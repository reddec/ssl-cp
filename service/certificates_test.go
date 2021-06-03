package service_test

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/db"
	"github.com/reddec/ssl-cp/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

type fastRandom struct{}

func (fr *fastRandom) Read(p []byte) (n int, err error) {
	return rand.Read(p)
}

func TestServiceImpl_CreateCertificate(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}
	t.Run("simple create should work", func(t *testing.T) {
		ca, err := srv.CreateCertificate(ctx, api.Subject{
			Name:    "example.com",
			Days:    365,
			Domains: []string{"abc.local", "xyz.local"},
			Ips:     []string{"127.0.0.1", "1.2.3.4"},
			Units:   []string{"john dow", "john@example.com"},
		})
		assert.NoError(t, err)
		assert.Equal(t, "example.com", ca.Name)
		assert.True(t, ca.Ca) // force CA
		assert.Empty(t, ca.Issuer)
		assert.NotEmpty(t, ca.CreatedAt)
		assert.Contains(t, ca.Domains, "abc.local")
		assert.Contains(t, ca.Domains, "xyz.local")
		assert.Contains(t, ca.Ips, "127.0.0.1")
		assert.Contains(t, ca.Ips, "1.2.3.4")
		assert.Contains(t, ca.Units, "john dow")
		assert.Contains(t, ca.Units, "john@example.com")
	})
	t.Run("multiple create should return same cert", func(t *testing.T) {
		ca, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
		})
		assert.NoError(t, err)

		secondary, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
		})
		assert.NoError(t, err)
		assert.Equal(t, ca, secondary)
	})
	t.Run("empty name not allowed", func(t *testing.T) {
		_, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "        ",
			Days: 365,
		})
		assert.Error(t, err)
	})
	t.Run("CA create should work", func(t *testing.T) {
		ca, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
			Ca:   true,
		})
		assert.NoError(t, err)
		assert.Equal(t, "example.com", ca.Name)
		assert.True(t, ca.Ca)
		assert.Empty(t, ca.Issuer)
		assert.NotEmpty(t, ca.CreatedAt)
	})

	t.Run("nested simple create should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
		})
		assert.NoError(t, err)

		ca, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)
		assert.Equal(t, "child", ca.Name)
		assert.False(t, ca.Ca)
		assert.Equal(t, root.Id, ca.Issuer)
		assert.NotEmpty(t, ca.CreatedAt)
	})

	t.Run("nested simple create should not work for parent non-CA", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
		})
		assert.NoError(t, err)

		child, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)

		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name:   "child.child",
			Issuer: child.Id,
			Days:   365,
		})
		assert.Error(t, err)
	})

	t.Run("nested simple create should not work for unknown parent", func(t *testing.T) {
		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name:   "root",
			Issuer: 100500,
			Days:   365,
		})
		assert.Error(t, err)
	})

	t.Run("nested CA create should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: t.Name(),
			Days: 365,
		})
		assert.NoError(t, err)

		ca, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child",
			Issuer: root.Id,
			Days:   365,
			Ca:     true,
		})
		assert.NoError(t, err)
		assert.Equal(t, "child", ca.Name)
		assert.True(t, ca.Ca)
		assert.Equal(t, root.Id, ca.Issuer)
		assert.NotEmpty(t, ca.CreatedAt)
	})

	t.Run("same name should work on different level", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example.com",
			Days: 365,
		})
		assert.NoError(t, err)

		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name:   "child",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)

		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name: "child",
			Days: 365,
		})
		assert.NoError(t, err)
	})
}

func TestService_ListCertificates(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	// root (ca)
	//   |- child.root
	//   |- ca.root (ca)
	//     |- child.ca.root
	// omega (ca)
	//   |- child.root

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "root",
		Days: 365,
	})
	assert.NoError(t, err)

	omega, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "omega",
		Days: 365,
	})
	assert.NoError(t, err)

	childRoot, err := srv.CreateCertificate(ctx, api.Subject{
		Name:   "child.root",
		Days:   365,
		Issuer: root.Id,
	})
	assert.NoError(t, err)

	caRoot, err := srv.CreateCertificate(ctx, api.Subject{
		Name:   "ca.root",
		Days:   365,
		Issuer: root.Id,
		Ca:     true,
	})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name:   "child.ca.root",
		Days:   365,
		Issuer: caRoot.Id,
	})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name:   "child.omega",
		Days:   365,
		Issuer: omega.Id,
	})
	assert.NoError(t, err)

	t.Run("list root certificates", func(t *testing.T) {
		list, err := srv.ListCertificates(ctx, 0)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.NotEmpty(t, findCert(root.Name, list))
		assert.NotEmpty(t, findCert(omega.Name, list))

		// check alias
		alias, err := srv.ListRootCertificates(ctx)
		assert.NoError(t, err)
		assert.Equal(t, list, alias)
	})

	t.Run("list nested certificates", func(t *testing.T) {
		list, err := srv.ListCertificates(ctx, root.Id)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.NotEmpty(t, findCert(caRoot.Name, list))
		assert.NotEmpty(t, findCert(childRoot.Name, list))
	})
}

