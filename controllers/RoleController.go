package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"fmt"
	"time"
	"encoding/json"
)

type Role struct {
	Id int64
	Name string
	Created string
}

//RoleController 角色管理
type RoleController struct {
	beego.Controller
}

type customPermission struct {
	Id int
	Name string
	Comment string
	Rw int
}
//Index 角色管理首页
func (c *RoleController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.Data["GroupMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "role/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "role/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.TplName = "role/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func(c *RoleController) GetRoleList(){
	var RoleArr []Role
	var role Role
	var message Message
	var permissionArr []customPermission
	var permission customPermission
	keyword := c.GetString("keyword","")
	role_id,err := c.GetInt("role_id")
	role_maps,err2 := models.GetRoles(keyword)
	if err == nil {
		ownedPermissions,unownedPermissions,err := models.GetPermissionsByRoleId(role_id)
		if err == nil {
			if len(ownedPermissions) != 0 {
				for _,v := range ownedPermissions {
					permission.Id = v.Id
					permission.Name = v.Name
					permission.Comment = v.Comment
					permission.Rw = 1
					permissionArr = append(permissionArr,permission)
				}
			}
			if len(unownedPermissions) != 0 {
				for _,v := range unownedPermissions {
					permission.Id = v.Id
					permission.Name = v.Name
					permission.Comment = v.Comment
					permission.Rw = 0
					permissionArr = append(permissionArr,permission)
				}
			}
			c.Data["json"] = &permissionArr
			c.ServeJSON()
		}
	}
	if err2 != nil {
		message.Msg = fmt.Sprintf("查询失败:%v",err2)
		message.Status = "fail"
		c.Data["json"] = &message
	} else {
		if role_maps != nil {
			for _, m := range role_maps {
				role.Id = m["Id"].(int64)
				role.Name = m["Name"].(string)
				role.Created = m["Created"].(time.Time).Format("2006-01-02 15:04:05")
				RoleArr = append(RoleArr, role)
			}
			c.Data["json"] = RoleArr
		} else {
			message.Msg = "查无此结果"
			message.Status = "fail"
			c.Data["json"] = &message
		}
	}
	c.ServeJSON()
}

func (c *RoleController) DoRole() {
	var message Message
	var ids []int
	id, err := c.GetInt("id")
	grant_list := c.GetString("grant_list","")
	role_name := c.GetString("role_name","")
	if err == nil && grant_list != "" {
		json.Unmarshal([]byte(grant_list), &ids)
		isModify := models.ModifyPermissionsForRole(id,ids)
		if isModify {
			message.Msg = "添加成功"
			message.Status = "ok"
		} else {
			message.Msg = "修改权限失败"
			message.Status = "fail"
		}
	}
	if role_name != "" && grant_list == "" {
		isAdd,err := models.AddRole(role_name)
		if isAdd && err == nil {
			message.Msg = "添加成功"
			message.Status = "ok"
		} else {
			message.Msg = err.Error()
			message.Status = "fail"
		}
	}
	if err == nil && grant_list == "" {
		isDeleted := models.DeleteRole(id)
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
