package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"log"
	"reflect"
	"strings"
	"task-queue-asynq/configs"
	"task-queue-asynq/controllers"
	"task-queue-asynq/exceptions"
)

var envConf configs.EnvConf

func init() {
	envConf = configs.NewEnv()
}

func main() {
	redisAddr := envConf.RedisUrl

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	asynqmonUI := asynqmon.New(asynqmon.Options{
		RootPath:     "/dashboard", // RootPath specifies the root for asynqmonUI app
		RedisConnOpt: asynq.RedisClientOpt{Addr: redisAddr},
	})

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	emailController := controllers.NewEmailController(validate, asynqClient)
	imageController := controllers.NewImageController(validate, asynqClient)

	router := gin.Default()
	//http.Handle(asynqmonUI.RootPath()+"/", asynqmonUI)
	//
	//// Go to http://localhost:8080/monitoring to see asynqmon homepage.
	//log.Fatal(http.ListenAndServe(":8080", nil))

	router.GET(asynqmonUI.RootPath()+"/*path", gin.WrapH(asynqmonUI))

	router.Use(gin.CustomRecovery(exceptions.ErrorHandler))

	imagePath := router.Group("image")
	imagePath.POST("resize", imageController.Resize)

	emailPath := router.Group("email")
	emailPath.POST("delivery", emailController.Delivery)

	err := router.Run()
	if err != nil {
		log.Fatalf("failed start server: %v", err)
	}
}
