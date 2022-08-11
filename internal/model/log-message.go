package model

import (
	"gaofamily/syslog/internal/config"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type LogMessage struct {
	Tag       string
	Content   string
	Facility  int
	Severity  int
	Client    string
	TlsPeer   string
	Hostname  string
	Timestamp time.Time
}

const TAG_KEY_NAME = "tag"
const CONTENT_KEY_NAME = "content"
const FACILITY_KEY_NAME = "facility"
const SEVERITY_KEY_NAME = "severity"
const CLIENT_KEY_NAME = "client"
const TLS_PEER_KEY_NAME = "tls_peer"
const HOSTNAME_KEY_NAME = "hostname"
const TIMESTAMP_KEY_NAME = "timestamp"

var timeZone *time.Location

func SetTimestampConvertToLocal(c config.TimeLocationConfiguration) {
	if c.Convert {
		if c.Timezone != "" {
			tz, err := time.LoadLocation(c.Timezone)
			if err != nil {
				log.Warnf("Invalid timezone: %v, use local timezone.", c.Timezone)
				timeZone = time.Now().Location()
			} else {
				timeZone = tz
			}
		} else {
			timeZone = time.Now().Location()
		}
		log.Infof("Time zone will be set to %v\n", timeZone)
	} else {
		timeZone = nil
		log.Infof("Time zone will not change")
	}
}

func ConvertLogMessage(logParts format.LogParts) LogMessage {
	logMessage := LogMessage{
		extractStringPart(logParts, TAG_KEY_NAME),
		extractStringPart(logParts, CONTENT_KEY_NAME),
		extractIntPart(logParts, FACILITY_KEY_NAME),
		extractIntPart(logParts, SEVERITY_KEY_NAME),
		extractStringPart(logParts, CLIENT_KEY_NAME),
		extractStringPart(logParts, TLS_PEER_KEY_NAME),
		extractStringPart(logParts, HOSTNAME_KEY_NAME),
		extractTimePart(logParts, TIMESTAMP_KEY_NAME),
	}
	return logMessage
}

func extractStringPart(logParts format.LogParts, keyName string) string {
	if tmp, ok := logParts[keyName]; ok {
		return tmp.(string)
	}
	return ""
}

func extractIntPart(logParts format.LogParts, keyName string) int {
	if tmp, ok := logParts[keyName]; ok {
		return tmp.(int)
	}
	return 0
}

func extractTimePart(logParts format.LogParts, keyName string) time.Time {
	var timestamp time.Time
	if tmp, ok := logParts[keyName]; ok {
		timestamp = tmp.(time.Time)
	} else {
		timestamp = time.Now()
	}
	if timeZone != nil {
		timestamp.In(timeZone)
	}
	return timestamp
}
