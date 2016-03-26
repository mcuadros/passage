package commands

import (
	"fmt"

	"github.com/mcuadros/passage/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
	"gopkg.in/inconshreveable/log15.v2"
)

type ServerCommand struct {
	LogLevel   string
	LogFile    string
	ConfigFile string
	Config     *server.Config
	Server     *server.Server
}

func NewServerCommand() *ServerCommand {
	return &ServerCommand{
		Server: server.NewServer(),
	}
}

func (c *ServerCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "A brief description of your command",
		Long:  `A longeraaaa description .`,
		RunE:  c.Execute,
	}

	cmd.Flags().StringVar(&c.ConfigFile,
		"config", "", "config file (default is $HOME/.passage.yaml)",
	)

	cmd.Flags().StringVar(&c.LogLevel,
		"level", "info", "max log level enabled",
	)

	return cmd
}

func (c *ServerCommand) Execute(cmd *cobra.Command, args []string) error {
	if err := c.setupLogging(); err != nil {
		return err
	}

	if err := c.readConfig(); err != nil {
		return err
	}

	log15.Info("configuration file loaded", "file", viper.ConfigFileUsed())
	if err := c.loadConfig(); err != nil {
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log15.Info("configuration file re-loaded", "file", e.Name)
		if err := c.loadConfig(); err != nil {
			log15.Error("unable to read/load config", "error", err.Error())
			return
		}
	})

	select {}
}

func (c *ServerCommand) readConfig() error {
	if c.ConfigFile != "" {
		viper.SetConfigFile(c.ConfigFile)
	}

	viper.WatchConfig()
	viper.SetConfigName(".passage")
	viper.AddConfigPath("$HOME")

	return viper.ReadInConfig()
}

func (c *ServerCommand) loadConfig() error {
	if err := viper.Unmarshal(&c.Config); err != nil {
		return err
	}

	if err := c.Server.Load(c.Config); err != nil {
		return err
	}

	return nil
}

func (c *ServerCommand) setupLogging() error {
	lvl, err := log15.LvlFromString(c.LogLevel)
	if err != nil {
		return fmt.Errorf("unknown log level name %q", c.LogLevel)
	}

	handler := log15.StdoutHandler
	format := log15.LogfmtFormat()

	if c.LogFile != "" {
		handler = log15.MultiHandler(
			handler,
			log15.Must.FileHandler(c.LogFile, format),
		)
	}

	log15.Root().SetHandler(log15.CallerFileHandler(
		log15.LvlFilterHandler(lvl, handler),
	))

	return nil
}
