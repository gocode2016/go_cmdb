package controllers

import (
	"cmdb/models"
	"github.com/astaxie/beego"
	_ "cmdb/utils"
	"encoding/json"
	"fmt"
	"strings"
	"cmdb/utils"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Index() {
	//var isExist bool
	//var username string
	//name := c.GetSession("name")
	//if name != nil {
	//	username = name.(string)
	//	isExist = true
	//} else {
	//	username = c.GetString("username","")
	//	sid := c.GetString("sid","")
	//	if username == "" || sid == "" {
	//		c.Redirect("http://aiops.jollycorp.net/",302)
	//	}
	//	conn,err := utils.Connect()
	//	if err != nil {
	//		c.Ctx.Output.Body([]byte(fmt.Sprintf("Redis err:%v",err)))
	//	}
	//	isExist = utils.GetKey(conn,sid)
	//}
	username := c.GetSession("name").(string)
	if username != "" {
		c.SetSession("name",username)
		c.Data["user"],_ = models.GetUserInfo(username)
		permissionMap := models.GetPermissionByUsername(username)
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["footerjs"] = "home/index_footerjs.html"
		c.Data["UserTotal"] = models.GetTotalUsers()
		c.Data["AssetTotal"] = models.GetTotalAssets()
		AddAssetList,UpdateAssetList := models.GetStatisticNumber()
		c.Data["AddAssetList"] = AddAssetList
		c.Data["UpdateAssetList"] = UpdateAssetList
		c.Layout = "shared/layout_page.html"
		c.TplName = "home/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *HomeController) Login() {
	isLogin := c.GetSession("isLogin")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "home/login_headcssjs.html"
	c.LayoutSections["footerjs"] = "home/login_footerjs.html"
	c.Layout = "shared/layout_base.html"
	if isLogin == "1" {
		username := c.GetSession("name")
		passwd := c.GetSession("passwd")
		c.Data["Username"] = username
		c.Data["Passwd"] = passwd
	} else {
		c.Data["Username"] = ""
		c.Data["Passwd"] = ""
	}
	c.TplName = "home/login.html"
}

func (c *HomeController) DoLogin() {
	username := strings.TrimSpace(c.GetString("UserName"))
	passwd := strings.TrimSpace(c.GetString("Passwd"))
	remeber := c.GetString("Remeber")
	if len(username) == 0 || len(passwd) == 0 {
		message := Message{"用户名或密码不能为空","fail"}
		c.Data["json"] = &message
		c.ServeJSON()
	}
	isExist := utils.UserAuthLdap(username,passwd)
	if isExist {
		user, err := models.GetUser(username)
		if err == nil {
			if ! user.IsEnabled {
				message := Message{"用户被禁用，请联系管理员", "fail"}
				c.Data["json"] = &message
				c.ServeJSON()
			}
			if user.IsLocked {
				message := Message{"账户已被锁定，请联系管理员", "fail"}
				c.Data["json"] = &message
				c.ServeJSON()
			}
			if len(remeber) != 0 {
				c.SetSession("isLogin", "1")
			}
		}
		c.SetSession("name",user.UserName)
		c.SetSession("passwd",passwd)
		if models.CheckRoleExist(user.UserName) {
			models.AddRoleForUser(user.UserName)
		}
		message := Message{"登录成功","ok"}
		c.Data["json"] = &message
		c.ServeJSON()
	} else {
		message := Message{"用户名或者密码错误","fail"}
		c.Data["json"] = &message
		c.ServeJSON()
	}
}

func (c *HomeController) Logout() {
	c.Redirect(c.URLFor("HomeController.Login"),302)
	//c.DelSession("name")
	//c.Redirect("http://aiops.jollycorp.net/",302)
}
