package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"polar/keyboard"

	"github.com/getlantern/systray"
	"go.bug.st/serial"
)

//go:embed shig.ico
var shigLove []byte

func main() {
	systray.Run(onReady, onExit)
	onReady()
}

func onReady() {
	systray.SetTitle("UiShigureLove")
	systray.SetIcon(shigLove)
	systray.SetTooltip("uisgrsuki")
	mquit := systray.AddMenuItem("Quitter", "Quitter")
	go func() {
		<-mquit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		os.Exit(1)
		fmt.Println("Finished quitting")
	}()

	portInfo := promptPort()
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portInfo.Name, mode)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Opened port", portInfo.Name)

	listenString(port, LF, keyboard.SendString)
}

func onExit() {
	return
}
