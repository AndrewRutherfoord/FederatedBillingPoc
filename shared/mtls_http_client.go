package shared

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// Creates the HTTP client with mTLS authentication configured. Used for communication between billing provider and CSP billing APIs
func NewMtlsHttpClient(mtlsKeyPath string, mtlsCertPath string, serverCertPaths []string) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(mtlsCertPath, mtlsKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server certificate: %v", err)
	}

	trustedPool := x509.NewCertPool()
	for _, path := range serverCertPaths {
		trustedCertPEM, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read trusted certificate: %w", err)
		}
		if !trustedPool.AppendCertsFromPEM(trustedCertPEM) {
			return nil, fmt.Errorf("failed to append trusted certificate")
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      trustedPool,
			},
		},
	}, nil
}
