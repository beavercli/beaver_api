package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	db_uri := "postgres://beaver_api:beaver_api@localhost:5432/beaver_api"
	srv_addr := "127.0.0.1:8080"
	dbg := true

	os.Setenv("DATABASE_URI", db_uri)
	os.Setenv("SERVER_ADDR", srv_addr)
	os.Setenv("DEBUG", "true")

	cfg := New()

	assert.Equal(t, db_uri, cfg.DB.URI)
	assert.Equal(t, srv_addr, cfg.Server.Addr)
	assert.Equal(t, dbg, cfg.DebugMode)
}
