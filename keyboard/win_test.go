package keyboard

import "testing"

func TestStringString(t *testing.T) {
	SendString(string([]byte{65, 30, 80, 13, 10})) // Should write A9P\n
}
