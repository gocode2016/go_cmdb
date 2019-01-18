package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
	"strings"
	"cmdb/utils"
)

type ZabbixTemplateController struct {
	beego.Controller
}

type zabbixTemplate struct {
	Id int64
	TemplateName string
	TemplateId int64
	Type string
	Count int64
}

func (c *ZabbixTemplateController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixTemplateMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixTemplate/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixTemplate/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalTemplates()
		c.TplName = "zabbixTemplate/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *ZabbixTemplateController) GetTemplateList(){
	var zabbixTemplateArr []zabbixTemplate
	var zabbixTemplate zabbixTemplate
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	environment := c.GetString("environment","")
	if err != nil {
		page = 1
	}
	groupMaps,count,err := models.GetTemplateList(page,keyword,environment)
	if err == nil {
		for _,m := range groupMaps {
			zabbixTemplate.Id = m["Id"].(int64)
			zabbixTemplate.TemplateName = m["TemplateName"].(string)
			zabbixTemplate.TemplateId = m["TemplateId"].(int64)
			zabbixTemplate.Type = m["Type"].(string)
			zabbixTemplate.Count = count
			zabbixTemplateArr = append(zabbixTemplateArr,zabbixTemplate)
		}
		c.Data["json"] = &zabbixTemplateArr
		c.ServeJSON()
	}
}

func (c *ZabbixTemplateController) DoZabbixTemplate(){
	var message Message
	id,err := c.GetInt("Id")
	templateId := c.GetString("TemplateId")
	templateName := strings.Trim(c.GetString("TemplateName","")," ")
	environment := c.GetString("environment","")
	groupIds := c.GetStrings("GroupIds")
	templateIds := c.GetStrings("TemplateIds")
	if templateName != "" && environment != "" && len(groupIds) != 0 {
		isExist,err := utils.CheckTemplateNameExist(templateName, environment)
		if isExist && err == nil {
			templateId,isAdd,err := utils.AddTemplate(templateName, environment, groupIds, templateIds)
			if isAdd && err == "" {
				models.AddTemplate(templateName, environment, templateId)
				message.Status = "ok"
				message.Msg = ""
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Error:%v",err)
			}
		} else {
			message.Status = "fail"
			message.Msg = err.Error()
		}
	}
	if err == nil && templateId != "" && environment != "" {
		isDelete,err := utils.DeleteTemplate(templateId, environment)
		if isDelete && err == "" {
			models.DeleteTemplate(int64(id))
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if templateId != "" && templateName != "" && environment != "" {
		isUpdate,err := utils.UpdateTemplate(templateId,templateName,environment)
		if isUpdate && err == "" {
			models.UpdateTemplate(templateId, templateName, environment)
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("%v",err)
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}