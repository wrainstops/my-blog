package common

import (
	"io"
	"log/slog"
	"os"
	"time"
)

func InitSlogLogger() *slog.Logger {
	option := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}

	var handler slog.Handler
	handler = slog.NewJSONHandler(io.MultiWriter(Global_LogFile, os.Stdout), option)

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
