// Code generated by simple-swagger  DO NOT EDIT.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	api "github.com/reddec/ssl-cp/api"
)

type RequestHook func(req *http.Request) error

func New(baseURL string, options ...Option) *Client {
	cl := &Client{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
	for _, opt := range options {
		opt(cl)
	}
	return cl
}

type Option func(cl *Client)

func HTTPClient(client *http.Client) Option {
	return func(cl *Client) {
		cl.client = client
	}
}

func Before(hook RequestHook) Option {
	return func(cl *Client) {
		cl.beforeHooks = append(cl.beforeHooks, hook)
	}
}

type Client struct {
	baseURL     string
	client      *http.Client
	beforeHooks []RequestHook
}

// Creates full copy of original client and applies new options.
func (client *Client) With(options ...Option) *Client {
	cp := &Client{
		baseURL:     client.baseURL,
		client:      client.client,
		beforeHooks: make([]RequestHook, len(client.beforeHooks)),
	}
	copy(cp.beforeHooks, client.beforeHooks)
	for _, opt := range options {
		opt(cp)
	}

	return cp
}
func (client *Client) GetStatus(ctx context.Context) (out api.Status, err error) {
	requestURL := client.baseURL + api.Prefix + "/status"
	var body bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) ListRootCertificates(ctx context.Context) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificates"
	var body bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) CreateCertificate(ctx context.Context, subject api.Subject) (out api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificates"
	var body bytes.Buffer

	if err = json.NewEncoder(&body).Encode(subject); err != nil {
		err = fmt.Errorf("encode subject: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) BatchCreateCertificate(ctx context.Context, batch []api.Batch) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificates"
	var body bytes.Buffer

	if err = json.NewEncoder(&body).Encode(batch); err != nil {
		err = fmt.Errorf("encode batch: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) ListExpiredCertificates(ctx context.Context) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificates/expired"
	var body bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) ListSoonExpireCertificates(ctx context.Context) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificates/soon-expire"
	var body bytes.Buffer

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) GetCertificate(ctx context.Context, certificateId uint) (out api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) RenewCertificate(ctx context.Context, certificateId uint, renewal api.Renewal) (out api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))
	if err = json.NewEncoder(&body).Encode(renewal); err != nil {
		err = fmt.Errorf("encode renewal: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) RevokeCertificate(ctx context.Context, certificateId uint) (err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}

	return
}

func (client *Client) GetPublicCert(ctx context.Context, certificateId uint) (out string, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}/cert"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) GetPrivateKey(ctx context.Context, certificateId uint) (out string, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}/key"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) ListCertificates(ctx context.Context, certificateId uint) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}/issued"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) ListRevokedCertificates(ctx context.Context, certificateId uint) (out []api.Certificate, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}/revoked"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func (client *Client) GetRevokedCertificatesList(ctx context.Context, certificateId uint) (out string, err error) {
	requestURL := client.baseURL + api.Prefix + "/certificate/{certificate_id}/revoked/crl"
	var body bytes.Buffer

	requestURL = strings.ReplaceAll(requestURL, "{certificate_id}", url.PathEscape(strconv.FormatUint(uint64(certificateId), 10)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, &body)
	if err != nil {
		err = fmt.Errorf("prepare request: %w", err)
		return
	}

	for _, hook := range client.beforeHooks {
		if err = hook(req); err != nil {
			return
		}
	}

	res, err := client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("execute request: %w", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		err = getError(res)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		return
	}

	return
}

func getError(res *http.Response) error {
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		return errors.New(string(payload))
	}
	var msg api.Error
	err = json.Unmarshal(payload, &msg)
	if err != nil {
		// fallback
		return errors.New(string(payload))
	}
	msg.Status = res.StatusCode
	return &msg
}