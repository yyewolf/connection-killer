package windivertwrapper

import "github.com/williamfhe/godivert"

func (w *WinDivertHandle) checkPacket(wd *godivert.WinDivertHandle, packetChan <-chan *godivert.Packet) {
	for packet := range packetChan {
		// Only send those who passthrough when the filter is OFF
		if !w.FilterAll {
			packet.Send(wd)
		}
	}
}

func (w *WinDivertHandle) Close() {
	w.Handler.Close()
}

func (w *WinDivertHandle) EnableFilter() {
	w.FilterAll = true
}

func (w *WinDivertHandle) DisableFilter() {
	w.FilterAll = false
}
