package service

import (
	"context"
	"fmt"
	"time"

	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/db"
)

func (srv *Service) GetStatus(ctx context.Context) (api.Status, error) {
	var total int64
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Unscoped().Count(&total).Error; err != nil {
		return api.Status{}, fmt.Errorf("coutn total number of certificates: %w", err)
	}

	now := time.Now()

	var expired int64
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Where("expire_at <= ?", now).Count(&expired).Error; err != nil {
		return api.Status{}, fmt.Errorf("coutn total number of expired certificates: %w", err)
	}

	var soonExpired int64
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Where("expire_at >= ? AND expire_at <= ?", now, now.Add(srv.soon)).Count(&soonExpired).Error; err != nil {
		return api.Status{}, fmt.Errorf("coutn total number of expired certificates: %w", err)
	}

	var ca int64
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Where("ca").Count(&ca).Error; err != nil {
		return api.Status{}, fmt.Errorf("coutn total number of CA: %w", err)
	}

	var revoked int64
	if err := srv.db.WithContext(ctx).Model(&db.Certificate{}).Unscoped().Where("deleted_at IS NOT NULL").Count(&revoked).Error; err != nil {
		return api.Status{}, fmt.Errorf("coutn total number of revoked: %w", err)
	}

	return api.Status{
		Total:      uint(total),
		Expired:    uint(expired),
		SoonExpire: uint(soonExpired),
		Ca:         uint(ca),
		Revoked:    uint(revoked),
	}, nil
}
