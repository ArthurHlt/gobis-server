package cli

import (
	"fmt"
	"github.com/cloudfoundry-community/gautocloud"
	"github.com/cloudfoundry-community/gautocloud/cloudenv"
	"github.com/cloudfoundry-community/gautocloud/connectors/generic"
	"github.com/cloudfoundry-community/gautocloud/interceptor/cli/urfave"
	"github.com/cloudfoundry-community/gautocloud/interceptor/configfile"
	"github.com/cloudfoundry-community/gautocloud/loader"
	"github.com/orange-cloudfoundry/gobis-server/server"
	"github.com/urfave/cli"
	"os"
)

var cliInterceptor *urfave.CliInterceptor
var confFileIntercept *configfile.ConfigFileInterceptor

func init() {
	confFileIntercept = configfile.NewConfigFile()
	cliInterceptor = urfave.NewCli()
	gautocloud.RegisterConnector(generic.NewConfigGenericConnector(
		server.GobisServerConfig{},
		confFileIntercept,
		cliInterceptor,
	))
}

type GobisServerApp struct {
	*cli.App
}

func NewApp() *GobisServerApp {
	app := &GobisServerApp{cli.NewApp()}
	app.Name = "gobis-server"
	app.Version = "1.4.0"
	app.Usage = "Create a gobis server based on a config file"
	app.ErrWriter = os.Stderr
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config-path, c",
			Value:  cloudenv.DEFAULT_CONFIG_PATH,
			Usage:  "Path to the config file (This file will not be used in a cloud env like Cloud Foundry, Heroku or kubernetes)",
			EnvVar: cloudenv.LOCAL_CONFIG_ENV_KEY,
		},
		cli.StringFlag{
			Name:  "cert",
			Value: "server.crt",
			Usage: "Path to a cert file or a cert content to enable https server",
		},
		cli.StringFlag{
			Name:  "key",
			Value: "server.key",
			Usage: "Path to a key file or a key content to enable https server",
		},
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: "Log level to use",
		},
		cli.StringFlag{
			Name:  "forward-url, f",
			Usage: "If set all non-found url by gobis will be forwarded to this url",
		},
		cli.BoolFlag{
			Name:  "log-json, j",
			Usage: "Write log in json",
		},
		cli.BoolFlag{
			Name:  "no-color",
			Usage: "Logger will not display colors",
		},
		cli.StringFlag{
			Name:  "lets-encrypt-domains, led",
			Usage: "If set server will use a certificate generated with let's encrypt, value should be your domain(s) (e.g.: --lets-encrypt=example.com[,seconddomain.com]). Host and port will be overwritten to use 0.0.0.0:443",
		},
	}
	return app
}

func (a *GobisServerApp) Run(arguments []string) (err error) {
	a.Action = a.RunServer
	return a.App.Run(arguments)
}
func (a *GobisServerApp) RunServer(c *cli.Context) error {
	cliInterceptor.SetContext(c)

	confPath := c.GlobalString("config-path")
	confFileIntercept.SetConfigPath(confPath)

	var config server.GobisServerConfig
	err := gautocloud.Inject(&config)
	if err != nil {
		if _, ok := err.(loader.ErrGiveService); ok {
			return fmt.Errorf("Configuration cannot be found")
		}
		return err
	}

	gobisServer, err := server.NewGobisServer(&config)
	if err != nil {
		return err
	}
	return gobisServer.Run()
}
