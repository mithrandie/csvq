package color

import (
	"strconv"
	"strings"
)

var UseEscapeSequences = false

const PlainColor = -1

const (
	FGBlack = iota + 30
	FGRed
	FGGreen
	FGYellow
	FGBlue
	FGMagenta
	FGCyan
	FGWhite
)

const (
	FGBrightBlack = iota + 90
	FGBrightRed
	FGBrightGreen
	FGBrightYellow
	FGBrightBlue
	FGBrightMagenta
	FGBrightCyan
	FGBrightWhite
)

func StripEscapeSequence(s string) string {
	if strings.HasPrefix(s, "\033[") && strings.HasSuffix(s, "\033[0m") {
		startIdx := 5
		if strings.Index(s, ";1m") == 4 {
			startIdx = 7
		}

		runes := []rune(s)
		s = string(runes[startIdx : len(runes)-4])
	}
	return s
}

func Colorize(s string, color int, bold bool) string {
	if !UseEscapeSequences || color == PlainColor {
		return s
	}
	boldStr := ""
	if bold {
		boldStr = ";1"
	}
	return "\033[" + strconv.Itoa(color) + boldStr + "m" + s + "\033[0m"
}

func Black(s string) string {
	return Colorize(s, FGBlack, false)
}

func Red(s string) string {
	return Colorize(s, FGRed, false)
}

func Green(s string) string {
	return Colorize(s, FGGreen, false)
}

func Yellow(s string) string {
	return Colorize(s, FGYellow, false)
}

func Blue(s string) string {
	return Colorize(s, FGBlue, false)
}

func Magenta(s string) string {
	return Colorize(s, FGMagenta, false)
}

func Cyan(s string) string {
	return Colorize(s, FGCyan, false)
}

func White(s string) string {
	return Colorize(s, FGWhite, false)
}

func BrightBlack(s string) string {
	return Colorize(s, FGBrightBlack, false)
}

func BrightRed(s string) string {
	return Colorize(s, FGBrightRed, false)
}

func BrightGreen(s string) string {
	return Colorize(s, FGBrightGreen, false)
}

func BrightYellow(s string) string {
	return Colorize(s, FGBrightYellow, false)
}

func BrightBlue(s string) string {
	return Colorize(s, FGBrightBlue, false)
}

func BrightMagenta(s string) string {
	return Colorize(s, FGBrightMagenta, false)
}

func BrightCyan(s string) string {
	return Colorize(s, FGBrightCyan, false)
}

func BrightWhite(s string) string {
	return Colorize(s, FGBrightWhite, false)
}

func BlackB(s string) string {
	return Colorize(s, FGBlack, true)
}

func RedB(s string) string {
	return Colorize(s, FGRed, true)
}

func GreenB(s string) string {
	return Colorize(s, FGGreen, true)
}

func YellowB(s string) string {
	return Colorize(s, FGYellow, true)
}

func BlueB(s string) string {
	return Colorize(s, FGBlue, true)
}

func MagentaB(s string) string {
	return Colorize(s, FGMagenta, true)
}

func CyanB(s string) string {
	return Colorize(s, FGCyan, true)
}

func WhiteB(s string) string {
	return Colorize(s, FGWhite, true)
}

func BrightBlackB(s string) string {
	return Colorize(s, FGBrightBlack, true)
}

func BrightRedB(s string) string {
	return Colorize(s, FGBrightRed, true)
}

func BrightGreenB(s string) string {
	return Colorize(s, FGBrightGreen, true)
}

func BrightYellowB(s string) string {
	return Colorize(s, FGBrightYellow, true)
}

func BrightBlueB(s string) string {
	return Colorize(s, FGBrightBlue, true)
}

func BrightMagentaB(s string) string {
	return Colorize(s, FGBrightMagenta, true)
}

func BrightCyanB(s string) string {
	return Colorize(s, FGBrightCyan, true)
}

func BrightWhiteB(s string) string {
	return Colorize(s, FGBrightWhite, true)
}

func Error(s string) string {
	return RedB(s)
}

func Warn(s string) string {
	return YellowB(s)
}

func Info(s string) string {
	return GreenB(s)
}
