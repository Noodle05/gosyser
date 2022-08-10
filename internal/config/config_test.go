package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("APP_CONFIG_FILE", "testdata/test-config.yml")

	c, err := LoadConfig()
	if err != nil {
		t.Error("Load config failed.", err)
	}

	expected := &Configuration{
		Server: ServerConfiguration{
			Udp: ListenConfiguration{true, "0.0.0.0", 514},
			Tcp: ListenConfiguration{true, "0.0.0.0", 514},
		},
		Database: DatabaseConfiguration{
			"localhost", 27017, "syslog-user", "password", "syslog", "logmessages", "SCRAM-SHA-256", 100, 2,
		},
		TimeStamp: TimeLocationConfiguration{true, ""}}

	assert.Equal(t, expected, c)
}

func TestLoadConfigWithOverrideByEnv(t *testing.T) {
	os.Setenv("APP_CONFIG_FILE", "testdata/test-config.yml")

	os.Setenv("SJ_DATABASE_HOST", "my-mongo")
	os.Setenv("SJ_TIMESTAMP_TIMEZONE", "America/Los_Angeles")

	c, err := LoadConfig()
	if err != nil {
		t.Error("Load config failed.", err)
	}

	expected := &Configuration{
		Server: ServerConfiguration{
			Udp: ListenConfiguration{true, "0.0.0.0", 514},
			Tcp: ListenConfiguration{true, "0.0.0.0", 514},
		},
		Database: DatabaseConfiguration{
			"my-mongo", 27017, "syslog-user", "password", "syslog", "logmessages", "SCRAM-SHA-256", 100, 2,
		},
		TimeStamp: TimeLocationConfiguration{true, "America/Los_Angeles"}}

	assert.Equal(t, expected, c)
}
