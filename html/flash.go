package html

import "context"

func NewFlashContext(parent context.Context, s string) context.Context {
	return context.WithValue(parent, flashCtxKey, s)
}

func FlashFromContext(ctx context.Context) string {
	f, ok := ctx.Value(flashCtxKey).(string)
	if ok {
		return f
	}
	return ""
}

func NewErrorContext(parent context.Context, err error) context.Context {
	return context.WithValue(parent, errorCtxKey, err)
}

func ErrorFromContext(ctx context.Context) error {
	e, ok := ctx.Value(errorCtxKey).(error)
	if ok {
		return e
	}
	return nil
}
