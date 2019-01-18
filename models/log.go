package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego"
)

type Log struct {
	Id int
	Name string `orm:"size(16)"`
	Ip string `orm:"size(64)"`
	Content string `orm:"size(256)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func AddLog(username,Ip,content string){
	O := orm.NewOrm()
	var log Log
	log.Name = username
	log.Ip = Ip
	log.Content = content
	if _,err := O.Insert(&log); err != nil {
		fmt.Println("Insert log err:",err)
	}
}

func GetLogList(page int,keyword,sort string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	number,_ := beego.AppConfig.Int("number")
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("log")
	if keyword != "" {
		cond1 := cond.Or("Content__icontains",keyword).Or("Name__icontains",keyword).Or("Ip__contains",keyword)
		count,_ = qs.SetCond(cond1).Count()
		qs = qs.SetCond(cond1)
	}
	if sort != "" {
		qs = qs.OrderBy("-create_time")
	} else {
		qs = qs.OrderBy("create_time")
	}
	offset := (page-1) * number
	_,err := qs.Limit(number,offset).Values(&maps)
	if err == nil {
		return maps,count,nil
	}
	return maps,count,err
}

func GetTotalLogs()(int64){
	count,err := orm.NewOrm().QueryTable("log").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func GetStatisticNumber()([]int64,[]int64){
	var AddAssetList,UpdateAssetList []int64
	var TimeList []string
	now := time.Now()
	yesterdayTime := now.AddDate(0,0,-1).Format("2006-01-02")
	beforeYesterdayTime := now.AddDate(0,0,-2).Format("2006-01-02")
	threeDaysAgoTime := now.AddDate(0,0,-3).Format("2006-01-02")
	fourDaysAgoTime := now.AddDate(0,0,-4).Format("2006-01-02")
	fiveDaysAgoTime := now.AddDate(0,0,-5).Format("2006-01-02")
	sixDaysAgoTime := now.AddDate(0,0,-6).Format("2006-01-02")
	sevenDaysAgoTime := now.AddDate(0,0,-7).Format("2006-01-02")
	TimeList = append(TimeList,sevenDaysAgoTime)
	TimeList = append(TimeList,sixDaysAgoTime)
	TimeList = append(TimeList,fiveDaysAgoTime)
	TimeList = append(TimeList,fourDaysAgoTime)
	TimeList = append(TimeList,threeDaysAgoTime)
	TimeList = append(TimeList,beforeYesterdayTime)
	TimeList = append(TimeList,yesterdayTime)
	o := orm.NewOrm()

	for _,time := range TimeList {
		var log1 []Log
		var log2 []Log
		AddCount,err1 := o.Raw("select * from log where content like ? and create_time like ?","%Type:add%","%"+time+"%").QueryRows(&log1)
		UpdateCount,err2 := o.Raw("select * from log where content like ? and create_time like ?","%Type:update%","%"+time+"%").QueryRows(&log2)
		if err1 == nil {
			AddAssetList = append(AddAssetList,AddCount)
		}
		if err2 == nil {
			UpdateAssetList = append(UpdateAssetList,UpdateCount)
		}
	}
	return AddAssetList,UpdateAssetList
}