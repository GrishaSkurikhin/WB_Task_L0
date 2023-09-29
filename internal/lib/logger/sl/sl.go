package sl

import (
	"golang.org/x/exp/slog"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Debug(msg string) slog.Attr {
	return slog.Attr{
		Key:   "message",
		Value: slog.StringValue(msg),
	}
}

func Warn(msg string) slog.Attr {
	return slog.Attr{
		Key:   "warning",
		Value: slog.StringValue(msg),
	}
}
