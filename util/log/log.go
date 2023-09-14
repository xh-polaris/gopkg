package log

import (
	"context"
	"io"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/zeromicro/go-zero/core/logx"
)

type logxLogger struct {
}

type HlogLogger struct {
	*logxLogger
}

type KlogLogger struct {
	*logxLogger
}

func NewHlogLogger() *HlogLogger {
	return &HlogLogger{
		&logxLogger{},
	}
}

func NewKlogLogger() *KlogLogger {
	return &KlogLogger{
		&logxLogger{},
	}
}

func (l *logxLogger) Debug(v ...interface{}) {
	getLogger().Debug(v...)
}

func (l *logxLogger) Info(v ...interface{}) {
	getLogger().Info(v...)
}

func (l *logxLogger) Error(v ...interface{}) {
	getLogger().Error(v...)
}

func (l *logxLogger) Trace(v ...interface{}) {
	getLogger().Debug(v...)
}

func (l *logxLogger) Notice(v ...interface{}) {
	getLogger().Info(v...)
}

func (l *logxLogger) Warn(v ...interface{}) {
	getLogger().Error(v...)
}

func (l *logxLogger) Fatal(v ...interface{}) {
	getLogger().Error(v...)
}

func (l *logxLogger) Debugf(format string, v ...interface{}) {
	getLogger().Debugf(format, v...)
}

func (l *logxLogger) Infof(format string, v ...interface{}) {
	getLogger().Infof(format, v...)
}

func (l *logxLogger) Errorf(format string, v ...interface{}) {
	getLogger().Errorf(format, v...)
}

func (l *logxLogger) Tracef(format string, v ...interface{}) {
	getLogger().Debugf(format, v...)
}

func (l *logxLogger) Noticef(format string, v ...interface{}) {
	getLogger().Infof(format, v...)
}

func (l *logxLogger) Warnf(format string, v ...interface{}) {
	getLogger().Errorf(format, v...)
}

func (l *logxLogger) Fatalf(format string, v ...interface{}) {
	getLogger().Errorf(format, v...)
}

func (l *logxLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Debugf(format, v...)
}

func (l *logxLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Debugf(format, v...)
}

func (l *logxLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Infof(format, v...)
}

func (l *logxLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Infof(format, v...)
}

func (l *logxLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Errorf(format, v...)
}

func (l *logxLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Errorf(format, v...)
}

func (l *logxLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	getLoggerCtx(ctx).Errorf(format, v...)
}

func (l *logxLogger) SetOutput(writer io.Writer) {
}

func (l *HlogLogger) SetLevel(level hlog.Level) {
}

func (l *KlogLogger) SetLevel(level klog.Level) {
}

func getLoggerCtx(ctx context.Context) logx.Logger {
	return logx.WithContext(ctx).WithCallerSkip(1)
}

func getLogger() logx.Logger {
	return logx.WithCallerSkip(1)
}

func CtxInfo(ctx context.Context, format string, v ...any) {
	getLoggerCtx(ctx).Infof(format, v...)
}

func Info(format string, v ...any) {
	getLogger().Infof(format, v...)
}

func CtxError(ctx context.Context, format string, v ...any) {
	getLoggerCtx(ctx).Errorf(format, v...)
}

func Error(format string, v ...any) {
	getLogger().Errorf(format, v...)
}

func CtxDebug(ctx context.Context, format string, v ...any) {
	getLoggerCtx(ctx).Debugf(format, v...)
}
