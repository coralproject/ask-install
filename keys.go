package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// DERtoPEM encodes a DER encoded key into a PEM format with the specified
// header.
func DERtoPEM(header string, data []byte) ([]byte, error) {

	// Create the PEM block.
	block := pem.Block{
		Type:  header,
		Bytes: data,
	}

	var buf bytes.Buffer

	if err := pem.Encode(&buf, &block); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateKeys will output the private key and public key of a newly created
// ecdsa keypair that has been encoded into PEM format + base64'd.
func GenerateKeys() (string, string, error) {

	// Generate the elliptic curve certificate.
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// Marshal the certificate into a ASN.1, DER format.
	privDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	// Encode the DER format into PEM.
	privPEM, err := DERtoPEM("EC PRIVATE KEY", privDER)
	if err != nil {
		return "", "", err
	}

	// Marshal the public key into a DER-encoded PKIX format.
	pubDER, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return "", "", err
	}

	// Encode the DER format into PEM.
	pubPEM, err := DERtoPEM("PUBLIC KEY", pubDER)
	if err != nil {
		return "", "", err
	}

	// Encode the PEM certificates into base64 strings.
	return base64.StdEncoding.EncodeToString(privPEM), base64.StdEncoding.EncodeToString(pubPEM), nil
}
