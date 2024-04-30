package vinculum_agollo

import (
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/go-kid/ioc/app"
	"github.com/go-kid/ioc/syslog"
)

func Plugin(cfg *config.AppConfig, configPath string, marshaller ...Marshaller) app.SettingOption {
	client, namespaces, err := NewAgolloClient(cfg, configPath)
	if err != nil {
		syslog.Panicf("create agollo client failed: %+v", err)
	}
	return app.Options(
		app.SetConfigLoader(NewConfigLoader(client, namespaces, marshaller...)),
		app.SetComponents(NewSpy(client)),
	)
}
