// Code generated by simple-swagger  DO NOT EDIT.
package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	api "github.com/reddec/ssl-cp/api"

	"github.com/julienschmidt/httprouter"
)

func New(impl api.API) http.Handler {
	router := httprouter.New()

	router.GET("/status", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		ctx := r.Context()
		res, err := impl.GetStatus(ctx)

		if err != nil {
			log.Println("getStatus: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificates", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		ctx := r.Context()
		res, err := impl.ListRootCertificates(ctx)

		if err != nil {
			log.Println("listRootCertificates: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.POST("/certificates", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramSubject api.Subject // in body

		if err := json.NewDecoder(r.Body).Decode(&paramSubject); err != nil {
			log.Println("createCertificate: decode subject from body:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := impl.CreateCertificate(ctx, paramSubject)

		if err != nil {
			log.Println("createCertificate: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.PUT("/certificates", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramBatch []api.Batch // in body

		if err := json.NewDecoder(r.Body).Decode(&paramBatch); err != nil {
			log.Println("batchCreateCertificate: decode batch from body:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := impl.BatchCreateCertificate(ctx, paramBatch)

		if err != nil {
			log.Println("batchCreateCertificate: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificates/expired", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		ctx := r.Context()
		res, err := impl.ListExpiredCertificates(ctx)

		if err != nil {
			log.Println("listExpiredCertificates: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificates/soon-expire", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		ctx := r.Context()
		res, err := impl.ListSoonExpireCertificates(ctx)

		if err != nil {
			log.Println("listSoonExpireCertificates: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificate/:certificate_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("getCertificate: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("getCertificate: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.GetCertificate(ctx, paramCertificateId)

		if err != nil {
			log.Println("getCertificate: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.PUT("/certificate/:certificate_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint  // in path
		var paramRenewal api.Renewal // in body

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("renewCertificate: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("renewCertificate: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&paramRenewal); err != nil {
			log.Println("renewCertificate: decode renewal from body:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := impl.RenewCertificate(ctx, paramCertificateId, paramRenewal)

		if err != nil {
			log.Println("renewCertificate: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.DELETE("/certificate/:certificate_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("revokeCertificate: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("revokeCertificate: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		err := impl.RevokeCertificate(ctx, paramCertificateId)

		if err != nil {
			log.Println("revokeCertificate: execute:", err)
			autoError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	})

	router.GET("/certificate/:certificate_id/cert", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("getPublicCert: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("getPublicCert: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.GetPublicCert(ctx, paramCertificateId)

		if err != nil {
			log.Println("getPublicCert: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificate/:certificate_id/key", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("getPrivateKey: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("getPrivateKey: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.GetPrivateKey(ctx, paramCertificateId)

		if err != nil {
			log.Println("getPrivateKey: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificate/:certificate_id/issued", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("listCertificates: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("listCertificates: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.ListCertificates(ctx, paramCertificateId)

		if err != nil {
			log.Println("listCertificates: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificate/:certificate_id/revoked", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("listRevokedCertificates: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("listRevokedCertificates: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.ListRevokedCertificates(ctx, paramCertificateId)

		if err != nil {
			log.Println("listRevokedCertificates: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	router.GET("/certificate/:certificate_id/revoked/crl", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer r.Body.Close()

		var paramCertificateId uint // in path

		if v, err := url.PathUnescape(ps.ByName("certificate_id")); err != nil {
			log.Println("getRevokedCertificatesList: decode certificate_id from path:", err)
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			if value, err := strconv.ParseUint(v, 10, 64); err == nil {
				paramCertificateId = uint(value)
			} else {
				log.Println("getRevokedCertificatesList: decode certificate_id from path:", err)
				jsonError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := impl.GetRevokedCertificatesList(ctx, paramCertificateId)

		if err != nil {
			log.Println("getRevokedCertificatesList: execute:", err)
			autoError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		w.WriteHeader(http.StatusOK)
		_ = encoder.Encode(res)

	})

	return router
}

func autoError(w http.ResponseWriter, err error) {
	if apiError, ok := api.AsAPIError(err); ok {
		jsonError(w, apiError.Message, apiError.Status)
		return
	}
	jsonError(w, err.Error(), http.StatusInternalServerError)
}

func jsonError(w http.ResponseWriter, err string, code int) {
	type errMessage struct {
		Message string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	w.WriteHeader(code)
	_ = encoder.Encode(&errMessage{Message: err})
}