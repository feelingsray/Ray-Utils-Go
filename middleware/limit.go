package middleware

import (
    "log"
    "time"
    
    "github.com/gin-gonic/gin"
    "go.uber.org/ratelimit"
)

func Ratelimit(limit ratelimit.Limiter) gin.HandlerFunc {
    prev := time.Now()
    return func(ctx *gin.Context) {
        now := limit.Take()
        log.Printf("%v\n", now.Sub(prev))
        prev = now
    }
}
