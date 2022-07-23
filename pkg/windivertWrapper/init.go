package windivertwrapper

import (
	"github.com/tawesoft/golib/v2/dialog"
	"github.com/williamfhe/godivert"
)

func NewHandle() *WinDivertHandle {
	// Just loading the DLLs here.
	godivert.LoadDLL(WinDivertDLLPath, WinDivertDLL32Path)

	// Starts the new handle, we want every packet here.
	winDivert, err := godivert.NewWinDivertHandle("true")
	if err != nil {
		dialog.Error("There has been an error : %s\nCheck for admin permissions.", err.Error())
		panic(err.Error())
	}

	// Sends the packet to the filter.
	packetChan, err := winDivert.Packets()
	if err != nil {
		panic(err)
	}

	handle := &WinDivertHandle{
		FilterAll: false,
		Handler:   winDivert,
	}

	go handle.checkPacket(winDivert, packetChan)

	return handle
}
