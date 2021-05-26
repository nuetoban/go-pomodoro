// +build !windows

package terminal

import "fmt"

func ClearLine() {
	fmt.Print("\033[F\033[K")
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func Bell() {
	fmt.Print("\a")
}
