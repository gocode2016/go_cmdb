package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"encoding/json"
	"fmt"
	"cmdb/utils"
	"strings"
	"strconv"
)

type ZabbixHostController struct {
	beego.Controller
}

type ZabbixHost struct {
	Id int64
	HostName string
	VisibleName string
	HostId int64
	Ip string
	GroupName string
	GroupId string
	TemplateName string
	TemplateId string
	Type string
	HostStatus bool
	InterfaceType string
	Count int64
}

func (c *ZabbixHostController) Index() {
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["zabbixHostMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "zabbixHost/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "zabbixHost/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalHosts()
		c.TplName = "zabbixHost/index.html"
	} else {
		c.Redirect("http://aiops.jollycorp.net/",302)
		//c.Redirect(c.URLFor("HomeController.Login"),302)
	}
}

func (c *ZabbixHostController) GetHostList(){
	var zabbixHostArr []ZabbixHost
	var zabbixHost ZabbixHost
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	hostMaps,count,err := models.GetHostList(page,keyword)
	if err == nil {
		for _,m := range hostMaps {
			zabbixHost.Id = m["Id"].(int64)
			zabbixHost.HostId = m["HostId"].(int64)
			zabbixHost.HostName = m["HostName"].(string)
			zabbixHost.Ip = m["Ip"].(string)
			if m["VisibleName"] == nil {
				zabbixHost.VisibleName = ""
			} else {
				zabbixHost.VisibleName = m["VisibleName"].(string)
			}
			if m["GroupName"] == nil {
				zabbixHost.GroupName = ""
			} else {
				zabbixHost.GroupName = m["GroupName"].(string)
			}
			if m["GroupId"] == nil {
				zabbixHost.GroupId = ""
			} else {
				zabbixHost.GroupId = m["GroupId"].(string)
			}
			if m["TemplateName"] == nil {
				zabbixHost.TemplateName = ""
			} else {
				zabbixHost.TemplateName = m["TemplateName"].(string)
			}
			if m["TemplateId"] == nil {
				zabbixHost.TemplateId = ""
			} else {
				zabbixHost.TemplateId = m["TemplateId"].(string)
			}
			if m["InterfaceType"] == nil {
				zabbixHost.InterfaceType = ""
			} else {
				zabbixHost.InterfaceType = m["InterfaceType"].(string)
			}
			zabbixHost.Type = m["Type"].(string)
			zabbixHost.Count = count
			zabbixHost.HostStatus = m["HostStatus"].(bool)
			zabbixHostArr = append(zabbixHostArr,zabbixHost)
		}
		c.Data["json"] = &zabbixHostArr
		c.ServeJSON()
	}
}

