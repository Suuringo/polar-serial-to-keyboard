package main

import (
	"fmt"
	"log"

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

func promptPort() (port *enumerator.PortDetails) {
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

	return ports[choice]
}

func listenString(port serial.Port, delimiter byte, callback func(string)) {
	var n int
	var err error
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
				callback(string(stringBytes))
				stringBytes = nil
			}
		}
	}
}