func TestService_ListExpiredCertificates(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "expired",
		Days: 0,
	})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name: "ok",
		Days: 12,
	})
	assert.NoError(t, err)

	list, err := srv.ListExpiredCertificates(ctx)
	assert.NoError(t, err)

	assert.Len(t, list, 1)

	assert.NotEmpty(t, findCert(root.Name, list))
}

func TestService_ListSoonExpireCertificates(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name: "nosoon",
		Days: 40,
	})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name: "expired",
		Days: 0,
	})
	assert.NoError(t, err)

	soon, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "soon",
		Days: 12,
	})
	assert.NoError(t, err)

	list, err := srv.ListSoonExpireCertificates(ctx)
	assert.NoError(t, err)

	assert.Len(t, list, 1)

	assert.NotEmpty(t, findCert(soon.Name, list))
}

func TestServiceImpl_GetCertificate(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name:    "root",
		Days:    365,
		Domains: []string{"xyz", "abc"},
	})
	assert.NoError(t, err)
	sort.Strings(root.Domains)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name: "omega",
		Days: 365,
	})
	assert.NoError(t, err)

	t.Run("get known", func(t *testing.T) {
		saved, err := srv.GetCertificate(ctx, root.Id)
		assert.NoError(t, err)
		sort.Strings(saved.Domains)
		assert.True(t, root.CreatedAt.Equal(saved.CreatedAt))
		assert.True(t, root.ExpireAt.Equal(saved.ExpireAt))
		assert.True(t, root.RevokedAt.Equal(saved.RevokedAt))
		assert.True(t, root.UpdatedAt.Equal(saved.UpdatedAt))

		root.CreatedAt, root.ExpireAt, root.RevokedAt, root.UpdatedAt = saved.CreatedAt, saved.ExpireAt, saved.RevokedAt, saved.UpdatedAt
		assert.Equal(t, root, saved)
	})

	t.Run("get unknown", func(t *testing.T) {
		saved, err := srv.GetCertificate(ctx, 100500)
		assert.Error(t, err)
		assert.Empty(t, saved)
	})

}

