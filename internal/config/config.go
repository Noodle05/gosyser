package config

import (
    "github.com/knadh/koanf"
    "github.com/knadh/koanf/parsers/yaml"
    "github.com/knadh/koanf/providers/env"
    "github.com/knadh/koanf/providers/file"
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
)

type Configuration struct {
    Server    ServerConfiguration
    Database  DatabaseConfiguration
    TimeStamp TimeLocationConfiguration
}

type ServerConfiguration struct {
    Udp ListenConfiguration
    Tcp ListenConfiguration
}

type ListenConfiguration struct {
    Enabled bool
    Address string
    Port    int
}

type DatabaseConfiguration struct {
    Host          string
    Port          int
    Username      string
    Password      string
    Database      string
    Table         string
    AuthMechanism string
    BatchSize     int
}

type TimeLocationConfiguration struct {
    Convert  bool
    Timezone string
}

func LoadConfig() (*Configuration, error) {
    configurationFile := os.Getenv("APP_CONFIG_FILE")
    if configurationFile == "" {
        configurationFile = "configs/application.yml"
    }
    log.Infof("Configuration file location: %v", configurationFile)

    k := koanf.New(".")
    if err := k.Load(file.Provider(configurationFile), yaml.Parser()); err != nil {
        return nil, err
    }
    k.Load(env.Provider("SJ_", ".", func(s string) string {
        return strings.Replace(strings.ToLower(
            strings.TrimPrefix(s, "SJ_")), "_", ".", -1)
    }), nil)

    c := &Configuration{}
    if err := k.Unmarshal("", c); err != nil {
        return nil, err
    }

    return c, nil
}
