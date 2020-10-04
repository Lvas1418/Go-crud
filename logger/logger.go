package logger

import (
  "fmt"
  "log"
  "os"
)

var Log *log.Logger

func init() {
  // set location of log file
  var logpath = "./logger"

  logfile, err := os.OpenFile(logpath+"/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if err != nil {
    fmt.Println(`Package logger. Error opening file: info.log. Error: `, err)
    os.Exit(1)
  }
  Log = log.New(logfile, "", log.LstdFlags|log.Lshortfile)
  //Log.Println("LogFile : " + logpath)
}