func TestServiceImpl_GetPublicCert(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	assert.NoError(t, database.Raw("TRUNCATE certificates").Error)

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "root",
		Days: 365,
	})
	assert.NoError(t, err)

	child, err := srv.CreateCertificate(ctx, api.Subject{
		Name:   "child.root",
		Days:   365,
		Issuer: root.Id,
	})
	assert.NoError(t, err)

	childCa, err := srv.CreateCertificate(ctx, api.Subject{
		Name:   "ca.root",
		Days:   365,
		Issuer: root.Id,
		Ca:     true,
	})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{
		Name: "omega",
		Days: 365,
	})
	assert.NoError(t, err)

	t.Run("get known cert", func(t *testing.T) {
		cert, err := srv.GetPublicCert(ctx, root.Id)
		assert.NoError(t, err)
		assert.NotEmpty(t, cert)
		assert.Contains(t, cert, "-----BEGIN CERTIFICATE-----")
	})

	t.Run("get unknown cert", func(t *testing.T) {
		assert.NoError(t, err)
		cert, err := srv.GetPublicCert(ctx, 100500)
		assert.Error(t, err)
		assert.Empty(t, cert)
	})

	t.Run("non-ca certificate should not have sign capability but can can authorize", func(t *testing.T) {
		cert, err := srv.GetPublicCert(ctx, child.Id)
		assert.NoError(t, err)

		block, _ := pem.Decode([]byte(cert))
		assert.Equal(t, "CERTIFICATE", block.Type)
		parsed, err := x509.ParseCertificate(block.Bytes)
		assert.NoError(t, err)
		assert.False(t, parsed.IsCA)
		assert.True(t, parsed.KeyUsage&x509.KeyUsageCertSign == 0)
		assert.Contains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		assert.Contains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})

	t.Run("child certificate should be signed by parent CA", func(t *testing.T) {
		cert, err := srv.GetPublicCert(ctx, child.Id)
		assert.NoError(t, err)

		block, _ := pem.Decode([]byte(cert))
		assert.Equal(t, "CERTIFICATE", block.Type)
		parsed, err := x509.ParseCertificate(block.Bytes)
		assert.NoError(t, err)

		assert.Equal(t, root.Name, parsed.Issuer.CommonName)
		assert.Equal(t, root.Serial, parsed.Issuer.SerialNumber)

		pool := x509.NewCertPool()

		ca, err := srv.GetPublicCert(ctx, root.Id)
		assert.NoError(t, err)

		added := pool.AppendCertsFromPEM([]byte(ca))
		assert.True(t, added)

		_, err = parsed.Verify(x509.VerifyOptions{
			DNSName:     parsed.Subject.CommonName,
			Roots:       pool,
			CurrentTime: time.Now(),
		})
		assert.NoError(t, err)
	})

	t.Run("ca certificate should have sign capability but can not authorize", func(t *testing.T) {
		cert, err := srv.GetPublicCert(ctx, root.Id)
		assert.NoError(t, err)

		block, _ := pem.Decode([]byte(cert))
		assert.Equal(t, "CERTIFICATE", block.Type)
		parsed, err := x509.ParseCertificate(block.Bytes)
		assert.NoError(t, err)
		assert.True(t, parsed.IsCA)
		assert.False(t, parsed.KeyUsage&x509.KeyUsageCertSign == 0)
		assert.NotContains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		assert.NotContains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})

	t.Run("nested ca certificate should have sign capability but can not authorize and has ref to parent CA", func(t *testing.T) {
		cert, err := srv.GetPublicCert(ctx, childCa.Id)
		assert.NoError(t, err)

		block, _ := pem.Decode([]byte(cert))
		assert.Equal(t, "CERTIFICATE", block.Type)
		parsed, err := x509.ParseCertificate(block.Bytes)
		assert.NoError(t, err)
		assert.True(t, parsed.IsCA)
		assert.False(t, parsed.KeyUsage&x509.KeyUsageCertSign == 0)
		assert.NotContains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		assert.NotContains(t, parsed.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		assert.Equal(t, root.Name, parsed.Issuer.CommonName)
	})
}

func TestServiceImpl_GetAuthorityKey(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "root",
		Days: 365,
	})
	assert.NoError(t, err)

	omega, err := srv.CreateCertificate(ctx, api.Subject{
		Name:   "omega",
		Days:   365,
		Issuer: root.Id,
	})
	assert.NoError(t, err)

	t.Run("can get key", func(t *testing.T) {
		key, err := srv.GetPrivateKey(ctx, omega.Id)
		assert.NoError(t, err)
		assert.Contains(t, key, "PRIVATE KEY")
	})

	t.Run("can not get CA key by default", func(t *testing.T) {
		key, err := srv.GetPrivateKey(ctx, root.Id)
		assert.Error(t, err)
		assert.Empty(t, key)
	})

	t.Run("can get CA key if allowed", func(t *testing.T) {
		public, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512), service.ExposeCAKey(true))
		if !assert.NoError(t, err) {
			return
		}
		key, err := public.GetPrivateKey(ctx, root.Id)
		assert.NoError(t, err)
		assert.Contains(t, key, "PRIVATE KEY")
	})

	t.Run("stored key encrypted", func(t *testing.T) {
		key, err := srv.GetPrivateKey(ctx, omega.Id)
		assert.NoError(t, err)
		var saved db.Certificate
		err = database.Model(&db.Certificate{}).Take(&saved, omega.Id).Error
		assert.NoError(t, err)
		assert.NotEmpty(t, key)
		assert.NotEmpty(t, saved.Key)
		assert.NotEqual(t, key, saved.Key)
	})

	t.Run("access stored key with wrong encryption key should not expose original key", func(t *testing.T) {
		key, err := srv.GetPrivateKey(ctx, omega.Id)
		assert.NoError(t, err)

		fake, err := service.New(database, "wrong-key")
		if !assert.NoError(t, err) {
			return
		}

		invalid, err := fake.GetPrivateKey(ctx, omega.Id)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.NotEqual(t, key, invalid)
	})

	t.Run("get unknown key", func(t *testing.T) {
		key, err := srv.GetPrivateKey(ctx, 100500)
		assert.Error(t, err)
		assert.Empty(t, key)
	})
}

