package log

import (
	"fmt"
	"time"
	"turutupa/gsnake/fsutil"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[91m"
	colorYellow = "\033[93m"
	colorCyan   = "\033[96m"
	colorWhite  = "\033[97m"
)

const LOG_NAME string = "gsnake."

var enableStorage = false

func EnableStorage() {
	Info("Log persistence enabled")
	enableStorage = true
}

func Info(message string) {
	printLog("INFO", colorCyan, message)
}

func Warn(message string) {
	printLog("WARN", colorYellow, message)
}

func Error(message string, error error) {
	if error != nil {
		printLog("ERROR", colorRed, message+colorRed+" ERROR_MSG "+colorReset+error.Error())
	} else {
		printLog("ERROR", colorRed, message)
	}
}

func Log(message string) {
	printLog("LOG", colorWhite, message)
}

func printLog(logType, color, message string) {
	timestamp := time.Now().Format("[2006-01-02][15:04:05]")
	logTypeBold := fmt.Sprintf("\033[1m%s\033[0m", logType)
	logLine := fmt.Sprintf("%s %s%s%s %s\n", timestamp, color, logTypeBold, colorReset, message)
	fmt.Printf(logLine)

	if enableStorage {
		fileName := getLogFileName()
		logRaw := fmt.Sprintf("%s %s %s\n", timestamp, logType, message)
		writeLogToFile(fileName, logRaw)
	}
}

func getLogFileName() string {
	currentTime := time.Now()
	dateString := currentTime.Format("20060102")
	fileName := LOG_NAME + dateString + ".log"
	return fileName
}

func writeLogToFile(filename, logLine string) {
	fsutil.AppendToFile(filename, logLine)
}
