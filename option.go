package vinculum_agollo

import (
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/go-kid/ioc/app"
	"github.com/go-kid/vinculum"
)

func Plugin(cfg *config.AppConfig, configPath string, configJson []byte) app.SettingOption {
	agolloLoader := NewConfigLoader(cfg, configPath, configJson)
	client := agolloLoader.(*loader).client
	spy := NewSpy(client, nil)
	return app.Options(
		vinculum.Refresher,
		app.SetConfig("x"),
		app.SetConfigLoader(agolloLoader),
		app.SetComponents(spy),
	)
}
