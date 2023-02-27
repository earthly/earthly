package parser

// GetLexerModeNames returns the generated mode names.
func GetLexerModeNames() []string {
	EarthLexerInit()
	return earthlexerLexerStaticData.modeNames
}

// GetLexerSymbolicNames returns the generated token names.
func GetLexerSymbolicNames() []string {
	EarthLexerInit()
	return earthlexerLexerStaticData.symbolicNames
}

// GetLexerLiteralNames returns the generated literal names.
func GetLexerLiteralNames() []string {
	EarthLexerInit()
	return earthlexerLexerStaticData.literalNames
}
