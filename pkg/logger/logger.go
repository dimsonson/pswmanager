package logger

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/dimsonson/pswmanager/internal/masterserver/settings"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Log struct {
	*zerolog.Logger
}

func (ll *Log) L() *Log {
	l := zerolog.New(
	zerolog.ConsoleWriter{
		Out:          os.Stderr,
		TimeFormat:   "2006/01/02 15:04:05",
		FormatCaller: func(i interface{}) string { return fmt.Sprintf("%s:", i) },
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s%s%s:", settings.ColorNewYellow, i, settings.ColorReset)
		},
		PartsOrder: []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.MessageFieldName, zerolog.CallerFieldName}}).
		With().Caller().Logger()

	stdlog.SetFlags(stdlog.Lshortfile)
	stdlog.SetOutput(log.Logger)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return &Log{&l}
}

//"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"

// type Logger interface {
// 	Debug(msg string, fields map[string]interface{})
// 	Info(msg string, fields map[string]interface{})
// 	Warn(msg string, fields map[string]interface{})
// 	Error(msg string, fields map[string]interface{})
// 	Fatal(msg string, fields map[string]interface{})
// }

// type ZeroLogger struct {
// 	logger *zerolog.Logger
// 	ctx    context.Context
// }

// func NewZeroLogger(loggerType string, ctx context.Context) *ZeroLogger {
// 	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

// 	log.Logger = log.Output(zerolog.ConsoleWriter{
// 		Out:          os.Stderr,
// 		TimeFormat:   "2006/01/02 15:04:05",
// 		FormatCaller: func(i interface{}) string { return fmt.Sprintf("%s:", i) },
// 		FormatMessage: func(i interface{}) string {
// 			return fmt.Sprintf("%s%s%s:", settings.ColorNewYellow, i, settings.ColorReset)
// 		},
// 		PartsOrder: []string{zerolog.TimestampFieldName, "level" /*zerolog.LevelFieldName*/, zerolog.MessageFieldName, zerolog.CallerFieldName}}).
// 		With().Caller().Logger()

// 	stdlog.SetFlags(stdlog.Lshortfile)
// 	stdlog.SetOutput(log.Logger)
// 	zerolog.SetGlobalLevel(zerolog.DebugLevel)

// 	return &ZeroLogger{logger: &logger, ctx: ctx}
// }

// func (l *ZeroLogger) Debug(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)
// 	l.logger.With().Fields(fields).Logger().Level(zerolog.DebugLevel)

// }

// func (l *ZeroLogger) Info(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)
// 	l.logger.With().Fields(fields).Logger().Level(zerolog.InfoLevel)
// }

// func (l *ZeroLogger) Warn(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)
// 	l.logger.With().Fields(fields).Logger().Level(zerolog.WarnLevel)
// }

// func (l *ZeroLogger) Error(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.With().Fields(fields).Logger().Level(zerolog.ErrorLevel)
// }

// func (l *ZeroLogger) Fatal(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)
// 	l.logger.With().Fields(fields).Logger().Level(zerolog.FatalLevel)
// }

// func (l *ZeroLogger) addContextCommonFields(fields map[string]interface{}) {
// 	if l.ctx != nil {
// 		for k, v := range l.ctx.Value("commonFields").(map[string]interface{}) {
// 			if _, ok := fields[k]; !ok {
// 				fields[k] = v
// 			}
// 		}
// 	}
// }

// type ZapLogger struct {
// 	logger *zap.Logger
// 	ctx    context.Context
// }

// func NewZapLogger(loggerType string, ctx context.Context) *ZapLogger {
// 	logger, _ := zap.NewProduction()

// 	return &ZapLogger{logger: logger, ctx: ctx}
// }

// func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.Debug("", zap.Any("args", fields))
// }

// func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.Info("", zap.Any("args", fields))
// }

// func (l *ZapLogger) Warn(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.Warn("", zap.Any("args", fields))
// }

// func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.Error("", zap.Any("args", fields))
// }

// func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
// 	l.addContextCommonFields(fields)

// 	l.logger.Fatal("", zap.Any("args", fields))
// }

// func (l *ZapLogger) addContextCommonFields(fields map[string]interface{}) {
// 	if l.ctx != nil {
// 		for k, v := range l.ctx.Value("commonFields").(map[string]interface{}) {
// 			if _, ok := fields[k]; !ok {
// 				fields[k] = v
// 			}
// 		}
// 	}
// }

// type LoggerWrapper struct {
// 	logger Logger
// }

// func NewLoggerWrapper(loggerType string, ctx context.Context) *LoggerWrapper {
// 	var logger Logger

// 	switch loggerType {
// 	case "zerolog":
// 		logger = NewZeroLogger(loggerType, ctx)
// 	case "zap":
// 		logger = NewZapLogger(loggerType, ctx)
// 	default:
// 		logger = NewZeroLogger(loggerType, ctx)
// 	}

// 	return &LoggerWrapper{logger: logger}
// }

// func (lw *LoggerWrapper) Debug(msg string, fields map[string]interface{}) {
// 	lw.logger.Debug(msg, fields)
// }

// func (lw *LoggerWrapper) Info(msg string, fields map[string]interface{}) {
// 	lw.logger.Info(msg, fields)
// }

// func (lw *LoggerWrapper) Warn(msg string, fields map[string]interface{}) {
// 	lw.logger.Warn(msg, fields)
// }

// func (lw *LoggerWrapper) Error(msg string, fields map[string]interface{}) {
// 	lw.logger.Error(msg, fields)
// }

// func (lw *LoggerWrapper) Fatal(msg string, fields map[string]interface{}) {
// 	lw.logger.Fatal(msg, fields)
// }
