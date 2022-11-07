package progressbar

import "strings"

var progressChars = []string{
	" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

// ProgressBar returns a progress bar as a string.
func ProgressBar(progress, width int) string {
	if progress > 100 {
		progress = 100
	}
	if progress < 0 {
		progress = 0
	}
	builder := make([]string, 0, width)
	fullChars := progress * width / 100
	blankChars := width - fullChars - 1
	deltaProgress := ((progress * width) % 100) * len(progressChars) / 100
	for i := 0; i < fullChars; i++ {
		builder = append(builder, progressChars[len(progressChars)-1])
	}
	if progress != 100 {
		builder = append(builder, progressChars[deltaProgress])
	}
	for i := 0; i < blankChars; i++ {
		builder = append(builder, progressChars[0])
	}
	return strings.Join(builder, "")
}
