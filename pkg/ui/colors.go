package ui

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

const AsciiLogo = `
  ____  ____    ____ _     ___ 
 |  _ \/ ___|  / ___| |   |_ _|
 | |_) \___ \ | |   | |    | | 
 |  __/ ___) || |___| |___ | | 
 |_|   |____/  \____|_____|___|
`

func PrintBanner(subtitle string) {
	fmt.Println("==================================================================")
	fmt.Printf("%s%s%s\n", ColorCyan, AsciiLogo, ColorReset)
	fmt.Printf("  %s%s%s %s(ThingsBoard Stack Phase 1)%s\n", ColorBold, subtitle, ColorReset, ColorDim, ColorReset)
	fmt.Println("==================================================================")
}

func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s%sℹ INFO:%s %s\n", ColorBold, ColorBlue, ColorReset, msg)
}

func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s%s✔ SUCCESS:%s %s\n", ColorBold, ColorGreen, ColorReset, msg)
}

func Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s%s▲ WARNING:%s %s\n", ColorBold, ColorYellow, ColorReset, msg)
}

func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s%s✖ ERROR:%s %s\n", ColorBold, ColorRed, ColorReset, msg)
}
