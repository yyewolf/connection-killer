package main

import (
	"os"

	"github.com/getlantern/systray"
)

func QuitButton() {
	go func() {
		<-config.quitBtn.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()
}
