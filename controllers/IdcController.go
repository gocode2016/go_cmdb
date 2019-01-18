package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"time"
	"fmt"
	"strings"
	"encoding/json"
)

type IdcController struct {
	beego.Controller
}

type Idc struct {
	Id int64
	Name string
	Bandwidth string
	Contacts string
	Mobile string
	Address string
	Network string
	Operator string
	Comment string
	CreateTime time.Time
}

func (c *IdcController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.Data["IdcMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "idc/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "idc/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.TplName = "idc/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *IdcController) GetIdcList(){
	var idcArr []Idc
	var idc Idc
	keyword := c.GetString("keyword","")
	idcMaps,err := models.GetIdcList(keyword)
	if err == nil {
		for _,m := range idcMaps {
			idc.Id = m["Id"].(int64)
			idc.Name = m["Name"].(string)
			idc.Bandwidth = m["Bandwidth"].(string)
			idc.Contacts = m["Contacts"].(string)
			idc.Mobile = m["Mobile"].(string)
			if m["Address"] == nil {
				idc.Address = ""
			} else {
				idc.Address = m["Address"].(string)
			}
			if m["Network"] == nil {
				idc.Network = ""
			} else {
				idc.Network = m["Network"].(string)
			}
			if m["Operator"] == nil {
				idc.Operator = ""
			} else {
				idc.Operator = m["Operator"].(string)
			}
			if m["Comment"] == nil {
				idc.Comment = ""
			} else {
				idc.Comment = m["Comment"].(string)
			}
			idc.CreateTime = m["CreateTime"].(time.Time)
			idcArr = append(idcArr,idc)
		}
		c.Data["json"] = idcArr
		c.ServeJSON()
	}
}

func (c *IdcController) DoIdc(){
	var message Message
	idc_id,err := c.GetInt("idc_id")
	Name := strings.Trim(c.GetString("Name","")," ")
	Bandwidth := strings.Trim(c.GetString("Bandwidth","")," ")
	Contacts := strings.Trim(c.GetString("Contacts","")," ")
	Mobile := strings.Trim(c.GetString("Mobile","")," ")
	Address := strings.Trim(c.GetString("Address","")," ")
	Network := strings.Trim(c.GetString("Network","")," ")
	Operator := strings.Trim(c.GetString("Operator","")," ")
	Comment := strings.Trim(c.GetString("Comment","")," ")
	FormType := strings.Trim(c.GetString("FormType","")," ")
	IdcNameStatus := strings.Trim(c.GetString("IdcNameStatus","")," ")
	if FormType == "modify_idc" && err == nil{
		IsExist,err := models.CheckIdcExist(Name,"modify",IdcNameStatus)
		if IsExist && err == nil {
			IsModify,err := models.ModifyIdc(idc_id,Name,Bandwidth,Contacts,Mobile,Address,Network,Operator,Comment)
			if IsModify && err == nil {
				message.Status = "ok"
				message.Msg = ""
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Modify idc error:%v",err)
			}
		} else {
			message.Status = "fail"
			message.Msg = "IDC名称已经存在！"
		}
	}
	if FormType == "" && Name != "" && Bandwidth != "" && Contacts != "" && Mobile != ""{
		IsExist,err := models.CheckIdcExist(Name,"add","")
		if IsExist && err == nil {
			IsAdd, err := models.AddIdc(Name, Bandwidth, Contacts, Mobile, Address, Network, Operator, Comment)
			if IsAdd && err == nil {
				message.Status = "ok"
				message.Msg = ""
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Add idc error:%v", err)
			}
		} else {
			message.Status = "fail"
			message.Msg = "IDC名称已经存在！"
		}
	}
	if FormType == "" && err == nil {
		IsDelete,err := models.DeleteIdc(idc_id)
		if IsDelete && err == nil {
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Delete idc error:%v",err)
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}
