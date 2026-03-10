package logs

import (
	"encoding/json"
	"log"
	"os"
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

func InitFilesLogs() error {
	mu.Lock()
	defer mu.Unlock()
	if currentInfoFile != nil {
		currentInfoFile.Close()
	}
	if currentErrorFile != nil {
		currentErrorFile.Close()
	}
	if currentDebugFile != nil {
		currentDebugFile.Close()
	}
	suffix := time.Now().Format("2006-01-02")
	currentInfoFile, err := os.OpenFile("../internal/logs/info_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	currentErrorFile, err = os.OpenFile("../internal/logs/error_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	currentDebugFile, err = os.OpenFile("../internal/logs/debug_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	InfoLogger = log.New(currentInfoFile, "", 0)
	ErrorLogger = log.New(currentErrorFile, "", 0)
	DebugLogger = log.New(currentDebugFile, "", 0)
	if err == nil {
		Debug("local.logs.go", "Archivos log creados correctamente")
	} else {
		println("Erro al escribir")
		panic("No se pudieron crear los archivos de log: " + err.Error())

	}

	return err
}
func ResetLogs() error {
	mu.Lock()
	defer mu.Unlock()

	if currentInfoFile != nil {
		currentInfoFile.Close()
	}
	if currentErrorFile != nil {
		currentErrorFile.Close()
	}
	if currentDebugFile != nil {
		currentDebugFile.Close()
	}

	suffix := time.Now().Format("2006-01-02")

	var err error
	currentInfoFile, err = os.OpenFile("../internal/logs/info_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	currentErrorFile, err = os.OpenFile("../internal/logs/error_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	currentDebugFile, err = os.OpenFile("../internal/logs/debug_"+suffix+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		panic("No se pudieron crear los archivos de log: " + err.Error())
	}

	InfoLogger = log.New(currentInfoFile, "", 0)
	ErrorLogger = log.New(currentErrorFile, "", 0)
	DebugLogger = log.New(currentDebugFile, "", 0)
	return err
}
func Info(action, user, location string) {
	entry := LogInfo{
		Timestamp: time.Now().Format(time.RFC3339),
		Action:    action,
		User:      user,
		Location:  location,
	}
	data, _ := json.Marshal(entry)
	InfoLogger.Println(string(data))
}

func Error(location, errMsg string) {
	entry := LogError{
		Timestamp: time.Now().Format(time.RFC3339),
		Location:  location,
		Error:     errMsg,
	}
	data, _ := json.Marshal(entry)
	ErrorLogger.Println(string(data))
}

func Debug(location, msg string) {
	entry := LogDebug{
		Timestamp: time.Now().Format(time.RFC3339),
		Location:  location,
		Msg:       msg,
	}
	data, _ := json.Marshal(entry)
	DebugLogger.Println(string(data))
}
