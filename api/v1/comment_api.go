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

// CreateComment 新增评论
func CreateComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, response.FailMsg("参数错误："+err.Error()))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.FailMsg("未登录"))
		return
	}

	// 评论必须关联文章
	if comment.PostID == 0 {
		c.JSON(http.StatusBadRequest, response.FailMsg("缺少文章ID"))
		return
	}

	comment.UserID = userID.(uint)
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	if err := pool.GetDB().Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("评论失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkData(comment))
}

// GetCommentsByPostID 查询文章下的评论列表
func GetCommentsByPostID(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []model.Comment

	if err := pool.GetDB().Preload("User").Where("post_id = ?", postID).
		Order("created_at asc").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkData(comments))
}

// DeleteComment 删除评论（作者或管理员）
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment model.Comment

	if err := pool.GetDB().First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, response.FailMsg("评论不存在"))
		} else {
			c.JSON(http.StatusInternalServerError, response.FailMsg("查询失败："+err.Error()))
		}
		return
	}

	userID := c.GetUint("user_id")
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, response.FailMsg("无权限删除该评论"))
		return
	}

	if err := pool.GetDB().Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.FailMsg("删除失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OkMsg("删除成功"))
}
