package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIntQuery(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func getBoolQuery(c *gin.Context, key string, defaultVal bool) bool {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}
