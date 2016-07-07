package http

import (
	"fmt"
	"strings"
)

func TestHost() string {
	return "http://localhost"
}

func TestPort() int {
	return 8181
}

func CreateFullTestURLForEndpoint(endpoint string) string {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	url := fmt.Sprintf("%s:%v%s", TestHost(), TestPort(), endpoint)
	return url
}
