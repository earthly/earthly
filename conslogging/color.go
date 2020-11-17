package conslogging

import "github.com/fatih/color"

var noColor = makeNoColor()
var cachedColor = makeColor(color.FgHiGreen)
var paramsColor = makeColor(color.FgHiBlack)
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

func makeColor(att color.Attribute) *color.Color {
	c := color.New()
	c.Add(att)
	return c
}

func makeNoColor() *color.Color {
	c := color.New()
	c.DisableColor()
	return c
}
