/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package validatingwebhook

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// code in this file is borrowed from https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09,
// modified as per our need

// GenerateCertificates generates tls certificate files with the path specified
func GenerateCertificates(certFilePath, privateKeyPath string) error {
	// ip address of the machine would be required to be added as
	// subject alternate name
	ipAddr, err := GetIP()
	if err != nil {
		return err
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"acme.org"},
			Country:      []string{"IN"},
		},
		IPAddresses:           []net.IP{ipAddr},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 1),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate, err: %s", err.Error())
	}

	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile, err := os.Create(certFilePath)
	if err != nil {
		return err
	}
	defer certFile.Close()
	certFile.Write(out.Bytes())

	out.Reset()

	pem.Encode(out, pemBlockForKey(priv))
	privKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privKeyFile.Close()
	privKeyFile.Write(out.Bytes())

	return nil
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	default:
		return nil
	}
}
