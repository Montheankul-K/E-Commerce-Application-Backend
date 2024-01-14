package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Montheankul-K/E-Commerce-Application-Backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ILogger interface {
	Print() ILogger
	Save()
	SetQuery(c *fiber.Ctx)
	SetBody(c *fiber.Ctx)
	SetResponse(res any)
}

type logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitLogger(c *fiber.Ctx, res any) ILogger {
	log := &logger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(), // จะแสดงหรือไม่แสดงขึ้นอยู่กับ reverse proxy
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *logger) Print() ILogger {
	utils.Debug(l)
	return l
}

func (l *logger) Save() {
	data := utils.Output(l)

	filename := fmt.Sprintf("./assets/logs/log_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

func (l *logger) SetQuery(c *fiber.Ctx) {
	var body any
	if err := c.QueryParser(&body); err != nil {
		log.Printf("Query parser error: %v", err)
	}
	l.Query = body
}

func (l *logger) SetBody(c *fiber.Ctx) {
	// message เป็น patch, put, post
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("Body parser error: %v", err)
	}

	switch l.Path {
	case "v1/users/signup":
		l.Body = "Signup"
	default:
		l.Body = body
	}
}

func (l *logger) SetResponse(res any) {
	l.Response = res
}
