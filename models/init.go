package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init(){
	//数据库类别
	dbType := beego.AppConfig.String("db_type")
	//连接名称
	dbAlias := beego.AppConfig.String(dbType + "::db_alias")
	//数据库名称
	dbName := beego.AppConfig.String(dbType + "::db_name")
	//数据库连接用户名
	dbUser := beego.AppConfig.String(dbType + "::db_user")
	//数据库连接用户名
	dbPwd := beego.AppConfig.String(dbType + "::db_pwd")
	//数据库IP（域名）
	dbHost := beego.AppConfig.String(dbType + "::db_host")
	//数据库端口
	dbPort := beego.AppConfig.String(dbType + "::db_port")

	dbCharset := beego.AppConfig.String(dbType + "::db_charset")
	orm.RegisterDataBase(dbAlias, dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+ dbPort+")/"+dbName+"?charset="+dbCharset+"&loc=Asia%2FShanghai", 30)
	orm.RegisterModel(new(Users),new(Roles),new(Asset),new(AssetGroup),new(Idc),new(Log),new(Permission),new(ZabbixHost),new(ZabbixHostGroup),new(ZabbixMedia),new(ZabbixTemplate),new(ZabbixUser),new(ZabbixUserGroup))
	//如果是开发模式，则显示命令信息
	isDev := (beego.AppConfig.String("runmode") == "dev")
	beego.BConfig.WebConfig.Session.SessionOn = true
	//自动建表
	//orm.RunSyncdb("default", false, isDev)
	if isDev {
		orm.Debug = isDev
	}

}
