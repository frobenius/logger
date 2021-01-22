package logger

import (
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LogLevel is type to represent log level
type LogLevel int

// Log levels
const (
	Fatal LogLevel = iota
	Error
	Warning
	Info
	Debug
	Trace
)

type queueNode struct {
	ts       time.Time
	format   string
	logLevel LogLevel
	args     []interface{}
	next     *queueNode
}

// Logger object
type Logger struct {
	logfile        string
	maxFileSize    int
	maxNumFiles    int
	perm           os.FileMode
	mutex          sync.Mutex
	mutexMsg       sync.Mutex
	enableStdout   bool
	enableMsec     bool
	enableDay      bool
	enableLevel    bool
	level          LogLevel
	head           *queueNode
	tail           *queueNode
	autoFlush      bool
	colorsOnStdout bool
	colorsOnFile   bool
	compression    bool
}

// StringToLogLevel convert string to LogLevel
func StringToLogLevel(s string) LogLevel {
	level := strings.TrimSpace(s)
	level = strings.ToLower(level)
	if level == "debug" {
		return Debug
	} else if level == "trace" {
		return Trace
	}
	return Info
}

func (l *Logger) addMessage(ts time.Time, logLevel LogLevel, format string, a ...interface{}) {
	var m queueNode
	m.ts = ts
	m.format = format
	m.logLevel = logLevel
	m.args = a
	l.mutex.Lock()
	m.next = nil
	if l.head == nil {
		l.head = &m
	} else {
		l.tail.next = &m
	}
	l.tail = &m
	l.mutex.Unlock()
}

func compressLogFile(src string, dst string) {
	f, err := os.Open(src)
	if err != nil {
		return
	}
	data := make([]byte, 64*1024)
	defer f.Close()
	compressedFile, err := os.Create(dst)
	w := gzip.NewWriter(compressedFile)
	defer w.Close()
	if err != nil {
		return
	}
	n, err := f.Read(data)
	for n > 0 && err == nil {
		w.Write(data[:n])
		n, err = f.Read(data)
	}
}

// NewLogger function istance an object Logger
func NewLogger(filename string, maxFileSize int, maxNumFiles int, level LogLevel, perm os.FileMode) *Logger {
	logger := new(Logger)
	logger.logfile = filename
	logger.maxFileSize = maxFileSize
	if logger.maxFileSize < 0 {
		logger.maxFileSize = 0
	}
	logger.maxNumFiles = maxNumFiles
	if logger.maxNumFiles < 1 {
		logger.maxNumFiles = 1
	}
	logger.perm = perm
	logger.mutex = sync.Mutex{}
	logger.mutexMsg = sync.Mutex{}
	logger.enableStdout = false
	logger.level = level
	logger.head = nil
	logger.enableMsec = true
	logger.enableDay = true
	logger.enableLevel = true
	logger.autoFlush = false
	logger.colorsOnStdout = true
	logger.compression = false
	return logger
}

// logMessage function is the core of logger functionalities
func (l *Logger) logMessage() bool {
	var node *queueNode
	var ts time.Time
	var level string
	var format string
	var a []interface{}

	l.mutexMsg.Lock()
	defer l.mutexMsg.Unlock()
	l.mutex.Lock()
	if l.head == nil {
		l.mutex.Unlock()
		return false
	}
	node = l.head
	l.head = node.next
	if node == l.tail {
		l.tail = nil
	}
	l.mutex.Unlock()

	ts = node.ts
	if !l.enableLevel {
		level = ""
	} else if node.logLevel == Info {
		level = "[INFO ] "
	} else if node.logLevel == Debug {
		level = "[DEBUG] "
	} else if node.logLevel == Warning {
		level = "[WARN ] "
	} else if node.logLevel == Error {
		level = "[ERROR] "
	} else if node.logLevel == Fatal {
		level = "[FATAL] "
	} else if node.logLevel == Trace {
		level = "[TRACE] "
	}
	color, _ := colors[node.logLevel]
	levelColor := fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, level)
	format = node.format
	a = node.args
	var msg string

	if !l.enableDay {
		if !l.enableMsec {
			msg = fmt.Sprintf("%02d:%02d:%02d.%03d",
				ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond()/1000000)
		} else {
			msg = fmt.Sprintf("%02d:%02d:%02d",
				ts.Hour(), ts.Minute(), ts.Second())
		}
	} else {
		if !l.enableMsec {
			msg = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				ts.Year(), ts.Month(), ts.Day(),
				ts.Hour(), ts.Minute(), ts.Second())
		} else {
			msg = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%03d",
				ts.Year(), ts.Month(), ts.Day(),
				ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond()/1000000)
		}
	}

	// Build the log row to write on file
	extraMsg := fmt.Sprintf(format, a...) + "\n"
	msgColor := msg + " " + levelColor + extraMsg
	msg += " " + level + extraMsg

	// If stdout enable print message also on standard output
	if l.enableStdout {
		if l.colorsOnStdout {
			fmt.Print(msgColor)
		} else {
			fmt.Print(msg)
		}
	}

	if l.logfile == "" {
		return true
	}

	// Perform eventually rotation of log files
	fileinfo, err := os.Stat(l.logfile)
	if err == nil {
		maxSize := int64(l.maxFileSize)
		msgLen := len(msg)
		if l.colorsOnFile {
			msgLen = len(msgColor)
		}
		if l.maxFileSize > 0 && fileinfo.Size()+int64(msgLen) > maxSize {
			// Remove older log file
			if l.maxNumFiles > 1 {
				os.Remove(l.logfile + "." + strconv.Itoa(l.maxNumFiles-1))
				os.Remove(l.logfile + "." + strconv.Itoa(l.maxNumFiles-1) + ".gz")
			}

			// Rotate files
			for i := l.maxNumFiles - 2; i >= 1; i-- {
				if l.compression {
					os.Rename(l.logfile+"."+strconv.Itoa(i)+".gz", l.logfile+"."+strconv.Itoa(i+1)+".gz")
				} else {
					os.Rename(l.logfile+"."+strconv.Itoa(i), l.logfile+"."+strconv.Itoa(i+1))
				}
			}
			if l.maxNumFiles > 1 {
				if l.compression {
					compressLogFile(l.logfile, l.logfile+".1.gz")
					os.Remove(l.logfile)
				} else {
					os.Rename(l.logfile, l.logfile+".1")
				}
			} else {
				os.Remove(l.logfile)
			}
		}
	}

	// Open log file in append
	f, err := os.OpenFile(l.logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, l.perm)
	if err != nil {
		return true
	}

	// Close file when function return
	defer f.Close()

	// Write on file
	if l.colorsOnFile {
		f.Write([]byte(msgColor))
	} else {
		f.Write([]byte(msg))
	}
	return true
}

