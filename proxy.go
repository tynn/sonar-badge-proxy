package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy extends the httputil.ReverseProxy to add
// some config to access a Sonar service securly
type Proxy struct {
	httputil.ReverseProxy
	addr          string
	remote        *url.URL
	authorization string
	metric        map[string]string
	secret        string
	json 		  bool
}

func (x *Proxy) requested(u *url.URL) (string, string) {
	p := strings.Split(u.EscapedPath(), "/")
	if len(p) != 3 {
		panic("Invalid Request=" + u.String())
	}
	if p[2] == "" {
		panic("Invalid Project=")
	}
	m, ok := x.metric[p[1]]
	if !ok {
		panic("Invalid Metric=" + p[1])
	}
	return m, p[2]
}

func (x *Proxy) buildQuery(m string, p string, b string) url.Values {
	q := make(url.Values)
	if x.json {
		q.Set("metricKeys", m)
		q.Set("component", p)
	} else {
		q.Set("metric", m)
		q.Set("project", p)
	}
	if b != "" {
		q.Set("branch", b)
	}
	return q
}

func (x *Proxy) authorize(h *http.Header) {
	if x.authorization != "" {
		h.Add("Authorization", x.authorization)
	}
}

func (x *Proxy) rewriteURL(r *http.Request, q url.Values) {
	u := x.remote
	r.URL.Scheme = u.Scheme
	r.URL.Host = u.Host
	r.URL.Path = u.Path
	r.URL.RawQuery = q.Encode()
	r.Host = u.Host
}

func (x *Proxy) verifyToken(p string, t ...string) {
	if s := x.secret; s != "" {
		m := md5.New()
		m.Write([]byte(p))
		m.Write([]byte{':'})
		m.Write([]byte(s))
		h := hex.EncodeToString(m.Sum(nil))
		for _, e := range t {
			if h == e {
				return
			}
		}
		panic(http.StatusUnauthorized)
	}
}

func (x *Proxy) director(r *http.Request) {
	m, p := x.requested(r.URL)
	q := r.URL.Query()
	x.verifyToken(p, q["token"]...)
	q = x.buildQuery(m, p, q.Get("branch"))
	x.authorize(&r.Header)
	x.rewriteURL(r, q)
}

func modifyResponse(r *http.Response) error {
	c := r.StatusCode
	switch {
	case c < http.StatusOK:
		panic(http.StatusBadGateway)
	case c < http.StatusMultipleChoices:
		return nil
	case c < http.StatusBadRequest:
		return nil
	case c < http.StatusInternalServerError:
		panic(http.StatusNotFound)
	default:
		panic(http.StatusBadGateway)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, e error) {
	panic(http.StatusBadGateway)
}

func handlePanic(w http.ResponseWriter) {
	switch err := recover(); err {
	case nil:
	case http.StatusBadGateway:
		http.Error(w, "502 bad gateway", http.StatusBadGateway)
	case http.StatusUnauthorized:
		http.Error(w, "401 request unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

func serveFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func (x *Proxy) serveHTTP(w http.ResponseWriter, r *http.Request) {
	defer handlePanic(w)
	x.ServeHTTP(w, r)
}

// Server sets up the proxy handlers for
// the favicon.ico and the badges
func (x *Proxy) Server() *http.Server {
	m := http.NewServeMux()
	m.HandleFunc("/favicon.ico", serveFavicon)
	m.HandleFunc("/", x.serveHTTP)
	return &http.Server{Addr: x.addr, Handler: m}
}

func basicAuthorization(t string) string {
	if t == "" {
		return ""
	}
	a := []byte(t + ":")
	b := base64.StdEncoding.EncodeToString(a)
	return "Basic " + b
}

// NewProxy creates a new proxy from config
func NewProxy(c *Config) *Proxy {
	p := new(Proxy)
	p.Director = p.director
	p.ModifyResponse = modifyResponse
	p.ErrorHandler = errorHandler
	p.addr = ":" + c.Port
	p.remote = c.Remote
	p.authorization = basicAuthorization(c.Authorization)
	p.metric = c.Metric
	p.secret = c.Secret
	p.json = c.Json
	return p
}
