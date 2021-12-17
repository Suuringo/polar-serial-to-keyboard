package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"polar/keyboard"
	"strings"

	"github.com/getlantern/systray"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type Config struct {
	Port       string `json:"port",omitempty`
	DeviceName string `json:"deviceName",omitempty`
}

//go:embed shig.ico
var shigLove []byte

func main() {
	systray.Run(onReady, onExit)
}

func readConfig() (conf Config) {
	buf, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error opening config.json ", err)
	}
	err = json.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatal("Error parsing config.json ", err)
	}
	return
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

	var portName string

	conf := readConfig()
	if conf.DeviceName != "" {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			log.Fatal("Error getting serial port list! ", err)
		}
		if len(ports) == 0 {
			log.Fatal("No serial ports found! ")
		}

		for _, port := range ports {
			if strings.HasPrefix(strings.ToLower(port.Product), strings.ToLower(conf.DeviceName)) {
				portName = port.Name
				break
			}
		}
	} else if conf.Port != "" {
		portName = conf.Port
	} else {
		log.Fatal("deviceName and port not specified in conf.json!")
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal("Error opening serial port ", portName, err)
	}

	log.Println("Opened port", portName)

	listenString(port, LF, keyboard.SendString)
}

func onExit() {
	return
}
