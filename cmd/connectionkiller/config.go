package main

import (
	"connection-killer/pkg/keylistener"
	windivertwrapper "connection-killer/pkg/windivertWrapper"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
	"github.com/spf13/viper"
	"github.com/tawesoft/golib/v2/dialog"
)

type Config struct {
	parsedMacro []int
	refreshRate int
	enabled     bool

	quitBtn *systray.MenuItem
	macro   *systray.MenuItem
	state   *systray.MenuItem
}

var config Config
var kl keylistener.Keylogger
var wd *windivertwrapper.WinDivertHandle

func (c *Config) parseConfig() {
	macro := viper.GetStringSlice("macro")
	c.parsedMacro = make([]int, len(macro))
	for i, m := range macro {
		c.parsedMacro[i] = keylistener.NameToCode(m)
	}
	kl.Macro = c.parsedMacro
	config.refreshRate = viper.GetInt("refreshRate")
}

func (c *Config) macroString() string {
	var s string
	for _, m := range c.parsedMacro {
		s += keylistener.CodeToName(m) + "+"
	}
	s = s[:len(s)-1]
	return "Macro : " + s
}

func (c *Config) stateString() string {
	if c.enabled {
		return "STATE : ON"
	}
	return "STATE : OFF"
}

func init() {
	kl = keylistener.NewKeylogger()
	wd = windivertwrapper.NewHandle()

	viper.SetDefault("macro", []string{"ctrl", "alt", "F"})
	viper.SetDefault("refreshRate", 500)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil && err != err.(viper.ConfigFileNotFoundError) {
		dialog.Error("There has been an error : %s", err.Error())
		panic(err.Error())
	}
	err = viper.SafeWriteConfig()
	if err != nil && err != err.(viper.ConfigFileAlreadyExistsError) {
		dialog.Error("There has been an error : %s", err.Error())
		panic(err.Error())
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		config.parseConfig()
	})
	viper.WatchConfig()

	config.parseConfig()
}
