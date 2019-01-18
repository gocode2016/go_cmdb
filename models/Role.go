package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"github.com/pkg/errors"
)

//用户角色
type Roles struct {
	Id int
	Name string `orm:"size(16);unique"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	User []*Users `orm:"reverse(many)"`
	Permission []*Permission `orm:"rel(m2m)"`
}

func GetRoleList()([]orm.Params,error){
	var maps []orm.Params
	_,err := orm.NewOrm().QueryTable("roles").Values(&maps)
	if err == nil {
		return maps,nil
	}
	return maps,err
}

func ModifyRole(id int,role_list []string)(bool){
	var roles []*Roles
	o := orm.NewOrm()
	user := Users{Id:id}
	m2m := o.QueryM2M(&user,"Role")
	num,err := orm.NewOrm().QueryTable("roles").Filter("User__Users__Id",id).All(&roles)
	if num >0 && err == nil {
		_, err2 := m2m.Remove(roles)
		if err2 == nil {
			for _,v := range role_list {
				v1,_ := strconv.Atoi(v)
				role := &Roles{Id:v1}
				_,err3 := m2m.Add(role)
				fmt.Println("err:",err3)
			}
			return true
		}
	} else if num == 0 && err == nil {
		for _,v := range role_list {
			v1,_ := strconv.Atoi(v)
			role := &Roles{Id:v1}
			_,err3 := m2m.Add(role)
			fmt.Println("err:",err3)
		}
		return true
	}
	return false
}

func GetRoles(keyword string)([]orm.Params,error){
	var maps []orm.Params
	var err error
	query := orm.NewOrm().QueryTable("roles")
	if keyword != "" {
		query = query.Filter("name__icontains",keyword)
	}
	_,err = query.OrderBy("-created").Values(&maps)
	if err != nil {
		return nil,err
	}
	return maps,nil
}

func DeleteRole(id int) bool {
	var permissions []*Permission
	var users []*Users
	o := orm.NewOrm()
	role := Roles{Id:id}
	p_m2m := o.QueryM2M(&role,"Permission")
	u_m2m := o.QueryM2M(&role,"User")
	_,err := orm.NewOrm().QueryTable("Permission").All(&permissions)
	_,err1 := orm.NewOrm().QueryTable("Users").All(&users)
	if err == nil && err1 == nil {
		_, err2 := p_m2m.Remove(permissions)
		_, err3 := u_m2m.Remove(users)
		if err2 != nil || err3 != nil {
			return false
		}
	}
	_, err4 := o.Delete(&role)
	if err4 == nil {
		return true
	} else {
		return false
	}
}

func AddRole(role_name string)(bool,error){
	o := orm.NewOrm()
	role := Roles{Name:role_name}
	if created,_,err := o.ReadOrCreate(&role,"Name"); err == nil {
		if created {
			return true,nil
		} else {
			return false,errors.New("该角色已存在")
		}
	} else {
		return false,err
	}
}
