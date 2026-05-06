package logger

import (
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path/filepath"
)

type Handler struct {
	handler slog.Handler
}

var Log *slog.Logger

func init() {
	Log = slog.New(newZlogHandler(os.Getenv("GO_CDN_CHECK_LOG_DIR")))
}

func newZlogHandler(dir string) slog.Handler {
	logger := newZapLogger(dir)
	handler := slogzap.Option{Level: slog.LevelDebug, AddSource: true, Logger: logger}.NewZapHandler()
	return handler
}

// var Log *slog.Logger
func newZapLogger(dir string) *zap.Logger {
	var coreArr []zapcore.Core

	format := "2006-01-02 15:04:05.000"
	encoderConfig := zap.NewProductionEncoderConfig()              // NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(format) // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder        // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder        // 显示完整文件路径
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	coreArr = append(coreArr, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel))
	if dir != "" {
		errLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		})
		infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev < zap.ErrorLevel && lev >= zap.InfoLevel
		})
		debugLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev == zap.DebugLevel
		})

		debugFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "debug.log"),
			MaxSize:    256,
			MaxBackups: 2,
			MaxAge:     2,
			Compress:   false,
		})
		infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "info.log"),
			MaxSize:    256,
			MaxBackups: 5,
			MaxAge:     1,
			Compress:   false,
		})
		errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "error.log"),
			MaxSize:    256,
			MaxBackups: 5,
			MaxAge:     1,
			Compress:   false,
		})
		coreArr = append(coreArr,
			zapcore.NewCore(encoder, debugFileWriteSyncer, debugLevel),
			zapcore.NewCore(encoder, infoFileWriteSyncer, infoLevel),
			zapcore.NewCore(encoder, errorFileWriteSyncer, errLevel),
		)
	}
	return zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) // zap.AddCaller()为显示文件名和行号，可省略
}

//var Log *zap.SugaredLogger
//
//const (
//	output_dir = "./logs/"
//	out_path   = "normal.log"
//	err_path   = "err.log"
//)
//
//func init() {
//	_, err := os.Stat(output_dir)
//	if err != nil {
//		if os.IsNotExist(err) {
//			err := os.Mkdir(output_dir, os.ModePerm)
//			if err != nil {
//				fmt.Printf("mkdir failed![%v]\n", err)
//			}
//		}
//	}
//
//	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
//	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
//		MessageKey:    "msg",
//		LevelKey:      "level",
//		TimeKey:       "ts",
//		CallerKey:     "caller",
//		StacktraceKey: "trace",
//		LineEnding:    zapcore.DefaultLineEnding,
//		EncodeLevel:   zapcore.LowercaseLevelEncoder,
//		EncodeCaller:  zapcore.ShortCallerEncoder,
//		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//			enc.AppendString(t.Format("2006-01-02 15:04:05"))
//		},
//		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
//			enc.AppendInt64(int64(d) / 1000000)
//		},
//	})
//
//	// 实现两个判断日志等级的interface
//	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
//		return true
//	})
//
//	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
//		return lvl >= zapcore.InfoLevel
//	})
//
//	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
//		return lvl >= zapcore.WarnLevel
//	})
//
//	infoHook_1 := os.Stdout
//	infoHook_2 := getWriter(out_path)
//	errorHook := getWriter(err_path)
//
//	// 最后创建具体的Logger
//	// INFO级别的日志同时写入文件和stdout, ERROR级别的日志单独写入文件, DEBUG级别的日志写入stdout
//	core := zapcore.NewTee(
//		zapcore.NewCore(encoder, zapcore.AddSync(infoHook_1), debugLevel),
//		zapcore.NewCore(encoder, zapcore.AddSync(infoHook_2), infoLevel),
//		zapcore.NewCore(encoder, zapcore.AddSync(errorHook), warnLevel),
//	)
//
//	// 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
//	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
//	Log = logger.Sugar()
//	defer logger.Sync()
//}
//
//func getWriter(filename string) io.Writer {
//	// 生成rotatelogs的Logger 实际生成的文件名 demo.log-YY-mm-dd
//	// 保存7天内的日志，每24小时分割一次日志
//	hook, err := rotatelogs.New(
//		// 没有使用go风格反人类的format格式
//		output_dir+filename+"-%Y-%m-%d",
//		//rotatelogs.WithLinkName(filename),
//		rotatelogs.WithMaxAge(time.Hour*24*7),
//		rotatelogs.WithRotationTime(time.Hour*24),
//	)
//	if err != nil {
//		panic(err)
//	}
//	return hook
//}
