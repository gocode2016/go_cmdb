package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"strings"
)

type Idc struct {
	Id int
	Name string `orm:"size(32)"`
	Bandwidth string `orm:"size(32);null"`  //带宽
	Contacts string `orm:"size(16);null"`
	Mobile string `orm:"size(16);null"`
	Address string `orm:"size(64);null"`
	Network string `orm:"size(64);null"`   //IP地址段
	Operator string `orm:"size(16);null"`  //运营商
	Comment string `orm:"size(256);null"`  //备注
	Asset []*Asset `orm:"reverse(many)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func GetIdcList(keyword string)([]orm.Params,error){
	var maps []orm.Params
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("idc")
	if keyword != "" {
		cond1 := cond.Or("Name__icontains",keyword).Or("Bandwidth__icontains",keyword).Or("Contacts__icontains",keyword).Or("Mobile__icontains",keyword)
		qs = qs.SetCond(cond1)
	}
	_,err := qs.Values(&maps)
	if err == nil {
		return maps,nil
	}
	return maps,err
}

func GetIdcName(id int64)(string){
	o := orm.NewOrm()
	var idc Idc
	err := o.QueryTable("idc").Filter("Asset__Id",id).Limit(1).One(&idc)
	if err == nil {
		return idc.Name
	} else {
		return ""
	}
}

func AddIdc(Name,Bandwidth,Contacts,Mobile,Address,Network,Operator,Comment string)(bool,error){
	var idc Idc
	idc.Name = Name
	idc.Bandwidth = Bandwidth
	idc.Contacts = Contacts
	idc.Mobile = Mobile
	idc.Address = Address
	idc.Network = Network
	idc.Operator = Operator
	idc.Comment = Comment
	_,err := orm.NewOrm().Insert(&idc)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteIdc(id int)(bool,error){
	if _,err := orm.NewOrm().Delete(&Idc{Id:id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func CheckIdcExist(Name,Operation,IdcNameStatus string)(bool,error){
	num,err := orm.NewOrm().QueryTable("idc").Filter("Name",strings.ToLower(Name)).Count()
	if(Operation == "add") {
		if num == 0 && err == nil {
			return true,nil
		} else {
			return false,err
		}
	} else if (Operation == "modify") {
		if IdcNameStatus == "Change" {
			if num == 0 && err == nil {
				return true, nil
			} else {
				return false, err
			}
		} else if IdcNameStatus == "NotChange" {
			if num == 1 && err == nil {
				return true, nil
			} else {
				return false, err
			}
		}
	}
	return false,nil
}

func ModifyIdc(id int,Name,Bandwidth,Contacts,Mobile,Address,Network,Operator,Comment string)(bool,error){
	var idc_map orm.Params
	idc_map = make(orm.Params)
	idc_map["Name"] = Name
	idc_map["Bandwidth"] = Bandwidth
	idc_map["Contacts"] = Contacts
	idc_map["Mobile"] = Mobile
	idc_map["Address"] = Address
	idc_map["Network"] = Network
	idc_map["Operator"] = Operator
	idc_map["Comment"] = Comment
	_,err := orm.NewOrm().QueryTable("idc").Filter("Id",id).Update(idc_map)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}