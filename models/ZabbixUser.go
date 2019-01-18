package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixUser struct {
	Id int64
	UserName string `orm:"size(32)"`
	UserId int64
	GroupName string `orm:"size(64);null"`
	GroupId string `orm:"size(16);null"`
	MediaTypeId string `orm:"size(16);null"`
	SendtoPeriod string `orm:"size(256);null"`
	Type string `orm:"size(16)"`
}

func GetUserList(page int,keyword string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_user")
	number,_ := beego.AppConfig.Int("number")
	if keyword != "" {
		cond1 := cond.Or("UserName__icontains",keyword).Or("GroupName__icontains",keyword).Or("SendtoPeriod__icontains",keyword).Or("Type__icontains",keyword)
		qs = qs.SetCond(cond1)
	}
	count,_ = qs.Count()
	offset = (page-1) * number
	_,err := qs.Limit(number,offset).Values(&maps)
	if err == nil {
		return maps,count,nil
	}
	return maps,count,err
}

func GetTotalZabbixUsers()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_user").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func AddZabbixUser(username,groupName,groupId,sendto,period,Type string,userId int64)(bool,error){
	var zabbixUser ZabbixUser
	var sendtoPeriod string
	zabbixUser.UserName = username
	zabbixUser.UserId = userId
	zabbixUser.GroupName = groupName
	zabbixUser.GroupId = groupId
	if sendto != "" {
		sendtoPeriod = sendto+":"+period
		if Type == "cn" {
			zabbixUser.MediaTypeId = beego.AppConfig.String("media::cn_media_type_id")
		} else if Type == "foreign" {
			zabbixUser.MediaTypeId = beego.AppConfig.String("media::foreign_media_type_id")
		}
	}
	zabbixUser.Type = Type
	zabbixUser.SendtoPeriod = sendtoPeriod
	_,err := orm.NewOrm().Insert(&zabbixUser)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteZabbixUser(Id int64)(bool,error){
	if _,err := orm.NewOrm().Delete(&ZabbixUser{Id:Id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateZabbixUser(userId, groupId, groupName, environment, sendto,period string)(bool) {
	var sendtoPeriod = ""
	var num1,num2,num3 int64
	var err1,err2,err3 error
	var mediaTypeId string
	if sendto != "" {
		if environment == "cn" {
			mediaTypeId = beego.AppConfig.String("media::cn_media_type_id")
		} else {
			mediaTypeId = beego.AppConfig.String("media::foreign_media_type_id")
		}
		sendtoPeriod = sendto+":"+period
	}
	if groupId != "" && groupName != "" && sendto == "" {
		num1, err1 = orm.NewOrm().QueryTable("zabbix_user").Filter("UserId", userId).Filter("type", environment).Update(orm.Params{"GroupName": groupName,
			"GroupId": groupId})
	} else if groupId == "" && groupName == "" && sendto != "" {
		num2, err2 = orm.NewOrm().QueryTable("zabbix_user").Filter("UserId", userId).Filter("type", environment).Update(orm.Params{"sendtoPeriod": sendtoPeriod,
			"MediaTypeId": mediaTypeId})
	} else if groupId != "" && groupName != "" && sendto != "" {
		num3, err3 = orm.NewOrm().QueryTable("zabbix_user").Filter("UserId", userId).Filter("type", environment).Update(orm.Params{"GroupName": groupName,
			"GroupId": groupId, "sendtoPeriod": sendtoPeriod, "MediaTypeId": mediaTypeId})
	}
	if (num1 > 0 && err1 == nil) || (num2 > 0 && err2 == nil) || (num3 >0 && err3 == nil) {
		return true
	} else {
		return false
	}
}