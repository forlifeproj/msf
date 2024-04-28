package fllog

type BuilderInterface interface {
	SetConsole(status bool) BuilderInterface
	GetConsole() bool
	SetLine(line bool)
	GetLine() bool
	SetLogLevel(logLevel string) BuilderInterface
	GetLogLevel() string
	SetLogFileName(logFileName string) BuilderInterface
	GetLogFileName() string
	SetMaxSize(MaxSize int) BuilderInterface
	SetMaxAge(MaxAge int) BuilderInterface
	SetMaxBackups(MaxBackups int) BuilderInterface
}
