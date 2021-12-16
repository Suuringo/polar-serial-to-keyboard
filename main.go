package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"polar/keyboard"

	"github.com/getlantern/systray"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const LF = 10

func isByteInBuff(buf []byte, bt byte) bool {
	for _, b := range buf {
		if bt == b {
			return true
		}
	}

	return false
}

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

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	fmt.Println("Choose port:")
	for i, port := range ports {
		fmt.Printf("%d: %v, (%v)\n", i, port.Name, port.Product)
	}

	var choice int
	fmt.Scanf("%d", &choice)
	for choice < 0 || choice >= len(ports) {
		fmt.Println("Choose port:")
		for i, port := range ports {
			fmt.Printf("%d: %v\n", i, port)
		}
		fmt.Scanf("%d", &choice)
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(ports[choice].Name, mode)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Opened port", ports[choice].Name)

	var n int
	buff := make([]byte, 16)
	var stringBytes []byte
	for {
		port.ResetOutputBuffer()
		n, err = port.Read(buff)
		if err != nil {
			log.Println("Error ", err)
		} else {
			stringBytes = append(stringBytes, buff[0:n]...)
			if isByteInBuff(buff[0:n], LF) {
				keyboard.SendString(string(stringBytes))
				// fmt.Println(string(stringBytes))
				stringBytes = nil
			}
		}
	}
}

func onExit() {
	return
}
