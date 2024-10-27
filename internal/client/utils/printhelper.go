package utils

import (
	"fmt"
	"strings"
	"time"
)

const terminalColor = "\033[32m"

// Printer – структура для работы с текстом.
type Printer struct{}

// Print – метод вывода текстовых данных на экран.
func (p *Printer) Print(s string) {
	roboPrint(s)
}

// Scan – обёртка для fmt.Scan.
func (p *Printer) Scan(a ...interface{}) (int, error) {
	n, err := fmt.Scanln(a...)
	return n, err
}

func roboPrint(text string) {
	arr := strings.Split(text, "")
	for _, char := range arr {
		fmt.Print(terminalColor, char)
		<-time.After(25 * time.Millisecond)
	}
	fmt.Println()
}

func (p *Printer) PrintLogo() {
	fmt.Print(terminalColor, " ██ ▄█▀▓█████ ▓█████  ██▓███  ▓█████  ██▀███  \n ██▄█▒ ▓█   ▀ ▓█   ▀ ▓██░  ██▒▓█   ▀ ▓██ ▒ ██▒\n▓███▄░ ▒███   ▒███   ▓██░ ██▓▒▒███   ▓██ ░▄█ ▒\n▓██ █▄ ▒▓█  ▄ ▒▓█  ▄ ▒██▄█▓▒ ▒▒▓█  ▄ ▒██▀▀█▄  \n▒██▒ █▄░▒████▒░▒████▒▒██▒ ░  ░░▒████▒░██▓ ▒██▒\n▒ ▒▒ ▓▒░░ ▒░ ░░░ ▒░ ░▒▓▒░ ░  ░░░ ▒░ ░░ ▒▓ ░▒▓░\n░ ░▒ ▒░ ░ ░  ░ ░ ░  ░░▒ ░      ░ ░  ░  ░▒ ░ ▒░\n░ ░░ ░    ░      ░   ░░          ░     ░░   ░ \n░  ░      ░  ░   ░  ░            ░  ░   ░     \n                                              \n")
}
