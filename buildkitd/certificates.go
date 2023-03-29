package buildkitd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type certificateData struct {
	Key       *rsa.PrivateKey
	CertBytes []byte
	Cert      *x509.Certificate
}

const (
	buildkit = "buildkit"
	earthly  = "earthly"
)

// GenerateCertificates creates and saves a CA and certificates for both sides of an mTLS TCP connection.
func GenerateCertificates(dir, hostname string) error {
	ca, err := createAndSaveCA(dir)
	if err != nil {
		return errors.Wrap(err, "create CA")
	}

	err = createAndSaveCertificate(ca, buildkit, dir, hostname)
	if err != nil {
		return errors.Wrap(err, "create buildkit certificate")
	}

	err = createAndSaveCertificate(ca, earthly, dir, hostname)
	if err != nil {
		return errors.Wrap(err, "create earthly certificate")
	}

	return nil
}

func createAndSaveCertificate(ca *certificateData, role, dir, hostname string) error {
	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return errors.Wrapf(err, "generate %s key", role)
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return errors.Wrapf(err, "generate %s serial", role)
	}

	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{fmt.Sprintf("Earthly GRPC: %v side", role)},
		},
		DNSNames:     []string{hostname},
		IPAddresses:  []net.IP{net.IPv6loopback, net.ParseIP("127.0.0.1")},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte(role),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &certKey.PublicKey, ca.Key)
	if err != nil {
		return errors.Wrapf(err, "generate %s certificate", role)
	}

	err = saveCertAndKeyAsPEM(role, dir, certBytes, certKey)
	if err != nil {
		return errors.Wrap(err, "save role")
	}

	return nil
}

func createAndSaveCA(dir string) (*certificateData, error) {
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization: []string{"Earthly Buildkit GRPC CA"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	err = saveCertAndKeyAsPEM("ca", dir, caBytes, caKey)
	if err != nil {
		return nil, errors.Wrap(err, "save CA")
	}

	return &certificateData{
		Key:       caKey,
		CertBytes: caBytes,
		Cert:      ca,
	}, nil
}

func saveCertAndKeyAsPEM(prefix, dir string, certBytes []byte, key *rsa.PrivateKey) error {
	certFile := filepath.Join(dir, fmt.Sprintf("%s_cert.pem", prefix))
	err := saveAsPEM(certFile, "CERTIFICATE", certBytes)
	if err != nil {
		return errors.Wrapf(err, "failed to save certificate %s", certFile)
	}

	keyFile := filepath.Join(dir, fmt.Sprintf("%s_key.pem", prefix))
	err = saveAsPEM(keyFile, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key))
	if err != nil {
		return errors.Wrapf(err, "failed to save key %s", keyFile)
	}

	return nil
}

func saveAsPEM(fn, typ string, bytes []byte) error {
	err := os.MkdirAll(path.Dir(fn), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	err = pem.Encode(f, &pem.Block{
		Type:  typ,
		Bytes: bytes,
	})
	if err != nil {
		return err
	}

	err = f.Chmod(0444)
	if err != nil {
		return err
	}

	return nil
}
