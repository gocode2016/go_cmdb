package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"strings"
	"time"
	"encoding/json"
	"fmt"
)

type Permission struct {
	Id int64
	Name string
	Created string
	Count int64
	Comment string
}
//PermissionController 角色权限管理
type PermissionController struct {
	beego.Controller
}

//Index 角色权限管理首页
func (c *PermissionController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		c.Data["PermissionMenu"] = true
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "permission/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "permission/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalPermissions()
		c.TplName = "permission/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *PermissionController) GetPermissionList() {
	var permissionArr []Permission
	var permission Permission
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	keyword := strings.Trim(c.GetString("keyword",""),"")
	sort := c.GetString("sort","")
	permissionMaps,count,err := models.GetPermissionList(page,keyword,sort)
	if err == nil {
		for _,m := range permissionMaps {
			permission.Id = m["Id"].(int64)
			permission.Name = m["Name"].(string)
			permission.Comment = m["Comment"].(string)
			permission.Count = count
			permission.Created = m["Created"].(time.Time).Format("2006-01-02 15:04:05")
			permissionArr = append(permissionArr,permission)
		}
		c.Data["json"] = permissionArr
		c.ServeJSON()
	}
}

func (c *PermissionController) DoPermission() {
	var message Message
	id, err := c.GetInt("id")
	permission_name := strings.Trim(c.GetString("permission_name",""),"")
	permission_comment := strings.Trim(c.GetString("permission_comment",""),"")
	if permission_name != "" && permission_comment != ""{
		isExist,err := models.CheckPermissionExist(permission_name)
		if isExist && err == nil {
			isAdd, err := models.AddPermission(permission_name, permission_comment)
			if isAdd && err == nil {
				message.Msg = "添加成功"
				message.Status = "ok"
			} else {
				message.Msg = err.Error()
				message.Status = "fail"
			}
		} else {
			message.Msg = "权限名称已存在"
			message.Status = "fail"
		}
	}
	if id > 0 && err == nil {
		isDeleted := models.DeletePermission(id)
		if isDeleted {
			message.Msg = "删除成功"
			message.Status = "ok"
		} else {
			message.Msg = "删除失败"
			message.Status = "fail"
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}

