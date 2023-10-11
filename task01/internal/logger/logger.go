package logger

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
)

func Init() error {
	f, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	InfoLog = log.New(f, "[INFO]\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(f, "[ERROR]\t", log.Ldate|log.Ltime)

	return nil
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		InfoLog.Printf("Request :%s %s", c.Request.Method, c.Request.URL)
		c.Next()
	}
}
