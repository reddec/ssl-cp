package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/db"
	"gorm.io/gorm"
)

const (
	defaultKeySize = 2048 // NIST recommendation
	serialNumSize  = 128
)

var (
	ErrInvalidName    = errors.New("subject name invalid")
	ErrNotIssuerNotCA = errors.New("issuer certificate is not CA")
	ErrKeyNotExposed  = errors.New("exposition CA key disabled by policy")
	ErrInvalidSerial  = errors.New("invalid serial number")
)

type Option func(impl *Service)

// Random source for key generation. By default: crypto/rand.
func Random(reader io.Reader) Option {
	return func(impl *Service) {
		impl.cryptoSource = reader
	}
}

// Key size in bits. Default 4096. Higher number is safer but much longer to generate and process.
func Key(bits int) Option {
	return func(impl *Service) {
		impl.keySize = bits
	}
}

// ExposeCAKey set policy to about retrieval private CA key. By default - disabled.
// There are no real reasons to expose CA private key except you want to sign certificates by yourself.
func ExposeCAKey(allow bool) Option {
	return func(impl *Service) {
		impl.exposeCa = allow
	}
}

// New implementation of certificate management API.
// DB - initialized connection to database.
func New(db *gorm.DB, secret string, options ...Option) (api.API, error) {
	tokenData := sha256.Sum256([]byte(secret)) // normalize to 32 bytes

	crypter, err := aes.NewCipher(tokenData[:])
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	srv := &Service{
		encryptor:    crypter,
		db:           db,
		cryptoSource: rand.Reader,
		keySize:      defaultKeySize,
		exposeCa:     false,
		soon:         time.Hour * 24 * 30,
	}
	for _, opt := range options {
		opt(srv)
	}
	return srv, nil
}

type Service struct {
	encryptor    cipher.Block
	db           *gorm.DB
	cryptoSource io.Reader
	keySize      int
	exposeCa     bool
	soon         time.Duration
}

func (srv *Service) TX(db *gorm.DB) *Service {
	cp := *srv
	cp.db = db
	return &cp
}

func (srv *Service) BatchCreateCertificate(ctx context.Context, list []api.Batch) ([]api.Certificate, error) {
	var ans = make([]api.Certificate, 0, len(list))
	for _, batch := range list {
		batch.Certificate.Ca = batch.Certificate.Ca || len(batch.Nested) > 0 || batch.Certificate.Issuer == 0
		cert, err := srv.CreateCertificate(ctx, batch.Certificate)
		if err != nil {
			return nil, fmt.Errorf("create certificate in batch: %w", err)
		}
		ans = append(ans, cert)
		for i := range batch.Nested {
			batch.Nested[i].Certificate.Issuer = cert.Id
		}
		result, err := srv.BatchCreateCertificate(ctx, batch.Nested)
		if err != nil {
			return nil, fmt.Errorf("create nested certificate in batch: %w", err)
		}
		ans = append(ans, result...)
	}
	return ans, nil
}

func (srv *Service) GetRevokedCertificatesList(ctx context.Context, certificateId uint) (string, error) {
	caKey, caCert, err := srv.getCA(ctx, certificateId)
	if err != nil {
		return "", fmt.Errorf("get CA: %w", err)
	}

	var list []db.Certificate
	err = srv.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL AND issuer_id = ?", certificateId).Find(&list).Error
	if err != nil {
		return "", fmt.Errorf("generate list revoked certificates: %w", err)
	}

	var revoked = make([]pkix.RevokedCertificate, len(list))

	for i, cert := range list {
		serial, ok := new(big.Int).SetString(cert.Serial, 10)
		if !ok {
			return "", fmt.Errorf("parse serial number for certificate %d: %w", cert.ID, ErrInvalidSerial)
		}
		revoked[i] = pkix.RevokedCertificate{
			SerialNumber:   serial,
			RevocationTime: cert.DeletedAt.Time,
		}
	}

	raw, err := x509.CreateRevocationList(srv.cryptoSource, &x509.RevocationList{
		RevokedCertificates: revoked,
		Number:              big.NewInt(time.Now().UnixNano()), // it is violation of x509 monotonic increasing number but in real scenario it is OK
	}, caCert, caKey)
	if err != nil {
		return "", fmt.Errorf("create CRL list: %w", err)
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "X509 CRL",
		Bytes: raw,
	})), nil
}

