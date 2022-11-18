package main

import (
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var metric = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "pk_http_requests",
		Buckets: []float64{.1, .25, 1, 2.5, 5, 20},
	},
	[]string{"domain", "method", "status", "route"},
)

func proxyTo(host string) *httputil.ReverseProxy {
	rp := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme:   "http",
		Host:     host,
		RawQuery: "",
	})
	rp.ModifyResponse = logTimeElapsed
	return rp
}

var systemsRegex = regexp.MustCompile("systems/[^/]+")
var membersRegex = regexp.MustCompile("members/[^/]+")
var groupsRegex = regexp.MustCompile("groups/[^/]+")
var switchesRegex = regexp.MustCompile("switches/[^/]+")
var guildsRegex = regexp.MustCompile("guilds/[^/]+")
var messagesRegex = regexp.MustCompile("messages/[^/]+")

func cleanPath(host, path string) string {
	if host != "api.pluralkit.me" {
		return ""
	}

	if !(strings.HasPrefix(host, "/v2") || strings.HasPrefix(host, "/private")) {
		return ""
	}

	path = systemsRegex.ReplaceAllString(path, "{systemRef}")
	path = membersRegex.ReplaceAllString(path, "{memberRef}")
	path = groupsRegex.ReplaceAllString(path, "{groupRef}")
	path = switchesRegex.ReplaceAllString(path, "{switchRef}")
	path = guildsRegex.ReplaceAllString(path, "{guild_id}")
	path = messagesRegex.ReplaceAllString(path, "{message_id}")

	return path
}

func requireEnv(key string) string {
	if val, ok := os.LookupEnv(key); !ok {
		panic("missing `" + key + "` in environment")
	} else {
		return val
	}
}