func TestService_RevokeCertificate(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	root, err := srv.CreateCertificate(ctx, api.Subject{
		Name: "root",
		Days: 365,
	})
	assert.NoError(t, err)

	err = srv.RevokeCertificate(ctx, root.Id)
	assert.NoError(t, err)

	_, err = srv.GetCertificate(ctx, root.Id)
	assert.Error(t, err)

	list, err := srv.ListRootCertificates(ctx)
	assert.NoError(t, err)
	assert.Empty(t, list)
}

func TestService_ListRevokedCertificates(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	t.Run("revoke should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example",
			Days: 365,
			Ca:   true,
		})
		assert.NoError(t, err)
		child1, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child1",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)
		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name:   "child2",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)

		err = srv.RevokeCertificate(ctx, child1.Id)
		assert.NoError(t, err)

		list, err := srv.ListRevokedCertificates(ctx, root.Id)
		assert.NoError(t, err)

		assert.Len(t, list, 1)

		assert.Equal(t, child1.Id, list[0].Id)
		assert.NotEmpty(t, list[0].RevokedAt)
	})
}

func TestService_GetRevokedCertificatesList(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	t.Run("revoke should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example",
			Days: 365,
			Ca:   true,
		})
		assert.NoError(t, err)

		child1, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child1",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)

		child1cert, err := getCert(ctx, srv, child1.Id)
		assert.NoError(t, err)

		rootCert, err := getCert(ctx, srv, root.Id)
		assert.NoError(t, err)

		_, err = srv.CreateCertificate(ctx, api.Subject{
			Name:   "child2",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)

		err = srv.RevokeCertificate(ctx, child1.Id)
		assert.NoError(t, err)

		pemCRL, err := srv.GetRevokedCertificatesList(ctx, root.Id)
		assert.NoError(t, err)

		block, _ := pem.Decode([]byte(pemCRL))
		assert.Equal(t, "X509 CRL", block.Type)
		assert.NotEmpty(t, block.Bytes)

		crl, err := x509.ParseCRL(block.Bytes)
		assert.NoError(t, err)

		pool := x509.NewCertPool()
		pool.AddCert(rootCert)

		_, err = child1cert.Verify(x509.VerifyOptions{
			DNSName:     child1cert.Subject.CommonName,
			Roots:       pool,
			CurrentTime: time.Now(),
		})
		assert.NoError(t, err)

		hasRevoked := false
		for _, revoked := range crl.TBSCertList.RevokedCertificates {
			hasRevoked = revoked.SerialNumber.Cmp(child1cert.SerialNumber) == 0
			if hasRevoked {
				break
			}
		}

		assert.True(t, hasRevoked)
	})
}

func TestService_BatchCreateCertificate(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	batch := api.Batch{
		Certificate: api.Subject{
			Name: "root1",
			Days: 365,
		},
		Nested: []api.Batch{
			{Certificate: api.Subject{
				Name: "child1",
				Days: 365,
			}},
			{Certificate: api.Subject{
				Name: "child2",
				Days: 365,
			}},
		},
	}

	list, err := srv.BatchCreateCertificate(ctx, []api.Batch{batch})
	assert.NoError(t, err)
	assert.Len(t, list, 3)
	root1 := findCert("root1", list)
	assert.NotEmpty(t, root1)
	assert.True(t, root1.Ca)

	child1 := findCert("child1", list)
	assert.NotEmpty(t, child1)
	assert.False(t, child1.Ca)
	assert.Equal(t, root1.Id, child1.Issuer)

	assert.NotEmpty(t, findCert("child2", list))
}

