package main

import (
	"github.com/envoyproxy/envoy/contrib/golang/filters/http/source/go/pkg/http"

	_ "example.com/filtermanager/plugins"
)

func init() {
	http.RegisterHttpFilterConfigFactoryAndParser("test", http.FilterManagerConfigFactory, &http.FilterManagerConfigParser{})
}

func main() {}
