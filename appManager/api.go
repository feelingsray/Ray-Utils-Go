package appManager

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

/*
 * 启动(重启)服务
 */
func (p *AppManager)procRestart(c *gin.Context)  {
    code := c.Query("code")
    name := c.Query("name")
    if code == "" {
        // 全部重启
        data,err := p.AllProcRestart()
        if err != nil {
            c.JSON(500,err.Error())
            return
        }
        c.JSON(200,data)
        return
    }
    c.JSON(200,p.ProcRestart(code,name))
    return
}

/*
 * 停止服务
 */
func (p *AppManager)procStop(c *gin.Context)  {
    code := c.Query("code")
    if code == "" {
        c.JSON(500,"服务编码不能为空")
    }
    p.ProcStop(code)
    c.JSON(200,fmt.Sprintf("发送停止指令:%s",code))
    return
}

/*
 * 停止服务
 */
func (p *AppManager)procList(c *gin.Context)  {
    serviceList := make(map[string]*Proc)
    p.ProcStore.Range(func(key, value interface{}) bool {
        serviceList[key.(string)]=value.(*Proc)
        return true
    })
    c.JSON(200,serviceList)
    return
}
