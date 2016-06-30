package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
)

func identifyPanic() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	switch {
	case name != "":
		return fmt.Sprintf("%v:%v", name, line)
	case file != "":
		return fmt.Sprintf("%v:%v", file, line)
	}

	return fmt.Sprintf("pc:%x", pc)
}

func usage() {
	log.Println("[command] -c=[config file path]")
	flag.PrintDefaults()
}

func checkSettings(config *Settings) error {
	if (config.Email == "" || config.Password == "") && config.LoginToken == "" {
		return errors.New("Input email/password or login token cannot be empty!")
	}

	return nil
}
