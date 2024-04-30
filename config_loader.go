package vinculum_agollo

import (
	"github.com/apolloconfig/agollo/v4"
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/ioc/util/fas"
	"github.com/go-kid/properties"
	"gopkg.in/yaml.v3"
)

type Marshaller func(in interface{}) (out []byte, err error)

type loader struct {
	client     agollo.Client
	marshal    Marshaller
	namespaces []string
}

func (l *loader) LoadConfig() ([]byte, error) {
	prop := properties.New()
	for _, namespace := range l.namespaces {
		cache := l.client.GetConfigCache(namespace)
		cache.Range(func(key, value interface{}) bool {
			prop.Set(key.(string), value)
			return true
		})
	}

	return l.marshal(prop)
}

func NewConfigLoader(client agollo.Client, namespaces []string, marshaller ...Marshaller) configure.Loader {
	return &loader{
		client:     client,
		marshal:    fas.TernaryOp(len(marshaller) != 0, marshaller[0], yaml.Marshal),
		namespaces: namespaces,
	}
}
