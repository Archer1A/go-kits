package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
	_multipleErrMsg     = "Multiple errors without a key."
)

var zapLogger *zap.Logger

func Init(logPath string, enab zapcore.Level) {
	jack := zapcore.AddSync(&lumberjack.Logger{
		Filename: logPath,
		Compress: true,
	})
	stdout := zapcore.AddSync(os.Stdout)
	multiSync := zapcore.NewMultiWriteSyncer(jack, stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//encoder := zapcore.NewConsoleEncoder(encoderConfig)

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(jsonEncoder, multiSync, enab)

	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	defer zapLogger.Sync()

}

type logger struct {
	ctx     context.Context
	traceId string
	_logger *zap.Logger
	spanId  string
	pspanId string
}

func NewLogger(ctx context.Context) *logger {
	return newLogger(ctx)
}

func newLogger(ctx context.Context) *logger {
	var traceId, spanId, pspanId string
	if value := ctx.Value("traceId"); value != nil {
		traceId = value.(string)
	}
	if value := ctx.Value("spanId"); value != nil {
		spanId = value.(string)
	}
	if value := ctx.Value("pspanId"); value != nil {
		pspanId = value.(string)
	}
	return &logger{
		ctx:     ctx,
		traceId: traceId,
		_logger: zapLogger,
		spanId:  spanId,
		pspanId: pspanId,
	}
}

func (l *logger) Infow(msg string, field ...interface{}) {
	l.logger(zapcore.InfoLevel, msg, field)
}

func (l *logger) Debugw(msg string, field ...interface{}) {
	l.logger(zapcore.DebugLevel, msg, field)
}

func (l *logger) Errorw(msg string, field ...interface{}) {
	l.logger(zapcore.ErrorLevel, msg, field)
}

func Infow(ctx context.Context, msg string, field ...interface{}) {
	l := newLogger(ctx)
	l.logger(zapcore.InfoLevel, msg, field)
}

func Debugw(ctx context.Context, msg string, field ...interface{}) {
	l := newLogger(ctx)
	l.logger(zapcore.DebugLevel, msg, field)
}

func Errorw(ctx context.Context, msg string, field ...interface{}) {
	l := newLogger(ctx)
	l.logger(zapcore.ErrorLevel, msg, field)
}

func Infof(ctx context.Context, fmtTemplate string, args ...interface{}) {
	l := newLogger(ctx)
	message := getFormatMessage(fmtTemplate, args)
	l.logger(zapcore.InfoLevel, message, nil)

}

func getFormatMessage(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}
	if msg != "" {
		return fmt.Sprintf(msg, args)
	}
	if len(args) == 1 {
		if str, ok := args[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(args...)
}

func (log *logger) logger(level zapcore.Level, msg string, fields []interface{}) {
	if log.spanId != "" {
		fields = append(fields, zap.String("spanid", log.spanId))
	}
	if log.traceId != "" {
		fields = append(fields, zap.String("traceid", log.traceId))
	}
	if log.pspanId != "" {
		fields = append(fields, zap.String("pspanid", log.pspanId))
	}
	if check := log._logger.Check(level, msg); check != nil {
		check.Write(log.sweetenFields(fields)...)
	}
}

func (log *logger) sweetenFields(fields []interface{}) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	var (
		zapFields = make([]zap.Field, 0, len(fields))
		seenError bool
	)

	for i := 0; i < len(fields); {
		if field, ok := fields[i].(zap.Field); ok {
			zapFields = append(zapFields, field)
			i++
			continue

		}
		if err, ok := fields[i].(error); ok {
			if !seenError {
				zapFields = append(zapFields, zap.Error(err))
			} else {
				log._logger.Error(_multipleErrMsg, zap.Error(err))
			}
			i++
			continue
		}

		if i == len(fields)-1 {
			log._logger.Error(_nonStringKeyErrMsg, zap.Any("ignored", fields[i]))
			i++
			continue
		}
		key, val := fields[i], fields[i+1]
		if s, ok := key.(string); ok {
			zapFields = append(zapFields, zap.Any(s, val))
		} else {
			zapFields = append(zapFields, zap.Any(fmt.Sprintf("%v", key), val))
		}
		i += 2
	}
	return zapFields
}

func (l *logger) Info() {
	l._logger.Info("hello", zap.String("name", "vic"))
}
