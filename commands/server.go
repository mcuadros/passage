package commands

import (
	"fmt"
	"net"
	"os"
	"os/signal"

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
	RPCAddr    string
	RPCServer  *server.RPCServer

	done chan bool
}

func NewServerCommand() *ServerCommand {
	return &ServerCommand{
		Server: server.NewServer(),
		done:   make(chan bool),
	}
}

func (c *ServerCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "runs a passage server",
		RunE:  c.Execute,
	}

	cmd.Flags().StringVar(&c.ConfigFile, "config", "", "config file (default is $HOME/.passage.yaml)")
	cmd.Flags().StringVar(&c.LogFile, "log-file", "", "log file")
	cmd.Flags().StringVar(&c.LogLevel, "log-level", "info", "max log level enabled")
	cmd.Flags().StringVar(&c.RPCAddr, "rpc-addr", "/tmp/passage.sock", "passage rpc server address, is an unix socket.")

	return cmd
}

func (c *ServerCommand) Execute(cmd *cobra.Command, args []string) error {
	c.handleSignals()
	if err := c.setupLogging(); err != nil {
		return err
	}

	if err := c.setupServer(); err != nil {
		return err
	}

	if err := c.setupRPCServer(); err != nil {
		return err
	}

	<-c.done
	log15.Info("server stopped successfully")
	return nil
}

func (c *ServerCommand) setupServer() error {
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

	return nil
}

func (c *ServerCommand) setupRPCServer() error {
	log15.Debug("rpc server started", "addr", c.RPCAddr)
	a, err := net.ResolveUnixAddr("unix", c.RPCAddr)
	if err != nil {
		return err
	}

	c.RPCServer = server.NewRPCServer(c.Server)
	if err := c.RPCServer.Listen(a); err != nil {
		return err
	}

	return nil
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
		handler = log15.MultiHandler(handler, log15.Must.FileHandler(c.LogFile, format))
	}

	handler = log15.LvlFilterHandler(lvl, handler)

	if lvl == log15.LvlDebug {
		handler = log15.CallerFileHandler(log15.LvlFilterHandler(lvl, handler))
	}

	log15.Root().SetHandler(handler)
	return nil
}

func (c *ServerCommand) handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		if err := c.stop(); err != nil {
			log15.Error("error stopping services", "err", err.Error())
		}
	}()
}

func (c *ServerCommand) stop() error {
	if err := c.Server.Close(); err != nil {
		return err
	}

	if err := c.RPCServer.Close(); err != nil {
		return err
	}

	c.done <- true
	return nil
}
