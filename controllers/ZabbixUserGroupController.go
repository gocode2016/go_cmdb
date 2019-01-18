package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
	"cmdb/utils"
	"strings"
)

type ZabbixUserGroupController struct {
	beego.Controller
}

type ZabbixUserGroup struct {
	Id int64
	GroupName string
	GroupId int64
	GroupStatus bool
	Type string
	Count int64
}

func (c *ZabbixUserGroupController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixUserGroupMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixUserGroup/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixUserGroup/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalZabbixUserGroups()
		c.TplName = "zabbixUserGroup/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *ZabbixUserGroupController) GetUserGroupList(){
	var zabbixUserGroupArr []ZabbixUserGroup
	var zabbixUserGroup ZabbixUserGroup
	keyword := c.GetString("keyword","")
	environment := c.GetString("environment","")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	groupMaps,count,err := models.GetUserGroupList(page,keyword,environment)
	if err == nil {
		for _,m := range groupMaps {
			zabbixUserGroup.Id = m["Id"].(int64)
			zabbixUserGroup.GroupName = m["GroupName"].(string)
			zabbixUserGroup.GroupId = m["GroupId"].(int64)
			zabbixUserGroup.Type = m["Type"].(string)
			zabbixUserGroup.GroupStatus = m["GroupStatus"].(bool)
			zabbixUserGroup.Count = count
			zabbixUserGroupArr = append(zabbixUserGroupArr,zabbixUserGroup)
		}
		c.Data["json"] = &zabbixUserGroupArr
		c.ServeJSON()
	}
}

func (c *ZabbixUserGroupController) DoZabbixUserGroup(){
	var message Message
	id,err := c.GetInt("Id")
	groupId := c.GetString("GroupId")
	environment := c.GetString("environment")
	groupStatus := c.GetString("GroupStatus")
	groupName := strings.Trim(c.GetString("GroupName")," ")
	if groupId == "" && environment != "" && groupName != "" {
		isExist,err := utils.CheckUserGroupNameExist(groupName, environment)
		if isExist && err == nil {
			groupId,isAdd,err := utils.AddUserGroup(groupName, environment)
			if isAdd && err == "" {
				models.AddUserGroup(groupName,environment,groupId)
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
		isDelete,err := utils.DeleteUserGroup(groupId, environment)
		if isDelete && err == "" {
			models.DeleteUserGroup(int64(id))
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if groupId != "" && environment != "" && groupStatus != "" {
		isChange,err := utils.ChangeUserGroupStatus(groupId,environment,groupStatus)
		if isChange && err == "" {
			models.UpdateUserGroupStatus(groupId,groupStatus,environment)
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if groupId != "" && environment != "" && groupName != "" {
		isExist,err := utils.CheckUserGroupNameExist(groupName, environment)
		if isExist && err == nil {
			isUpdate,err := utils.UpdateUserGroupName(groupId, groupName, environment)
			if isUpdate && err == "" {
				models.UpdateUserGroup(groupId, groupName, environment)
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
