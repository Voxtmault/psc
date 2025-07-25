package db

import (
	"testing"

	"github.com/voxtmault/psc/config"
)

func TestInitConnection(t *testing.T) {
	config := config.DBConfig{
		DBHost: "localhost",
		DBPort: "3306",
		DBUser: "user",
		DBPass: "pass",
		DBName: "psc",
	}

	if err := InitConnection(&config); err != nil {
		t.Errorf("failed to init connection: %v", err)
	}
}
