package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gin-api/pkg/app"
	"gin-api/pkg/config"
	"path"
	"runtime"
	"strings"
	"time"
)

//RuntimeLog 记录错误日志
func RuntimeLog(text string) {
	filename := app.RuntimeLogFile()
	GetLogger(filename).Error(text)
}

//AccessLog  记录访问日志
func AccessLog(name string, text interface{}) {
	var txt string
	//将 text 处理成 string 类型
	if str , ok := text.(string); ok {
		txt = str
	} else {
		txt = fmt.Sprintf("%+v",text)
	}

	filename := app.AccessLogFile()
	//应用如果是cli这种形式,因为没有 url, 因此访问日志就不能基于url来构造了
	if filename == "" {
		filename = getFilename()
	}

	if len(txt) > 0 {
		GetLogger(filename).Info(name, zap.String("info", txt))
	} else {
		GetLogger(filename).Info(name)
	}
}

//Log 记录普通日志(可自定义存放日志的目录)
func Log(name string, text interface{}, logFile... string) {
	filename := getFilename(logFile...)
	GetLogger(filename).Sugar().Info(name + " : ", text)
}

//LogIf 记录错误日志
func LogIf(name string, err error, logFile... string) bool {
	if err != nil {
		filename := getFilename(logFile...)
		GetLogger(filename).Sugar().Info(name + " : ", err.Error())
		return false
	}
	return true
}

func getFilename(logFile... string) string {
	filename := ""
	if len(logFile) > 0 {
		filename = logFile[0]
	} else {
		//pc, file, line, _ := runtime.Caller(1)
		//funcName := path.Base(runtime.FuncForPC(pc).Name())
		_, file, _, _ := runtime.Caller(2)
		filename      = strings.Replace(file, app.GetRootPath(), "", -1)
		filename      = strings.Replace(filename, path.Ext(filename), ".log", -1)
	}
	return app.UserLogFile(filename)
}

//GetLogger 初始化日志驱动
func GetLogger(filename string) *zap.Logger {
	// 获取日志写入介质
	writeSyncer := getLogWriter(filename)

	// 设置日志等级，具体请见 config/log.go 文件
	logLevel := new(zapcore.Level)
	if logLevel.UnmarshalText([]byte(config.GetString("log.level"))) != nil {
		*logLevel = zapcore.DebugLevel
	}

	// 初始化 core
	core := zapcore.NewCore(getEncoder(), writeSyncer, zapcore.DebugLevel)

	// 初始化 Logger
	logger := zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),        // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才打印调用栈
	)

	// 将自定义的 logger 替换为全局的 logger
	// zap.L().Fatal() 调用时，就会使用我们自定的 Logger
	zap.ReplaceGlobals(logger)

	return logger
}

// getEncoder 设置日志存储格式
func getEncoder() zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用，如 paginator/paginator.go:148
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoder,              // 时间格式，我们自定义为 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	// 使用内置的 Console 编码器(支持换行)
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// customTimeEncoder 自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 日志记录
func getLogWriter(filename string) zapcore.WriteSyncer {
	maxSize   := config.GetInt("log.max_size")
	maxBackup := config.GetInt("log.max_backup")
	maxAge    := config.GetInt("log.max_age")
	compress  := config.GetBool("log.compress")
	logType   := config.GetString("log.type")

	//按日期记录日志
	if logType == "daily" {
		current := time.Now().Format("2006-01-02")
		ext := path.Ext(filename)
		if ext == "" {
			filename = fmt.Sprintf("%s-%s.log", filename, current)
		} else {
			newExt   := fmt.Sprintf("-%s%s", current, ext)
			filename = strings.Replace(filename, ext, newExt, -1)
		}
	}

	// 滚动日志，详见 config/log.go
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	return zapcore.AddSync(lumberJackLogger)
}
