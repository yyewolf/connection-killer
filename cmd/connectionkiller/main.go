package main

import (
	"os"
	"time"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Connection Killer")
	systray.SetTooltip("I kill your connection when you want to.")

	config.state = systray.AddMenuItem(config.stateString(), config.stateString())
	config.macro = systray.AddMenuItem(config.macroString(), config.macroString())
	config.quitBtn = systray.AddMenuItem("Quit", "Quit the whole app")
	QuitButton()

	// Sets the icon of a menu item. Only available on Mac and Windows.
	config.quitBtn.SetIcon(icon.Data)

	for true {
		if kl.CheckMacro() {
			config.enabled = !config.enabled
			switch config.enabled {
			case true:
				wd.EnableFilter()
			case false:
				wd.DisableFilter()
			}
		}

		// Update the systray
		config.macro.SetTitle(config.macroString())
		config.state.SetTitle(config.stateString())
		time.Sleep(time.Duration(config.refreshRate) * time.Millisecond)
	}
}

func onExit() {
	os.Exit(0)
}
