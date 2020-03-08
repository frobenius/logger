package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	var log = NewLogger("test_logger.log", 5000, 5, Trace, 0666)
	SetLevel(Trace)
	Infof("%06d - This is a global log message", 1)
	Debugf("%06d - This is a global log message", 2)
	Warningf("%06d - This is a global log message", 3)
	Tracef("%06d - This is a global log message", 4)
	Errorf("%06d - This is a global log message", 5)
	Fatalf("%06d - This is a global log message", 6)
	for i := 1; i <= 30; i++ {
		log.Infof("%06d - InfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfoInfo", i)
		log.Debugf("%06d - DebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebugDebug", i)
		log.Tracef("%06d - TraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTraceTrace", i)
		log.Errorf("%06d - ErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorErrorError", i)
		log.Fatalf("%06d - FatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatalFatal", i)
		log.Warningf("%06d - WarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarnWarn", i)
		if i == 5 {
			log.Flush()
			log.EnableDay(false)
		}
		if i == 10 {
			log.Flush()
			log.EnableMsec(false)
		}
		if i == 15 {
			log.Flush()
			log.EnableDay(true)
		}
		if i == 20 {
			log.Flush()
			log.SetLevel(Info)
		}
		if i == 25 {
			log.Flush()
			log.EnableLevel(false)
		}
	}
	log.Flush()
}
