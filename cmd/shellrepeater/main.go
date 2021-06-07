package main

import (
	"context"

	"github.com/earthly/earthly/debugger/server"
	"github.com/earthly/earthly/slog"

	"github.com/sirupsen/logrus"
)

const addr = "0.0.0.0:8373"

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()
	log := slog.GetLogger(ctx).With("app", "shellrepeater")

	x := server.NewServer(addr, log)
	x.Start()
}
