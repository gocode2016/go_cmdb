package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixUserGroup struct {
	Id int64
	GroupName string `orm:"size(32)"`
	GroupId int64
	GroupStatus bool `orm:"default(true)"`
	Type string `orm:"size(16)"`
}

func GetUserGroupList(page int,keyword,environment string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_user_group").Filter("GroupStatus", 0)
	number,_ := beego.AppConfig.Int("number")
	if keyword != "" {
		cond1 := cond.Or("GroupName__icontains",keyword).Or("Type__icontains",keyword)
		qs = qs.SetCond(cond1)
	}
	if environment != "" {
		qs = qs.Filter("Type",environment)
		_,err := qs.Values(&maps)
		if err == nil {
			return maps,count,nil
		}
		return maps,count,err
	}
	count,_ = qs.Count()
	offset = (page-1) * number
	_,err := qs.Limit(number,offset).Values(&maps)
	if err == nil {
		return maps,count,nil
	}
	return maps,count,err
}

func GetTotalZabbixUserGroups()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_user_group").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func AddUserGroup(groupName,Type string,groupId int64)(bool,error){
	var zabbixUserGroup ZabbixUserGroup
	zabbixUserGroup.GroupName = groupName
	zabbixUserGroup.GroupId = groupId
	zabbixUserGroup.GroupStatus = false
	zabbixUserGroup.Type = Type
	_,err := orm.NewOrm().Insert(&zabbixUserGroup)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteUserGroup(Id int64)(bool,error){
	if _,err := orm.NewOrm().Delete(&ZabbixUserGroup{Id:Id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateUserGroupStatus(groupId,groupStatus,Type string)(bool) {
	var status string
	if groupStatus == "false" {
		status = "1"
	} else {
		status = "0"
	}
	num,err := orm.NewOrm().QueryTable("zabbix_user_group").Filter("GroupId",groupId).Filter("type",Type).Update(orm.Params{"GroupStatus": status})
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}

func UpdateUserGroup(groupId, groupName, environment string)(bool) {
	num,err := orm.NewOrm().QueryTable("zabbix_user_group").Filter("GroupId",groupId).Filter("type",environment).Update(orm.Params{"GroupName": groupName})
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}