package conslogging

import "github.com/fatih/color"

var noColor = makeNoColor()
var cachedColor = makeColor(color.FgHiGreen)
var metadataModeColor = makeColor(color.FgHiWhite, color.BgHiBlack)
var successColor = makeColor(color.FgHiGreen)
var warnColor = makeColor(color.FgHiRed)

var availablePrefixColors = []*color.Color{
	makeColor(color.FgBlue),
	makeColor(color.FgMagenta),
	makeColor(color.FgCyan),
	makeColor(color.FgYellow),
	makeColor(color.FgGreen),
	makeColor(color.FgHiBlue),
	makeColor(color.FgHiMagenta),
	makeColor(color.FgHiCyan),
	makeColor(color.FgHiYellow),
	makeColor(color.FgHiWhite),
}

func makeColor(attrs ...color.Attribute) *color.Color {
	c := color.New()
	for _, attr := range attrs {
		c.Add(attr)
	}
	return c
}

func makeNoColor() *color.Color {
	c := color.New()
	c.DisableColor()
	return c
}
