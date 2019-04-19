package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	Addr          string
	Authorization string
	Metric        []string
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

func metric() []string {
	s := os.Getenv("METRIC")
	m := strings.Split(s, ",")
	for _, k := range m {
		if _, ok := metricMapping[k]; !ok {
			panic("Invalid METRIC=" + s)
		}
	}
	return m
}

func remote(c *http.Client) *url.URL {
	s := os.Getenv("REMOTE")
	u, err := url.Parse(s + "/api/project_badges/measure")
	if err != nil {
		panic("Invalid REMOTE=" + s)
	}
	u.Scheme = "https"
	r, err := c.Head(u.String())
	if err != nil {
		panic("Invalid REMOTE=" + s)
	}
	switch r.StatusCode {
	case http.StatusUnauthorized:
		log.Print("AUTHORIZATION required for REMOTE=" + s)
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
	return &Config{
		Addr:          ":" + port(),
		Authorization: authorization(),
		Metric:        metric(),
		Remote:        remote(http.DefaultClient),
		Secret:        secret(),
	}
}
