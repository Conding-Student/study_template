package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	Separator     *log.Logger
)

func SystemLoggerAPI(url string, body interface{}, class string, resp *http.Response, ret interface{}, ip string) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_API.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strBody := fmt.Sprintf("%#v", body)
	strRet := fmt.Sprintf("%#v", ret)
	strStatus := fmt.Sprintf("%v", resp.Status)

	Separator.Println("")
	InfoLogger.Println(class + ": - - - - - - - - - - - - - - -")
	InfoLogger.Println(class + ": FROM: " + ip)
	InfoLogger.Println(class + ": URL REQUEST: " + url)
	InfoLogger.Println(class + ": BODY: " + strBody)
	InfoLogger.Println(class + ": RESPONSE: " + strRet)
	InfoLogger.Println(class + ": STATUS: " + strStatus)
	file.Close()
}

func SystemLoggerErrorAPI(url string, body interface{}, class string, resp *http.Response, ret interface{}, ip string) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_API.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strBody := fmt.Sprintf("%#v", body)
	strRet := fmt.Sprintf("%#v", ret)
	strStatus := fmt.Sprintf("%v", resp.Status)

	Separator.Println("")
	InfoLogger.Println(class + ": - - - - - - - - - - - - - - -")
	InfoLogger.Println(class + ": FROM: " + ip)
	InfoLogger.Println(class + ": URL REQUEST: " + url)
	InfoLogger.Println(class + ": BODY: " + strBody)
	InfoLogger.Println(class + ": RESPONSE: " + strRet)
	InfoLogger.Println(class + ": STATUS: " + strStatus)
	file.Close()
}

func SystemLoggerDB(body interface{}, class string, status int, ret string, ip string) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_DB.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strBody := fmt.Sprintf("%#v", body)
	strRet := fmt.Sprintf("%#v", ret)

	Separator.Println("")
	InfoLogger.Println(class + ": - - - - - - - - - - - - - - -")
	InfoLogger.Println(class + ": FROM: " + ip)
	InfoLogger.Println(class + ": BODY: " + strBody)
	InfoLogger.Println(class + ": RESPONSE: " + strRet)
	InfoLogger.Println(class + ": STATUS: " + strconv.Itoa(status))
	file.Close()
}

func SystemLoggerErrorDB(body interface{}, class string, status int, ret string, ip string) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_DB.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strBody := fmt.Sprintf("%#v", body)
	strRet := fmt.Sprintf("%#v", ret)

	Separator.Println("")
	ErrorLogger.Println(class + ": - - - - - - - - - - - - - - -")
	ErrorLogger.Println(class + ": FROM: " + ip)
	ErrorLogger.Println(class + ": BODY: " + strBody)
	ErrorLogger.Println(class + ": RESPONSE: " + strRet)
	ErrorLogger.Println(class + ": STATUS: " + strconv.Itoa(status))
	file.Close()
}

func SystemLoggerError(class string, proccess string, errorData error) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_SYSTEM_ERROR.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strError := fmt.Sprintf("%#v", errorData.Error())

	Separator.Println("")
	ErrorLogger.Println(class + ": - - - - - - - - - - - - - - -")
	ErrorLogger.Println(class + ": ERROR AT: " + proccess)
	ErrorLogger.Println(class + ": ERROR: " + strError)
	file.Close()
}

func SystemLogger(class string, data string, username string) {

	currentTime := time.Now()
	file, err := os.OpenFile("logs/"+currentTime.Format("01022006")+"_SYSTEM_ERROR.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	Separator = log.New(file, "", log.Ldate|log.Ltime)

	strError := fmt.Sprintf("%#v", username)

	Separator.Println("")
	InfoLogger.Println(class + ": - - - - - - - - - - - - - - -")
	InfoLogger.Println(class + ": UPDATED RECORD is_active column: value(" + data + ")")
	InfoLogger.Println(class + ": ERROR: " + strError)
	file.Close()
}
