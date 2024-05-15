package logger

import (
	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path"
)

type Handler struct {
	handler slog.Handler
}

var Log *slog.Logger

func init() {
	Log = slog.New(newZlogHandler("./log"))
}
func newZlogHandler(dir string) slog.Handler {
	logger := newZapLogger(dir)
	handler := slogzap.Option{Level: slog.LevelDebug, AddSource: true, Logger: logger}.NewZapHandler()
	return handler
}

// var Log *slog.Logger
func newZapLogger(dir string) *zap.Logger {
	if dir == "" {
		dir = path.Join("logs", "kernel")
	}
	var coreArr []zapcore.Core

	format := "2006-01-02 15:04:05.000"
	encoderConfig := zap.NewProductionEncoderConfig()              // NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(format) // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder        // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder        // 显示完整文件路径
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 日志级别
	errLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.InfoLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})

	// info文件writeSyncer
	debugFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(dir, "debug.log"), // 日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    256,                         // 文件大小限制,单位MB
		MaxBackups: 2,                           // 最大保留日志文件数量
		MaxAge:     2,                           // 日志文件保留天数
		Compress:   false,                       // 是否压缩处理
	})
	debugFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(debugFileWriteSyncer, zapcore.AddSync(os.Stdout)), debugLevel) // 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(dir, "info.log"), // 日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    256,                        // 文件大小限制,单位MB
		MaxBackups: 5,                          // 最大保留日志文件数量
		MaxAge:     1,                          // 日志文件保留天数
		Compress:   false,                      // 是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), infoLevel) // 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(dir, "error.log"), // 日志文件存放目录
		MaxSize:    256,                         // 文件大小限制,单位MB
		MaxBackups: 5,                           // 最大保留日志文件数量
		MaxAge:     1,                           // 日志文件保留天数
		Compress:   false,                       // 是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), errLevel) // 第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, debugFileCore)
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
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
