package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixHost struct {
	Id int64
	HostName string `orm:"size(64)"`
	VisibleName string `orm:"size(64)"`
	HostId int64
	HostStatus bool `orm:"default(true)"`
	Ip string `orm:"size(32);null"`
	Dns string `orm:"size(64);null"`
	InterfaceType string `orm:"size(32);null"`
	InterfaceAddr string `orm:"size(16);null"`
	GroupName string `orm:"size(64);null"`
	GroupId string `orm:"size(32);null"`
	TemplateName string `orm:"size(256);null"`
	TemplateId string `orm:"size(64);null"`
	Type string `orm:"size(16)"`   //国内还是国外
}

func GetHostList(page int,keyword string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_host")
	number,_ := beego.AppConfig.Int("number")
	if keyword != "" {
		cond1 := cond.Or("HostName__icontains",keyword).Or("Ip__icontains",keyword).Or("GroupName__icontains",keyword).
			Or("TemplateName__icontains",keyword).Or("Type__icontains",keyword)
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

func GetTotalHosts()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_host").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func AddHost(HostName,visibleName,interfaceAddr,interfaceValue,interfaceType,GroupName,GroupId,TemplateName,TemplateId,Type string,HostId int64)(bool,error){
	var zabbixHost ZabbixHost
	zabbixHost.HostName = HostName
	zabbixHost.VisibleName = visibleName
	zabbixHost.HostId = HostId
	if interfaceAddr == "IP" {
		zabbixHost.Ip = interfaceValue
	} else {
		zabbixHost.Dns = interfaceValue
	}
	if interfaceType == "1" {
		zabbixHost.InterfaceType = "Agent"
	} else if interfaceType == "2" {
		zabbixHost.InterfaceType = "SNMP"
	} else if interfaceType == "3" {
		zabbixHost.InterfaceType = "IPMI"
	} else if interfaceType == "4" {
		zabbixHost.InterfaceType = "JMX"
	}
	zabbixHost.GroupName = GroupName
	zabbixHost.GroupId = GroupId
	zabbixHost.TemplateName = TemplateName
	zabbixHost.TemplateId = TemplateId
	zabbixHost.Type = Type
	_,err := orm.NewOrm().Insert(&zabbixHost)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteHost(Id int64)(bool,error){
	if _,err := orm.NewOrm().Delete(&ZabbixHost{Id:Id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateHostStatus(HostId int64,HostStatus,Type string)(bool) {
	var host_status = 0
	if HostStatus == "false" {
		host_status = 1
	} else {
		host_status = 0
	}
	num,err := orm.NewOrm().QueryTable("zabbix_host").Filter("HostId",HostId).Filter("type",Type).Update(orm.Params{"HostStatus": host_status})
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}

func UpdateHost(HostId int64,updateHostGroups,updateTemplates,allocateTemplateNames,allocateGroupNames,environment,hostname,visibleName string)(bool) {
	params := make(orm.Params)
	if updateHostGroups != "" && updateTemplates != ""{
		params["group_id"] = updateHostGroups
		params["template_id"] = updateTemplates
		params["group_name"] = allocateGroupNames
		params["template_name"] = allocateTemplateNames
	}
	if updateHostGroups != "" && updateTemplates == ""{
		params["group_id"] = updateHostGroups
		params["group_name"] = allocateGroupNames
	}
	if updateHostGroups == "" && updateTemplates != ""{
		params["template_id"] = updateTemplates
		params["template_name"] = allocateTemplateNames
	}
	params["host_name"] = hostname
	params["visible_name"] = visibleName
	num,err := orm.NewOrm().QueryTable("zabbix_host").Filter("HostId",HostId).Filter("type",environment).Update(params)
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}

func MassDeleteForeignHost(foreignHostIdArr []string)(bool) {
	count,err := orm.NewOrm().QueryTable("zabbix_host").Filter("Type","foreign").Filter("HostId__in",foreignHostIdArr).Delete()
	if count > 0 && err == nil {
		return true
	} else {
		return false
	}
}

func MassDeleteCnHost(cnHostIdArr []string)(bool) {
	count,err := orm.NewOrm().QueryTable("zabbix_host").Filter("Type","cn").Filter("HostId__in",cnHostIdArr).Delete()
	if count > 0 && err == nil {
		return true
	} else {
		return false
	}
}

