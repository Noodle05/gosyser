package main

import (
    "context"
    "gaofamily/syslog/internal/config"
    "gaofamily/syslog/internal/model"
    "gaofamily/syslog/internal/store"
    syslogServer "gaofamily/syslog/internal/syslog-server"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func init() {
    log.SetFormatter(&log.TextFormatter{
        DisableColors: true,
        FullTimestamp: true,
    })
}

func main() {
    appConfig, err := config.LoadConfig()
    if err != nil {
        panic(err)
    }

    model.SetTimestampConvertToLocal(appConfig.TimeStamp)

    persistChannel := make(chan model.LogMessage)

    datastore, err := store.NewDatastore(appConfig.Database, persistChannel)
    if err != nil {
        panic(err)
    }
    defer func() {
        err = datastore.Stop()
        if err != nil {
            log.Error("Error happened when close datastore.", err)
        }
    }()
    if err := datastore.Start(); err != nil {
        panic(err)
    }

    server, err := syslogServer.NewServer(appConfig.Server, persistChannel)
    if err != nil {
        panic(err)
    }
    defer func() {
        err = server.Stop()
        if err != nil {
            log.Error("Error happened when close syslog server.", err)
        }
    }()
    if err := server.Start(); err != nil {
        panic(err)
    }

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

    <-stop

    if err := server.Stop(); err != nil {
        log.Warn("Error happened when stop syslog server.", err)
    }
    if err := datastore.Stop(); err != nil {
        log.Warn("Error happened when stop datastore.", err)
    }

    _, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
}
