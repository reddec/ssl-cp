package db

import (
	"time"

	"gorm.io/gorm"
)

type Certificate struct {
	gorm.Model
	IssuerID    *uint
	Issued      []Certificate `gorm:"foreignkey:IssuerID"`
	Name        string
	Serial      string
	CA          bool
	Description string
	ExpireAt    time.Time
	Key         string
	Cert        string
	Domains     []Domain
	Units       []Unit
}

func (crt *Certificate) RefID() *uint {
	if crt == nil {
		return nil
	}
	return &crt.ID
}

func (crt *Certificate) Issuer() uint {
	if crt.IssuerID == nil {
		return 0
	}
	return *crt.IssuerID
}

type Domain struct {
	CertificateID uint   `gorm:"primaryKey"`
	Name          string `gorm:"primaryKey"`
}

type Unit struct {
	CertificateID uint   `gorm:"primaryKey"`
	Name          string `gorm:"primaryKey"`
}
