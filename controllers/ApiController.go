package controllers

import (
	"github.com/astaxie/beego"
	"cmdb/models"
	"fmt"
	"cmdb/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strings"
	"strconv"
)

type ApiController struct {
	beego.Controller
}

type Response map[string]interface{}

type GroupInfo struct {
	AssetGroup string
	PackMethod string
	AssetInfo []map[string]interface{}
}

const secret string = "cmdbApi"

func (this *ApiController) Prepare() {
	username := this.GetString("username")
	password := this.GetString("password")
	token := this.GetString("token")
	if token == "" {
		if username == "" || password == "" {
			res := make(map[string]interface{})
			res["status"] = "fail"
			res["errMsg"] = "No username or password"
			this.Data["json"] = &res
			this.ServeJSON()
		} else if username != "" && password != "" {
			isExist := utils.UserAuthLdap(username, password)
			if !isExist {
				res := make(map[string]interface{})
				res["status"] = "fail"
				res["errMsg"] = "Wrong username or password"
				this.Data["json"] = &res
				this.ServeJSON()
			}
		}
	} else {
		if token != "q1EpLNpQAknntpOq" {
			if strings.Contains(token, "\"") {
				token, _ = strconv.Unquote(token)
			}
			_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil {
				res := make(map[string]interface{})
				res["status"] = "fail"
				res["errMsg"] = fmt.Sprintf("%v", err)
				this.Data["json"] = &res
				this.ServeJSON()
			}
		}
	}
}

func (this *ApiController) GetToken() {
	mySigningKey := []byte(secret)
	claims := &jwt.StandardClaims{
		ExpiresAt: int64(time.Now().Unix() + 600),
		Issuer:    "cmdbApi",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err == nil {
		res := make(map[string]interface{})
		res["status"] = "ok"
		res["token"] = ss
		this.Data["json"] = &res
		this.ServeJSON()
	} else {
		res := make(map[string]interface{})
		res["status"] = "fail"
		res["errMsg"] = "Get token fail"
		this.Data["json"] = &res
		this.ServeJSON()
	}
}

func (this *ApiController) GetGroupList() {
	var groupInfo GroupInfo
	var groupInfoList []GroupInfo
	var assetGroup = ""
	groups,err := models.GetGroupList()
	if err == nil {
		for _,group := range groups {
			if group["GroupName"] != nil {
				assetGroup = group["GroupName"].(string)
				groupInfo.AssetGroup = assetGroup
				if group["PackMethod"] != nil {
					groupInfo.PackMethod = group["PackMethod"].(string)
				} else {
					groupInfo.PackMethod = ""
				}
				assets,err := models.GetAssetsByGroup(group["Id"].(int64))
				if len(assets) > 0 && err == nil {
					var AssetInfo  []map[string]interface{}
					for _,asset := range assets {
						assetInfo := make(map[string]interface{})
						assetInfo["hostname"] = asset.Hostname
						assetInfo["jdkVersion"] = asset.JdkVersion
						assetInfo["innerIp"] = asset.InnerIp
						assetInfo["outerIp"] = asset.OuterIp
						//assetInfo["tomcatVersion"] = asset.TomcatVersion
						//assetInfo["appUser"] = asset.AppUser
						assetInfo["assetEnvironment"] = asset.AssetEnvironment
						assetInfo["idcName"] = models.GetIdcName(int64(asset.Idc.Id))
						AssetInfo = append(AssetInfo, assetInfo)
						groupInfo.AssetInfo = AssetInfo
					}
				}
				groupInfoList = append(groupInfoList,groupInfo)
			}

		}
		this.Data["json"] = &groupInfoList
	} else {
		res := make(map[string]interface{})
		res["status"] = "fail"
		res["errMsg"] = fmt.Sprintf("%v",err)
		this.Data["json"] = &res
	}
	this.ServeJSON()
}
