package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
	"strings"
	"cmdb/utils"
	"strconv"
)

type ZabbixUserController struct {
	beego.Controller
}

type ZabbixUser struct {
	Id int64
	UserName string
	UserId int64
	GroupName string
	GroupId string
	MediaName string
	MediaTypeId string
	SendtoPeriod string
	Type string
	Count int64
}

func (c *ZabbixUserController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixUserMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixUser/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixUser/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalZabbixUsers()
		c.TplName = "zabbixUser/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *ZabbixUserController) GetUserList(){
	var zabbixUserArr []ZabbixUser
	var zabbixUser ZabbixUser
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	userMaps,count,err := models.GetUserList(page,keyword)
	if err == nil {
		for _,m := range userMaps {
			zabbixUser.Id = m["Id"].(int64)
			zabbixUser.UserId = m["UserId"].(int64)
			zabbixUser.UserName = m["UserName"].(string)
			if m["GroupName"] == nil {
				zabbixUser.GroupName = ""
			} else {
				zabbixUser.GroupName = m["GroupName"].(string)
			}
			if m["GroupId"] == nil {
				zabbixUser.GroupId = ""
			} else {
				zabbixUser.GroupId = m["GroupId"].(string)
			}
			if m["MediaTypeId"] == nil {
				zabbixUser.MediaTypeId = ""
			} else {
				zabbixUser.MediaTypeId = m["MediaTypeId"].(string)
			}
			mediaName := models.GetMediaNameByMediaTypeId(m["MediaTypeId"].(string), m["Type"].(string))
			if mediaName == "" {
				zabbixUser.MediaName = "-"
			} else {
				zabbixUser.MediaName = mediaName
			}
			zabbixUser.Type = m["Type"].(string)
			if m["SendtoPeriod"].(string) == "" {
				zabbixUser.SendtoPeriod = "-"
			} else {
				zabbixUser.SendtoPeriod = m["SendtoPeriod"].(string)
			}
			zabbixUser.Count = count
			zabbixUserArr = append(zabbixUserArr,zabbixUser)
		}
		c.Data["json"] = &zabbixUserArr
		c.ServeJSON()
	}
}

func (c *ZabbixUserController) DoZabbixUser(){
	var message Message
	id,err := c.GetInt("id")
	userId := c.GetString("userId","")
	//username := strings.Trim(c.GetString("UserName",""),"")
	//groupIds := c.GetStrings("GroupIds")
	sendto := strings.Trim(c.GetString("Sendto",""),"")
	period := strings.Trim(c.GetString("Period","1-7,00:00-24:00"),"")
	environment := c.GetString("environment","")
	//groupNames := c.GetString("GroupNames","")
	allocateUserGroups := c.GetString("allocate_user_group_list","")
	userGroupNames := c.GetString("userGroupNames","")
	//if username != "" && len(groupIds) != 0 && environment != "" && groupNames != "" {
	//	userId,isAdd,err := utils.AddUser(username,environment,sendto,period,groupIds)
	//	if isAdd && err == "" {
	//		var groupId string
	//		for _,id := range groupIds {
	//			if groupId == "" {
	//				groupId = id
	//			} else {
	//				groupId = groupId+","+id
	//			}
	//		}
	//		models.AddZabbixUser(username,groupNames,groupId,sendto,period,environment,userId)
	//		message.Status = "ok"
	//		message.Msg = ""
	//	} else {
	//		message.Status = "fail"
	//		message.Msg = fmt.Sprintf("Error:%v",err)
	//	}
	//}
	if err == nil && userId != "" && environment != ""{
		isDelete,err := utils.DeleteUser(userId,environment)
		if isDelete && err == "" {
			models.DeleteZabbixUser(int64(id))
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if userId != "" && environment != "" && ((allocateUserGroups != "" && userGroupNames != "") || sendto != ""){
		var allocateUserGroupList []int64
		if allocateUserGroups != "" {
			json.Unmarshal([]byte(allocateUserGroups), &allocateUserGroupList)
		}
		isUpdate,err := utils.UpdateUserGroup(userId,environment,sendto,period,allocateUserGroupList)
		if isUpdate && err == "" {
			var userGroupIds string
			if len(allocateUserGroupList) != 0 {
				for _,groupId := range allocateUserGroupList {
					groupIdStr := strconv.Itoa(int(groupId))
					if userGroupIds != "" {
						userGroupIds = userGroupIds+","+groupIdStr
					} else {
						userGroupIds = groupIdStr
					}
				}
			}
			models.UpdateZabbixUser(userId, userGroupIds, userGroupNames, environment, sendto,period)
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}
