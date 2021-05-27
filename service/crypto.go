package service

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/db"
)

var ErrEncryptedTextTooSmall = errors.New("encrypted text too small")

func (srv *Service) encrypt(text string) (string, error) {
	data := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("encrypt: read iv: %w", err)
	}
	cfb := cipher.NewCFBEncrypter(srv.encryptor, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (srv *Service) decrypt(text string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", fmt.Errorf("decrypt: decode b64 value: %w", err)
	}

	if len(data) < aes.BlockSize {
		return "", ErrEncryptedTextTooSmall
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(srv.encryptor, iv)
	cfb.XORKeyStream(data, data)
	return string(data), nil
}

type certificate struct {
	Subject             api.Subject
	ExpiredAt           time.Time
	Issuer              *db.Certificate
	Serial              string
	CertPEM             string
	PrivateEncryptedPEM string
}

func (srv *Service) generateSignedCertificate(ctx context.Context, subject api.Subject, id uint) (*certificate, error) {
	subject.Name = strings.TrimSpace(subject.Name)
	if len(subject.Name) == 0 {
		return nil, ErrInvalidName
	}
	serial, err := srv.generateSerial(serialNumSize)
	if err != nil {
		return nil, fmt.Errorf("generate serial: %w", err)
	}

	issuer, err := srv.validateIssuer(subject.Issuer)
	if err != nil {
		return nil, fmt.Errorf("validate issuer: %w", err)
	}

	if len(subject.Domains) == 0 {
		subject.Domains = append(subject.Domains, subject.Name)
	}

	if issuer == nil {
		subject.Ca = true
	}

	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:         subject.Name,
			OrganizationalUnit: []string{strconv.FormatUint(uint64(id), 10)},
			SerialNumber:       serial.String(),
			ExtraNames:         []pkix.AttributeTypeAndValue{},
		},
		DNSNames:              subject.Domains,
		IsCA:                  subject.Ca,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, int(subject.Days)),
		BasicConstraintsValid: true,
	}

	if subject.Ca {
		cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	} else {
		cert.KeyUsage = x509.KeyUsageDigitalSignature
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	}

	certKey, err := rsa.GenerateKey(srv.cryptoSource, srv.keySize)
	if err != nil {
		return nil, fmt.Errorf("generate private key: %w", err)
	}

	var caCert *x509.Certificate
	var caKey interface{}
	if subject.Issuer == 0 {
		caCert = cert
		caKey = certKey
	} else if issuerKey, issuerCert, err := srv.getCA(ctx, subject.Issuer); err == nil {
		caCert = issuerCert
		caKey = issuerKey
	} else {
		return nil, fmt.Errorf("get CA: %w", err)
	}

	signedCert, err := x509.CreateCertificate(srv.cryptoSource, cert, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("generate certificate: %w", err)
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signedCert,
	})
	if err != nil {
		return nil, fmt.Errorf("encode certificate to PEM: %w", err)
	}

	keyPEM := new(bytes.Buffer)
	err = pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})
	if err != nil {
		return nil, fmt.Errorf("encode key to PEM: %w", err)
	}

	cryptedKey, err := srv.encrypt(keyPEM.String())
	if err != nil {
		return nil, fmt.Errorf("encrypt key: %w", err)
	}

	return &certificate{
		Subject:             subject,
		ExpiredAt:           cert.NotAfter,
		Issuer:              issuer,
		Serial:              serial.String(),
		CertPEM:             certPEM.String(),
		PrivateEncryptedPEM: cryptedKey,
	}, nil
}
