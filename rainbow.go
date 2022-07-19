package main

import "github.com/gookit/color"

func Rainbowifier(s string) string {
	color.Style{color.FgGreen, color.OpBold}.String()
	const (
		RED           byte = 0x0
		LIGHT_RED     byte = 0x1
		LIGHT_YELLOW  byte = 0x2
		LIGHT_GREEN   byte = 0x3
		LIGHT_BLUE    byte = 0x4
		LIGHT_MAGENTA byte = 0x5
	)
	colorMap := map[byte]color.Style{
		0x0: {color.FgRed},
		0x1: {color.FgLightRed},
		0x2: {color.FgLightYellow},
		0x3: {color.FgLightGreen},
		0x4: {color.FgLightBlue},
		0x5: {color.FgLightMagenta},
	}
	output := ""
	color := byte(0x0)
	for _, c := range s {
		output += colorMap[color].Sprint(string(c))
		if c == ' ' {
			continue
		}
		if color == 0x5 {
			color = 0x0
		} else {
			color++
		}
	}
	return output
}