func (c *ZabbixHostController) DoZabbixHost() {
	var message Message
	hostname := strings.Trim(c.GetString("HostName","")," ")
	interfaceAddr := c.GetString("InterfaceAddr","")
	interfaceValue := strings.Trim(c.GetString("InterfaceValue","")," ")
	interfaceType := c.GetString("InterfaceType")
	environment := c.GetString("environment","")
	groupIdList := c.GetStrings("GroupIds")
	templateIdList := c.GetStrings("TemplateIds")
	groupNames := c.GetString("GroupNames","")
	templateNames := c.GetString("TemplateNames","")
	HostId,err2 := c.GetInt("HostId")
	HostStatus := c.GetString("HostStatus","")
	Id,err3 := c.GetInt("Id")
	Type := c.GetString("Type")
	allocateHostGroups := c.GetString("allocate_host_group_list")
	allocateTemplates := c.GetString("allocate_template_list")
	editedEnvironment := c.GetString("EditedEnvironment","")
	allocateTemplateNames := c.GetString("AllocateTemplateNames","")
	allocateGroupNames := c.GetString("AllocateGroupNames","")
	cnHostId := c.GetString("cnHostId")
	foreignHostId := c.GetString("foreignHostId")
	visibleName := strings.Trim(c.GetString("VisibleName",""),"")
	if err2 == nil && editedEnvironment != "" && allocateHostGroups != "" || allocateTemplates != "" || hostname != "" ||  visibleName != ""{
		var allocateHostGroupList []int64
		var allocateTemplateList []int64
		if allocateHostGroups != "" {
			json.Unmarshal([]byte(allocateHostGroups), &allocateHostGroupList)
		}
		if allocateTemplates != "" {
			json.Unmarshal([]byte(allocateTemplates), &allocateTemplateList)
		}
		isUpdate,err := utils.UpdateHost(editedEnvironment,hostname,visibleName,int64(HostId),allocateHostGroupList,allocateTemplateList)
		if isUpdate && err == "" {
			var updateHostGroups,updateTemplates string
			if allocateHostGroups != "" {
				for _,group := range allocateHostGroupList {
					groupStr := strconv.Itoa(int(group))
					if updateHostGroups != "" {
						updateHostGroups = updateHostGroups+","+groupStr
					} else {
						updateHostGroups = groupStr
					}
				}
			}
			if allocateTemplates != "" {
				for _,template := range allocateTemplateList {
					templateStr := strconv.Itoa(int(template))
					if updateTemplates != "" {
						updateTemplates = updateTemplates+","+templateStr
					} else {
						updateTemplates = templateStr
					}
				}
			}
			models.UpdateHost(int64(HostId),updateHostGroups,updateTemplates,allocateTemplateNames,allocateGroupNames,editedEnvironment,hostname,visibleName)
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if hostname != "" && visibleName != "" && interfaceAddr != "" && interfaceValue != "" && environment != "" {
		isExist,err := utils.CheckHostExist(hostname,interfaceAddr,interfaceValue,environment)
		if isExist && err == nil {
			hostid,IsAdd,err := utils.AddHost(hostname,visibleName,interfaceAddr,interfaceValue,interfaceType,environment,groupIdList,templateIdList)
			if IsAdd && err == "" {
				var templateids = ""
				for _,templateid := range templateIdList {
					if templateids != "" {
						templateids = templateids+","+templateid
					} else {
						templateids = templateid
					}
				}
				var groupids = ""
				for _,groupid := range groupIdList {
					if groupids != "" {
						groupids = groupids+","+groupid
					} else {
						groupids = groupid
					}
				}
				IsAdd,err := models.AddHost(hostname,visibleName,interfaceAddr,interfaceValue,interfaceType,groupNames,groupids,templateNames,templateids,environment,hostid)
				if IsAdd && err == nil {
					message.Status = "ok"
					message.Msg = ""
				} else {
					message.Status = "fail"
					message.Msg = err.Error()
				}
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Error:%v",err)
			}
		} else {
			message.Status = "fail"
			message.Msg = err.Error()
		}
	}
	if err2 == nil && HostStatus != "" && Type != "" {
		IsUpdate,err := utils.ChangeStatus(int64(HostId), HostStatus, Type)
		if IsUpdate && err == "" {
			models.UpdateHostStatus(int64(HostId),HostStatus, Type)
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if err2 == nil && err3 == nil && Type != ""{
		IsDelete,err := utils.DeleteHost(int64(HostId),Type)
		if IsDelete && err == "" {
			models.DeleteHost(int64(Id))
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err)
		}
	}
	if cnHostId != "" || foreignHostId != "" {
		var cnHostIdArr,foreignHostIdArr []string
		var isCnMassDelete = false
		var isForeignMassDelete = false
		var err1,err2 string
		if cnHostId != "" {
			json.Unmarshal([]byte(cnHostId), &cnHostIdArr)
			isCnMassDelete,err1 = utils.MassDeleteCnHost(cnHostIdArr)
		}
		if foreignHostId != "" {
			json.Unmarshal([]byte(foreignHostId), &foreignHostIdArr)
			isForeignMassDelete,err2 = utils.MassDeleteForeignHost(foreignHostIdArr)
		}
		if (isCnMassDelete && err1 == "") || (isForeignMassDelete && err2 == "") {
			if len(foreignHostIdArr) != 0 {
				models.MassDeleteForeignHost(foreignHostIdArr)
			}
			if len(cnHostIdArr) != 0 {
				models.MassDeleteCnHost(cnHostIdArr)
			}
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Error:%v",err1+","+err2)
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}