package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
)

type ZabbixMediaController struct {
	beego.Controller
}

type ZabbixMedia struct {
	Id int64
	MediaName string
	MediaTypeId int64
	MediaStatus bool
	Description string
	Type string
	Count int64
}

func (c *ZabbixMediaController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixMediaMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixMedia/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixMedia/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalZabbixMedias()
		c.TplName = "zabbixMedia/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *ZabbixMediaController) GetMediaList(){
	var zabbixMediaArr []ZabbixMedia
	var zabbixMedia ZabbixMedia
	keyword := c.GetString("keyword","")
	environment := c.GetString("environment","")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	mediaMaps,count,err := models.GetMediaList(page,keyword,environment)
	if err == nil {
		for _,m := range mediaMaps {
			zabbixMedia.Id = m["Id"].(int64)
			zabbixMedia.MediaName = m["MediaName"].(string)
			zabbixMedia.MediaTypeId = m["MediaTypeId"].(int64)
			zabbixMedia.Type = m["Type"].(string)
			zabbixMedia.MediaStatus = m["MediaStatus"].(bool)
			zabbixMedia.Description = m["Description"].(string)
			zabbixMedia.Count = count
			zabbixMediaArr = append(zabbixMediaArr,zabbixMedia)
		}
		c.Data["json"] = &zabbixMediaArr
		c.ServeJSON()
	}
}
