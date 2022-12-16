package middleware

import (
  "bufio"
  "bytes"
  "net"
  "net/http"
  "strings"
  "time"
  
  "github.com/gin-gonic/gin"
  "github.com/tidwall/sjson"
)

const (
  noWritten     = -1
  defaultStatus = 200
)

// Cost 接口耗时统计中间件
func Cost() gin.HandlerFunc {
  return func(c *gin.Context) {
    if strings.Contains(c.Request.URL.Path, "download") {
      c.Next()
    }
    start := time.Now()
    var wb *ResponseBuffer
    if w, ok := c.Writer.(gin.ResponseWriter); ok {
      wb = NewResponseBuffer(w)
      c.Writer = wb
      c.Next()
    }
    data := wb.Body.Bytes()
    wb.Body.Reset()
    body, _ := sjson.SetBytes(data, "cost", time.Since(start).Milliseconds())
    wb.Body.Write(body)
    wb.Flush()
  }
}

type ResponseBuffer struct {
  Response gin.ResponseWriter
  status   int
  Body     *bytes.Buffer
  Flushed  bool
}

func NewResponseBuffer(w gin.ResponseWriter) *ResponseBuffer {
  return &ResponseBuffer{
    Response: w, status: defaultStatus, Body: &bytes.Buffer{},
  }
}

func (w *ResponseBuffer) Header() http.Header {
  return w.Response.Header() // use the actual response header
}

func (w *ResponseBuffer) Pusher() http.Pusher {
  return w.Response.Pusher()
}

func (w *ResponseBuffer) Write(buf []byte) (int, error) {
  w.Body.Write(buf)
  return len(buf), nil
}

func (w *ResponseBuffer) WriteString(s string) (n int, err error) {
  n, err = w.Write([]byte(s))
  return
}

func (w *ResponseBuffer) Written() bool {
  return w.Body.Len() != noWritten
}

func (w *ResponseBuffer) WriteHeader(status int) {
  w.status = status
}

func (w *ResponseBuffer) WriteHeaderNow() {
  //if !w.Written() {
  //	w.size = 0
  //	w.ResponseWriter.WriteHeader(w.status)
  //}
}

func (w *ResponseBuffer) Status() int {
  return w.status
}

func (w *ResponseBuffer) Size() int {
  return w.Body.Len()
}

func (w *ResponseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
  return w.Response.(http.Hijacker).Hijack()
}

func (w *ResponseBuffer) CloseNotify() <-chan bool {
  return w.Response.(http.CloseNotifier).CloseNotify()
}

func (w *ResponseBuffer) Flush() {
  w.realFlush()
}

func (w *ResponseBuffer) realFlush() {
  if w.Flushed {
    return
  }
  w.Response.WriteHeader(w.status)
  if w.Body.Len() > 0 {
    _, err := w.Response.Write(w.Body.Bytes())
    if err != nil {
      panic(err)
    }
    w.Body.Reset()
  }
  w.Flushed = true
}
