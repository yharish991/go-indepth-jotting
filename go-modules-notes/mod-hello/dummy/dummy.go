package dummy

import "fmt"

// Dummy is simple string
type Dummy string

// DummyValue simply prints the value on screeen
func (d Dummy) DummyValue() {
	fmt.Println(d)
}
