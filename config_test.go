package main

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestPort(t *testing.T) {
	os.Setenv("PORT", "4000")
	AssertEqual(t, "4000", port(), "Wrong PORT=%v")
}

func TestPort_panic(t *testing.T) {
	os.Setenv("PORT", "port")
	defer AssertPanic(t, "PORT")
	port()
}

func TestAuthorization(t *testing.T) {
	os.Setenv("AUTHORIZATION", "BASIC foo:bar")
	AssertEqual(t, "BASIC foo:bar", authorization(), "Wrong AUTHORIZATION=%v")
}

func TestAuthorization_empty(t *testing.T) {
	os.Setenv("AUTHORIZATION", "")
	AssertEqual(t, "", authorization(), "Wrong AUTHORIZATION=%v")
}

func TestAuthorization_unset(t *testing.T) {
	os.Unsetenv("AUTHORIZATION")
	AssertEqual(t, "", authorization(), "Wrong AUTHORIZATION=%v")
}

func TestMetric(t *testing.T) {
	os.Setenv("METRIC", "bugs,lines")
	e := map[string]string{
		"bugs":  metricMapping["bugs"],
		"lines": metricMapping["lines"],
	}
	AssertDeepEqual(t, e, metric(), "Wrong METRIC=%v")
}

func TestMetric_empty(t *testing.T) {
	os.Setenv("METRIC", "")
	defer AssertPanic(t, "METRIC")
	metric()
}

func TestMetric_blank(t *testing.T) {
	os.Setenv("METRIC", "bugs,,lines")
	defer AssertPanic(t, "METRIC")
	metric()
}

func TestMetric_unkown(t *testing.T) {
	os.Setenv("METRIC", "bugs,lines,stats")
	defer AssertPanic(t, "METRIC")
	metric()
}

func TestRemote_domain(t *testing.T) {
	os.Setenv("REMOTE", "\n")
	defer AssertPanic(t, "Domain")
	remote(new(http.Client))
}

func TestRemote_request(t *testing.T) {
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL + "/api/project_badges/measure")
	AssertNoError(t, err)

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	os.Setenv("REMOTE", u.Host)
	AssertDeepEqual(t, u, remote(c), "Wrong REMOTE=%v")
}

func TestRemote_unauthorized(t *testing.T) {
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL + "/api/project_badges/measure")
	AssertNoError(t, err)

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	os.Setenv("REMOTE", u.Host)
	AssertDeepEqual(t, u, remote(c), "Wrong REMOTE=%v")
}

func TestRemote_none(t *testing.T) {
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "OK", http.StatusOK)
	}))
	defer s.Close()

	u, err := url.Parse(s.URL + "/api/project_badges/measure")
	AssertNoError(t, err)

	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	os.Setenv("REMOTE", u.Host)
	defer AssertPanic(t, "REMOTE")
	remote(c)
}

func TestRemote_empty(t *testing.T) {
	os.Setenv("REMOTE", "")
	defer AssertPanic(t, "REMOTE")
	remote(http.DefaultClient)
}

func TestRemote_unset(t *testing.T) {
	os.Unsetenv("REMOTE")
	defer AssertPanic(t, "REMOTE")
	remote(http.DefaultClient)
}

func TestSecret(t *testing.T) {
	os.Setenv("SECRET", "0123456789abcdef")
	AssertEqual(t, "0123456789abcdef", secret(), "Wrong SECRET=%v")
}

func TestSecret_empty(t *testing.T) {
	os.Setenv("SECRET", "")
	AssertEqual(t, "", secret(), "Wrong SECRET=%v")
}

func TestSecret_unset(t *testing.T) {
	os.Unsetenv("SECRET")
	AssertEqual(t, "", secret(), "Wrong SECRET=%v")
}
