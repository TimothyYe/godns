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
	L_INFO int = iota
	L_WARNING
	L_DEBUG
	PRE_INFO    = "[   INFO]"
	PRE_WARNING = "[WARNING]"
	PRE_DEBUG   = "[  DEBUG]"
)

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

func InitLogger(logfile string, size, num int) (err error) {
	logger, err := NewLogger(logfile, size, num, L_INFO, -1, -1)
	if logger != nil {
		logger.level = L_INFO - 1
	}
	return
}

//immplement write
func (this *Logger) Write(p []byte) (n int, err error) {
	if this.DEV_MODE {
		n, err = os.Stdout.Write(p)
		return
	}
	n, err = this.fd.Write(p)
	if err == nil {
		fi, e := this.fd.Stat()
		if e != nil {
			err = e
			return
		}
		if fi.Size() > int64(this.size) {
			this.muSplit.Lock()
			defer this.muSplit.Unlock()

			fname := fi.Name()
			strings.HasSuffix(fname, ".log")
			fbase := fname[:len(fname)-3]

			oldBs := make([]byte, 0, this.size)
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
				if len(oldBs)+len(line) > this.size {
					newBs = append(newBs, line...)
				} else {
					oldBs = append(oldBs, line...)
				}
			}
			fd.Close()

			_, err = this.saveLog(1, fbase, oldBs)
			if err != nil {
				return
			}
			err = this.fd.Close()
			if err != nil {
				return
			}
			err = os.Remove(fname)
			if err != nil {
				return
			}
			this.fd, err = os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModeAppend|0666)
			if err != nil {
				return
			}
			_, err = this.fd.Write(newBs)
			if err != nil {
				return
			}
		}
	}
	return
}

func (this *Logger) saveLog(index int, fbase string, data []byte) (n int, err error) {
	fn := fbase + strconv.Itoa(index) + ".log"
	_, err = os.Stat(fn)
	if index < this.num && err == nil {
		var b []byte
		b, err = ioutil.ReadFile(fn)
		if err != nil {
			return
		}
		n, err = this.saveLog(index+1, fbase, b)
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

//flush buf data to std log
func (this *Logger) Flush() {
	if this.buf.Len() > 0 {
		this.mu.Lock()
		defer this.mu.Unlock()

		log.SetFlags(0)
		log.Print(this.buf)
		log.SetFlags(log.LstdFlags)
		this.buf.Reset()
	}
}

//clean prefix and check buf size
func (this *Logger) clean() {
	this.log.SetPrefix("")
	if this.buf.Len()/1024 > this.flushSize {
		go this.Flush()
	}
}

func (this *Logger) setPrefix(lv int) bool {
	if lv > this.level {
		return false
	}

	switch lv {
	case L_INFO:
		this.log.SetPrefix(PRE_INFO)
	case L_WARNING:
		this.log.SetPrefix(PRE_WARNING)
	case L_DEBUG:
		this.log.SetPrefix(PRE_DEBUG)
	default:
		return false
	}
	return true
}

func (this *Logger) logPrint(lv int, args ...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if !this.setPrefix(lv) {
		return
	}
	this.log.Print(args...)
	this.clean()
}

func (this *Logger) logPrintln(lv int, args ...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if !this.setPrefix(lv) {
		return
	}
	this.log.Println(args...)
	this.clean()
}

func (this *Logger) logPrintf(lv int, format string, args ...interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if !this.setPrefix(lv) {
		return
	}
	this.log.Printf(format, args...)
	this.clean()
}

//close fd
func (this *Logger) Close() {
	if this.fd != nil {
		this.Flush()
		this.fd.Close()
	}
}

func (this *Logger) Info(args ...interface{}) {
	this.logPrint(L_INFO, args...)
}

func (this *Logger) Infoln(args ...interface{}) {
	this.logPrintln(L_INFO, args...)
}

func (this *Logger) Infof(format string, args ...interface{}) {
	this.logPrintf(L_INFO, format, args...)
}

func (this *Logger) Warning(args ...interface{}) {
	this.logPrint(L_WARNING, args...)
}

func (this *Logger) Warningln(args ...interface{}) {
	this.logPrintln(L_WARNING, args...)
}

func (this *Logger) Warningf(format string, args ...interface{}) {
	this.logPrintf(L_WARNING, format, args...)
}

func (this *Logger) Debug(args ...interface{}) {
	this.logPrint(L_DEBUG, args...)
	this.Flush()
}

func (this *Logger) Debugln(args ...interface{}) {
	this.logPrintln(L_DEBUG, args...)
	this.Flush()
}

func (this *Logger) Debugf(format string, args ...interface{}) {
	this.logPrintf(L_DEBUG, format, args...)
	this.Flush()
}
