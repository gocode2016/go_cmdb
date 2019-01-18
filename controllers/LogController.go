package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"time"
	"strings"
	"encoding/json"
	"fmt"
)

type LogController struct {
	beego.Controller
}

type Log struct {
	Id int64
	Name string
	Ip string
	Content string
	Count int64
	CreateTime string
}

func (c *LogController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.Data["LogMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "log/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "log/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalLogs()
		c.TplName = "log/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *LogController) GetLogList(){
	var logArr []Log
	var log Log
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	keyword := strings.Trim(c.GetString("keyword",""),"")
	sort := c.GetString("sort","")
	logMaps,count,err := models.GetLogList(page,keyword,sort)
	if err == nil {
		for _,m := range logMaps {
			log.Id = m["Id"].(int64)
			log.Name = m["Name"].(string)
			log.Ip = m["Ip"].(string)
			content := m["Content"].(string)
			if len(content) > 128 {
				content = fmt.Sprintf("%v....",content[0:129])
			}
			log.Content = content
			log.Count = count
			log.CreateTime = m["CreateTime"].(time.Time).Format("2006-01-02 15:04:05")
			logArr = append(logArr,log)
		}
		c.Data["json"] = logArr
		c.ServeJSON()
	}
}
