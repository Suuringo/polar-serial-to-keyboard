package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"polar-serial-to-usb/keyboard"
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
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("UiShigureLove")
	systray.SetIcon(shigLove)
	systray.SetTooltip("uisgrsuki")
	mquit := systray.AddMenuItem("Quitter", "Quitter")
	go func() {
		<-mquit.ClickedCh
		systray.Quit()
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

func onExit() {
	os.Exit(0)
	return
}
