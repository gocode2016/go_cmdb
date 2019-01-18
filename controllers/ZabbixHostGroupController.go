package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
	"strings"
	"cmdb/utils"
)

type ZabbixHostGroupController struct {
	beego.Controller
}

type zabbixHostGroup struct {
	Id int64
	GroupName string
	GroupId int64
	Type string
	Count int64
}

func (c *ZabbixHostGroupController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixHostGroupMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixHostGroup/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixHostGroup/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalHostGroups()
		c.TplName = "zabbixHostGroup/index.html"
	} else {
		//c.Redirect("http://aiops.jollycorp.net/",302)
		c.Redirect(c.URLFor("HomeController.Login"),302)
	}
}

func (c *ZabbixHostGroupController) GetHostGroupList(){
	var zabbixHostGroupArr []zabbixHostGroup
	var zabbixHostGroup zabbixHostGroup
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	environment := c.GetString("environment","")
	if err != nil {
		page = 1
	}
	groupMaps,count,err := models.GetHostGroupList(page,keyword,environment)
	if err == nil {
		for _,m := range groupMaps {
			zabbixHostGroup.Id = m["Id"].(int64)
			zabbixHostGroup.GroupName = m["GroupName"].(string)
			zabbixHostGroup.GroupId = m["GroupId"].(int64)
			zabbixHostGroup.Type = m["Type"].(string)
			zabbixHostGroup.Count = count
			zabbixHostGroupArr = append(zabbixHostGroupArr,zabbixHostGroup)
		}
		c.Data["json"] = &zabbixHostGroupArr
		c.ServeJSON()
	}
}

func (c *ZabbixHostGroupController) DoZabbixHostGroup(){
	var message Message
	id,err := c.GetInt("Id")
	groupId := c.GetString("GroupId")
	environment := c.GetString("environment")
	groupName := strings.Trim(c.GetString("GroupName")," ")
	if groupId == "" && environment != "" && groupName != "" {
		isExist,err := utils.CheckHostGroupNameExist(groupName, environment)
		if isExist && err == nil {
			groupId,isAdd,err := utils.AddHostGroupName(groupName, environment)
			if isAdd && err == "" {
				models.AddHostGroup(groupName,environment,groupId)
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
	if err == nil && groupId != "" && environment != "" {
		isDelete,err := utils.DeleteHostGroupName(groupId, environment)
		if isDelete && err == "" {
			models.DeleteHostGroup(int64(id))
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if groupId != "" && environment != "" && groupName != "" {
		isExist,err := utils.CheckHostGroupNameExist(groupName, environment)
		if isExist && err == nil {
			isUpdate,err := utils.UpdateHostGroupName(groupId, groupName, environment)
			if isUpdate && err == "" {
				models.UpdateHostGroup(groupId, groupName, environment)
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
	c.Data["json"] = &message
	c.ServeJSON()
}