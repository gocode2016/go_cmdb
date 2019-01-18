package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type Users struct {
	Id                 int
	UserName           string `orm:"size(24);unique;column(username)"`
	GivenName          string `orm:"size(24);column(givenname)"`
	IsSuper            bool `orm:"default(0)"`
	IsLocked		   bool `orm:"default(0)"`
	IsEnabled          bool `orm:"default(1)"`
	Mobile             string `orm:"size(16);null"`
	Email              string `orm:"size(256);null"`
	Avatar             string `orm:"size(256);null"`
	Role               []*Roles `orm:"rel(m2m)"`
}

func GetUser(username string) (*Users,error){
	var user Users
	err := orm.NewOrm().QueryTable("users").Filter("username",username).One(&user)
	if err != nil {
		return nil,err
	}
	return &user,nil
}

func GetRole(username string)([]*Roles,error) {
	var roles []*Roles
	_,err := orm.NewOrm().QueryTable("roles").Filter("User__Users__UserName",username).All(&roles)
	if err != nil {
		return nil,err
	}
	return roles,nil
}

func GetUserInfo(username string) (*Users,error){
	var user Users
	err := orm.NewOrm().QueryTable("users").Filter("username",username).One(&user)
	if err != nil {
		return nil,err
	}
	return &user,nil
}

func GetUsers(status int,keyword string,page int)([]orm.Params,int64,error){
	var maps []orm.Params
	var err error
	var count int64 = 0
	query := orm.NewOrm().QueryTable("users")
	number,_ := beego.AppConfig.Int("number")
	if status != 2 {
		query = query.Filter("is_enabled",status)
	}
	if keyword != "" {
		query = query.Filter("username__icontains",keyword)
	}
	count,_ = query.Count()
	offset := (page-1) * number
	_,err = query.Limit(number,offset).Values(&maps)
	if err != nil {
		return nil,count,err
	}
	return maps,count,nil
}

func DeleteUser(id int) bool {
	var roles []*Roles
	o := orm.NewOrm()
	user := Users{Id:id}
	m2m := o.QueryM2M(&user,"Role")
	_,err := orm.NewOrm().QueryTable("roles").Filter("User__Users__Id",id).All(&roles)
	if err == nil {
		_, err1 := m2m.Remove(roles)
		_, err2 := o.Delete(&user)
		if err1 == nil && err2 == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func ChangeStatus(id int,status bool) bool {
	user := Users{Id:id}
	user.IsEnabled = !status
	if _,err := orm.NewOrm().Update(&user,"IsEnabled"); err == nil {
		return true
	} else {
		return false
	}
}

func GetTotalUsers()(int64){
	count,err := orm.NewOrm().QueryTable("users").Count()
	if err == nil {
		return count
	} else {
		return 0
	}
}


func CheckRoleExist(username string)(bool) {
	var role []*Roles
	num,err := orm.NewOrm().QueryTable("roles").Filter("User__Users__UserName",username).All(&role)
	if num == 0 && err == nil {
		return true
	} else {
		return false
	}
}

func AddRoleForUser(username string){
	var user Users
	err := orm.NewOrm().QueryTable("Users").Filter("UserName",username).One(&user)
	if err == nil {
		u := Users{Id:user.Id}
		m2m := orm.NewOrm().QueryM2M(&u, "Role")
		var role Roles
		err := orm.NewOrm().QueryTable("Roles").Filter("Name","user").One(&role)
		if err == nil {
			r := &Roles{Id: role.Id}
			m2m.Add(r)
		}
	}

}
