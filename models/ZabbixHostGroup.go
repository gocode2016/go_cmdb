package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixHostGroup struct {
	Id int64
	GroupName string `orm:"size(64)"`
	GroupId int64
	Type string `orm:"size(16)"`
}

func GetHostGroupList(page int,keyword,environment string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_host_group")
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

func GetTotalHostGroups()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_host_group").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func AddHostGroup(groupName,Type string,groupId int64)(bool,error){
	var zabbixHostGroup ZabbixHostGroup
	zabbixHostGroup.GroupName = groupName
	zabbixHostGroup.GroupId = groupId
	zabbixHostGroup.Type = Type
	_,err := orm.NewOrm().Insert(&zabbixHostGroup)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteHostGroup(Id int64)(bool,error){
	if _,err := orm.NewOrm().Delete(&ZabbixHostGroup{Id:Id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateHostGroup(groupId, groupName, environment string)(bool) {
	num,err := orm.NewOrm().QueryTable("zabbix_host_group").Filter("GroupId",groupId).Filter("type",environment).Update(orm.Params{"GroupName": groupName})
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}