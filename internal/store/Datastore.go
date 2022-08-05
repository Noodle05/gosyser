package store

import (
    "context"
    "errors"
    "fmt"
    "gaofamily/syslog/internal/config"
    "gaofamily/syslog/internal/model"
    log "github.com/sirupsen/logrus"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "time"
)

type Datastore struct {
    configuration     config.DatabaseConfiguration
    logMessageChannel chan model.LogMessage
    inited            bool
    dbClient          *mongo.Client
    dbContext         context.Context
    flushTimer        *time.Timer
}

func NewDatastore(configuration config.DatabaseConfiguration, logMessageChannel chan model.LogMessage) (*Datastore, error) {
    if configuration.Host == "" {
        return nil, errors.New("invalid database host")
    }
    if configuration.Port <= 0 || configuration.Port > 65536 {
        return nil, errors.New("invalid database port")
    }
    if configuration.Username == "" {
        return nil, errors.New("invalid database username")
    }
    if configuration.Password == "" {
        return nil, errors.New("invalid database password")
    }
    if configuration.Database == "" {
        return nil, errors.New("invalid database name")
    }
    if configuration.AuthMechanism == "" {
        return nil, errors.New("invalid auth mechanism")
    }
    if configuration.BatchSize <= 0 || configuration.BatchSize > 50000 {
        return nil, errors.New("batch size must greater than 0 and less than 50000")
    }

    d := &Datastore{configuration: configuration, logMessageChannel: logMessageChannel}
    return d, nil
}

func (d *Datastore) Start() error {
    ctx := context.TODO()
    d.dbContext = ctx

    credential := options.Credential{
        AuthMechanism: d.configuration.AuthMechanism,
        AuthSource:    d.configuration.Database,
        Username:      d.configuration.Username,
        Password:      d.configuration.Password,
    }
    uri := fmt.Sprintf("mongodb://%v:%d", d.configuration.Host, d.configuration.Port)
    log.Infof("Mongodb URI: %v", uri)
    clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)
    log.Trace("Connection to mongodb")
    client, err := mongo.Connect(ctx, clientOpts)
    log.Debug("Connected to mongodb")
    if err != nil {
        return err
    }
    d.dbClient = client
    log.Trace("Try to ping mongodb")
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        return err
    }
    log.Debug("Ping mongodb success")
    d.inited = true
    syslogDatabase := d.dbClient.Database(d.configuration.Database)
    syslogCollection := syslogDatabase.Collection(d.configuration.Table)
    index := []mongo.IndexModel{
        {
            Keys: bson.M{"facility": -1},
        },
        {
            Keys: bson.M{"severity": -1},
        },
        {
            Keys: bson.M{"hostname": 1},
        },
        {
            Keys: bson.M{"timestamp": -1},
        },
        {
            Keys: bson.M{"client": 1},
        },
        {
            Keys: bson.M{"tag": 1},
        },
    }
    opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
    log.Trace("Creating indexes")
    _, err = syslogCollection.Indexes().CreateMany(ctx, index, opts)
    if err != nil {
        return err
    }
    log.Debug("Create indexes success")
    go func(channel chan model.LogMessage) {
        var logMessages []interface{}
        for logMessage := range channel {
            if d.flushTimer == nil {
                d.flushTimer = time.NewTimer(2 * time.Second)
                go func() {
                    <-d.flushTimer.C
                    if len(logMessages) > 0 {
                        log.Trace("Timeout, flush log messages")
                        d.bulkSave(syslogCollection, ctx, logMessages)
                        logMessages = nil
                    }
                }()
            }
            logMessages = append(logMessages, logMessage)
            if len(logMessages) >= d.configuration.BatchSize {
                log.Trace("Reach batch size, flush log messages")
                d.bulkSave(syslogCollection, ctx, logMessages)
                logMessages = nil
            }
        }
        if len(logMessages) > 0 {
            log.Debug("Has remaining log messages, flush log messages")
            d.bulkSave(syslogCollection, ctx, logMessages)
            logMessages = nil
        }
    }(d.logMessageChannel)

    log.Info("Datastore started.")

    return nil
}

func (d *Datastore) bulkSave(collection *mongo.Collection, ctx context.Context, data []interface{}) {
    log.Debugf("Bulk inserting %d log messages", len(data))
    if d.flushTimer != nil {
        d.flushTimer.Stop()
        d.flushTimer = nil
    }
    _, err := collection.InsertMany(ctx, data)
    if err != nil {
        log.Warn("Error happened when insert log message. ", err)
    }
}

func (d *Datastore) Stop() error {
    if d.inited {
        if err := d.dbClient.Disconnect(d.dbContext); err != nil {
            return err
        }
        d.inited = false
        log.Info("Datastore stopped.")
    }
    return nil
}
