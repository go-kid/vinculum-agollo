package vinculum_agollo

import (
	"encoding/json"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func NewAgolloClient(appConfig *config.AppConfig, configPath string) (agollo.Client, []string, error) {
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		if appConfig != nil {
			return appConfig, nil
		}
		if configPath != "" {
			jsonConfig, err := os.ReadFile(configPath)
			if err != nil {
				return nil, errors.Wrap(err, "read agollo config error")
			}
			appConfig = &config.AppConfig{}
			err = json.Unmarshal(jsonConfig, appConfig)
			if err != nil {
				return nil, errors.Wrap(err, "")
			}

		}
		if appConfig == nil {
			return nil, errors.New("agollo loader need a config")
		}
		return appConfig, nil
	})
	if err != nil {
		return nil, nil, err
	}
	namespaces := strings.Split(appConfig.NamespaceName, ",")
	return client, namespaces, nil
}
