package earthfile2llb

import (
	"log"

	"github.com/earthly/earthly/earthfile2llb/parser"
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

func (l *debugListener) EnterStmtWords(ctx *parser.StmtWordsContext) {
	log.Printf("StmtWords: %s\n", ctx.GetText())
}
