package pacific

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// Certificates holds the PEM-encoded keys and certificates needed to
// establish a mutually authenticated TLS connection.
type Certificates struct {
	CACertPEM     []byte
	CAKeyPEM      []byte
	ServerCertPEM []byte
	ServerKeyPEM  []byte
	ClientCertPEM []byte
	ClientKeyPEM  []byte
}

// GenerateCertificates creates a self-signed CA along with server and client
// certificates signed by that CA, suitable for mutual TLS.
func GenerateCertificates() (*Certificates, error) {
	caCertPEM, caKeyPEM, err := generateCA()
	if err != nil {
		return nil, err
	}

	caCert, caKey, err := parseCA(caCertPEM, caKeyPEM)
	if err != nil {
		return nil, err
	}

	serverCertPEM, serverKeyPEM, err := generateSignedCert(caCert, caKey, x509.ExtKeyUsageServerAuth, []string{"localhost", "127.0.0.1"})
	if err != nil {
		return nil, err
	}

	clientCertPEM, clientKeyPEM, err := generateSignedCert(caCert, caKey, x509.ExtKeyUsageClientAuth, nil)
	if err != nil {
		return nil, err
	}

	return &Certificates{
		CACertPEM:     caCertPEM,
		CAKeyPEM:      caKeyPEM,
		ServerCertPEM: serverCertPEM,
		ServerKeyPEM:  serverKeyPEM,
		ClientCertPEM: clientCertPEM,
		ClientKeyPEM:  clientKeyPEM,
	}, nil
}

// ServerTLSConfig returns a tls.Config that presents the server certificate
// and requires clients to present a certificate signed by the CA.
func (c *Certificates) ServerTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(c.ServerCertPEM, c.ServerKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("load server key pair: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(c.CACertPEM) {
		return nil, fmt.Errorf("append CA certificate to pool")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// ClientTLSConfig returns a tls.Config that presents the client certificate
// and verifies the server against the CA.
func (c *Certificates) ClientTLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair(c.ClientCertPEM, c.ClientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("load client key pair: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(c.CACertPEM) {
		return nil, fmt.Errorf("append CA certificate to pool")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// generateCA creates a self-signed CA certificate and returns its PEM-encoded
// certificate and private key.
func generateCA() (certPEM, keyPEM []byte, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate CA key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber:          randomSerial(),
		Subject:               pkix.Name{CommonName: "abyss-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, fmt.Errorf("create CA certificate: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM, err = marshalECKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal CA key: %w", err)
	}

	return certPEM, keyPEM, nil
}

// generateSignedCert creates a certificate signed by the given CA for the
// provided usage and DNS names.
func generateSignedCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, usage x509.ExtKeyUsage, dnsNames []string) (certPEM, keyPEM []byte, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generate key: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: randomSerial(),
		Subject:      pkix.Name{CommonName: "abyss"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{usage},
		DNSNames:     dnsNames,
	}

	if usage == x509.ExtKeyUsageServerAuth {
		template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	}

	der, err := x509.CreateCertificate(rand.Reader, template, caCert, &key.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("create certificate: %w", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM, err = marshalECKey(key)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal key: %w", err)
	}

	return certPEM, keyPEM, nil
}

// parseCA decodes PEM-encoded CA certificate and key material.
func parseCA(certPEM, keyPEM []byte) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, nil, fmt.Errorf("decode CA certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA certificate: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("decode CA key PEM")
	}

	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA key: %w", err)
	}

	return cert, key, nil
}

// marshalECKey encodes an ECDSA private key as PEM.
func marshalECKey(key *ecdsa.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}), nil
}

// randomSerial returns a random positive serial number for a certificate.
func randomSerial() *big.Int {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		panic(fmt.Sprintf("generate serial number: %v", err))
	}
	return serial
}
