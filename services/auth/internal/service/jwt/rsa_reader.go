package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func (j *jwtImpl) readRSAPublicKeyFromString(rsaPublic string) (*rsa.PublicKey, error) {
	data, _ := pem.Decode([]byte(rsaPublic))
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("readRSAPublicKeyFromString: imposible type assertion")
	}

	return publicKey, nil
}

func (j *jwtImpl) readRSAPrivateKeyFromString(rsaPrivate string) (*rsa.PrivateKey, error) {
	data, _ := pem.Decode([]byte(rsaPrivate))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKeyImported, nil
}
