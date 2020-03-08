package logger

// gLog is the static logger
var gLog *Logger

func init() {
	gLog = NewLogger("", 0, 0, Info, 0666)
	gLog.EnableStdOut(true)
	gLog.EnableAutoFlush(true)
}

// SetLevel function set a new log level
func SetLevel(level LogLevel) {
	gLog.level = level
}

// GetDefaultLogger get the global logger to use in other packages
func GetDefaultLogger() *Logger {
	return gLog
}

// EnableAutoFlush function set autoflush flag
func EnableAutoFlush(flag bool) {
	gLog.EnableAutoFlush(flag)
}

// EnableStdOut function set or unset the output also in standard output
func EnableStdOut(flag bool) {
	gLog.EnableStdOut(flag)
}

// EnableMsec function set or unset milliseconds in log messages
func EnableMsec(flag bool) {
	gLog.EnableMsec(flag)
}

// EnableDay function set or unset day (YYYY-MM-DD) in log messages
func EnableDay(flag bool) {
	gLog.EnableDay(flag)
}

// EnableLevel function set or unset log level in log messages
func EnableLevel(flag bool) {
	gLog.EnableLevel(flag)
}

// Infof function emit a log at Info level with an interface as fmt.Printf
func Infof(format string, a ...interface{}) {
	gLog.Infof(format, a...)
}

// Debugf function emit a log at Debug level with an interface as fmt.Printf
func Debugf(format string, a ...interface{}) {
	gLog.Debugf(format, a...)
}

// Warningf function emit a log at Warning level with an interface as fmt.Printf
func Warningf(format string, a ...interface{}) {
	gLog.Warningf(format, a...)
}

// Errorf function emit a log at Error level with an interface as fmt.Printf
func Errorf(format string, a ...interface{}) {
	gLog.Errorf(format, a...)
}

// Fatalf function emit a log at Fatal level with an interface as fmt.Printf
func Fatalf(format string, a ...interface{}) {
	gLog.Fatalf(format, a...)
}

// Tracef function emit a log at Trace level with an interface as fmt.Printf
func Tracef(format string, a ...interface{}) {
	gLog.Tracef(format, a...)
}

// FlushLog function process all messages in queue
func FlushLog() {
	gLog.Flush()
}
