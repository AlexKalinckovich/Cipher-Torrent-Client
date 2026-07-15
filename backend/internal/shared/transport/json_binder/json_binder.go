package json_binder

import "github.com/gin-gonic/gin"

type BindResult[T any] struct {
	Value T
	Err   error
}

func BindJSON[T any](c *gin.Context) BindResult[T] {
	var value T
	err := c.ShouldBindJSON(&value)
	return BindResult[T]{Value: value, Err: err}
}
