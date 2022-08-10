package model

import (
	"gaofamily/syslog/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"testing"
	"time"
)

func TestSetTimestampConvertToLocal(t *testing.T) {
	config := config.TimeLocationConfiguration{Convert: true, Timezone: ""}

	SetTimestampConvertToLocal(config)

	assert.True(t, config.Convert)
	assert.Empty(t, config.Timezone)
}

func TestConvertLogMessage(t *testing.T) {
	config := config.TimeLocationConfiguration{Convert: true, Timezone: ""}

	SetTimestampConvertToLocal(config)

	timestamp := time.Now()

	logParts := format.LogParts{
		"facility":  3,
		"severity":  5,
		"client":    "dummy client",
		"content":   "dummy content",
		"hostname":  "some host",
		"timestamp": timestamp,
		"tag":       "Some tag",
		"tls_peer":  "",
	}
	actual := ConvertLogMessage(logParts)

	expected := LogMessage{Client: "dummy client", Content: "dummy content", Hostname: "some host",
		Facility: 3, Severity: 5, Timestamp: timestamp, Tag: "Some tag"}

	assert.Equal(t, expected, actual)
}