func (srv *Service) ListRevokedCertificates(ctx context.Context, certificateId uint) ([]api.Certificate, error) {
	var list []db.Certificate
	err := srv.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL AND issuer_id = ?", certificateId).Preload("Domains").Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return make([]api.Certificate, 0), nil
	}
	if err != nil {
		return nil, fmt.Errorf("list revoked certificates: %w", err)
	}
	ans := make([]api.Certificate, len(list))
	for i, c := range list {
		ans[i] = mapCert(&c)
	}
	return ans, nil
}

func (srv *Service) RevokeCertificate(ctx context.Context, certificateId uint) error {
	return srv.db.WithContext(ctx).Delete(&db.Certificate{}, certificateId).Error
}

func (srv *Service) CreateCertificate(ctx context.Context, subject api.Subject) (api.Certificate, error) {
	log.Println("creating certificate", subject.Name)

	// check maybe we already have - generate key quite slow
	var savedCert = db.Certificate{IssuerID: nullId(subject.Issuer), Name: subject.Name}
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Where(&savedCert).Take(&savedCert).Error; err == nil {
		return mapCert(&savedCert), nil
	}

	var dbCA db.Certificate

	err := srv.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		dbCA = db.Certificate{
			Name: subject.Name,
			CA:   subject.Ca,
		}

		if err := tx.Create(&dbCA).Error; err != nil {
			return fmt.Errorf("create draft certificate: %w", err)
		}

		cert, err := srv.TX(tx).generateSignedCertificate(ctx, subject, dbCA.ID)
		if err != nil {
			return fmt.Errorf("generate certificate: %w", err)
		}

		subject = cert.Subject

		dbCA = db.Certificate{
			Model:    dbCA.Model,
			IssuerID: cert.Issuer.RefID(),
			Name:     subject.Name,
			Serial:   cert.Serial,
			CA:       subject.Ca,
			ExpireAt: cert.ExpiredAt,
			Key:      cert.PrivateEncryptedPEM,
			Cert:     cert.CertPEM,
		}

		for _, domain := range subject.Domains {
			dbCA.Domains = append(dbCA.Domains, db.Domain{Name: strings.TrimSpace(domain)})
		}

		for _, ip := range subject.Ips {
			dbCA.Addresses = append(dbCA.Addresses, db.Address{IP: strings.TrimSpace(ip)})
		}

		for _, field := range subject.Units {
			dbCA.Units = append(dbCA.Units, db.Unit{Name: field})
		}

		if err := tx.Save(&dbCA).Error; err != nil {
			return fmt.Errorf("save CA in DB: %w", err)
		}

		return nil
	})
	if err != nil {
		return api.Certificate{}, fmt.Errorf("create certificate: %w", err)
	}
	return mapCert(&dbCA), nil
}

func (srv *Service) validateIssuer(issuer uint) (*db.Certificate, error) {
	if issuer == 0 {
		return nil, nil
	}
	var ca db.Certificate
	err := srv.db.Model(&db.Certificate{}).Take(&ca, issuer).Error
	if err != nil {
		return nil, fmt.Errorf("get issuer certificate: %w", err)
	}
	if !ca.CA {
		return nil, ErrNotIssuerNotCA
	}
	ca.Key, err = srv.decrypt(ca.Key)
	if err != nil {
		return nil, fmt.Errorf("decrypt issuer key: %w", err)
	}
	return &ca, nil
}

func (srv *Service) ListRootCertificates(ctx context.Context) ([]api.Certificate, error) {
	return srv.ListCertificates(ctx, 0)
}

func (srv *Service) GetCertificate(ctx context.Context, certificateId uint) (api.Certificate, error) {
	var ca db.Certificate
	err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Preload("Units").Preload("Addresses").Preload("Domains").Take(&ca, certificateId).Error
	if err != nil {
		return api.Certificate{}, fmt.Errorf("get certificate %d: %w", certificateId, err)
	}
	return mapCert(&ca), nil
}

func (srv *Service) GetPublicCert(ctx context.Context, certificateId uint) (string, error) {
	var ca db.Certificate
	err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Select("cert").Take(&ca, certificateId).Error
	if err != nil {
		return "", fmt.Errorf("get certificate key: %w", err)
	}
	return ca.Cert, nil
}

