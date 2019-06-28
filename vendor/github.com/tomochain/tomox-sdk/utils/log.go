package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/op/go-logging"
)

// TODO: add log prefix to params
const (
	LogPrefix = "tomochain"
)

var Logger = NewLogger("main", "./logs/main.log")

var StdoutLogger = NewStandardOutputLogger()
var TerminalLogger = NewColoredLogger()

// NewFileLogger creates a logging utility that outputs to the file passed as argument and also output to stdout.
func NewLogger(module string, logFile string) *logging.Logger {
	_, fileName, _, _ := runtime.Caller(1)
	logDir := path.Join(path.Dir(fileName), "../logs/")
	logFile = path.Join(path.Dir(fileName), "../", logFile)

	logger, err := logging.GetLogger(module)
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{message}`,
	)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	log, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, log)
	backend := logging.NewLogBackend(writer, LogPrefix, 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}

// NewFileLogger creates a logging utility that only output to stdout.
func NewStandardOutputLogger() *logging.Logger {
	_, fileName, _, _ := runtime.Caller(1)
	logDir := path.Join(path.Dir(fileName), "../logs/")

	logger, err := logging.GetLogger("main")
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{message}`,
	)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	writer := io.MultiWriter(os.Stdout)
	backend := logging.NewLogBackend(writer, LogPrefix, 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}

// NewFileLogger creates a logging utility that outputs to the file passed as argument but
// but does not output to stdout.
func NewFileLogger(module string, logFile string) *logging.Logger {
	_, fileName, _, _ := runtime.Caller(1)
	logDir := path.Join(path.Dir(fileName), "../logs/")
	logFile = path.Join(path.Dir(fileName), "../", logFile)

	logger, err := logging.GetLogger(module)
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{message}`,
	)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	log, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(log)
	backend := logging.NewLogBackend(writer, LogPrefix, 0)
	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}

func NewErrorLogger() *logging.Logger {
	return NewLogger("error", "./logs/errors.log")
}

func NewColoredLogger() *logging.Logger {
	logger, err := logging.GetLogger("colored")
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{color}%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{color:reset} %{message}`,
	)

	writer := io.MultiWriter(os.Stdout)
	backend := logging.NewLogBackend(writer, LogPrefix, 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}

func NewNoopLogger() *logging.Logger {
	logger, err := logging.GetLogger("noop")
	if err != nil {
		panic(err)
	}
	noopBackend := logging.NewLogBackend(ioutil.Discard, "", 0)
	formattedBackend := logging.NewBackendFormatter(noopBackend, logging.DefaultFormatter)
	leveledBackend := logging.AddModuleLevel(formattedBackend)
	logger.SetBackend(leveledBackend)
	return logger
}
