package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Debug(msg string, fields ...zap.Field) {
	defaultLogger.getLogger().Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	defaultLogger.getLogger().Info(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	defaultLogger.getLogger().Error(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	defaultLogger.getLogger().Warn(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.getLogger().Fatal(msg, fields...)
}

func DebugF(msg string, args ...interface{}) {
	defaultLogger.getSugarLogger().Debugf(msg, args...)
}
func InfoF(msg string, args ...interface{}) {
	defaultLogger.getSugarLogger().Infof(msg, args...)
}
func ErrorF(msg string, args ...interface{}) {
	defaultLogger.getSugarLogger().Errorf(msg, args...)
}
func WarnF(msg string, args ...interface{}) {
	defaultLogger.getSugarLogger().Warnf(msg, args...)
}
func FatalF(msg string, args ...interface{}) {
	defaultLogger.getSugarLogger().Fatalf(msg, args...)
}

type LogConfig struct {
	Level      string `json:"level"`       // Level 最低日志等级，DEBUG<INFO<WARN<ERROR<FATAL 例如：info-->收集info等级以上的日志
	FileName   string `json:"file_name"`   // FileName 日志文件位置
	MaxSize    int    `json:"max_size"`    // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge     int    `json:"max_age"`     // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups int    `json:"max_backups"` // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
}

var (
	defaultLogger = NewLogger(&LogConfig{
		Level:      "DEBUG",
		MaxAge:     1,
		MaxBackups: 1,
		FileName:   "xx.log",
	})
)

type Logger struct {
	cfg *LogConfig
	*zap.Logger
}

func NewLogger(cfg *LogConfig) *Logger {
	l := &Logger{cfg: cfg}
	l.Init()
	return l
}

func (l *Logger) Init() {
	fileWriter := l.getLogWriter(l.cfg.FileName, l.cfg.MaxSize, l.cfg.MaxBackups, l.cfg.MaxAge)
	consoleWriter := zapcore.AddSync(os.Stdout)
	writerSyncer := zapcore.NewMultiWriteSyncer(fileWriter, consoleWriter)
	// 获取日志编码格式
	encoder := l.getEncoder()
	level := l.getLevel()
	// 创建一个将日志写入 WriteSyncer 的核心。
	core := zapcore.NewCore(encoder, writerSyncer, level)
	l.Logger = zap.New(core, zap.AddCaller())
}

// 负责日志写入的位置
func (l *Logger) getLogWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer。
	// 它试图变得智能：如果 io.Writer 的具体类型实现了 WriteSyncer，我们将使用现有的 Sync 方法。
	// 如果没有，我们将添加一个无操作同步。

	return zapcore.AddSync(lumberJackLogger)
}

// 负责设置 encoding 的日志格式
func (l *Logger) getEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConfig := zap.NewProductionEncoderConfig()

	// 设置每个日志条目使用的键。如果有任何键为空，则省略该条目的部分。

	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "time"
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 以 package/file:行 的格式 序列化调用程序，从完整路径中删除除最后一个目录外的所有目录。
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func (l *Logger) getLevel() *zapcore.Level {
	// 获取日志最低等级，即>=该等级，才会被写入。
	var level = new(zapcore.Level)
	err := level.UnmarshalText([]byte(l.cfg.Level))
	if err != nil {
		return level
	}
	return level
}

func (l *Logger) getSugarLogger() *zap.SugaredLogger {
	return l.Sugar().WithOptions(zap.AddCallerSkip(1))
}

func (l *Logger) getLogger() *zap.Logger {
	return l.WithOptions(zap.AddCallerSkip(1))
}
