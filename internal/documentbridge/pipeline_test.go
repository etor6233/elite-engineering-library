package documentbridge

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentIntakeBoundsAndFixtureIsolation(t *testing.T) {
	b := []byte("%PDF-1.7 fixture")
	o := Original{Name: "invoice.pdf", SHA256: Hash(b), Bytes: b}
	if o.Validate("PROVIDER") != nil {
		t.Fatal("reference structural PDF rejected")
	}
	if o.Validate("FIXTURE") == nil {
		t.Fatal("arbitrary bytes reached public fixture lane")
	}
	for _, name := range []string{"../invoice.pdf", "C:invoice.pdf", "folder\\invoice.pdf", "invoice.jpg", "invoice.pdf ", strings.Repeat("a", 129) + ".pdf"} {
		v := o
		v.Name = name
		if v.Validate("PROVIDER") == nil {
			t.Fatal("unbounded filename or mismatched type", name)
		}
	}
	o.SHA256 = strings.Repeat("0", 64)
	if o.Validate("PROVIDER") == nil {
		t.Fatal("unbound bytes")
	}
	if FixtureSecurity(context.Background(), "absent", filepath.Join(t.TempDir(), "receipt")) == nil {
		t.Fatal("fixture security accepted missing source")
	}
	if _, e := (SecurityCommand{Python: os.Args[0], PythonSHA: strings.Repeat("0", 64)}).Stage(); e == nil {
		t.Fatal("unbound security process selected")
	}
}
func FuzzDocumentOriginalBoundary(f *testing.F) {
	f.Add("invoice.pdf", []byte("%PDF-1.7 reference"))
	f.Add("../invoice.jpg", []byte{255, 216, 255, 224})
	f.Add("file.zip", []byte("PK"))
	f.Fuzz(func(t *testing.T, name string, b []byte) {
		if len(b) > MaxBytes+1 {
			return
		}
		o := Original{Name: name, SHA256: Hash(b), Bytes: b}
		e := o.Validate("PROVIDER")
		if e == nil {
			if len(b) == 0 || len(b) > MaxBytes || strings.ContainsAny(name, "/\\:") {
				t.Fatal("accepted path or byte-budget violation")
			}
			mutated := append([]byte(nil), b...)
			mutated[0] ^= 1
			o.Bytes = mutated
			if o.Validate("PROVIDER") == nil {
				t.Fatal("content mutation preserved authority")
			}
		}
	})
}
