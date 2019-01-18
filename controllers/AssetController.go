package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"strings"
	"fmt"
	"time"
	"encoding/json"
	"github.com/astaxie/beego/orm"
)

type AssetController struct {
	beego.Controller
}

type AssetMap map[string] interface{}
type ModifyAssetMap map[string]AssetMap

type Asset struct {
	Id int64
	Hostname string
	InnerIp string
	OuterIp string
	Os string
	OsRelease string
	Cpu string
	Memory string
	Disk string
	Sn string
	AssetType string
	AssetEnvironment string
	Manufacturer string
	Model string
	AssetNumber string
	CabinetNumber string
	CabinetPosition string
	IsPublisher bool
	AssetStatus string
	Applications string
	Comment string
	OpManager string
	IdcName string
	AssetGroupName string
	JdkVersion string
	TomcatVersion string
	Count int64
	ExpireTime string
	CreateTime string
	UpdateTime string
}
var AssetNameList = [...]string{"Hostname","InnerIp","OuterIp","Os","OsRelease","Cpu","Memory","Disk","AssetType","AssetEnvironment","CabinetNumber","CabinetPosition","IsPublisher","AssetStatus","Applications","Comment","OpManager","ExpireTime","JdkVersion"}

func (c *AssetController) Index() {
	//username := c.Ctx.GetCookie("username")
	username := c.GetSession("name")
	if username != nil {
		c.Data["user"], _ = models.GetUserInfo(username.(string))
		permissionMap := models.GetPermissionByUsername(username.(string))
		permissionJson, _ := json.Marshal(permissionMap)
		c.Data["permissionMap"] = fmt.Sprintf("%s", string(permissionJson))
		c.Data["AssetMenu"] = true
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["headcssjs"] = "asset/index_headcssjs.html"
		c.LayoutSections["footerjs"] = "asset/index_footerjs.html"
		c.Layout = "shared/layout_page.html"
		c.Data["number"], _ = beego.AppConfig.Int("number")
		c.Data["total"] = models.GetTotalAssets()
		c.TplName = "asset/index.html"
	} else {
		//c.Redirect(c.URLFor("HomeController.Login"),302)
		c.Redirect("http://aiops.jollycorp.net/",302)
	}
}

