package vinculum_agollo

import (
	"fmt"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type tCmp struct {
	Config *Config `refreshScope:""`
	AB     int     `prop:"AB" refreshScope:""`
}

func (c *tCmp) OnScopeChange(path string) error {
	fmt.Println("onChange", path, c.Config, c.AB)
	return nil
}

type Config struct {
	Value  int    `yaml:"value"`
	Value2 string `yaml:"value2"`
	Value3 []int  `yaml:"value3"`
}

func (c *Config) Prefix() string {
	return "Test"
}

func TestInit(t *testing.T) {
	iocApp := ioc.RunTest(t,
		Plugin(nil, "agollo.json", nil),
		app.SetComponents(
			&tCmp{},
		),
	)
	defer iocApp.Close()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGHUP,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc
}
