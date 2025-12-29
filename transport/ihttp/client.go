package ihttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Verthandii/spring/encoding"
	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
)

var (
	Logger ilogger.Logger // 日志
)

func init() {
	path := fmt.Sprintf("%s/runtime/logs/request.log", env.GetRootPath())
	Logger = ilogger.NewFileZapLogger(path)
}

type DecodeErrorFunc func(ctx context.Context, res *http.Response) error

type EncodeRequestFunc func(ctx context.Context, contentType string, in any) (body []byte, err error)

type DecodeResponseFunc func(ctx context.Context, res *http.Response, out any) error

type ClientOption func(*clientOptions)

type clientOptions struct {
	tlsConf      *tls.Config
	timeout      time.Duration
	timeoutTry   int
	endpoint     string
	userAgent    string
	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
	transport    http.RoundTripper
}

func WithTlsConfig(cfg *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = cfg
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(o *clientOptions) {
		o.userAgent = userAgent
	}
}

func WithRequestEncoder(encoder EncodeRequestFunc) ClientOption {
	return func(o *clientOptions) {
		o.encoder = encoder
	}
}

func WithResponseDecoder(decoder DecodeResponseFunc) ClientOption {
	return func(o *clientOptions) {
		o.decoder = decoder
	}
}

func WithErrorDecoder(errorDecoder DecodeErrorFunc) ClientOption {
	return func(o *clientOptions) {
		o.errorDecoder = errorDecoder
	}
}

func WithTransport(trans *http.Transport) ClientOption {
	return func(o *clientOptions) {
		o.transport = trans
	}
}

func WithTimeoutTry(try int) ClientOption {
	return func(o *clientOptions) {
		o.timeoutTry = try
	}
}

type Client struct {
	opts clientOptions
	cc   *http.Client
}

func NewClient(opts ...ClientOption) (*Client, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	options := clientOptions{
		timeout:      3000 * time.Millisecond,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		errorDecoder: DefaultErrorDecoder,
		transport: otelhttp.NewTransport(roundTripperMetrics(&http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   50, // 尝试解决 context deadline exceeded (Client.Timeout exceeded while awaiting headers)
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		})),
	}

	for _, option := range opts {
		option(&options)
	}

	if options.tlsConf != nil {
		if t, ok := options.transport.(*http.Transport); ok {
			t.TLSClientConfig = options.tlsConf
		}
	}
	if options.timeoutTry < 1 {
		options.timeoutTry = 1
	}

	return &Client{
		opts: options,
		cc: &http.Client{
			Timeout:   options.timeout,
			Transport: options.transport,
		},
	}, nil
}

func (client *Client) Get(ctx context.Context, path string, reply any) error {
	return client.Invoke(ctx, http.MethodGet, path, "", nil, reply)
}

func (client *Client) Post(ctx context.Context, path string, contentType string, args, reply any) error {
	return client.Invoke(ctx, http.MethodPost, path, contentType, args, reply)
}

// Invoke makes an rpc call procedure for remote service.
func (client *Client) Invoke(ctx context.Context, method, path, contentType string, args any, reply any) error {
	req, err := client.BuildRequest(ctx, method, path, contentType, args)
	if err != nil {
		return err
	}
	return client.invoke(ctx, req, reply)
}

func (client *Client) invoke(ctx context.Context, req *http.Request, reply any) error {
	res, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	err = client.BuildResponse(ctx, res, reply)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) BuildRequest(ctx context.Context, method, path, contentType string, args any) (*http.Request, error) {
	var body io.Reader
	if args != nil {
		if contentType == "" {
			contentType = "application/json"
		}
		data, err := client.opts.encoder(ctx, contentType, args)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	url := fmt.Sprintf("%s%s", client.opts.endpoint, path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if client.opts.userAgent != "" {
		req.Header.Set("User-Agent", client.opts.userAgent)
	}

	return req, nil
}

func (client *Client) BuildResponse(ctx context.Context, res *http.Response, reply any) error {
	defer res.Body.Close()
	if err := client.opts.decoder(ctx, res, reply); err != nil {
		return err
	}
	return nil
}

func (client *Client) Do(req *http.Request) (*http.Response, error) {
	return client.do(req.Context(), req)
}

func (client *Client) do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	var (
		reqData  []byte
		respData []byte
	)
	t := time.Now()
	defer func() {
		if resp != nil && resp.Request != nil {
			req = resp.Request
			ctx = resp.Request.Context()
		}
		kvs := []any{
			"URL", req.URL.String(),
			"Duration", time.Since(t).Milliseconds(),
			"RequestHeader", req.Header,
			"RequestBody", string(reqData),
		}
		if err != nil {
			kvs = append(kvs,
				"Error", err.Error(),
			)
		}
		if resp != nil {
			kvs = append(kvs,
				"ResponseHeader", resp.Header,
				"ResponseBody", string(respData),
				"StatusCode", resp.StatusCode,
			)
		}
		Logger.InfowCtx(ctx, "RequestLog", kvs...)
	}()

	if req.Body != nil {
		reqData, err = io.ReadAll(req.Body)
		if err != nil {
			return
		}
		req.Body = io.NopCloser(bytes.NewReader(reqData))
	}
	var dataBuffer *bytes.Reader
	for i := 0; i < client.opts.timeoutTry; i++ {
		if req.GetBody != nil {
			bodyReadCloser, _ := req.GetBody()
			req.Body = bodyReadCloser
		} else if req.Body != nil {
			if dataBuffer == nil {
				data, err := io.ReadAll(req.Body)
				req.Body.Close()
				if err != nil {
					return nil, err
				}
				dataBuffer = bytes.NewReader(data)
				req.ContentLength = int64(dataBuffer.Len())
				req.Body = io.NopCloser(dataBuffer)
			}
			_, _ = dataBuffer.Seek(0, io.SeekStart)
		}
		resp, err = client.cc.Do(req)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				// 重试
				continue
			}
			return
		}
		if resp.Body != nil {
			respData, err = io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			resp.Body = io.NopCloser(bytes.NewReader(respData))
		}
		if err = client.opts.errorDecoder(ctx, resp); err != nil {
			return
		}
		return
	}
	return
}

func DefaultRequestEncoder(_ context.Context, contentType string, in any) ([]byte, error) {
	name := ContentSubtype(contentType)
	body, err := encoding.GetCodec(name).Marshal(in)
	if err != nil {
		return nil, err
	}
	return body, err
}

func DefaultResponseDecoder(_ context.Context, res *http.Response, v any) error {
	defer res.Body.Close()
	if v == nil {
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return CodecForResponse(res).Unmarshal(data, v)
}

func DefaultErrorDecoder(_ context.Context, res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err == nil {
		return errors.New(fmt.Sprintf("http status code:%d.", res.StatusCode) + string(data))
	}
	return err
}

func CodecForResponse(r *http.Response) encoding.Codec {
	codec := encoding.GetCodec(ContentSubtype(r.Header.Get("Content-Type")))
	if codec != nil {
		return codec
	}
	return encoding.GetCodec("json")
}
