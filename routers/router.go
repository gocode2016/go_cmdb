package routers

import (
	"github.com/astaxie/beego"
	"cmdb/controllers"
)

func init() {
	//登录
	beego.Router("/",&controllers.HomeController{},"Get:Login")
	beego.Router("/home/dologin", &controllers.HomeController{}, "Post:DoLogin")
	beego.Router("/home/index", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/logout", &controllers.HomeController{}, "*:Logout")
	//用户模块
	beego.Router("/user/index", &controllers.UserController{}, "*:Index")
	beego.Router("/user/list", &controllers.UserController{}, "*:GetUserList")
	beego.Router("/user/douser",&controllers.UserController{},"*:DoUser")
	beego.Router("/role/index",&controllers.RoleController{},"*:Index")
	beego.Router("/role/list", &controllers.RoleController{}, "*:GetRoleList")
	beego.Router("/role/dorole",&controllers.RoleController{},"*:DoRole")
	//资产模块
	beego.Router("/asset/index", &controllers.AssetController{}, "*:Index")
	beego.Router("/asset/doasset", &controllers.AssetController{}, "*:DoAsset")
	beego.Router("/asset/list", &controllers.AssetController{}, "*:GetAssetList")
	beego.Router("/asset/expirelist", &controllers.AssetController{}, "*:GetExpiredAssetList")
	//资产组模块
	beego.Router("/group/index", &controllers.AssetGroupController{}, "*:Index")
	beego.Router("/group/list",&controllers.AssetGroupController{},"*:GetAssetGroupList")
	beego.Router("/group/dogroup",&controllers.AssetGroupController{},"*:DoAssetGroup")
	//IDC模块
	beego.Router("/idc/index", &controllers.IdcController{}, "*:Index")
	beego.Router("/idc/list",&controllers.IdcController{},"*:GetIdcList")
	beego.Router("/idc/doidc",&controllers.IdcController{},"*:DoIdc")
	//log模块
	beego.Router("/log/index", &controllers.LogController{}, "*:Index")
	beego.Router("/log/list",&controllers.LogController{},"*:GetLogList")
	//权限路由
	beego.Router("/permission/index", &controllers.PermissionController{}, "*:Index")
	beego.Router("/permission/list",&controllers.PermissionController{},"*:GetPermissionList")
	beego.Router("/permission/doidc",&controllers.PermissionController{},"*:DoPermission")
	//API路由
	beego.Router("/v1/group/list", &controllers.ApiController{}, "get:GetGroupList")
	beego.Router("/v1/token", &controllers.ApiController{}, "get:GetToken")
	//zabbix主机路由
	beego.Router("/host/index", &controllers.ZabbixHostController{}, "*:Index")
	beego.Router("/host/list", &controllers.ZabbixHostController{}, "*:GetHostList")
	beego.Router("/host/dohost", &controllers.ZabbixHostController{}, "*:DoZabbixHost")
	//zabbix主机组路由
	beego.Router("/hostGroup/index", &controllers.ZabbixHostGroupController{}, "*:Index")
	beego.Router("/hostGroup/list", &controllers.ZabbixHostGroupController{}, "*:GetHostGroupList")
	beego.Router("/hostGroup/dohostgroup", &controllers.ZabbixHostGroupController{}, "*:DoZabbixHostGroup")
	//zabbix模板路由
	beego.Router("/template/index", &controllers.ZabbixTemplateController{}, "*:Index")
	beego.Router("/template/list", &controllers.ZabbixTemplateController{}, "*:GetTemplateList")
	beego.Router("/template/dotemplate", &controllers.ZabbixTemplateController{}, "*:DoZabbixTemplate")
	//zabbix用户路由
	beego.Router("/zabbixUser/index", &controllers.ZabbixUserController{}, "*:Index")
	beego.Router("/zabbixUser/list", &controllers.ZabbixUserController{}, "*:GetUserList")
	beego.Router("/zabbixUser/douser", &controllers.ZabbixUserController{}, "*:DoZabbixUser")
	//zabbix用户组路由
	beego.Router("/userGroup/index", &controllers.ZabbixUserGroupController{}, "*:Index")
	beego.Router("/userGroup/list", &controllers.ZabbixUserGroupController{}, "*:GetUserGroupList")
	beego.Router("/userGroup/dogroup", &controllers.ZabbixUserGroupController{}, "*:DoZabbixUserGroup")
	//zabbix媒介路由
	beego.Router("/media/index", &controllers.ZabbixMediaController{}, "*:Index")
	beego.Router("/media/list", &controllers.ZabbixMediaController{}, "*:GetMediaList")
}
