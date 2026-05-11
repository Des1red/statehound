package logger

import (
	"fmt"
	"strings"
)

func Status(s string) {
	fmt.Println("[*] " + s)
}

func Success(s string) {
	fmt.Println("[+] " + s)
}

func Failed(s string, err error) {
	if err == nil {
		fmt.Println("[!] " + s)
		return
	}

	if strings.TrimSpace(s) == "" {
		fmt.Println("[!] " + err.Error())
		return
	}

	fmt.Println("[!] " + s + ": " + err.Error())
}

func Warn(s string) {
	fmt.Println("[~] " + s)
}