func (c *AssetController) DoAsset(){
	var message Message
	var checkGroup bool = false
	var content string
	var Ip string
	data := c.GetString("data","")
	asset_id,err := c.GetInt("id")
	ip := c.GetString("ip","")
	username := c.GetSession("name")
	if data == "" && err == nil {
		asset,_ := models.GetAsset(asset_id)
		if IsDelete,err := models.DeleteAsset(asset_id);IsDelete {
			if asset.OuterIp != "" {
				Ip = asset.OuterIp
			} else {
				Ip = asset.InnerIp
			}
			content = fmt.Sprintf("Type:delete,Hostname:%s,InnerIp:%s",asset.Hostname,asset.InnerIp)
			models.AddLog(username.(string),Ip,content)
			message.Msg = ""
			message.Status = "ok"
		} else {
			message.Msg = fmt.Sprintf("Delete asset error:%v", err)
			message.Status = "fail"
		}
		c.Data["json"] = &message
		c.ServeJSON()
	}
	if data != "" && err == nil {
		var modify_asset_list []ModifyAssetMap
		var message Message
		var UpdateDataMaps= make(orm.Params)
		var isContinue = true
		json.Unmarshal([]byte(data), &modify_asset_list)
		for i := 0; i < len(modify_asset_list); i++ {
			if _, ok := modify_asset_list[i]["Hostname"]; ok {
				isExist, _ := models.CheckAssetsExist(modify_asset_list[i]["Hostname"]["newValue"].(string),nil,nil)
				if isExist {
					isContinue = false
					message.Msg = "主机名已存在"
					message.Status = "fail"
					c.Data["json"] = &message
					c.ServeJSON()
				}
			}
			if _, ok := modify_asset_list[i]["OuterIp"]; ok {
				isExist, _ := models.CheckAssetsExist(nil,nil,modify_asset_list[i]["OuterIp"]["newValue"].(string))
				if isExist {
					isContinue = false
					message.Msg = "外网IP已存在"
					message.Status = "fail"
					c.Data["json"] = &message
					c.ServeJSON()
				}
			}
		}
		if isContinue {
			for i := 0; i < len(modify_asset_list); i++ {
				if _, ok := modify_asset_list[i]["AssetGroupName"]; ok {
					oldVal := modify_asset_list[i]["AssetGroupName"]["oldValue"]
					newVal := modify_asset_list[i]["AssetGroupName"]["newValue"]
					models.UpdateAssetGroup(asset_id, oldVal, newVal)
					checkGroup = true
					content = fmt.Sprintf("Field:AssetGroupName,oldValue:%s,newValue:%s;", oldVal, newVal)
					if len(modify_asset_list) == 1 {
						content = "Type:update," + content
						models.AddLog(username.(string), ip, content)
						message.Msg = ""
						message.Status = "ok"
						c.Data["json"] = &message
						c.ServeJSON()
					}
				}
			}
			for i := 0; i < len(modify_asset_list); i++ {
				if _, ok := modify_asset_list[i]["IdcName"]; ok {
					oldVal := modify_asset_list[i]["IdcName"]["oldValue"]
					newVal := modify_asset_list[i]["IdcName"]["newValue"]
					models.UpdateIdc(asset_id, newVal.(string))
					if content == "" {
						content = fmt.Sprintf("Field:idc_name,oldValue:%s,newValue:%s;", oldVal, newVal)
					} else {
						content = fmt.Sprintf("%s,Field:idc_name,oldValue:%s,newValue:%s;", content, oldVal, newVal)
					}
					if (checkGroup && len(modify_asset_list) == 2) || (len(modify_asset_list) == 1) {
						if ! strings.Contains(content, "Type:update") {
							content = "Type:update," + content
						}
						models.AddLog(username.(string), ip, content)
						message.Msg = ""
						message.Status = "ok"
						c.Data["json"] = &message
						c.ServeJSON()
					}
				}
			}
			for i := 0; i < len(modify_asset_list); i++ {
				for j := 0; j < len(AssetNameList); j++ {
					if _, ok := modify_asset_list[i][AssetNameList[j]]; ok {
						if modify_asset_list[i][AssetNameList[j]]["newValue"] != modify_asset_list[i][AssetNameList[j]]["oldValue"] {
							UpdateDataMaps[AssetNameList[j]] = modify_asset_list[i][AssetNameList[j]]["newValue"]
							if content == "" {
								content = fmt.Sprintf("Type:update,Field:%s,oldValue:%v,newValue:%v;", AssetNameList[j], modify_asset_list[i][AssetNameList[j]]["oldValue"], modify_asset_list[i][AssetNameList[j]]["newValue"])
							} else {
								content = fmt.Sprintf("%s;Field:%s,oldValue:%v,newValue:%v;", content, AssetNameList[j], modify_asset_list[i][AssetNameList[j]]["oldValue"], modify_asset_list[i][AssetNameList[j]]["newValue"])
							}
						}
					}
				}
			}
			if ok, err := models.MultiUpdate(asset_id, UpdateDataMaps); ok {
				message.Msg = ""
				message.Status = "ok"
			} else {
				message.Msg = fmt.Sprintf("Update asset error:%v", err)
				message.Status = "ok"
			}
			if ! strings.Contains(content, "Type:update") {
				content = "Type:update," + content
			}
			models.AddLog(username.(string), ip, content)
			c.Data["json"] = &message
			c.ServeJSON()
		}
	}
	if data == "" && err != nil {
		asset_map := make(map[string]interface{})
		asset_map["hostname"] = strings.Trim(c.GetString("hostname"), " ")
		asset_map["inner_ip"] = strings.Trim(c.GetString("inner_ip"), " ")
		asset_map["outer_ip"] = strings.Trim(c.GetString("outer_ip"), " ")
		isExist, assetName := models.CheckAssetsExist(asset_map["hostname"], asset_map["inner_ip"], asset_map["outer_ip"])
		if isExist {
			var message Message
			message.Msg = fmt.Sprintf("%s已存在", assetName)
			message.Status = "fail"
			c.Data["json"] = &message
			c.ServeJSON()
		}
		asset_map["os"] = strings.Trim(c.GetString("os"), " ")
		asset_map["os_release"] = strings.Trim(c.GetString("os_release"), " ")
		asset_map["cpu"] = strings.Trim(c.GetString("cpu", ""), " ")
		asset_map["memory"] = strings.Trim(c.GetString("memory", ""), " ")
		asset_map["disk"] = strings.Trim(c.GetString("disk", ""), " ")
		asset_map["sn"] = strings.Trim(c.GetString("sn", ""), " ")
		asset_map["asset_type"] = strings.Trim(c.GetString("asset_type"), " ")
		asset_map["manufacturer"] = strings.Trim(c.GetString("manufacturer", ""), " ")
		asset_map["model"] = strings.Trim(c.GetString("model", ""), " ")
		asset_map["asset_number"] = strings.Trim(c.GetString("asset_number", ""), " ")
		asset_map["cabinet_number"] = strings.Trim(c.GetString("cabinet_number", ""), " ")
		asset_map["cabinet_position"] = strings.Trim(c.GetString("cabinet_position", ""), " ")
		asset_map["is_publisher"], _ = c.GetInt("is_publisher")
		asset_map["asset_environment"] = strings.Trim(c.GetString("asset_environment", ""), " ")
		asset_map["asset_status"] = strings.Trim(c.GetString("asset_status"), " ")
		asset_map["applications"] = strings.Trim(c.GetString("applications", ""), " ")
		asset_group_list := c.GetStrings("asset_group")
		asset_map["idc"] = strings.Trim(c.GetString("idc"), " ")
		asset_map["expire_time"] = c.GetString("expire_time", "")
		asset_map["comment"] = strings.Trim(c.GetString("comment", ""), " ")
		asset_map["op_manager"] = strings.Trim(c.GetString("op_manager", ""), " ")
		asset_map["jdk_version"] = strings.Trim(c.GetString("jdk_version", ""), " ")
		id, isAdd, err := models.AddAsset(asset_map)
		if isAdd && err == nil {
			if asset_group_list != nil {
				models.AddGroupForAsset(id, asset_group_list)
			}
			content = fmt.Sprintf("Type:Add,Hostname:%s,InnerIp:%s",asset_map["hostname"],asset_map["inner_ip"])
			if asset_map["outer_ip"].(string) != "" {
				Ip = asset_map["outer_ip"].(string)
			} else {
				Ip = asset_map["inner_ip"].(string)
			}
			models.AddLog(username.(string),Ip,content)
			message.Msg = ""
			message.Status = "ok"
		} else {
			message.Msg = fmt.Sprintf("添加失败:%v", err.Error())
			message.Status = "fail"
		}
		c.Data["json"] = &message
		c.ServeJSON()
	}
}

