package logger

import "os"

// gLog is the static logger
var gLog *Logger

var colors map[LogLevel]int

func init() {
	colors = make(map[LogLevel]int)
	colors[Info] = 32
	colors[Debug] = 34
	colors[Warning] = 33
	colors[Error] = 31
	colors[Trace] = 36
	colors[Fatal] = 35
	gLog = NewLogger("", 0, 0, Info, 0666)
	gLog.EnableStdOut(true)
	gLog.EnableAutoFlush(true)
}

// SetLevel function set a new log level
func SetLevel(level LogLevel) {
	gLog.level = level
}

// SetFileName set the log file name
func SetFileName(fileName string) {
	gLog.SetFileName(fileName)
}

// SetMaxFileSize set the maximum size of a single log file in bytes. 0 means no limit
func SetMaxFileSize(maxSize int) {
	gLog.SetMaxFileSize(maxSize)
}

// SetMaxNumFiles set the maximum number of log files
func SetMaxNumFiles(maxNum int) {
	gLog.SetMaxNumFiles(maxNum)
}

// SetPermissionsFile set the permission on log file
func SetPermissionsFile(perm os.FileMode) {
	gLog.SetPermissionsFile(perm)
}

// GetDefaultLogger get the global logger to use in other packages
func GetDefaultLogger() *Logger {
	return gLog
}

// SetDefaultLogger set the defalut logger
func SetDefaultLogger(logger *Logger) {
	gLog = logger
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

// EnableColorsOnFile function enable or disable colors in log written on log file
func EnableColorsOnFile(flag bool) {
	gLog.EnableColorsOnFile(flag)
}

// EnableColorsOnStdout function enable or disable colors in log written on stdout
func EnableColorsOnStdout(flag bool) {
	gLog.EnableColorsOnStdout(flag)
}

// EnableCompression enable compression of rolled file
func EnableCompression(flag bool) {
	gLog.EnableCompression(flag)
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
