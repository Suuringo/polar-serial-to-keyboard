package main

import (
	"log"

	"go.bug.st/serial"
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
