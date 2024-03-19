package vinculum_agollo

import (
	"encoding/json"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/vinculum"
	"os"
)

type spy struct {
	config *config.AppConfig
	client agollo.Client
	ch     chan vinculum.UpdateHandler
}

func NewSpy(client agollo.Client, c *config.AppConfig) vinculum.Spy {
	return &spy{
		config: c,
		client: client,
		ch:     make(chan vinculum.UpdateHandler),
	}
}

func (s *spy) Init() error {
	if s.client == nil {
		client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
			if s.config == nil {
				cfgBytes, err := os.ReadFile("agollo.json")
				if err != nil {
					return nil, err
				}
				s.config = &config.AppConfig{}
				err = json.Unmarshal(cfgBytes, s.config)
				if err != nil {
					return nil, err
				}
			}
			return s.config, nil
		})
		if err != nil {
			return err
		}
		s.client = client
	}
	s.client.AddChangeListener(s)
	return nil
}

func (s *spy) Change() <-chan vinculum.UpdateHandler {
	return s.ch
}

func (s *spy) Close() error {
	s.client.Close()
	return nil
}

func (s *spy) OnChange(event *storage.ChangeEvent) {
	s.ch <- func(binder configure.Binder) error {
		for path, change := range event.Changes {
			binder.Set(path, change.NewValue)
		}
		return nil
	}
}

func (s *spy) OnNewestChange(event *storage.FullChangeEvent) {
}