func (c *AssetController) GetAssetList() {
	var assetArr []Asset
	var asset Asset
	idc_id := c.GetString("idc","")
	asset_status := c.GetString("asset_status","")
	keyword := c.GetString("keyword","")
	page,err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	assetMaps,count,err := models.GetAssetList(page,idc_id,asset_status,keyword)
	if assetMaps != nil && err == nil {
		for _,m := range assetMaps {
			asset.Id = m["Id"].(int64)
			asset.Hostname = m["Hostname"].(string)
			asset.InnerIp = m["InnerIp"].(string)
			asset.OuterIp = m["OuterIp"].(string)
			asset.Os = m["Os"].(string)
			asset.OsRelease = m["OsRelease"].(string)
			asset.Cpu = m["Cpu"].(string)
			asset.Memory = m["Memory"].(string)
			asset.Disk = m["Disk"].(string)
			if m["Sn"] != nil {
				asset.Sn = m["Sn"].(string)
			} else {
				asset.Sn = ""
			}
			asset.AssetType = m["AssetType"].(string)
			asset.AssetEnvironment = m["AssetEnvironment"].(string)
			if m["Manufacturer"] != nil {
				asset.Manufacturer = m["Manufacturer"].(string)
			} else {
				asset.Manufacturer = ""
			}
			if m["Model"] != nil {
				asset.Model = m["Model"].(string)
			} else {
				asset.Model = ""
			}
			if m["AssetNumber"] != nil {
				asset.AssetNumber = m["AssetNumber"].(string)
			} else {
				asset.AssetNumber = ""
			}
			if m["CabinetNumber"] != nil {
				asset.CabinetNumber = m["CabinetNumber"].(string)
			} else {
				asset.CabinetNumber = ""
			}
			if m["CabinetPosition"] != nil {
				asset.CabinetPosition = m["CabinetPosition"].(string)
			} else {
				asset.CabinetPosition = ""
			}
			asset.IsPublisher,_ = m["IsPublisher"].(bool)
			asset.AssetStatus = m["AssetStatus"].(string)
			if m["Applications"] != nil {
				asset.Applications = m["Applications"].(string)
			} else {
				asset.Applications = ""
			}
			if m["Comment"] != nil {
				asset.Comment = m["Comment"].(string)
			} else {
				asset.Comment = ""
			}
			if m["OpManager"] != nil {
				asset.OpManager = m["OpManager"].(string)
			} else {
				asset.OpManager = ""
			}
			if m["JdkVersion"] != nil {
				asset.JdkVersion = m["JdkVersion"].(string)
			} else {
				asset.JdkVersion = ""
			}
			asset.IdcName = models.GetIdcName(m["Id"].(int64))
			groups,err := models.GetGroupNameByAsset(m["Id"].(int64))
			group_name := ""
			if groups != nil && err == nil {
				for _, v := range groups {
					group_name += "," + v.GroupName
				}
				group_name = strings.TrimLeft(group_name,",")
			}
			asset.AssetGroupName = group_name
			if m["ExpireTime"] != nil {
				asset.ExpireTime = m["ExpireTime"].(string)
			} else {
				asset.ExpireTime = ""
			}
			asset.CreateTime = m["CreateTime"].(time.Time).Format("2006-01-02 15:04:05")
			asset.UpdateTime = m["UpdateTime"].(time.Time).Format("2006-01-02 15:04:05")
			asset.Count = count
			assetArr = append(assetArr,asset)
		}
		c.Data["json"] = &assetArr
	} else {
		var message Message
		message.Status = "fail"
		message.Msg = fmt.Sprintf("获取失败:%v",err)
		c.Data["json"] = &message
	}
	c.ServeJSON()
}

func (c *AssetController) GetExpiredAssetList(){
	var ExpiredAssetList []string
	IsExist,assets := models.GetExpiredAsset()
	if IsExist {
		for _,asset := range assets {
			if asset.OuterIp != "" {
				ExpiredAssetList = append(ExpiredAssetList,asset.OuterIp)
			} else {
				ExpiredAssetList = append(ExpiredAssetList,asset.InnerIp)
			}
		}
	}
	c.Data["json"] = ExpiredAssetList
	c.ServeJSON()
}
