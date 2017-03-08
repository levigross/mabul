package http

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/lucas-clemente/quic-go/h2quic"
	"github.com/valyala/fasthttp"
)

// GetFloodAttack is our generic interface that allows us to execute a get flood
type GetFloodAttack interface {
	SetAttacker(*Attacker)
	Get() error
}

// FastHTTPGet implements our GetFloodAttack
type FastHTTPGet struct {
	attacker *Attacker
	client   *fasthttp.Client
}

var _ GetFloodAttack = &FastHTTPGet{}

// SetAttacker sets the attacker and instantiates the client
func (f *FastHTTPGet) SetAttacker(a *Attacker) {
	f.attacker = a
	f.client = &fasthttp.Client{
		TLSConfig:    f.attacker.Config.TLSConfig,
		ReadTimeout:  f.attacker.Config.Timeout,
		WriteTimeout: f.attacker.Config.Timeout,
	}
}

// Get does what it says it will
func (f *FastHTTPGet) Get() error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("GET")
	req.SetRequestURI(f.attacker.Config.url.String())
	req.Header.SetHost(f.attacker.Config.url.Hostname())
	req.SetConnectionClose()

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := f.client.Do(req, resp); err != nil {
		return err
	}
	if _, err := resp.WriteTo(ioutil.Discard); err != nil {
		return err
	}
	return nil
}

// RegHTTPGet implements our GET flood using the regular net/http
type RegHTTPGet struct {
	attacker *Attacker
	client   *http.Client
}

var _ GetFloodAttack = &RegHTTPGet{}

// SetAttacker ....
func (r *RegHTTPGet) SetAttacker(a *Attacker) {
	r.attacker = a
	if a.Config.Quic {
		r.client = &http.Client{
			Transport: &h2quic.QuicRoundTripper{},
		}
	} else {
		r.client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   r.attacker.Config.Timeout,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       r.attacker.Config.Timeout,
				TLSHandshakeTimeout:   r.attacker.Config.Timeout,
				ExpectContinueTimeout: r.attacker.Config.Timeout,
			},
		}
	}
}

// Get ...
func (r *RegHTTPGet) Get() error {
	ourReq, err := http.NewRequest("GET", r.attacker.Config.url.String(), nil)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	var cancelfunc context.CancelFunc
	ctx, cancelfunc = context.WithTimeout(ctx, r.attacker.Config.Timeout)
	defer cancelfunc()
	ourReq = ourReq.WithContext(ctx)
	resp, err := r.client.Do(ourReq)
	if err != nil {
		return err
	}
	if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
		return err
	}
	return nil
}
