package utils

import (
	"github.com/astaxie/beego"
	"gopkg.in/ldap.v2"
	"fmt"
)

//LDAP 登录验证
func UserAuthLdap(username, password string) (bool){
	// The username and password we want to check
	bindUsername := beego.AppConfig.String("ldap::mdn")
	bindPassword := beego.AppConfig.String("ldap::mpw")
	ldapServer := beego.AppConfig.String("ldap::ldapServer")
	baseDn := beego.AppConfig.String("ldap::basedn")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s", ldapServer))
	if err != nil {
		fmt.Printf("Ldap server connect error!")
		return false
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		fmt.Printf("Ldap manager auth error!")
		return false
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(cn=%s))", username), //指定用户
		[]string{"dn", "cn", "givenname", "email"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		fmt.Printf("Ldap serach error!")
		return false
	}
	//fmt.Println(len(sr.Entries))
	if len(sr.Entries) != 1 {
		fmt.Printf("User does not exist or too many entries returned")
		return false
	}

	userdn := sr.Entries[0].DN
	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("Ldap user auth error!")
		return false
	}

	/*var getUser, getName, getMail string
	for _, entry := range sr.Entries {
		// fmt.Printf("%s\n", entry.DN)
		getUser =  entry.GetAttributeValue("cn")
		getName = entry.GetAttributeValue("givenName")
		getMail =  entry.GetAttributeValue("email")
	}*/
	return true
}


/*//LDAP 获取账号信息并返回map 和 json格式数据
func UserInfoLdap() ([]User) {
	// The username and password we want to check
	bindUsername := beego.AppConfig.String("ldap::mdn")
	bindPassword := beego.AppConfig.String("ldap::mpw")
	ldapServer := beego.AppConfig.String("ldap::ldapServer")
	basedn := beego.AppConfig.String("ldap::basedn")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s", ldapServer))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(cn=*))", //所有用户
		[]string{"dn", "cn", "givenname", "email"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	//遍历所有用户
	var user User
	var userList []User
	for _, entry := range sr.Entries {
		//userList.Dn = entry.DN
		user.UserName =  entry.GetAttributeValue("cn")
		user.GivenName = entry.GetAttributeValue("givenName")
		user.Email = entry.GetAttributeValue("email")
		userList = append(userList, user) //添加slice列表
	}

	return userList //返回用户列表map 和 json格式
}

type TestJob struct {

}

func (this TestJob) Run() {
	fmt.Println("testJob.................")
}
*/