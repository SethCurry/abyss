package runenv

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func listTar(t *testing.T, b []byte) {
	tr := tar.NewReader(bytes.NewReader(b))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("entry: %q type=%s mode=%s size=%d", hdr.Name, typeString(hdr.Typeflag), os.FileMode(hdr.Mode).String(), hdr.Size)
	}
}

func typeString(b byte) string {
	switch b {
	case tar.TypeDir:
		return "dir"
	case tar.TypeReg:
		return "reg"
	case tar.TypeSymlink:
		return "symlink"
	default:
		return fmt.Sprintf("0x%x", b)
	}
}

func TestBuildTarNested(t *testing.T) {
	root := t.TempDir()
	// root/top.txt
	if err := os.WriteFile(filepath.Join(root, "top.txt"), []byte("top"), 0644); err != nil {
		t.Fatal(err)
	}
	// root/sub/deep.txt and root/sub/nested.txt
	if err := os.MkdirAll(filepath.Join(root, "sub", "deepdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "nested.txt"), []byte("nested"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "deepdir", "a.txt"), []byte("a"), 0644); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if err := buildTar(&buf, root, info); err != nil {
		t.Fatal(err)
	}

	listTar(t, buf.Bytes())

	base := filepath.Base(root)
	want := map[string]string{
		base + "/":                       "dir",
		base + "/sub/":                   "dir",
		base + "/sub/deepdir/":            "dir",
		base + "/sub/deepdir/a.txt":       "reg",
		base + "/sub/nested.txt":          "reg",
		base + "/top.txt":                 "reg",
	}
	got := map[string]string{}
	tr := tar.NewReader(bytes.NewReader(buf.Bytes()))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		got[hdr.Name] = typeString(hdr.Typeflag)
	}
	for name, typ := range want {
		if g, ok := got[name]; !ok || g != typ {
			t.Errorf("entry %q: want %s, got %q (present=%v)", name, typ, g, ok)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Errorf("unexpected entry %q", name)
		}
	}
}