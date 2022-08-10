package syslog_server

import (
	"errors"
	"fmt"
	"gaofamily/syslog/internal/config"
	"gaofamily/syslog/internal/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mcuadros/go-syslog.v2"
)

type server struct {
	config            config.ServerConfiguration
	logMessageChannel chan model.LogMessage
	server            *syslog.Server
	started           bool
}

func NewServer(configuration config.ServerConfiguration, logMessageChannel chan model.LogMessage) (*server, error) {
	if !configuration.Udp.Enabled && !configuration.Tcp.Enabled {
		return nil, errors.New("neither UDP nor TCP enabled")
	}
	if err := validateListenConfiguration(configuration.Udp, "UDP"); err != nil {
		return nil, err
	}
	if err := validateListenConfiguration(configuration.Tcp, "TCP"); err != nil {
		return nil, err
	}

	server := &server{configuration, logMessageChannel, nil, false}
	return server, nil
}

func (s *server) Start() error {
	if s.started {
		return errors.New("syslog server already started")
	}
	logChannel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(logChannel)

	s.server = syslog.NewServer()
	s.server.SetFormat(syslog.Automatic)
	s.server.SetHandler(handler)
	if s.config.Udp.Enabled {
		addr := fmt.Sprintf("%v:%v", s.config.Udp.Address, s.config.Udp.Port)
		err := s.server.ListenUDP(addr)
		if err != nil {
			return err
		}
	}
	if s.config.Tcp.Enabled {
		addr := fmt.Sprintf("%v:%v", s.config.Tcp.Address, s.config.Tcp.Port)
		err := s.server.ListenTCP(addr)
		if err != nil {
			return err
		}
	}

	if err := s.server.Boot(); err != nil {
		return err
	}
	s.started = true

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			logMessage := model.ConvertLogMessage(logParts)

			s.logMessageChannel <- logMessage
		}
	}(logChannel)

	log.Info("Syslog server started.")

	return nil
}

func (s *server) Stop() error {
	if s.started {
		close(s.logMessageChannel)
		if err := s.server.Kill(); err != nil {
			return err
		}
		s.server.Wait()
		s.started = false
		log.Info("Syslog server stopped.")
	}
	return nil
}

func validateListenConfiguration(configuration config.ListenConfiguration, protocol string) error {
	invalidAddress := fmt.Sprintf("invalid %v listening address", protocol)
	invalidPort := fmt.Sprintf("invalid %v listening port", protocol)
	if configuration.Enabled {
		if configuration.Address == "" {
			return errors.New(invalidAddress)
		}
		if configuration.Port <= 0 || configuration.Port > 65536 {
			return errors.New(invalidPort)
		}
	}
	return nil
}
