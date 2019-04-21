package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestRequested(t *testing.T) {
	u, err := url.Parse("/metric/project")
	AssertNoError(t, err)

	x := &Proxy{
		metric: map[string]string{
			"metric": "mapped",
		},
	}

	m, p := x.requested(u)
	AssertEqual(t, "mapped", m, "Wrong Metric=%v")
	AssertEqual(t, "project", p, "Wrong Project=%v")
}

func TestRequested_request(t *testing.T) {
	u, err := url.Parse("metric/project")
	AssertNoError(t, err)

	defer AssertPanic(t, "Request")
	new(Proxy).requested(u)
}

func TestRequested_project(t *testing.T) {
	u, err := url.Parse("/metric/")
	AssertNoError(t, err)

	defer AssertPanic(t, "Project")
	new(Proxy).requested(u)
}

func TestRequested_metric(t *testing.T) {
	u, err := url.Parse("/metric/project")
	AssertNoError(t, err)

	defer AssertPanic(t, "Metric")
	new(Proxy).requested(u)
}

func TestBuildQuery(t *testing.T) {
	m := "metric"
	p := "project"
	b := "branch"
	q := make(url.Values)

	q.Set("metric", m)
	q.Set("project", p)
	AssertDeepEqual(t, q, buildQuery(m, p, ""), "Wrong Query=%v")

	q.Set("branch", b)
	AssertDeepEqual(t, q, buildQuery(m, p, b), "Wrong Query=%v")
}

func TestAuthorize(t *testing.T) {
	p := new(Proxy)
	h := make(http.Header)

	p.authorize(&h)
	AssertEqual(t, "", h.Get("Authorization"), "Wrong Authorization=%v")

	p.authorization = "Basic 0123456789"
	p.authorize(&h)
	AssertEqual(t, p.authorization, h.Get("Authorization"), "Wrong Authorization=%v")
}

func TestRewriteURL(t *testing.T) {
	u, err := url.Parse("https://localhost:4000/path/to")
	AssertNoError(t, err)

	q, err := url.ParseQuery("q1=foo&q2=bar")
	AssertNoError(t, err)

	p := &Proxy{remote: u}
	r := &http.Request{URL: new(url.URL)}

	p.rewriteURL(r, q)
	u.RawQuery = q.Encode()
	AssertDeepEqual(t, u, r.URL, "Wrong r.URL=%v")
	AssertEqual(t, u.Host, r.Host, "Wrong r.Host=%v")
}

func TestVerifyToken(t *testing.T) {
	p := &Proxy{secret: "abc"}

	p.verifyToken("def", "2ac8358c6394edc2831b81b24c70287d")
	p.verifyToken("def", "2ac8358c6394edc2831b81b24c70287d", "def")
	p.verifyToken("def", "def", "2ac8358c6394edc2831b81b24c70287d")

	defer AssertPanic(t, "Token")
	p.verifyToken("def", "def", "2ac8358c6394edc2", "831b81b24c70287d")
}

func TestVerifyToken_empty(t *testing.T) {
	new(Proxy).verifyToken("def", "abc")
}

func TestServer_Addr(t *testing.T) {
	p := &Proxy{addr: ":4000"}
	s := p.Server()
	AssertEqual(t, p.addr, s.Addr, "Wrong Addr=%v")
}

func TestServer_badge(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusOK, w.Code, "Wrong Code=%v")

	b, err := ioutil.ReadFile("favicon.ico")
	AssertNoError(t, err)
	AssertDeepEqual(t, b, w.Body.Bytes(), "Wrong Body=%v")
}

func TestServer_favicon(t *testing.T) {
	r, err := http.NewRequest("GET", "/favicon.ico", nil)
	AssertNoError(t, err)

	h := new(Proxy).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusOK, w.Code, "Wrong Code=%v")

	b, err := ioutil.ReadFile("favicon.ico")
	AssertNoError(t, err)
	AssertDeepEqual(t, b, w.Body.Bytes(), "Wrong Body=%v")
}

func TestServer_root(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	AssertNoError(t, err)

	h := new(Proxy).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNotFound, w.Code, "Wrong Code=%v")
}

func TestServer_authorization(t *testing.T) {
	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	p := &Proxy{
		secret: "abc",
		metric: map[string]string{
			"metric": "metric",
		},
	}
	p.Director = p.director

	h := p.Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusUnauthorized, w.Code, "Wrong Code=%v")
}

func TestServer_connection(t *testing.T) {
	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse("http://localhost/")
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusBadGateway, w.Code, "Wrong Code=%v")
}

func TestServer_informational(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusBadGateway, w.Code, "Wrong Code=%v")
}

func TestServer_success(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNoContent, w.Code, "Wrong Code=%v")
}

func TestServer_redirection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNotModified, w.Code, "Wrong Code=%v")
}

func TestServer_clienterror(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusNotFound, w.Code, "Wrong Code=%v")
}

func TestServer_servererror(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))
	defer s.Close()

	r, err := http.NewRequest("GET", "/metric/project", nil)
	AssertNoError(t, err)

	u, err := url.Parse(s.URL)
	AssertNoError(t, err)

	c := &Config{
		Port:          "4000",
		Remote:        u,
		Authorization: "",
		Secret:        "",
		Metric: map[string]string{
			"metric": "metric",
		},
	}

	h := NewProxy(c).Server().Handler
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)
	AssertEqual(t, http.StatusBadGateway, w.Code, "Wrong Code=%v")
}

func TestNewProxy(t *testing.T) {
	p := "4000"
	r := new(url.URL)
	a := "authorization"
	s := "secret"
	m := map[string]string{
		"metric": "metric",
	}

	c := &Config{
		Port:          p,
		Remote:        r,
		Authorization: a,
		Secret:        s,
		Metric:        m,
	}

	x := NewProxy(c)
	AssertEqual(t, ":"+p, x.addr, "Wrong Director=%v")
	AssertEqual(t, r, x.remote, "Wrong Remote=%v")
	AssertEqual(t, a, x.authorization, "Wrong Authorization=%v")
	AssertEqual(t, s, x.secret, "Wrong Secret=%v")
	AssertDeepEqual(t, m, x.metric, "Wrong Metric=%v")
	AssertPointerEqual(t, x.director, x.Director, "Wrong Director=%v")
}
