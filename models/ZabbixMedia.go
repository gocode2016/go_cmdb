package models

import (
	"strings"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type ZabbixMedia struct {
	Id int64
	MediaName string `orm:"size(32)"`
	MediaTypeId int64
	MediaStatus bool `orm:"default(true)"`
	Description string `orm:"size(128);null"`
	Type string `orm:"size(16)"`
}

func GetMediaList(page int,keyword,environment string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("zabbix_media")
	number,_ := beego.AppConfig.Int("number")
	if keyword != "" {
		cond1 := cond.Or("MediaName__icontains",keyword).Or("Type__icontains",keyword)
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

func GetMediaNameByMediaTypeId(mediaTypeId,environment string)(string){
	var mediaLists []orm.ParamsList
	var mediaName string
	qs := orm.NewOrm().QueryTable("zabbix_media")
	if strings.Contains(mediaTypeId, ",") {
		mediaTypeIdArr := strings.Split(mediaTypeId,",")
		num, err := qs.Filter("MediaTypeId__in", mediaTypeIdArr).Filter("Type", environment).ValuesList(&mediaLists,"MediaName")
		if num > 0 && err == nil {
			for _,row := range mediaLists {
				if mediaName == "" {
					mediaName = row[0].(string)
				} else {
					mediaName = mediaName+","+row[0].(string)
				}
			}
		}
	} else {
		num, err := qs.Filter("MediaTypeId", mediaTypeId).Filter("Type", environment).ValuesList(&mediaLists,"MediaName")
		if num > 0 && err == nil {
			for _,row := range mediaLists {
				mediaName = row[0].(string)
			}
		}
	}
	return mediaName
}

func GetTotalZabbixMedias()(int64){
	count,err := orm.NewOrm().QueryTable("zabbix_media").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}