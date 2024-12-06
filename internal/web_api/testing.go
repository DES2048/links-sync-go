package webapi

import (
	"links-sync-go/internal/config"
	"testing"
)

func TestServer(t *testing.T) IApiServer {
	t.Helper()

	config := &config.Config{
		Server: config.Server{
			ApiKey: "secret-key",
		},
		Db: config.Db{
			Url: ":memory:",
		},
	}
	server := NewApiServer(config)

	return server
}