func (srv *Service) GetPrivateKey(ctx context.Context, certificateId uint) (string, error) {
	var cert db.Certificate
	err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Select("key", "ca").Take(&cert, certificateId).Error
	if err != nil {
		return "", fmt.Errorf("get certificate key: %w", err)
	}
	if cert.CA && !srv.exposeCa {
		return "", ErrKeyNotExposed
	}
	return srv.decrypt(cert.Key)
}

func (srv *Service) ListCertificates(ctx context.Context, certificateId uint) ([]api.Certificate, error) {
	var nodes []db.Certificate
	req := srv.db.WithContext(ctx).Model(&db.Certificate{})
	if certificateId == 0 {
		// root
		req = req.Where("issuer_id IS NULL")
	} else {
		req = req.Where("issuer_id = ?", certificateId)
	}
	err := req.Preload("Units").Preload("Addresses").Preload("Domains").Find(&nodes).Error
	if err != nil {
		return nil, fmt.Errorf("list certificates: %w", err)
	}
	var ans = make([]api.Certificate, len(nodes))
	for i, n := range nodes {
		ans[i] = mapCert(&n)
	}
	return ans, nil
}

func (srv *Service) RenewCertificate(ctx context.Context, certificateId uint, renewal api.Renewal) (api.Certificate, error) {
	var old db.Certificate
	err := srv.db.WithContext(ctx).Preload("Units").Preload("Addresses").Preload("Domains").Preload("Issued").Model(&old).Take(&old, certificateId).Error
	if err != nil {
		return api.Certificate{}, fmt.Errorf("find certificate: %w", err)
	}

	cert, err := srv.generateSignedCertificate(ctx, api.Subject{
		Name:    old.Name,
		Issuer:  old.Issuer(),
		Days:    renewal.Days,
		Ca:      old.CA,
		Domains: renewal.Domains,
		Units:   renewal.Units,
		Ips:     renewal.Ips,
	}, old.ID)
	if err != nil {
		return api.Certificate{}, fmt.Errorf("regenerate certificate: %w", err)
	}

	var dbDomains = make([]db.Domain, 0, len(cert.Subject.Domains))
	for _, d := range cert.Subject.Domains {
		dbDomains = append(dbDomains, db.Domain{
			CertificateID: old.ID,
			Name:          d,
		})
	}

	var dbFields = make([]db.Unit, 0, len(cert.Subject.Units))
	for _, f := range cert.Subject.Units {
		dbFields = append(dbFields, db.Unit{
			CertificateID: old.ID,
			Name:          f,
		})
	}

	var dbIPs = make([]db.Address, 0, len(cert.Addresses))
	for _, f := range cert.Addresses {
		dbIPs = append(dbIPs, db.Address{
			CertificateID: old.ID,
			IP:            f.String(),
		})
	}

	err = srv.db.Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Model(&db.Domain{}).Where("certificate_id = ?", old.ID).Delete(db.Domain{}).Error
		if err != nil {
			return fmt.Errorf("remove old domains: %w", err)
		}

		err = tx.WithContext(ctx).Model(&db.Unit{}).Where("certificate_id = ?", old.ID).Delete(db.Unit{}).Error
		if err != nil {
			return fmt.Errorf("remove old fields: %w", err)
		}

		err = tx.WithContext(ctx).Model(&db.Address{}).Where("certificate_id = ?", old.ID).Delete(db.Address{}).Error
		if err != nil {
			return fmt.Errorf("remove old IP: %w", err)
		}

		err = tx.WithContext(ctx).Model(&db.Domain{}).Create(dbDomains).Error
		if err != nil {
			return fmt.Errorf("create new domains: %w", err)
		}

		if len(dbFields) > 0 {
			err = tx.WithContext(ctx).Model(&db.Unit{}).Create(dbFields).Error
			if err != nil {
				return fmt.Errorf("create new units: %w", err)
			}
		}

		if len(dbIPs) > 0 {
			err = tx.WithContext(ctx).Model(&db.Address{}).Create(dbIPs).Error
			if err != nil {
				return fmt.Errorf("create new IP: %w", err)
			}
		}
		err = tx.WithContext(ctx).Model(&db.Certificate{}).
			Where("id = ?", certificateId).
			Updates(&db.Certificate{
				ExpireAt: cert.ExpiredAt,
				Cert:     cert.CertPEM,
				Key:      cert.PrivateEncryptedPEM,
			}).Error
		if err != nil {
			return fmt.Errorf("update cert: %w", err)
		}
		err = tx.WithContext(ctx).Model(&db.Certificate{}).Preload("Units").Preload("Addresses").Preload("Domains").Take(&old, certificateId).Error
		if err != nil {
			return fmt.Errorf("get cert: %w", err)
		}
		return nil
	})
	if err != nil {
		return api.Certificate{}, fmt.Errorf("save renewed certificate: %w", err)
	}

	for _, issued := range old.Issued {
		if issued.DeletedAt.Valid || time.Now().After(issued.ExpireAt) {
			continue
		}
		days := uint(math.Ceil(issued.ExpireAt.Sub(issued.UpdatedAt).Hours() / 24))
		_, err = srv.RenewCertificate(ctx, issued.ID, api.Renewal{Days: days})
		if err != nil {
			return api.Certificate{}, fmt.Errorf("renew nested certificate %d: %w", issued.ID, err)
		}
	}
	return mapCert(&old), nil
}

