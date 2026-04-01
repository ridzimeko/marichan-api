package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, code int, success bool, message string, data any) {
	c.JSON(code, gin.H{
		"success": success,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, code int, message string, errs any) {
	c.JSON(code, gin.H{
		"success": false,
		"message": message,
		"errors":  errs,
	})
}
