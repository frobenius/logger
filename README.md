# logger ![Go](https://github.com/frobenius/logger/workflows/Go/badge.svg?branch=master) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![Build Status](https://travis-ci.com/frobenius/logger.svg?branch=master)](https://travis-ci.com/frobenius/logger)
Logging library for Go

## Features

- Very simple and intuitive api
- Log messages are written to file throught go routine in order to avoid delays
- Global logger ready to use to log on standard output
- Rolling files by size and/or date
- Log level colored

## Code samples

- Use global logger
  ````go
  Infof("This is a global log message at INFO level")
  Warningf("This is another global log message at WARN level")
  ````
  
- Log on file
  ````go
  var log = NewLogger("file.log", 1024*1024, 5, Trace, 0666)
  log.Infof("This log message will be written on file.log")
  log.Flush()
  ````
  
