package log

var (
	defaultLogger *Logger
	logBuilder    *Builder
	With          = defaultLogger.With
	Debug         = defaultLogger.Debug
	Info          = defaultLogger.Info
	Warn          = defaultLogger.Warn
	Error         = defaultLogger.Error
	Fatal         = defaultLogger.Fatal
	Debugf        = defaultLogger.Debugf
	Fatalf        = defaultLogger.Fatalf
	Infof         = defaultLogger.Infof
	Warnf         = defaultLogger.Warnf
	Errorf        = defaultLogger.Errorf
	Panic         = defaultLogger.Panic
	Debugw        = defaultLogger.Debugw
	Infow         = defaultLogger.Infow
	ErrorW        = defaultLogger.Errorw
	Panicf        = defaultLogger.Panicf
	Warnw         = defaultLogger.Warnw
)

// todo 实现，只需要基本的default实现即可，控制打印到某个默认文件
type Builder struct {
}

// 初始化default log
func init() {

}

func GetDefaultLogger() *Logger {
	return defaultLogger
}
