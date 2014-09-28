package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	. "modinputs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func generateSelfSignedCertificate(host string, rsaBits int, directory string, prefix string) (certPath string, keyPath string, retErr error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		retErr = fmt.Errorf("Failed to generate private key: %s", err)
		return
	}

	notBefore := time.Now().Add(-1 * 24 * time.Hour)
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		retErr = fmt.Errorf("Failed to generate serial number: %s", err)
		return
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		retErr = fmt.Errorf("Failed to create certificate: %s", err)
		return
	}

	certPath = filepath.Join(directory, fmt.Sprintf("%s_cert.pem", prefix))
	certOut, err := os.Create(certPath)
	if err != nil {
		retErr = fmt.Errorf("Failed to open %s for writing: %s", certPath, err)
		return
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyPath = filepath.Join(directory, fmt.Sprintf("%s_key.pem", prefix))
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		retErr = fmt.Errorf("failed to open key.pem for writing:", err)
		return
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	return
}
