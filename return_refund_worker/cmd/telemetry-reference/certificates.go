// AUTHORED ephemeral reference fixture; not a certificate authority service.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func referenceCertificates(directory string) error {
	info, e := os.Lstat(directory)
	if e != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("REFERENCE_CERTIFICATE_REJECTED")
	}
	write := func(name string, data []byte) error {
		f, e := os.OpenFile(filepath.Join(directory, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if e != nil {
			return e
		}
		_, e = f.Write(data)
		if e == nil {
			e = f.Sync()
		}
		ce := f.Close()
		if e != nil {
			return e
		}
		return ce
	}
	serial := func() (*big.Int, error) { return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128)) }
	pub, key, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		return e
	}
	id := sha256.Sum256(pub)
	sn, e := serial()
	if e != nil {
		return e
	}
	now := time.Now()
	ca := &x509.Certificate{SerialNumber: sn, Subject: pkix.Name{CommonName: "Elite ephemeral reference CA"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(12 * time.Hour), IsCA: true, BasicConstraintsValid: true, MaxPathLen: 0, MaxPathLenZero: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign, SubjectKeyId: id[:20]}
	der, e := x509.CreateCertificate(rand.Reader, ca, ca, pub, key)
	if e != nil {
		return e
	}
	if e = write("ca.pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})); e != nil {
		return e
	}
	for _, role := range []string{"server", "client", "rogue"} {
		lp, lk, e := ed25519.GenerateKey(rand.Reader)
		if e != nil {
			return e
		}
		sn, e = serial()
		if e != nil {
			return e
		}
		lid := sha256.Sum256(lp)
		leaf := &x509.Certificate{SerialNumber: sn, Subject: pkix.Name{CommonName: "elite-reference-" + role}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(2 * time.Hour), BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature, SubjectKeyId: lid[:20], AuthorityKeyId: ca.SubjectKeyId, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}}
		if role == "server" {
			leaf.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
			leaf.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
		}
		signer := key
		if role == "rogue" {
			signer = lk
		}
		// The rogue fixture is a valid self-issued leaf whose issuer name resembles
		// the real CA, but whose signing key is unrelated. It must never authenticate.
		parent := ca
		if role == "rogue" {
			copyCA := *ca
			copyCA.PublicKey = lp
			parent = &copyCA
		}
		der, e = x509.CreateCertificate(rand.Reader, leaf, parent, lp, signer)
		if e != nil {
			return e
		}
		if e = write(role+".pem", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})); e != nil {
			return e
		}
		kb, e := x509.MarshalPKCS8PrivateKey(lk)
		if e != nil {
			return e
		}
		if e = write(role+".key", pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: kb})); e != nil {
			return e
		}
	}
	return nil
}
