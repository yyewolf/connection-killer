package windivertwrapper

import "github.com/williamfhe/godivert"

var (
	WinDivertDLLPath   = "WinDivert.dll"
	WinDivertDLL32Path = "WinDivert32.dll"
)

type WinDivertHandle struct {
	FilterAll bool
	Handler   *godivert.WinDivertHandle
}
