package stringutil

// ProcessParamsAndQuotes takes in a slice of strings, and rearranges the slices
// depending on quotes and parenthesis.
//
// For example "hello ", "wor(", "ld)" becomes "hello ", "wor( ld)".
func ProcessParamsAndQuotes(args []string) []string {
	curQuote := rune(0)
	allowedQuotes := map[rune]rune{
		'"':  '"',
		'\'': '\'',
		'(':  ')',
	}
	ret := make([]string, 0, len(args))
	var newArg []rune
	for _, arg := range args {
		for _, char := range arg {
			newArg = append(newArg, char)
			if curQuote == 0 {
				_, isQuote := allowedQuotes[char]
				if isQuote {
					curQuote = char
				}
				continue
			}
			if char == allowedQuotes[curQuote] {
				curQuote = rune(0)
			}
		}
		if curQuote == 0 {
			ret = append(ret, string(newArg))
			newArg = []rune{}
			continue
		}
		// Unterminated quote - join up two args into one.
		// Add a space between joined-up args.
		newArg = append(newArg, ' ')
	}
	if curQuote != 0 {
		// Unterminated quote case.
		newArg = newArg[:len(newArg)-1] // remove last space
		ret = append(ret, string(newArg))
	}

	return ret
}
