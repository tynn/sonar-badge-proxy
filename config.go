package main

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var metricMapping = map[string]string{
	"status":          "alert_status",
	"bugs":            "bugs",
	"codesmells":      "code_smells",
	"coverage":        "coverage",
	"duplications":    "duplicated_lines_density",
	"lines":           "ncloc",
	"maintainability": "sqale_rating",
	"reliability":     "reliability_rating",
	"security":        "security_rating",
	"techdept":        "sqale_index",
	"vulnerabilities": "vulnerabilities",
}

// Config represents the environment configuration of the server
type Config struct {
	Port          string
	Authorization string
	Metric        map[string]string
	Remote        *url.URL
	Secret        string
}

func port() string {
	s := os.Getenv("PORT")
	_, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		panic("Invalid PORT=" + s)
	}
	return s
}

func authorization() string {
	return os.Getenv("AUTHORIZATION")
}

func metric() map[string]string {
	s := os.Getenv("METRIC")
	m := make(map[string]string)
	for _, k := range strings.Split(s, ",") {
		v, ok := metricMapping[k]
		if !ok {
			panic("Invalid METRIC=" + s)
		}
		m[k] = v
	}
	return m
}

func remote(c *http.Client) *url.URL {
	s := os.Getenv("REMOTE")
	u, err := url.Parse("https://" + s + "/api/project_badges/measure")
	if err != nil {
		panic("Invalid REMOTE=" + s)
	}
	r, err := c.Get(u.String())
	if err != nil {
		panic("Invalid REMOTE=" + s)
	}
	switch r.StatusCode {
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusBadRequest:
		return u
	default:
		panic("Invalid REMOTE=" + s)
	}
}

func secret() string {
	return os.Getenv("SECRET")
}

// LoadConfig prepares the Config from the environment
func LoadConfig() *Config {
	c := &http.Client{Timeout: 10 * time.Second}
	return &Config{
		Port:          port(),
		Authorization: authorization(),
		Metric:        metric(),
		Remote:        remote(c),
		Secret:        secret(),
	}
}
