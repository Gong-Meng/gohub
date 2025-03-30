// Package logger 处理日志相关逻辑
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gongmeng/gohub/app"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 全局 Logger 对象
var Logger *zap.Logger

// InitLogger 日志初始化
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {

	// 获取日志写入介质
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)

	// 设置日志等级，具体请见 config/log.go 文件
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置有误。请修改 config/log.go 文件中的 log.level 配置项")
	}

	// 初始化 core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 设置全局 Logger
	Logger = zap.New(
		core,
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)

	// 设置 Logger 为全局 Logger
	// 将自定义的 logger 替换为全局的 logger
	// zap.L().Fatal() 调用时，就会使用我们自定的 Logger
	zap.ReplaceGlobals(Logger)

}

// getEncoder 设置日志存储格式
func getEncoder() zapcore.Encoder {

	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 时间字段名
		LevelKey:       "level",                        // 日志级别字段名
		NameKey:        "logger",                       // 日志名称字段名
		CallerKey:      "caller",                       // 代码调用，如 paginator/paginator.go:148
		FunctionKey:    zapcore.OmitKey,                // 函数名字段名
		MessageKey:     "message",                      // 日志信息字段名
		StacktraceKey:  "stacktrace",                   // 堆栈字段名
		LineEnding:     zapcore.DefaultLineEnding,      // 换行符
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别编码器 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoder,              // 时间编码器 自定义时间格式 我们自定义为 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // 时长编码器 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	// 本地环境配置
	if app.IsLocal() {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// 本地设置内置的 Console 解码器（支持 stacktrace 换行）
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 线上环境使用 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder 自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 日志记录介质。Gohub 中使用了两种介质，os.Stdout 和文件
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 如果配置了按照日期记录日志文件

	if logType == "daily" {
		// 日志文件按日期分割
		logname := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logname)
	}

	// 滚动日志，详见 config/log.go
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件路径
		MaxSize:    maxSize,   // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackup, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,    // 文件最多保存多少天
		Compress:   compress,  // 是否压缩
	}

	// 配置输出介质
	if app.IsLocal() {
		// 本地开发终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产环境只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}
}