func (srv *Service) ListExpiredCertificates(ctx context.Context) ([]api.Certificate, error) {
	var nodes []db.Certificate
	err := srv.db.WithContext(ctx).
		Model(&db.Certificate{}).
		Where("expire_at <= ?", time.Now()).
		Preload("Domains").
		Preload("Addresses").
		Preload("Units").
		Find(&nodes).
		Error
	if err != nil {
		return nil, fmt.Errorf("list expired certificates: %w", err)
	}
	var ans = make([]api.Certificate, len(nodes))
	for i, n := range nodes {
		ans[i] = mapCert(&n)
	}
	return ans, nil
}

func (srv *Service) ListSoonExpireCertificates(ctx context.Context) ([]api.Certificate, error) {
	var nodes []db.Certificate
	now := time.Now()
	err := srv.db.WithContext(ctx).
		Model(&db.Certificate{}).
		Where("expire_at > ? AND expire_at <= ?", now, now.Add(srv.soon)).
		Preload("Domains").
		Preload("Addresses").
		Preload("Units").
		Find(&nodes).
		Error
	if err != nil {
		return nil, fmt.Errorf("list soon expire certificates: %w", err)
	}
	var ans = make([]api.Certificate, len(nodes))
	for i, n := range nodes {
		ans[i] = mapCert(&n)
	}
	return ans, nil
}

func (srv *Service) generateSerial(bits int) (*big.Int, error) {
	var payload = make([]byte, bits/8)
	_, err := io.ReadFull(srv.cryptoSource, payload)
	if err != nil {
		return nil, fmt.Errorf("read random source: %w", err)
	}
	v := new(big.Int)
	v.SetBytes(payload)
	return v, nil
}

func (srv *Service) getCA(ctx context.Context, id uint) (*rsa.PrivateKey, *x509.Certificate, error) {
	var ca db.Certificate
	err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Select("key", "ca", "cert").Take(&ca, id).Error
	if err != nil {
		return nil, nil, fmt.Errorf("get issuer key: %w", err)
	}

	if !ca.CA {
		return nil, nil, ErrNotIssuerNotCA
	}

	decrypted, err := srv.decrypt(ca.Key)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt issuer key: %w", err)
	}

	block, _ := pem.Decode([]byte(decrypted))

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse private key: %w", err)
	}

	block, _ = pem.Decode([]byte(ca.Cert))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse public cert: %w", err)
	}

	return priv, cert, nil
}

func mapCert(header *db.Certificate) api.Certificate {
	n := api.Certificate{
		Id:        header.ID,
		Issuer:    header.Issuer(),
		Serial:    header.Serial,
		Ca:        header.CA,
		Name:      header.Name,
		CreatedAt: header.CreatedAt,
		ExpireAt:  header.ExpireAt,
		UpdatedAt: header.UpdatedAt,
	}
	if header.DeletedAt.Valid {
		n.RevokedAt = header.DeletedAt.Time
	}
	for _, domain := range header.Domains {
		n.Domains = append(n.Domains, domain.Name)
	}
	for _, field := range header.Units {
		n.Units = append(n.Units, field.Name)
	}
	for _, addr := range header.Addresses {
		n.Ips = append(n.Ips, addr.IP)
	}
	return n
}
