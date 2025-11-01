package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/wjhcoding/metanode-task-go-blog/internal/dao/pool"
	"github.com/wjhcoding/metanode-task-go-blog/internal/model"
	"github.com/wjhcoding/metanode-task-go-blog/pkg/common/response"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("参数错误："+err.Error()))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.FailMsg("未登录"))
		return
	}

	post.UserID = userID.(uint)
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := pool.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("创建失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkData(post))
}

// GetPostList 获取所有文章列表
func GetPostList(c *gin.Context) {
	var posts []model.Post
	if err := pool.GetDB().Preload("User").Order("created_at desc").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkData(posts))
}

// GetPostByID 获取单篇文章详情
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := pool.GetDB().Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.FailMsg("文章不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.OkData(post))
}

// UpdatePost 更新文章（仅作者可操作）
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := pool.GetDB().First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.FailMsg("文章不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		}
		return
	}

	userID := c.GetUint("user_id")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, response.FailMsg("无权限修改该文章"))
		return
	}

	var req model.Post
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("参数错误："+err.Error()))
		return
	}

	post.Title = req.Title
	post.Content = req.Content
	post.UpdatedAt = time.Now()

	if err := pool.GetDB().Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("更新失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkMsg("更新成功"))
}

// DeletePost 删除文章（仅作者可操作）
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post
	if err := pool.GetDB().First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.FailMsg("文章不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		}
		return
	}

	userID := c.GetUint("user_id")
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, response.FailMsg("无权限删除该文章"))
		return
	}

	if err := pool.GetDB().Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("删除失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkMsg("删除成功"))
}
