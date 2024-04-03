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
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/util/fileutil"
)

type certData struct {
	Key  *rsa.PrivateKey
	Cert *x509.Certificate
}

const (
	buildkit = "buildkit"
	earthly  = "earthly"

	typeCert   = "CERTIFICATE"
	typeRSAKey = "RSA PRIVATE KEY"
)

// GenCerts creates and saves a CA and certificates for both sides of an mTLS TCP connection.
func GenCerts(cfg config.Config, hostname string) error {
	caKey, err := parseTLSKey(cfg.Global.TLSCAKey)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "failed reading CA key")
	}
	if errors.Is(err, os.ErrNotExist) {
		all := []string{
			cfg.Global.TLSCACert,
			cfg.Global.ServerTLSCert,
			cfg.Global.ServerTLSKey,
			cfg.Global.ClientTLSCert,
			cfg.Global.ClientTLSKey,
		}
		var missing []string
		for _, f := range all {
			if exists, _ := fileutil.FileExists(f); !exists {
				missing = append(missing, f)
			}
		}
		switch len(missing) {
		case 0:
			return nil
		case len(all):
			key, err := createTLSKey(cfg.Global.TLSCAKey)
			if err != nil {
				return errors.Wrap(err, "could not create CA")
			}
			caKey = key
		default:
			found := all
			for _, m := range missing {
				for i, f := range found {
					if f == m {
						found = append(found[:i], found[i+1:]...)
						break
					}
				}
			}
			return hint.Wrap(errors.New("cannot generate missing certificates"),
				fmt.Sprintf("missing certificates: %v", missing),
				fmt.Sprintf("found certificates: %v", found),
				"you may want to stop earthly-buildkitd, delete your certificates, and run 'earthly bootstrap' to regenerate certificates",
			)
		}
	}
	caCert, err := parseTLSCert(cfg.Global.TLSCACert)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "could not parse CA certificate")
	}
	if errors.Is(err, os.ErrNotExist) {
		caCert, err = createCACert(caKey, cfg.Global.TLSCACert)
		if err != nil {
			return errors.Wrap(err, "could not create CA certificate")
		}
	}
	ca := &certData{
		Key:  caKey,
		Cert: caCert,
	}

	if err := genCert(ca, buildkit, hostname, cfg.Global.ServerTLSKey, cfg.Global.ServerTLSCert); err != nil {
		return errors.Wrapf(err, "could not generate server TLS key/cert pair for %v", buildkit)
	}
	if err := genCert(ca, earthly, hostname, cfg.Global.ClientTLSKey, cfg.Global.ClientTLSCert); err != nil {
		return errors.Wrapf(err, "could not generate client TLS key/cert pair for %v", earthly)
	}

	return nil
}

func genCert(ca *certData, role, hostname, keyPath, certPath string) error {
	certExists, _ := fileutil.FileExists(certPath)
	key, err := parseTLSKey(keyPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "could not parse %v TLS key", role)
	}
	if errors.Is(err, os.ErrNotExist) {
		if certExists {
			return errors.Wrapf(err, "refusing to generate TLS key %q: TLS cert %q exists", keyPath, certPath)
		}
		key, err = createTLSKey(keyPath)
		if err != nil {
			return errors.Wrapf(err, "could not create %v TLS key", role)
		}
	}
	if !certExists {
		if _, err := createTLSCert(ca, key, role, certPath, hostname); err != nil {
			return errors.Wrapf(err, "could not create %v TLS cert", role)
		}
	}
	return nil
}

func parseTLSKey(path string) (*rsa.PrivateKey, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read private key %q", path)
	}
	dec, _ := pem.Decode(body)
	key, err := x509.ParsePKCS1PrivateKey(dec.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode %q as RSA private key", path)
	}
	return key, nil
}

func createTLSKey(path string) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Wrapf(err, "could not generate RSA key")
	}
	if err := savePEM(path, typeRSAKey, x509.MarshalPKCS1PrivateKey(key)); err != nil {
		return nil, errors.Wrapf(err, "saving private key to %q failed", path)
	}
	return key, nil
}

func parseTLSCert(path string) (*x509.Certificate, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read public cert %q", path)
	}
	dec, _ := pem.Decode(body)
	cert, err := x509.ParseCertificate(dec.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode %q as x509 certificate", path)
	}
	return cert, nil
}

func createTLSCert(ca *certData, key *rsa.PrivateKey, role, path, hostname string) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, errors.Wrapf(err, "could not generate serial for role %q", role)
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

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Cert, &key.PublicKey, ca.Key)
	if err != nil {
		return nil, errors.Wrapf(err, "could not generate certificate for role %q", role)
	}

	if err := savePEM(path, typeCert, certBytes); err != nil {
		return nil, errors.Wrapf(err, "could not save certificate for role %q to path %q", role, path)
	}

	return cert, nil
}

func createCACert(key *rsa.PrivateKey, path string) (*x509.Certificate, error) {
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

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		return nil, errors.Wrap(err, "creating CA certificate failed")
	}

	if err := savePEM(path, typeCert, caBytes); err != nil {
		return nil, errors.Wrapf(err, "saving CA certificate to %q failed", path)
	}

	return ca, nil
}

func savePEM(path, typ string, bytes []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
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

	if err := f.Chmod(0444); err != nil {
		return err
	}

	return nil
}

// ConfigureSatelliteTLS uses the CA cert and key associate with the satellite
// to generate a new certificate/key pair for use in client-side mTLS.
// The certificates are configured in the settings for a new buildkit client.
func ConfigureSatelliteTLS(settings *Settings, sat *cloud.SatelliteInstance) (cleanupFn func(), err error) {
	dir := filepath.Join(os.TempDir(), "earthly", "certs", uuid.NewString())
	caRootPath := filepath.Join(dir, "root_ca_cert.pem")
	caCertPath := filepath.Join(dir, "ca_cert.pem")
	caKeyPath := filepath.Join(dir, "ca_key.pem")
	earthlyCertPath := filepath.Join(dir, "earthly_cert.pem")
	earthlyKeyPath := filepath.Join(dir, "earthly_key.pem")

	settings.TLSCA = caRootPath
	settings.ClientTLSCert = earthlyCertPath
	settings.ClientTLSKey = earthlyKeyPath
	settings.ServerTLSCert = ""
	settings.ServerTLSKey = ""

	if err = os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.Wrap(err, "could not make temp tls dir")
	}

	cleanupFn = func() { _ = os.RemoveAll(dir) }

	// TODO consider using concurrent goroutines to speed this up.

	if err = os.WriteFile(caRootPath, sat.Certificate.RootCa, 0444); err != nil {
		return nil, errors.Wrap(err, "failed saving ca cert")
	}

	defer func() {
		if err != nil {
			cleanupFn()
		}
	}()

	if err = os.WriteFile(caCertPath, sat.Certificate.ClientCa, 0444); err != nil {
		return nil, errors.Wrap(err, "failed saving ca cert")
	}

	if err = os.WriteFile(caKeyPath, sat.Certificate.ClientCaKey, 0444); err != nil {
		return nil, errors.Wrap(err, "failed saving ca key")
	}

	caCert, err := parseTLSCert(caCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing ca cert")
	}

	caKey, err := parseTLSKey(caKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing ca key")
	}

	ca := &certData{
		Cert: caCert,
		Key:  caKey,
	}

	if err = genCert(ca, earthly, sat.Address, earthlyKeyPath, earthlyCertPath); err != nil {
		return nil, errors.Wrap(err, "could not generate client TLS key/cert pair")
	}

	return cleanupFn, nil
}
