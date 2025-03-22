package log

import "go.uber.org/zap"

type Logger struct {
	base  *zap.Logger
	sugar *zap.SugaredLogger
}

func NewLogger(base *zap.Logger) *Logger {
	return &Logger{
		base:  base,
		sugar: base.Sugar(),
	}
}

func (l *Logger) Base() *zap.Logger {
	return l.base
}

func (l *Logger) With(args ...interface{}) Logger {
	if l.sugar == nil {
		return Logger{}
	}
	sl := l.sugar.With(args...)
	return Logger{
		base:  sl.Desugar(),
		sugar: sl,
	}
}

func (l *Logger) Debugw(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Debugw(msg, args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Debug(args...)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Debugf(format, args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Info(args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Infof(format, args...)
	}
}

func (l *Logger) Infow(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Infow(msg, args...)
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Warn(args...)
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Warnf(format, args...)
	}
}

func (l *Logger) Warnw(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Warnw(msg, args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Error(args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Errorf(format, args...)
	}
}

func (l *Logger) Errorw(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Errorw(msg, args...)
	}
}

func (l *Logger) Panicw(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Panicw(msg, args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Panic(args...)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Panicf(format, args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Fatal(args...)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Fatalf(format, args...)
	}
}

func (l *Logger) Fatalw(msg string, args ...interface{}) {
	if l.sugar != nil {
		l.sugar.Fatalw(msg, args...)
	}
}