func TestService_RenewCertificate(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	t.Run("basic CA renew should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example",
			Days: 365,
			Ca:   true,
		})
		assert.NoError(t, err)
		rootCert, err := srv.GetPublicCert(ctx, root.Id)
		assert.NoError(t, err)

		child1, err := srv.CreateCertificate(ctx, api.Subject{
			Name:   "child1",
			Issuer: root.Id,
			Days:   365,
		})
		assert.NoError(t, err)
		child1Cert, err := srv.GetPublicCert(ctx, child1.Id)
		assert.NoError(t, err)

		newRoot, err := srv.RenewCertificate(ctx, root.Id, api.Renewal{Days: 400})
		assert.NoError(t, err)
		assert.Equal(t, root.Id, newRoot.Id)
		assert.Equal(t, root.Ca, newRoot.Ca)
		assert.True(t, root.CreatedAt.Equal(newRoot.CreatedAt))
		assert.True(t, newRoot.ExpireAt.After(root.ExpireAt))

		newChild, err := srv.GetCertificate(ctx, child1.Id)
		assert.NoError(t, err)
		assert.Equal(t, child1.Id, newChild.Id)
		assert.Equal(t, child1.Ca, newChild.Ca)
		assert.Equal(t, child1.Issuer, newChild.Issuer)
		assert.True(t, child1.CreatedAt.Equal(newChild.CreatedAt))
		assert.True(t, newChild.ExpireAt.After(child1.ExpireAt))

		child1NewCert, err := srv.GetPublicCert(ctx, child1.Id)
		assert.NoError(t, err)
		rootNewCert, err := srv.GetPublicCert(ctx, root.Id)
		assert.NoError(t, err)

		assert.NotEqual(t, rootCert, rootNewCert)
		assert.NotEqual(t, child1Cert, child1NewCert)
	})

	t.Run("basic CA renew should work", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name:    "example2",
			Days:    365,
			Ca:      true,
			Domains: []string{"abc", "def"},
		})
		assert.NoError(t, err)
		newRoot, err := srv.RenewCertificate(ctx, root.Id, api.Renewal{
			Days:    400,
			Domains: []string{"gamma"},
		})
		assert.NoError(t, err)

		assert.Len(t, newRoot.Domains, 1)
		assert.Contains(t, newRoot.Domains, "gamma")
	})

	t.Run("renew fields", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name:  "example2",
			Days:  365,
			Units: []string{"a", "b"},
		})
		assert.NoError(t, err)
		newRoot, err := srv.RenewCertificate(ctx, root.Id, api.Renewal{
			Days:  400,
			Units: []string{"c"},
		})
		assert.NoError(t, err)

		assert.Len(t, newRoot.Units, 1)
		assert.Equal(t, "c", newRoot.Units[0])
	})
	t.Run("renew IP", func(t *testing.T) {
		root, err := srv.CreateCertificate(ctx, api.Subject{
			Name: "example2",
			Days: 365,
			Ips:  []string{"127.0.0.1", "1.2.3.4"},
		})
		assert.NoError(t, err)
		newRoot, err := srv.RenewCertificate(ctx, root.Id, api.Renewal{
			Days: 400,
			Ips:  []string{"5.6.7.8"},
		})
		assert.NoError(t, err)

		assert.Len(t, newRoot.Ips, 1)
		assert.Equal(t, "5.6.7.8", newRoot.Ips[0])
	})
}

func TestService_GetStatus(t *testing.T) {
	database, err := db.New(sqlite.Open("file:" + t.Name() + "?mode=memory&cache=shared&_foreign_keys=on"))
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()

	srv, err := service.New(database, "test123", service.Random(&fastRandom{}), service.Key(512))
	if !assert.NoError(t, err) {
		return
	}

	// add nested
	root1, err := srv.CreateCertificate(ctx, api.Subject{Name: "root1", Days: 365, Ca: true})
	assert.NoError(t, err)

	_, err = srv.CreateCertificate(ctx, api.Subject{Name: "root2", Days: 365, Issuer: root1.Id})
	assert.NoError(t, err)

	// add expired
	_, err = srv.CreateCertificate(ctx, api.Subject{Name: "expired", Days: 0, Issuer: root1.Id})
	assert.NoError(t, err)

	// add soon expire
	_, err = srv.CreateCertificate(ctx, api.Subject{Name: "soon-expire", Days: 12, Issuer: root1.Id})
	assert.NoError(t, err)

	// add revoked
	revoked, err := srv.CreateCertificate(ctx, api.Subject{Name: "revoked", Days: 12, Issuer: root1.Id})
	assert.NoError(t, err)
	err = srv.RevokeCertificate(ctx, revoked.Id)
	assert.NoError(t, err)

	expected := api.Status{
		Total:      5,
		Expired:    1,
		SoonExpire: 1,
		Ca:         1,
		Revoked:    1,
	}

	status, err := srv.GetStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, status)

}

func getCert(ctx context.Context, client api.API, id uint) (*x509.Certificate, error) {
	pemCert, err := client.GetPublicCert(ctx, id)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(pemCert))
	return x509.ParseCertificate(block.Bytes)
}

func findCert(name string, list []api.Certificate) api.Certificate {
	for _, item := range list {
		if item.Name == name {
			return item
		}
	}
	return api.Certificate{}
}
