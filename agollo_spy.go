package vinculum_agollo

import (
	"encoding/json"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/go-kid/vinculum"
	"github.com/go-kid/vinculum-agollo/properties"
	"gopkg.in/yaml.v3"
	"os"
)

type spy struct {
	config *config.AppConfig
	client agollo.Client
	ch     chan []byte
}

func NewSpy(client agollo.Client, c *config.AppConfig) vinculum.Spy {
	return &spy{
		config: c,
		client: client,
		ch:     make(chan []byte),
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

func (s *spy) Change() <-chan []byte {
	return s.ch
}

func (s *spy) Close() error {
	s.client.Close()
	return nil
}

func (s *spy) OnChange(event *storage.ChangeEvent) {
}

func (s *spy) OnNewestChange(event *storage.FullChangeEvent) {
	expand := properties.PropMapExpand(event.Changes)
	out, _ := yaml.Marshal(expand)
	s.ch <- out
}
