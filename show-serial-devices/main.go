package main

import (
	"fmt"

	"go.bug.st/serial/enumerator"
)

// Prints all serial devices and their port name
func main() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		fmt.Println(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
	} else {
		for _, port := range ports {
			fmt.Printf(" - Name : %v | Port : %v\n", port.Product, port.Name)
		}
	}

	fmt.Println("Press the enter key to continue...")
	fmt.Scanln()
}
