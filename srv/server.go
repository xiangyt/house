package srv

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xiangyt/house/task"
)

// web服务器
type webService struct {
	Addr               string
	engine             *gin.Engine
	beforeStartHandler []HandlerFunc
	afterStartHandler  []func() error
	beforeStopHandler  []HandlerFunc
	afterStopHandler   []HandlerFunc
}

type HandlerFunc func(en *gin.Engine) error

func NewWebService(addr string) *webService {
	return &webService{
		Addr:   addr,
		engine: gin.Default(),
	}
}

func (web *webService) BeforeStart(hs ...HandlerFunc) {
	web.beforeStartHandler = append(web.beforeStartHandler, hs...)
}

func (web *webService) AfterStart(hs ...func() error) {
	web.afterStartHandler = append(web.afterStartHandler, hs...)
}

func (web *webService) BeforeStop(hs ...HandlerFunc) {
	web.beforeStopHandler = append(web.beforeStopHandler, hs...)
}

func (web *webService) AfterStop(hs ...HandlerFunc) {
	web.afterStopHandler = append(web.afterStopHandler, hs...)
}

func (web *webService) Run() {
	srv := &http.Server{
		Addr:    web.Addr,
		Handler: web.engine,
	}

	start := make(chan error)
	go func() {
		for _, handler := range web.beforeStartHandler {
			if err := handler(web.engine); err != nil {
				logrus.Fatal("beforeStart err: ", err)
			}
		}

		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			close(start)
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		select {
		case <-start:
			return
		case <-time.After(5 * time.Second):
			// 启动定时任务
			task.GetManager().Start()
			for _, handler := range web.afterStartHandler {
				if err := handler(); err != nil {
					logrus.Fatal("beforeStart err: ", err)
				}
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("Shutdown Server ...")

	for _, handler := range web.beforeStopHandler {
		if err := handler(web.engine); err != nil {
			logrus.Fatal("beforeStop err: ", err)
		}
	}

	// 等待定时任务结束
	task.GetManager().Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown: ", err)
	}

	for _, handler := range web.afterStopHandler {
		if err := handler(web.engine); err != nil {
			logrus.Fatal("afterStop err: ", err)
		}
	}
	logrus.Println("Server exiting")
}