// SetFileName set the log file name
func (l *Logger) SetFileName(fileName string) {
	l.logfile = fileName
}

// SetMaxFileSize set the maximum size of a single log file in bytes. 0 means no limit
func (l *Logger) SetMaxFileSize(maxSize int) {
	if maxSize < 0 {
		l.maxFileSize = 0
	} else {
		l.maxFileSize = maxSize
	}
}

// SetMaxNumFiles set the maximum number of log files
func (l *Logger) SetMaxNumFiles(maxNum int) {
	if maxNum < 1 {
		l.maxNumFiles = 1
	} else {
		l.maxNumFiles = maxNum
	}
}

// SetPermissionsFile set the permission on log file
func (l *Logger) SetPermissionsFile(perm os.FileMode) {
	l.perm = perm
}

// SetLevel function set a new log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// EnableAutoFlush function set autoflush flag
func (l *Logger) EnableAutoFlush(flag bool) {
	l.autoFlush = flag
}

// EnableStdOut function set or unset the output also in standard output
func (l *Logger) EnableStdOut(flag bool) {
	l.enableStdout = flag
}

// EnableMsec function set or unset milliseconds in log messages
func (l *Logger) EnableMsec(flag bool) {
	l.enableMsec = flag
}

// EnableDay function set or unset day (YYYY-MM-DD) in log messages
func (l *Logger) EnableDay(flag bool) {
	l.enableDay = flag
}

// EnableLevel function set or unset log level in log messages
func (l *Logger) EnableLevel(flag bool) {
	l.enableLevel = flag
}

// EnableCompression enable compression of rolled file
func (l *Logger) EnableCompression(flag bool) {
	l.compression = flag
}

// EnableColorsOnFile function enable or disable colors in log written on log file
func (l *Logger) EnableColorsOnFile(flag bool) {
	l.colorsOnFile = flag
}

// EnableColorsOnStdout function enable or disable colors in log written on stdout
func (l *Logger) EnableColorsOnStdout(flag bool) {
	l.colorsOnStdout = flag
}

// Flush function process all messages in queue
func (l *Logger) Flush() {
	for l.logMessage() {
	}
}

func (l *Logger) processMessage() {
	if l.autoFlush {
		l.logMessage()
	} else {
		go l.logMessage()
	}
}

// Infof function emit a log at Info level with an interface as fmt.Printf
func (l *Logger) Infof(format string, a ...interface{}) {
	if Info <= l.level {
		l.addMessage(time.Now(), Info, format, a...)
		l.processMessage()
	}
}

// Debugf function emit a log at Debug level with an interface as fmt.Printf
func (l *Logger) Debugf(format string, a ...interface{}) {
	if Debug <= l.level {
		l.addMessage(time.Now(), Debug, format, a...)
		l.processMessage()
	}
}

// Warningf function emit a log at Warning level with an interface as fmt.Printf
func (l *Logger) Warningf(format string, a ...interface{}) {
	if Warning <= l.level {
		l.addMessage(time.Now(), Warning, format, a...)
		l.processMessage()
	}
}

// Errorf function emit a log at Error level with an interface as fmt.Printf
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.addMessage(time.Now(), Error, format, a...)
	l.processMessage()
}

// Fatalf function emit a log at Fatal level with an interface as fmt.Printf
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.addMessage(time.Now(), Fatal, format, a...)
	l.processMessage()
}

// Tracef function emit a log at Trace level with an interface as fmt.Printf
func (l *Logger) Tracef(format string, a ...interface{}) {
	if Trace <= l.level {
		l.addMessage(time.Now(), Trace, format, a...)
		l.processMessage()
	}
}
