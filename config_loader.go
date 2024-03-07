package vinculum_agollo

import (
	"encoding/json"
	"errors"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/go-kid/ioc/configure"
	"github.com/go-kid/vinculum-agollo/properties"
	"gopkg.in/yaml.v3"
	"os"
)

type loader struct {
	configPath string
	configJson []byte
	cfg        *config.AppConfig
	client     agollo.Client
}

func (l *loader) LoadConfig(_ string) ([]byte, error) {
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		if l.cfg == nil {
			var cfg = &config.AppConfig{}
			var err error
			if l.configPath != "" {
				l.configJson, err = os.ReadFile(l.configPath)
				if err != nil {
					return nil, err
				}
			}
			if len(l.configJson) == 0 {
				return nil, errors.New("agollo loader need a config")
			}
			err = json.Unmarshal(l.configJson, cfg)
			if err != nil {
				return nil, err
			}
			l.cfg = cfg
		}
		if l.cfg == nil {
			return nil, errors.New("agollo loader need a config")
		}
		return l.cfg, nil
	})
	if err != nil {
		return nil, err
	}
	l.client = client
	cache := client.GetConfigCache(l.cfg.NamespaceName)
	var m = make(map[string]any)
	cache.Range(func(key, value interface{}) bool {
		m[key.(string)] = value
		return true
	})
	expand := properties.PropMapExpand(m)
	return yaml.Marshal(expand)
}

func NewConfigLoader(cfg *config.AppConfig, configPath string, configJson []byte) configure.Loader {
	return &loader{
		configPath: configPath,
		configJson: configJson,
		cfg:        cfg,
	}
}
