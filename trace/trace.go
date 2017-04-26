// Package trace is used to calculate page speed
package trace

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

type Latency map[string]time.Duration

const (
	maxRedirects = 10
)

func Trace(method, uri string, headers []string, body string, redirected int, latency map[string]time.Duration) error {

	u, err := parseUrl(uri)
	if err != nil {
		return err
	}

	req, err := newRequest(method, u, headers, body)
	if err != nil {
		return err
	}

	var t0, t1, t2, t3, t4 time.Time

	t := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t0 = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t1 = time.Now()
		},
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			// FIXME: 2017/4/25 return error
			if err != nil {
				log.Fatalf("E! unable to connect to host %s %v", addr, err)
			}
			t2 = time.Now()
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			t3 = time.Now()
		},
		GotFirstResponseByte: func() {
			t4 = time.Now()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), t))

	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if t0.IsZero() {
		t0 = t1
	}

	latency["dns_lookup"] += t1.Sub(t0)
	latency["tcp_connection"] += t2.Sub(t1)
	switch u.Scheme {
	case "https":
		latency["tls_handshake"] += t3.Sub(t2)
		latency["server_processing"] += t4.Sub(t3)
	case "http":
		latency["tls_handshake"] = 0
		latency["server_processing"] += t4.Sub(t2)
	}

	if isRedirect(resp) {
		loc, err := resp.Location()
		if err != nil {
			return err
		}
		redirected++
		if redirected > maxRedirects {
			return errors.New("got max redirects")
		}
		err = Trace(method, loc.String(), headers, body, redirected, latency)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseUrl(uri string) (*url.URL, error) {

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
		if !strings.HasSuffix(u.Host, ":80") {
			u.Scheme += "s"
		}
	}

	return u, nil
}

func newRequest(method string, u *url.URL, headers []string, body string) (*http.Request, error) {

	req, err := http.NewRequest(method, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	for _, h := range headers {
		k, v, err := headerKeyValue(h)
		if err != nil {
			return nil, err
		}
		req.Header.Add(k, v)
	}

	return req, nil
}

func headerKeyValue(h string) (string, string, error) {

	i := strings.Index(h, ":")
	if i == -1 {
		return "", "", errors.New("invalid header " + h)
	}

	return h[:i], strings.TrimLeft(h[i:], " "), nil
}

func isRedirect(resp *http.Response) bool {

	return resp.StatusCode > 299 && resp.StatusCode < 400
}
