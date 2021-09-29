package conslogging

import "github.com/fatih/color"

var (
	noColor            = makeNoColor()
	cachedColor        = makeColor(color.FgHiGreen)
	metadataModeColor  = makeColor(color.FgHiWhite, color.BgHiBlack)
	phaseColor         = makeNoColor()
	disabledPhaseColor = makeColor(color.FgHiBlack)
	specialPhaseColor  = makeColor(color.FgYellow)
	successColor       = makeColor(color.FgHiGreen)
	warnColor          = makeColor(color.FgHiRed)
	localColor         = makeColor(color.FgHiBlue)
)

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
