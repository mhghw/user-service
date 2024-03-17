package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], logger *logrus.Entry) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base:   handler,
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

func ApplyCommandWithResultDecorators[H any, R any](handler CommandWithResultHandler[H, R], logger *logrus.Entry) CommandWithResultHandler[H, R] {
	return commandWithResultLoggingDecorator[H, R]{
		base:   handler,
		logger: logger,
	}
}

type CommandWithResultHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
