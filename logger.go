package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// L_INFO log level
	L_INFO int = iota
	// L_WARNING log level
	L_WARNING
	// L_DEBUG log level
	L_DEBUG
	// PRE_INFO log level
	PRE_INFO = "[   INFO]"
	// PRE_WARNING log level
	PRE_WARNING = "[WARNING]"
	// PRE_DEBUG log level
	PRE_DEBUG = "[  DEBUG]"
)

// Logger struct
type Logger struct {
	DEV_MODE      bool
	fd            *os.File
	size          int
	num           int
	level         int
	mu            sync.Mutex
	muSplit       sync.Mutex
	flushInterval int64 //Second
	flushSize     int
	buf           *bytes.Buffer
	log           *log.Logger
}

// NewLogger returns a new created logger
func NewLogger(logfile string, size, num int, level int, flushInterval int64, flushSize int) (logger *Logger, err error) {
	if size < 1 || num < 1 || level < L_INFO || len(logfile) < 1 {
		err = errors.New("NewLogWriter:param error.")
		return
	}
	logger = &Logger{size: size * 1024, num: num, level: level, DEV_MODE: false}
	logger.fd, err = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0666)
	if err != nil {
		logger = nil
		return
	}
	log.SetOutput(logger)
	if flushInterval > 0 && flushSize > 0 {
		logger.buf = new(bytes.Buffer)
		logger.log = log.New(logger.buf, "", log.LstdFlags)

		go func(interval int64, logger *Logger) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("logger Tick, Recovered in %v:\n %s", r, debug.Stack())
				}
			}()
			c := time.Tick(time.Duration(interval) * time.Second)
			for _ = range c {
				logger.Flush()
			}
		}(flushInterval, logger)
	}
	return
}

// InitLogger initialize logger with specified log filename & size
func InitLogger(logfile string, size, num int) (err error) {
	logger, err := NewLogger(logfile, size, num, L_INFO, -1, -1)
	if logger != nil {
		logger.level = L_INFO - 1
	}
	return
}

// immplement write
func (logger *Logger) Write(p []byte) (n int, err error) {
	if logger.DEV_MODE {
		n, err = os.Stdout.Write(p)
		return
	}
	n, err = logger.fd.Write(p)
	if err == nil {
		fi, e := logger.fd.Stat()
		if e != nil {
			err = e
			return
		}
		if fi.Size() > int64(logger.size) {
			logger.muSplit.Lock()
			defer logger.muSplit.Unlock()

			fname := fi.Name()
			strings.HasSuffix(fname, ".log")
			fbase := fname[:len(fname)-3]

			oldBs := make([]byte, 0, logger.size)
			newBs := []byte{}
			fd, e := os.Open(fname)
			if e != nil {
				err = e
				return
			}
			rd := bufio.NewReader(fd)
			for {
				line, e := rd.ReadBytes('\n')
				if e == io.EOF {
					break
				}
				if e != nil {
					err = e
					return
				}
				if len(oldBs)+len(line) > logger.size {
					newBs = append(newBs, line...)
				} else {
					oldBs = append(oldBs, line...)
				}
			}
			fd.Close()

			_, err = logger.saveLog(1, fbase, oldBs)
			if err != nil {
				return
			}
			err = logger.fd.Close()
			if err != nil {
				return
			}
			err = os.Remove(fname)
			if err != nil {
				return
			}
			logger.fd, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0666)
			if err != nil {
				return
			}
			_, err = logger.fd.Write(newBs)
			if err != nil {
				return
			}
		}
	}
	return
}

func (logger *Logger) saveLog(index int, fbase string, data []byte) (n int, err error) {
	fn := fbase + strconv.Itoa(index) + ".log"
	_, err = os.Stat(fn)
	if index < logger.num && err == nil {
		var b []byte
		b, err = ioutil.ReadFile(fn)
		if err != nil {
			return
		}
		n, err = logger.saveLog(index+1, fbase, b)
		if err != nil {
			return
		}
	}

	fd, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm|0666)
	if err != nil {
		return
	}
	defer fd.Close()
	n, err = fd.Write(data)
	return
}

// Flush buf data to std log
func (logger *Logger) Flush() {
	if logger.buf.Len() > 0 {
		logger.mu.Lock()
		defer logger.mu.Unlock()

		log.SetFlags(0)
		log.Print(logger.buf)
		log.SetFlags(log.LstdFlags)
		logger.buf.Reset()
	}
}

// Clean prefix and check buf size
func (logger *Logger) clean() {
	logger.log.SetPrefix("")
	if logger.buf.Len()/1024 > logger.flushSize {
		go logger.Flush()
	}
}

func (logger *Logger) setPrefix(lv int) bool {
	if lv > logger.level {
		return false
	}

	switch lv {
	case L_INFO:
		logger.log.SetPrefix(PRE_INFO)
	case L_WARNING:
		logger.log.SetPrefix(PRE_WARNING)
	case L_DEBUG:
		logger.log.SetPrefix(PRE_DEBUG)
	default:
		return false
	}
	return true
}

func (logger *Logger) logPrint(lv int, args ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if !logger.setPrefix(lv) {
		return
	}
	logger.log.Print(args...)
	logger.clean()
}

func (logger *Logger) logPrintln(lv int, args ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if !logger.setPrefix(lv) {
		return
	}
	logger.log.Println(args...)
	logger.clean()
}

func (logger *Logger) logPrintf(lv int, format string, args ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if !logger.setPrefix(lv) {
		return
	}
	logger.log.Printf(format, args...)
	logger.clean()
}

// Close fd
func (logger *Logger) Close() {
	if logger.fd != nil {
		logger.Flush()
		logger.fd.Close()
	}
}

// Info output info log
func (logger *Logger) Info(args ...interface{}) {
	logger.logPrint(L_INFO, args...)
}

// Infoln output info log with newline
func (logger *Logger) Infoln(args ...interface{}) {
	logger.logPrintln(L_INFO, args...)
}

// Infof output formatted info log
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logPrintf(L_INFO, format, args...)
}

// Warning output warning log
func (logger *Logger) Warning(args ...interface{}) {
	logger.logPrint(L_WARNING, args...)
}

//Warningln output warning log with newline
func (logger *Logger) Warningln(args ...interface{}) {
	logger.logPrintln(L_WARNING, args...)
}

// Warningf output formatted warning log
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.logPrintf(L_WARNING, format, args...)
}

// Debug output debug log
func (logger *Logger) Debug(args ...interface{}) {
	logger.logPrint(L_DEBUG, args...)
	logger.Flush()
}

// Debugln output debug log with newline
func (logger *Logger) Debugln(args ...interface{}) {
	logger.logPrintln(L_DEBUG, args...)
	logger.Flush()
}

// Debugf output formatted debug log
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logPrintf(L_DEBUG, format, args...)
	logger.Flush()
}
