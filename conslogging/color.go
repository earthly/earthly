package conslogging

import "github.com/fatih/color"

var (
	noColor            = makeNoColor()
	cachedColor        = makeColor(color.FgGreen)
	metadataModeColor  = makeColor(color.FgHiWhite, color.BgHiBlack)
	phaseColor         = makeColor(color.FgCyan)
	disabledPhaseColor = makeColor(color.FgHiBlack)
	specialPhaseColor  = makeColor(color.FgYellow)
	successColor       = makeColor(color.FgGreen)
	warnColor          = makeColor(color.FgHiRed)
	localColor         = makeColor(color.FgHiBlue)
	helpColor          = makeColor(color.FgMagenta)
)

var availablePrefixColors = []*color.Color{
	makeColor(color.FgBlue),
	makeColor(color.FgMagenta),
	makeColor(color.FgCyan),
	makeColor(color.FgYellow),
	makeColor(color.FgHiGreen),
	makeColor(color.FgHiBlue),
	makeColor(color.FgHiMagenta),
	makeColor(color.FgHiCyan),
	makeColor(color.FgHiYellow),
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
