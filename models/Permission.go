package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"strings"
)

//用户权限
type Permission struct {
	Id int
	Name string `orm:"size(32);unique"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Comment string `orm:"size(128)"`
	Role []*Roles `orm:"reverse(many)"`
}

func GetPermissionList(page int,keyword,sort string)([]orm.Params,int64,error){
	var maps []orm.Params
	var count int64 = 0
	number,_ := beego.AppConfig.Int("number")
	cond := orm.NewCondition()
	qs := orm.NewOrm().QueryTable("permission")
	if keyword != "" {
		cond1 := cond.Or("Name__icontains",keyword).Or("Comment__icontains",keyword)
		qs = qs.SetCond(cond1)
	}
	if sort != "" {
		qs = qs.OrderBy("-created")
	} else {
		qs = qs.OrderBy("created")
	}
	count,_ = qs.Count()
	offset := (page-1) * number
	_,err := qs.Limit(number,offset).Values(&maps)
	if err == nil {
		return maps,count,nil
	}
	return maps,count,err
}

func AddPermission(permission_name,permission_comment string)(bool,error){
	o := orm.NewOrm()
	permission := Permission{Name:permission_name,Comment:permission_comment}
	if created,_,err := o.ReadOrCreate(&permission,"Name","Comment"); err == nil {
		if created {
			return true,nil
		} else {
			return false,errors.New("该权限已存在")
		}
	} else {
		return false,err
	}
}

func DeletePermission(id int) bool {
	o := orm.NewOrm()
	permission := Permission{Id:id}
	_, err := o.Delete(&permission)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetTotalPermissions()(int64){
	count,err := orm.NewOrm().QueryTable("permission").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}

func GetPermissionsByRoleId(role_id int)([]*Permission,[]*Permission,error){
	var ownedPermissions []*Permission
	var unownedPermissions []*Permission
	var ids []int
	num,err := orm.NewOrm().QueryTable("permission").Filter("Role__Roles__id",role_id).All(&ownedPermissions)
	if num > 0 && err == nil {
		for _,v := range ownedPermissions {
			ids = append(ids,v.Id)
		}
	}

	if len(ids) != 0 {
		num,err := orm.NewOrm().QueryTable("permission").Exclude("id__in",ids).All(&unownedPermissions)
		if num > 0 && err == nil {
			return ownedPermissions,unownedPermissions,nil
		} else {
			return ownedPermissions,unownedPermissions,err
		}
	} else {
		num,err := orm.NewOrm().QueryTable("permission").All(&unownedPermissions)
		if num > 0 && err == nil {
			return ownedPermissions,unownedPermissions,nil
		} else {
			return ownedPermissions,unownedPermissions,err
		}
	}
}

func ModifyPermissionsForRole(role_id int,ids []int)(bool) {
	var permissions []Permission
	o := orm.NewOrm()
	role := Roles{Id: role_id}
	m2m := o.QueryM2M(&role, "Permission")
	num, err := orm.NewOrm().QueryTable("permission").Filter("Role__Roles__Id", role_id).All(&permissions)
	if num > 0 && err == nil {
		_, err := m2m.Remove(permissions)
		if err != nil {
			return false
		}
	}
	if len(ids) != 0 {
		_, err := orm.NewOrm().QueryTable("permission").Filter("Id__in", ids).All(&permissions)
		if err == nil {
			_, err := m2m.Add(permissions)
			if err == nil {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return true
	}
}

func GetPermissionByUsername(username string)(map[string]int) {
	var roles []*Roles
	var role_list []int
	num,err := orm.NewOrm().QueryTable("roles").Filter("User__Users__UserName",username).All(&roles)
	permissionMap := make(map[string]int)
	if err == nil {
		if num > 0 {
			for _,role := range roles {
				role_list = append(role_list,role.Id)
			}
		}
		var permissions []*Permission
		for _,id := range role_list {
			num,err := orm.NewOrm().QueryTable("permission").Filter("Role__Roles__Id",id).All(&permissions)
			if num > 0 && err == nil {
				for _,permission := range permissions {
					permissionMap[permission.Name] = 1
				}
			}
		}
	}
	return permissionMap
}

func CheckPermissionExist(Name string)(bool,error){
	num,err := orm.NewOrm().QueryTable("permission").Filter("Name",strings.ToLower(Name)).Count()
	if num == 0 && err == nil {
		return true,nil
	} else {
		return false,err
	}
}