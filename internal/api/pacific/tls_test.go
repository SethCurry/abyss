package pacific

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestMutualTLS(t *testing.T) {
	certs, err := GenerateCertificates()
	if err != nil {
		t.Fatalf("GenerateCertificates: %v", err)
	}

	serverTLSConfig, err := certs.ServerTLSConfig()
	if err != nil {
		t.Fatalf("ServerTLSConfig: %v", err)
	}

	srv := NewServer(func(w http.ResponseWriter, r *http.Request) struct{} {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return struct{}{}
	})
	srv.AddRoute("GET", "/", func(struct{}) {})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	server := &http.Server{
		Handler:   srv.router,
		TLSConfig: serverTLSConfig,
	}

	go func() {
		_ = server.ServeTLS(listener, "", "")
	}()
	defer server.Close()

	clientTLSConfig, err := certs.ClientTLSConfig()
	if err != nil {
		t.Fatalf("ClientTLSConfig: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: clientTLSConfig,
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://" + listener.Addr().String() + "/")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(body) != "ok" {
		t.Fatalf("body = %q, want %q", body, "ok")
	}
}

func TestMutualTLSRejectsUntrustedClient(t *testing.T) {
	certs, err := GenerateCertificates()
	if err != nil {
		t.Fatalf("GenerateCertificates: %v", err)
	}

	serverTLSConfig, err := certs.ServerTLSConfig()
	if err != nil {
		t.Fatalf("ServerTLSConfig: %v", err)
	}

	srv := NewServer(func(w http.ResponseWriter, r *http.Request) struct{} {
		return struct{}{}
	})
	srv.AddRoute("GET", "/", func(struct{}) {})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	server := &http.Server{
		Handler:   srv.router,
		TLSConfig: serverTLSConfig,
	}

	go func() {
		_ = server.ServeTLS(listener, "", "")
	}()
	defer server.Close()

	// A client that trusts the CA but presents no client certificate must be
	// rejected by the server.
	pool := serverTLSConfig.ClientCAs
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:    pool,
				MinVersion: tls.VersionTLS13,
			},
		},
		Timeout: 5 * time.Second,
	}

	_, err = client.Get("https://" + listener.Addr().String() + "/")
	if err == nil {
		t.Fatal("expected connection to fail without a client certificate")
	}
}
