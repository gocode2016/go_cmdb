package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixTemplate struct {
	Id int64
	TemplateName string `orm:"size(64)"`
	TemplateId int64
	Type string `orm:"size(16)"`
}

func GetTemplateList(page int,keyword,environment string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_template")
	number,_ := beego.AppConfig.Int("number")
	if keyword != "" {
		cond1 := cond.Or("TemplateName__icontains",keyword).Or("Type__icontains",keyword)
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

func GetTotalTemplates()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_template").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func AddTemplate(templateName,Type string,templateId int64)(bool,error){
	var zabbixTemplate ZabbixTemplate
	zabbixTemplate.TemplateName = templateName
	zabbixTemplate.TemplateId = templateId
	zabbixTemplate.Type = Type
	_,err := orm.NewOrm().Insert(&zabbixTemplate)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteTemplate(Id int64)(bool,error){
	if _,err := orm.NewOrm().Delete(&ZabbixTemplate{Id:Id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateTemplate(templateId, templateName, environment string)(bool) {
	num,err := orm.NewOrm().QueryTable("zabbix_template").Filter("TemplateId",templateId).Filter("type",environment).Update(orm.Params{"TemplateName": templateName})
	if num > 0 && err == nil {
		return true
	} else {
		return false
	}
}
