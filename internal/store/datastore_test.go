package store

import (
	"gaofamily/syslog/internal/config"
	"gaofamily/syslog/internal/model"
	"reflect"
	"testing"
)

func TestNewDatastore(t *testing.T) {
	logMessageChannel := make(chan model.LogMessage)
	type args struct {
		configuration     config.DatabaseConfiguration
		logMessageChannel chan model.LogMessage
	}
	tests := []struct {
		name    string
		args    args
		want    *datastore
		wantErr bool
	}{
		{"success", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "password", "syslog", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, &datastore{config.DatabaseConfiguration{
			"localhost", 27017, "user", "password", "syslog", "logmessages", "SHA-1", 100, 2,
		}, logMessageChannel, false, nil, nil, nil, nil,
		}, false},
		{"missing database host", args{
			config.DatabaseConfiguration{
				"", 27017, "user", "password", "syslog", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"invalid database port", args{
			config.DatabaseConfiguration{
				"localhost", 0, "user", "password", "syslog", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"missing username", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "", "password", "syslog", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"missing password", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "", "syslog", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"missing database", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "password", "", "logmessages", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"missing table", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "password", "syslog", "", "SHA-1", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"missing auth mechanism", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "password", "syslog", "logmessages", "", 100, 2,
			},
			logMessageChannel,
		}, nil, true},
		{"invalid batch size", args{
			config.DatabaseConfiguration{
				"localhost", 27017, "user", "password", "syslog", "logmessages", "SHA-1", 0, 2,
			},
			logMessageChannel,
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDatastore(tt.args.configuration, tt.args.logMessageChannel)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDatastore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDatastore() got = %v, want %v", got, tt.want)
			}
		})
	}
}
