package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/jobx/api/middleware"
	v1 "github.com/sqc157400661/jobx/api/v1"
)

func InitRouter() *gin.Engine {

	//设置gin运行模式
	//gin.SetMode(conf.GetString("runMode"))

	r := gin.New()
	r.Use(middleware.Cors())
	gin.SetMode(gin.DebugMode)
	r.Any("health", func(c *gin.Context) {
		c.JSON(200, "ok")
	})
	//恢复中间件 默认可用：r.Use(gin.Recovery())
	r.Use(gin.Recovery())
	//设置路由分组
	CheckRoleRouter(r)
	return r
}

// 需认证的路由
func CheckRoleRouter(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	// 可根据业务需求来设置接口版本
	v1Group := r.Group("/api/v1", middlewares...)

	taskRouter := v1Group.Group("/job")
	{
		job := v1.Job{}
		// 获取任务列表
		taskRouter.GET("/list", job.GetPage)
		// 获取任务节点列表
		taskRouter.GET("/steps", job.TaskList)

		taskRouter.GET("/query", job.Get)
		// 任务重试
		taskRouter.POST("/retry", job.Retry)
		// 任务跳过
		taskRouter.POST("/skip", job.Skip)
		// 任务暂停
		taskRouter.POST("/pause", job.Pause)
		// 任务恢复
		taskRouter.POST("/restart", job.Restart)
		// 任务废弃
		taskRouter.POST("/discard", job.Discard)
		// 任务废弃
		taskRouter.POST("/force-discard", job.ForceDiscard)

		jobLog := v1.JobLog{}
		// 获取任务的日志
		taskRouter.GET("log", jobLog.Get)
	}

}
