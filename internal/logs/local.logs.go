package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	mu               sync.RWMutex
	currentInfoFile  *os.File
	currentErrorFile *os.File
	currentDebugFile *os.File

	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

type LogInfo struct {
	Timestamp string `json:"timestamp"`
	Action    string
	User      string
	Location  string
}

type LogError struct {
	Timestamp string `json:"timestamp"`
	Location  string
	Error     string
}

type LogDebug struct {
	Timestamp string `json:"timestamp"`
	Location  string
	Msg       string
}

func logDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("internal", "logs")
	}
	return filepath.Dir(file)
}

func openLogFiles(suffix string) (infoFile, errorFile, debugFile *os.File, err error) {
	dir := logDir()
	if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
		return nil, nil, nil, fmt.Errorf("no se pudo crear directorio de logs: %w", mkErr)
	}

	infoPath := filepath.Join(dir, "info_"+suffix+".log")
	errorPath := filepath.Join(dir, "error_"+suffix+".log")
	debugPath := filepath.Join(dir, "debug_"+suffix+".log")

	infoFile, err = os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, nil, nil, err
	}
	errorFile, err = os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		_ = infoFile.Close()
		return nil, nil, nil, err
	}
	debugFile, err = os.OpenFile(debugPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		_ = infoFile.Close()
		_ = errorFile.Close()
		return nil, nil, nil, err
	}

	return infoFile, errorFile, debugFile, nil
}

func setLoggers(infoFile, errorFile, debugFile *os.File) {
	InfoLogger = log.New(infoFile, "", 0)
	ErrorLogger = log.New(errorFile, "", 0)
	DebugLogger = log.New(debugFile, "", 0)
}

func InitFilesLogs() error {
	mu.Lock()
	defer mu.Unlock()

	if currentInfoFile != nil {
		_ = currentInfoFile.Close()
	}
	if currentErrorFile != nil {
		_ = currentErrorFile.Close()
	}
	if currentDebugFile != nil {
		_ = currentDebugFile.Close()
	}

	suffix := time.Now().Format("2006-01-02")
	infoFile, errorFile, debugFile, err := openLogFiles(suffix)
	if err != nil {
		return err
	}

	currentInfoFile = infoFile
	currentErrorFile = errorFile
	currentDebugFile = debugFile
	setLoggers(currentInfoFile, currentErrorFile, currentDebugFile)
	Debug("local.logs.go", "Archivos log creados correctamente")
	return nil
}

func ResetLogs() error {
	mu.Lock()
	defer mu.Unlock()

	if currentInfoFile != nil {
		_ = currentInfoFile.Close()
	}
	if currentErrorFile != nil {
		_ = currentErrorFile.Close()
	}
	if currentDebugFile != nil {
		_ = currentDebugFile.Close()
	}

	suffix := time.Now().Format("2006-01-02")
	infoFile, errorFile, debugFile, err := openLogFiles(suffix)
	if err != nil {
		return err
	}

	currentInfoFile = infoFile
	currentErrorFile = errorFile
	currentDebugFile = debugFile
	setLoggers(currentInfoFile, currentErrorFile, currentDebugFile)
	return nil
}

func Info(action, user, location string) {
	entry := LogInfo{
		Timestamp: time.Now().Format(time.RFC3339),
		Action:    action,
		User:      user,
		Location:  location,
	}
	data, _ := json.Marshal(entry)
	if InfoLogger == nil {
		log.Println(string(data))
		return
	}
	InfoLogger.Println(string(data))
}

func Error(location, errMsg string) {
	entry := LogError{
		Timestamp: time.Now().Format(time.RFC3339),
		Location:  location,
		Error:     errMsg,
	}
	data, _ := json.Marshal(entry)
	if ErrorLogger == nil {
		log.Println(string(data))
		return
	}
	ErrorLogger.Println(string(data))
}

func Debug(location, msg string) {
	entry := LogDebug{
		Timestamp: time.Now().Format(time.RFC3339),
		Location:  location,
		Msg:       msg,
	}
	data, _ := json.Marshal(entry)
	if DebugLogger == nil {
		log.Println(string(data))
		return
	}
	DebugLogger.Println(string(data))
}
