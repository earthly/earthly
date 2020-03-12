package earthfile2llb

import (
	"log"

	"github.com/vladaionescu/earthly/earthfile2llb/parser"
)

type debugListener struct {
	*parser.BaseEarthParserListener
}

func newDebugListener() *debugListener {
	return new(debugListener)
}

func (l *debugListener) EnterTarget(ctx *parser.TargetContext) {
	log.Printf("Target: %s\n", ctx.GetText())
}

func (l *debugListener) EnterStmt(ctx *parser.StmtContext) {
	log.Printf("Stmt: %s\n", ctx.GetText())
}

func (l *debugListener) EnterFromStmt(ctx *parser.FromStmtContext) {
	log.Printf("FromStmt: %s\n", ctx.GetText())
}

func (l *debugListener) EnterTargetName(ctx *parser.TargetNameContext) {
	log.Printf("TargetName: %s\n", ctx.GetText())
}

func (l *debugListener) EnterCommandName(ctx *parser.CommandNameContext) {
	log.Printf("CommandName: %s\n", ctx.GetText())
}

func (l *debugListener) EnterFlag(ctx *parser.FlagContext) {
	log.Printf("Flag: %s\n", ctx.GetText())
}

func (l *debugListener) EnterStmtWords(ctx *parser.StmtWordsContext) {
	log.Printf("StmtWords: %s\n", ctx.GetText())
}

func (l *debugListener) EnterArg(ctx *parser.ArgContext) {
	log.Printf("Arg: %s\n", ctx.GetText())
}

// func (l *debugListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
// 	log.Printf("Every: %s\n", ctx.GetText())
// }
