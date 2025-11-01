package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/wjhcoding/metanode-task-go-blog/internal/dao/pool"
	"github.com/wjhcoding/metanode-task-go-blog/internal/middleware"
	"github.com/wjhcoding/metanode-task-go-blog/internal/model"
	"github.com/wjhcoding/metanode-task-go-blog/pkg/common/response"
)

// 用户注册请求体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email"`
}

// 用户登录请求体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("参数错误："+err.Error()))
		return
	}

	// 检查用户名是否已存在
	var existing model.User
	if err := pool.GetDB().Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("用户名已存在"))
		return
	}

	// 密码加密
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	if err := pool.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("注册失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkMsg("注册成功"))
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("参数错误："+err.Error()))
		return
	}

	var user model.User
	if err := pool.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, response.FailMsg("用户不存在"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.FailMsg("数据库错误："+err.Error()))
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, response.FailMsg("密码错误"))
		return
	}

	// 生成 JWT Token
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("Token生成失败"))
		return
	}

	c.JSON(http.StatusOK, response.OkData(gin.H{
		"token":    token,
		"username": user.Username,
		"user_id":  user.ID,
	}))
}
