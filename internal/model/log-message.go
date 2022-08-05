package model

import (
    "gaofamily/syslog/internal/config"
    log "github.com/sirupsen/logrus"
    "gopkg.in/mcuadros/go-syslog.v2/format"
    "time"
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

    timestamp := logParts["timestamp"].(time.Time)
    if timeZone != nil {
        timestamp.In(timeZone)
    }

    logMessage := LogMessage{logParts["tag"].(string),
        logParts["content"].(string), logParts["facility"].(int), logParts["severity"].(int), logParts["client"].(string),
        logParts["tls_peer"].(string), logParts["hostname"].(string), timestamp}
    return logMessage
}
