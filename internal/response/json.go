package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonSuccess(c *gin.Context, data interface{}) {
	JsonSuccessWithMessage(c, data, "success")
}

func JsonSuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": message, "data": data})
}

func JsonFail(c *gin.Context, code int, message string) {
	JsonFailWithMessage(c, code, nil, message)
}

func JsonFailWithMessage(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message, "data": data})

}
