package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xiangyt/house/config"
	"github.com/xiangyt/house/db"
	"github.com/xiangyt/house/srv"
	"github.com/xiangyt/house/task"
	"net/http"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
	if config.Get().Mode == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	}

	web := srv.NewWebService(":30030")
	web.BeforeStart(func(en *gin.Engine) error {
		if err := InitDatabase(); err != nil {
			return err
		}
		if err := InitTask(); err != nil {
			return err
		}

		en.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		return nil
	})

	web.Run()

}

func InitDatabase() error {
	return db.GetInstance().Init()
}

func InitTask() error {
	err := task.GetManager().AddJob(task.Job{
		Name: "cron_refresh_house_info",
		Spec: "*/10 * * * *",
		Func: func() {
			logrus.Println("test job")
		},
	})
	if err != nil {
		return err
	}
	return nil
}
