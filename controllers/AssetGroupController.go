package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"strings"
	"fmt"
	"encoding/json"
)

type AssetGroupController struct {
	beego.Controller
}

type AssetGroup struct {
	Id int64
	GroupName string
	AliasName string
	DevManager string
	OpManager string
	Comment string
	Count int64
}

type AssetInfo struct {
	IP string
	Hostname string
}

func (c *AssetGroupController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil{
		c.Data["user"],_ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson,_ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s",string(permissionJson))
		c.Data["AssetGroupMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "assetGroup/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "assetGroup/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"],_ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalAssetGroups()
		c.TplName = "assetGroup/index.html"
	} else {
		c.Redirect(c.URLFor("HomeController.Login"),302)
		//c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *AssetGroupController) GetAssetGroupList(){
	var assetGroupArr []AssetGroup
	var assetGroup AssetGroup
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	selectType := c.GetString("select_type")
	if err != nil {
		page = 1
	}
	assetGroupMaps,count,err := models.GetAssetGroupList(page,keyword,selectType)
	if err == nil {
		for _,m := range assetGroupMaps {
			assetGroup.Id = m["Id"].(int64)
			assetGroup.GroupName = m["GroupName"].(string)
			if m["AliasName"] == nil {
				assetGroup.AliasName = ""
			} else {
				assetGroup.AliasName = m["AliasName"].(string)
			}
			if m["DevManager"] == nil {
				assetGroup.DevManager = ""
			} else {
				assetGroup.DevManager = m["DevManager"].(string)
			}
			if m["OpManager"] == nil {
				m["OpManager"] = ""
			} else {
				assetGroup.OpManager = m["OpManager"].(string)
			}
			if m["Comment"] == nil {
				m["Comment"] = ""
			} else {
				assetGroup.Comment = m["Comment"].(string)
			}
			assetGroup.Count = count
			assetGroupArr = append(assetGroupArr,assetGroup)
		}
		c.Data["json"] = assetGroupArr
		c.ServeJSON()
	}
}

func (c *AssetGroupController) DoAssetGroup(){
	var message Message
	var asset_list []AssetInfo
	var assetInfo AssetInfo
	group_id,err := c.GetInt("group_id")
	target := c.GetString("target")
	GroupName := strings.Trim(c.GetString("GroupName","")," ")
	AliasName := strings.Trim(c.GetString("AliasName","")," ")
	DevManager := strings.Trim(c.GetString("DevManager","")," ")
	OpManager := strings.Trim(c.GetString("OpManager","")," ")
	Comment := strings.Trim(c.GetString("Comment","")," ")
	FormType := strings.Trim(c.GetString("FormType","")," ")
	GroupNameStatus := strings.Trim(c.GetString("GroupNameStatus","")," ")
	if err == nil && target == "get_assets" {
		assets,err := models.GetAssetsByGroup(int64(group_id))
		if err == nil {
			for _,asset := range assets {
				assetInfo.Hostname = asset.Hostname
				if asset.OuterIp != "" {
					assetInfo.IP = asset.OuterIp
				} else {
					assetInfo.IP = asset.InnerIp
				}
				asset_list = append(asset_list,assetInfo)
			}
			c.Data["json"] = asset_list
		}
		c.ServeJSON()
	}
	if GroupName != "" && FormType == "modify_asset_group" && err == nil && target == "" && GroupNameStatus != "" {
		IsExist,err := models.CheckAssetGroupExist(GroupName,"modify",GroupNameStatus)
		if IsExist && err == nil {
			IsModify,err := models.ModifyAssetGroup(group_id,GroupName,AliasName,DevManager,OpManager,Comment)
			if IsModify && err == nil {
				message.Status = "ok"
				message.Msg = ""
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Modify asset error:%v",err)
			}
		} else {
			message.Status = "fail"
			message.Msg = "资产组名称已存在！"
		}
	}
	if GroupName != "" && FormType == "" && err != nil && target == "" {
		IsExist,err := models.CheckAssetGroupExist(GroupName,"add","")
		if IsExist && err == nil {
			IsAdd, err := models.AddAssetGroup(GroupName, AliasName, DevManager, OpManager, Comment)
			if IsAdd && err == nil {
				message.Status = "ok"
				message.Msg = ""
			} else {
				message.Status = "fail"
				message.Msg = fmt.Sprintf("Add asset error:%v", err)
			}
		} else {
			message.Status = "fail"
			message.Msg = "资产组名称已存在！"
		}
	}
	if GroupName == "" && FormType == "" && err == nil && target == ""{
		isDelete,err := models.DeleteAssetGroup(group_id)
		if isDelete && err == nil {
			message.Status = "ok"
			message.Msg = ""
		} else {
			message.Status = "fail"
			message.Msg = fmt.Sprintf("Delete asset error:%v",err)
		}
	}
	c.Data["json"] = &message
	c.ServeJSON()
}
