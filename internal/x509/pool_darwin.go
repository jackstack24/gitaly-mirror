package x509

import (
	"crypto/x509"
	"io/ioutil"
	"os"
	"path"
)

// SystemCertPool circumvents the fact that Go on macOS does not support
// SSL_CERT_{DIR,FILE}.
func SystemCertPool() (*x509.CertPool, error) {
	var certPem []byte

	if f := os.Getenv(SSLCertFile); len(f) > 0 {
		pem, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}

		pem = append(pem, '\n')
		certPem = append(certPem, pem...)
	}

	if d := os.Getenv(SSLCertDir); len(d) > 0 {
		entries, err := ioutil.ReadDir(d)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			pem, err := ioutil.ReadFile(path.Join(d, entry.Name()))
			if err != nil {
				return nil, err
			}

			pem = append(pem, '\n')
			certPem = append(certPem, pem...)
		}
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	pool.AppendCertsFromPEM(certPem)
	return pool, nil
}
