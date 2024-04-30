package vinculum_agollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/go-kid/properties"
	"github.com/go-kid/vinculum"
)

type spy struct {
	client agollo.Client
	ch     chan<- properties.Properties
}

func (s *spy) RegisterChannel(ch chan<- properties.Properties) {
	s.ch = ch
}

func NewSpy(client agollo.Client) vinculum.Spy {
	return &spy{
		client: client,
	}
}

func (s *spy) Init() error {
	s.client.AddChangeListener(s)
	return nil
}

func (s *spy) OnChange(event *storage.ChangeEvent) {
	p := properties.New()
	for path, change := range event.Changes {
		p.Set(path, change.NewValue)
	}
	s.ch <- p
}

func (s *spy) OnNewestChange(event *storage.FullChangeEvent) {
}
