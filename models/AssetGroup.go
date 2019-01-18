package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"github.com/astaxie/beego"
	"strings"
)

type AssetGroup struct {
	Id int
	GroupName string `orm:"size(64)"`
	AliasName string `orm:"size(64);null"`
	DevManager string `orm:"size(64);null"`    //开发负责人
	OpManager string `orm:"size(64);null"`     //运维负责人
	Comment string `orm:"size(256);null"`      //备注
	PackMethod string `orm:"size(16);null"`    //打包方式
	Asset []*Asset `orm:"reverse(many)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func GetAssetGroupList(page int,keyword,selectType string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	number,_ := beego.AppConfig.Int("number")
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("asset_group")
	if keyword != "" {
		cond1 := cond.Or("GroupName__icontains",keyword).Or("AliasName__icontains",keyword).Or("DevManager__icontains",keyword).Or("OpManager__icontains",keyword)
		qs = qs.SetCond(cond1)
	}
	if selectType == "all" {
		_,err := qs.Values(&maps)
		if err == nil {
			return maps,count,nil
		}
		return maps,count,err
	}
	count,_ = qs.Count()
	offset := (page-1) * number
	_,err := qs.Limit(number,offset).Values(&maps)
	if err == nil {
		return maps,count,nil
	}
	return maps,count,err
}

func AddGroupForAsset(id int64,asset_group_list []string){
	o := orm.NewOrm()
	asset_id := int(id)
	asset:= Asset{Id:asset_id}
	m2m := o.QueryM2M(&asset, "AssetGroup")
	for i := 0;i < len(asset_group_list);i++{
		id,_ := strconv.Atoi(asset_group_list[i])
		group := &AssetGroup{Id:id}
		_,err := m2m.Add(group)
		fmt.Println(err)
	}
}

func GetGroupNameByAsset(id int64)([]*AssetGroup,error){
	var assetGroup []*AssetGroup
	num,err := orm.NewOrm().QueryTable("asset_group").Filter("Asset__Asset__Id",id).All(&assetGroup)
	if num > 0 && err == nil {
		return assetGroup,nil
	}
	return nil,err
}

func AddAssetGroup(GroupName,AliasName,DevManager,OpManager,Comment string)(bool,error){
	var group AssetGroup
	group.GroupName = GroupName
	group.AliasName = AliasName
	group.DevManager = DevManager
	group.OpManager = OpManager
	group.Comment = Comment
	_,err := orm.NewOrm().Insert(&group)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func DeleteAssetGroup(id int)(bool,error){
	if _,err := orm.NewOrm().Delete(&AssetGroup{Id:id});err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func CheckAssetGroupExist(GroupName,Operation,GroupNameStatus string)(bool,error){
	num,err := orm.NewOrm().QueryTable("asset_group").Filter("GroupName",strings.ToLower(GroupName)).Count()
	if Operation == "add" {
		if num == 0 && err == nil {
			return true,nil
		} else {
			return false,err
		}
	} else if Operation == "modify" {
		if GroupNameStatus == "Change" {
			if num == 0 && err == nil {
				return true, nil
			} else {
				return false, err
			}
		} else if GroupNameStatus == "NotChange" {
			if num == 1 && err == nil {
				return true, nil
			} else {
				return false, err
			}
		}
	}
	return false,nil
}

func ModifyAssetGroup(id int,GroupName,AliasName,DevManager,OpManager,Comment string)(bool,error){
	var asset_map orm.Params
	asset_map = make(orm.Params)
	asset_map["GroupName"] = GroupName
	asset_map["AliasName"] = AliasName
	asset_map["DevManager"] = DevManager
	asset_map["OpManager"] = OpManager
	asset_map["Comment"] = Comment
	_,err := orm.NewOrm().QueryTable("asset_group").Filter("Id",id).Update(asset_map)
	if err == nil {
		return true,nil
	} else {
		return false,err
	}
}

func GetTotalAssetGroups()(int64){
	count,err := orm.NewOrm().QueryTable("asset_group").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func GetAssetsByGroup(id int64)([]*Asset,error){
	var asset []*Asset
	num,err := orm.NewOrm().QueryTable("asset").Filter("AssetGroup__AssetGroup__Id",id).All(&asset)
	if num > 0 && err == nil {
		return asset,nil
	} else {
		return asset,err
	}
}

func GetGroupList()([]orm.Params,error){
	var maps []orm.Params
	_,err := orm.NewOrm().QueryTable("asset_group").Values(&maps,"Id", "GroupName", "PackMethod")
	if err == nil {
		return maps,nil
	}
	return maps,err
}