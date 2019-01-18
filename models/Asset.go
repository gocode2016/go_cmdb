package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"fmt"
	"github.com/astaxie/beego"
)

type Asset struct {
	Id int
	Hostname string `orm:"size(64);"`
	InnerIp string `orm:"size(64);unique"`
	OuterIp string `orm:"size(64);null;unique"`
	Os string `orm:"size(16);null"`
	OsRelease string "size(16);null"
	Cpu string `orm:"size(32);null"`
	Memory string `orm:"size(16);null"`
	Disk string `orm:"size(128);null"`
	Sn string `orm:"size(32);null"`
	AssetType string `orm:"size(16)"`             //资产类型
	AssetEnvironment string `orm:"size(16);null"` //资产环境
	Manufacturer string `orm:"size(64);null"`
	Model string `orm:"size(64);null"`           //主机型号
	AssetNumber string `orm:"size(32);null"`     //资产编号
	CabinetNumber string `orm:"size(32);null"`   //机柜号
	CabinetPosition string `orm:"size(32);null"` //机柜位置
	IsPublisher bool `orm:"default(false)"`      //是否是发布机
	AssetStatus string `orm:"size(16)"`          //机器状态:线上，空闲，下架，报修
	Applications string `orm:"size(128);null"`    //应用
	JdkVersion string `orm:"size(16);null"`       //jdk版本
	Comment string `orm:"size(128);null"`         //备注
	OpManager string `orm:"size(16);null"`        //运维负责人
	TomcatVersion string `orm:"size(256);null"`    //tomcat版本
	AppUser string `orm:"size(128);null"`         //应用启动用户
	Idc *Idc `orm:"null;rel(fk);on_delete(set_null)"`
	AssetGroup []*AssetGroup `orm:"rel(m2m)"`
	ExpireTime string `orm:"null"`                //过期时间
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

type AssetMap map[string] interface{}

func CheckAssetsExist(hostname, inner_ip, outer_ip interface{}) (bool,string) {
	query := orm.NewOrm().QueryTable("asset")
	if hostname != nil {
		newHostname := strings.ToLower(hostname.(string))
		if query.Filter("hostname",newHostname).Exist() {
			return true,hostname.(string)
		}
	}
	if inner_ip != nil {
		if query.Filter("inner_ip",inner_ip).Exist() {
			return true,inner_ip.(string)
		}
	}
	if outer_ip != nil {
		if query.Filter("outer_ip",outer_ip).Exist() {
			return true,outer_ip.(string)
		}
	}
	return false,""
}

func AddAsset(asset_map AssetMap)(int64,bool,error){
	o := orm.NewOrm()
	var asset Asset
	idc_id,_ := strconv.Atoi(asset_map["idc"].(string))
	idc := Idc{Id:idc_id}
	asset.Hostname = asset_map["hostname"].(string)
	asset.InnerIp = asset_map["inner_ip"].(string)
	asset.OuterIp = asset_map["outer_ip"].(string)
	asset.Os = asset_map["os"].(string)
	asset.OsRelease = asset_map["os_release"].(string)
	asset.Cpu = asset_map["cpu"].(string)
	asset.Memory = asset_map["memory"].(string)
	asset.Disk = asset_map["disk"].(string)
	asset.Sn = asset_map["sn"].(string)
	asset.AssetType = asset_map["asset_type"].(string)
	asset.AssetEnvironment = asset_map["asset_environment"].(string)
	asset.Manufacturer = asset_map["manufacturer"].(string)
	asset.Model = asset_map["model"].(string)
	asset.AssetNumber = asset_map["asset_number"].(string)
	asset.CabinetNumber = asset_map["cabinet_number"].(string)
	asset.CabinetPosition = asset_map["cabinet_position"].(string)
	asset.IsPublisher,_ = strconv.ParseBool(strconv.Itoa(asset_map["is_publisher"].(int)))
	asset.AssetStatus = asset_map["asset_status"].(string)
	asset.Applications = asset_map["applications"].(string)
	asset.Comment = asset_map["comment"].(string)
	asset.OpManager = asset_map["op_manager"].(string)
	asset.JdkVersion = asset_map["jdk_version"].(string)
	asset.Idc = &idc
	asset.ExpireTime= asset_map["expire_time"].(string)
	id,err := o.Insert(&asset)
	if id > 0 && err == nil {
		return id,true,nil
	} else {
		return 0,false,err
	}
}

func GetAssetList(page int,idc_id,asset_status,keyword string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	var offset int
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("asset")
	number,_ := beego.AppConfig.Int("number")
	if idc_id != "" {
		idc_id,_ := strconv.Atoi(idc_id)
		qs = qs.Filter("Idc",idc_id)
	}
	if asset_status != "" {
		qs = qs.Filter("AssetStatus",asset_status)
	}
	if keyword != "" {
		cond1 := cond.Or("Hostname__icontains",keyword).Or("InnerIp__icontains",keyword).Or("OuterIp__icontains",keyword).Or("Applications__icontains",keyword).Or("OpManager__icontains",keyword)
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

func GetAsset(id int)(*Asset,error){
	var asset Asset
	err := orm.NewOrm().QueryTable("asset").Filter("Id",id).One(&asset)
	if err != nil {
		return nil,err
	}
	return &asset,nil
}

func MultiUpdate(id int,UpdateDataMaps orm.Params)(bool,error){
	num,err := orm.NewOrm().QueryTable("asset").Filter("Id",id).Update(UpdateDataMaps)
	if num > 0 && err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func UpdateAssetGroup(id int,oldVal,newVal interface{}){
	var group AssetGroup
	asset := Asset{Id:id}
	m2m := orm.NewOrm().QueryM2M(&asset,"AssetGroup")
	if strings.Contains(oldVal.(string),","){
		var groups []AssetGroup
		oldValList := strings.Split(oldVal.(string),",")
		for i := 0; i < len(oldValList);i++ {
			err := orm.NewOrm().QueryTable("AssetGroup").Filter("group_name",oldValList[i]).One(&group)
			if err == nil {
				groups = append(groups,group)
			}
		}
		m2m.Remove(groups)
	} else {
		groupName := oldVal.(string)
		var group AssetGroup
		err := orm.NewOrm().QueryTable("AssetGroup").Filter("group_name",groupName).One(&group)
		if err == nil {
			m2m.Remove(group)
		}
	}
	if strings.Contains(newVal.(string),","){
		var groups []AssetGroup
		newValList := strings.Split(newVal.(string),",")
		for i := 0; i < len(newValList);i++ {
			err := orm.NewOrm().QueryTable("AssetGroup").Filter("group_name",newValList[i]).One(&group)
			if err == nil {
				groups = append(groups,group)
			}
		}
		m2m.Add(groups)
	} else {
		groupName := newVal.(string)
		var group AssetGroup
		err := orm.NewOrm().QueryTable("AssetGroup").Filter("group_name",groupName).One(&group)
		if err == nil {
			m2m.Add(group)
		}
	}
}

func UpdateIdc(id int,IdcName string){
	O := orm.NewOrm()
	var idc Idc
	err := O.QueryTable("idc").Filter("name",IdcName).One(&idc,"Id")
	if err == nil {
		_, err := O.QueryTable("asset").Filter("id", id).Update(orm.Params{"Idc": idc.Id})
		if err != nil {
			fmt.Println("err:", err)
		}
	}
}

func DeleteAsset(id int)(bool,error){
	if _,err := orm.NewOrm().Delete(&Asset{Id:id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func GetTotalAssets()(int64){
	count,err := orm.NewOrm().QueryTable("asset").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func GetExpiredAsset()(bool,[]*Asset){
	var assets []*Asset
	now := time.Now()
	tomorrowTime := now.AddDate(0,0,1).Format("2006-01-02")
	num,err := orm.NewOrm().QueryTable("asset").Filter("expire_time",tomorrowTime).All(&assets)
	if num > 0 && err == nil {
		return true,assets
	} else {
		return false,assets
	}
}

//func GetGroupList()([]orm.Params,error){
//	var maps []orm.Params
//	_,err := orm.NewOrm().QueryTable("asset").Filter("asset_environment__icontains","product").Exclude("applications__icontains","redis").
//		Exclude("applications__icontains","localhost").Exclude("applications__icontains","db").Exclude("applications__icontains","backup").
//		Exclude("applications__icontains","apollo").Exclude("applications__icontains","log").Exclude("applications__icontains","ec2").Exclude("applications__icontains","bi").
//		Exclude("applications__icontains","balancer").Exclude("applications__icontains","rabbit").Exclude("applications__icontains","wlj").Exclude("applications__icontains","yuce").
//		Exclude("applications__icontains","mongodb").Exclude("applications__icontains","gpn").Exclude("applications__icontains","nginx").Exclude("applications__icontains","logstash").
//		Exclude("applications__icontains","salt").Exclude("applications__icontains","git").Exclude("applications__contains","AT").Exclude("applications__icontains","vpn").
//		Exclude("applications__icontains","slave").Exclude("applications__icontains","mysql").Exclude("applications__icontains","balncer").Exclude("applications__icontains","zookeeper").
//		Exclude("applications","master").Exclude("applications__icontains","image").Exclude("applications__icontains","hk").Exclude("applications__icontains","ldap").Exclude("applications__icontains","gateway").
//		Exclude("applications__icontains","img").Exclude("applications__icontains","ansible").Exclude("applications","es").Exclude("applications__icontains","dallas").Distinct().Values(&maps,"Applications","OpManager")
//	if err == nil {
//		return maps,nil
//	}
//	return maps,err
//}