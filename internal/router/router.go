package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/wjhcoding/metanode-task-go-blog/api/v1"
	"github.com/wjhcoding/metanode-task-go-blog/internal/middleware"
	"github.com/wjhcoding/metanode-task-go-blog/pkg/common/response"
	"github.com/wjhcoding/metanode-task-go-blog/pkg/global/log"
	"go.uber.org/zap"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	server := gin.Default()
	server.Use(Cors())
	server.Use(Recovery)
	// server.Use(gin.Recovery())

	// API v1 分组
	api := server.Group("/api/v1")
	{
		// 🧍 用户模块
		api.POST("/user/register", v1.Register)
		api.POST("/user/login", v1.Login)

		// 📰 文章模块（需要登录）
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.POST("/posts", v1.CreatePost)
			auth.GET("/posts", v1.GetPostList)
			auth.GET("/posts/:id", v1.GetPostByID)
			auth.PUT("/posts/:id", v1.UpdatePost)
			auth.DELETE("/posts/:id", v1.DeletePost)

			// 💬 评论模块
			auth.POST("/comments", v1.CreateComment)
			auth.GET("/comments/:post_id", v1.GetCommentsByPostID)
			auth.DELETE("/comments/:id", v1.DeleteComment)
		}
	}

	// 健康检测接口
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return server
}

// ----------------- 以下为通用中间件 -----------------

// Cors 跨域处理
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error("HttpError", zap.Any("HttpError", err))
			}
		}()

		c.Next()
	}
}

func Recovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("gin catch error: ", log.Any("gin catch error: ", r))
			c.JSON(http.StatusOK, response.FailMsg("系统内部错误"))
		}
	}()
	c.Next()
}
