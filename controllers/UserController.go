package controllers

import (
	"cmdb/models"
	"fmt"
	"strings"
	"github.com/astaxie/beego"
	"time"
	"encoding/json"
)

type UserController struct {
	beego.Controller
}

type Users struct {
	Id int64
	UserName string
	GivenName string
	Role string
	Email string
	IsEnabled bool
	Count int64
}

type Message struct {
	Msg string
	Status string
}

type Roles struct {
	Id int64
	Name string
	Created time.Time
}

/*func (c *BackendUserController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的多数Action都需要权限控制，则将验证放到Prepare
	c.checkAuthor("DataGrid")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//权限控制里会进行登录验证，因此这里不用再作登录验证
	//c.checkLogin()

}*/
func (c *UserController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.Data["UserMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "backenduser/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "backenduser/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalUsers()
		c.TplName = "backenduser/index.html"
	} else {
		//c.Redirect(c.URLFor("HomeController.Login"),302)
		c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func(c *UserController) GetUserList(){
	var UserArr []Users
	var user Users
	var message Message
	status,_ := c.GetInt("status",2)
	keyword := c.GetString("keyword")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	user_maps,count,err := models.GetUsers(status,keyword,page)
	if err != nil {
		message.Msg = fmt.Sprintf("查询失败:%v",err)
		message.Status = "fail"
		c.Data["json"] = &message
	} else {
		if user_maps != nil {
			for _, m := range user_maps {
				role := ""
				roles, err := models.GetRole(m["UserName"].(string))
				if err == nil && len(roles) != 0{
					for _, v := range roles {
						role += "," + v.Name
					}
					user.Role = strings.TrimLeft(role, ",")
				} else {
					user.Role = "-"
				}
				user.Id = m["Id"].(int64)
				user.UserName = m["UserName"].(string)
				if m["GivenName"] != nil {
					user.GivenName = m["GivenName"].(string)
					if user.GivenName == "" {
						user.GivenName = "-"
					}
				} else {
					user.GivenName = "-"
				}
				user.Email = m["Email"].(string)
				user.IsEnabled = m["IsEnabled"].(bool)
				user.Count = count
				UserArr = append(UserArr, user)
			}
			c.Data["json"] = &UserArr
		} else {
			message.Msg = ""
			message.Status = "fail"
			c.Data["json"] = &message
		}
	}
	c.ServeJSON()
}

func (c *UserController) DoUser(){
	id,err1 := c.GetInt("id")
	status,err2 := c.GetBool("status")
	ajax_type := c.GetString("type","")
	roles_id := c.GetString("roles_id","")
	var message Message
	if err1 == nil && roles_id != "" {
		new_roles_id := strings.TrimLeft(roles_id,",")
		role_id_list := strings.Split(new_roles_id,",")
		isModified := models.ModifyRole(id,role_id_list)
		if isModified {
			message.Msg = "修改成功"
			message.Status = "ok"
		} else {
			message.Msg = "修改失败"
			message.Status = "fail"
		}
		c.Data["json"] = &message
	}
	if err1 == nil && err2 == nil && ajax_type == "" && roles_id == ""{
		isChanged := models.ChangeStatus(id,status)
		if isChanged {
			message.Msg = "修改成功"
			message.Status = "ok"
		} else {
			message.Msg = "修改失败"
			message.Status = "fail"
		}
		c.Data["json"] = &message
	} else if err1 == nil && err2 != nil && ajax_type == "" && roles_id == ""{
		isDeleted := models.DeleteUser(id)
		if isDeleted {
			message.Msg = "删除成功"
			message.Status = "ok"
		} else {
			message.Msg = "删除失败"
			message.Status = "fail"
		}
		c.Data["json"] = &message
	}
	if err1 == nil && ajax_type != "" && roles_id == ""{
		var roleArr []Roles
		var role Roles
		roles,err := models.GetRoleList()
		if err != nil {
			message.Msg = "获取失败"
			message.Status = "fail"
			c.Data["json"] = &message
		}
		if roles == nil && err == nil {
			message.Msg = "角色内容为空"
			message.Status = "fail"
			c.Data["json"] = &message
		} else if roles != nil && err == nil {
			for _,m := range roles{
				role.Id = m["Id"].(int64)
				role.Name = m["Name"].(string)
				roleArr = append(roleArr,role)
			}
			c.Data["json"] = &roleArr
		}
	}
	c.ServeJSON()
}