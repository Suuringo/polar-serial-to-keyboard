package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"polar/keyboard"

	"github.com/getlantern/systray"
	"go.bug.st/serial"
)

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func main() {
	systray.Run(onReady, onExit)
	onReady()
}

func onReady() {
	systray.SetTitle("UiShigureLove")
	systray.SetIcon(getIcon("./shig.ico"))
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
