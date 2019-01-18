package utils

import (
	"github.com/astaxie/beego"
	"zabbix"
	"encoding/json"
	"strconv"
	"errors"
)

var username = beego.AppConfig.String("zabbix::username")
var	password = beego.AppConfig.String("zabbix::password")

func CheckHostExist(hostname,interfaceAddr,interfaceValue,environment string)(bool,error) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false,err
	}
	_, err = api.Login()
	if err != nil {
		return false,err
	}
	//判断主机名是否存在
	data := make(map[string]interface{})
	hostMap := make(map[string]interface{})
	hostMap["host"] = hostname
	data["output"] = "extend"
	data["filter"] = hostMap
	result,_ := api.ZabbixRequest("host.get",data)
	//判断IP或DNS是否存在
	data2 := make(map[string]interface{})
	ipMap := make(map[string]interface{})
	if interfaceAddr == "IP" {
		ipMap["ip"] = interfaceValue
	} else {
		ipMap["dns"] = interfaceValue
	}
	data2["output"] = "host"
	data2["filter"] = ipMap
	result2,_ := api.ZabbixRequest("host.get",data2)
	if len(result.Result.([]interface{})) != 0 {
		return false,errors.New("主机名已存在！")
	} else if len(result2.Result.([]interface{})) != 0 {
		return false,errors.New("接口方式的值已存在！")
	} else {
		return true,nil
	}
}

func AddHost(hostname,visibleName,interfaceAddr,interfaceValue,interfaceType,environment string,groupIdList,templateIdList []string)(int64,bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return 0,false,err.Error()
	}

	_, err = api.Login()
	if err != nil {
		return 0,false,err.Error()
	}
	data := make(map[string]interface{})
	interfaceMap := make(map[string]interface{})
	var interfaceArr []map[string]interface{}
	var groupArr []map[string]interface{}
	var templateArr []map[string]interface{}
	if interfaceAddr == "IP" {
		interfaceMap["ip"] = interfaceValue
		interfaceMap["dns"] = ""
		interfaceMap["useip"] = 1
	} else {
		interfaceMap["ip"] = ""
		interfaceMap["dns"] = interfaceValue
		interfaceMap["useip"] = 0
	}
	//1:agent,2:SNMP,3:IPMX,4:JMX
	if interfaceType == "1" {
		interfaceMap["type"] = 1
		interfaceMap["port"] = "10050"
	} else if interfaceType == "2" {
		interfaceMap["type"] = 2
		interfaceMap["port"] = "161"
	} else if interfaceType == "3" {
		interfaceMap["type"] = 3
		interfaceMap["port"] = "623"
	} else if interfaceType == "4" {
		interfaceMap["type"] = 4
		interfaceMap["port"] = "12345"
	}
	interfaceMap["main"] = 1
	interfaceMap["bulk"] = 1
	interfaceArr = append(interfaceArr,interfaceMap)
	for _,groupId := range groupIdList {
		groupMap := make(map[string]interface{})
		groupMap["groupid"] = groupId
		groupArr = append(groupArr,groupMap)
	}
	for _,templateId := range templateIdList {
		templateMap := make(map[string]interface{})
		templateMap["templateid"] = templateId
		templateArr = append(templateArr,templateMap)
	}
	data["host"] = hostname
	data["name"] = visibleName
	data["interfaces"] = interfaceArr
	data["groups"] = groupArr
	data["templates"] = templateArr
	result,err := api.ZabbixRequest("host.create",data)
	if err == nil && result.Result != nil {
		res, _ := json.Marshal(result.Result)  //{2.0 {0  } map[hostids:[10442]] 1}
		type ZabbixHost map[string][]string
		var ret ZabbixHost
		var hostids = ""
		err = json.Unmarshal(res, &ret)
		for _,hostid := range ret["hostids"] {
			if hostids != "" {
				hostids = hostids+","+hostid
			} else {
				hostids = hostid
			}
		}
		hostid,_ := strconv.Atoi(hostids)
		return int64(hostid),true,""
	} else {
		return 0,false,result.Error.Error()
	}
}

func DeleteHost(HostId int64,environment string)(bool,string){
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false,err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false,err.Error()
	}
	var hostIdsArr []int64
	hostIdsArr = append(hostIdsArr,HostId)
	data := hostIdsArr
	result,_ := api.ZabbixRequest("host.delete",data)
	if  result.Error.Code == 0 {
		return true,""
	} else {
		return false,result.Error.Error()
	}
}

func MassDeleteForeignHost(foreignHostIdArr []string)(bool,string) {
	var url string
	if len(foreignHostIdArr) != 0 {
		url = beego.AppConfig.String("zabbix::foreign_url")
		data := foreignHostIdArr
		api, err := zabbix.NewAPI(url, username, password)
		if err != nil {
			return false, err.Error()
		}
		_, err = api.Login()
		if err != nil {
			return false, err.Error()
		}
		result, _ := api.ZabbixRequest("host.delete", data)
		if result.Error.Code == 0 {
			return true, ""
		} else {
			return false, result.Error.Error()
		}
	} else {
		return false,"要删除的国内主机列表为空！"
	}
}

func MassDeleteCnHost(cnHostIdArr []string)(bool,string){
	var url string
	if len(cnHostIdArr) != 0{
		url = beego.AppConfig.String("zabbix::cn_url")
		data := cnHostIdArr
		api, err := zabbix.NewAPI(url, username, password)
		if err != nil {
			return false,err.Error()
		}
		_, err = api.Login()
		if err != nil {
			return false,err.Error()
		}
		result,_ := api.ZabbixRequest("host.delete",data)
		if result.Error.Code == 0 {
			return true,""
		} else {
			return false,result.Error.Error()
		}
	} else {
		return false,"要删除的国外主机列表为空！"
	}
}

func ChangeStatus(HostId int64,HostStatus,environment string)(bool,string){
	var url string
	var host_status = 0
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false,err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false,err.Error()
	}
	data := make(map[string]interface{})
	if HostStatus == "false" {
		host_status = 1
	} else {
		host_status = 0
	}
	data["hostid"] = HostId
	data["status"] = host_status
	result,_ := api.ZabbixRequest("host.update",data)
	if  result.Error.Code == 0 {
		return true,""
	} else {
		return false,result.Error.Error()
	}
}

func UpdateHost(editedEnvironment,hostname,visibleName string,hostId int64,allocateHostGroupList,allocateTemplateList []int64)(bool,string){
	var url string
	if editedEnvironment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false,err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false,err.Error()
	}
	data := make(map[string]interface{})
	if len(allocateTemplateList) != 0 {
		var addTemplateArr []map[string]int64
		for _,templateId := range allocateTemplateList {
			templateMaps := make(map[string]int64)
			templateMaps["templateid"] = templateId
			addTemplateArr = append(addTemplateArr,templateMaps)
		}
		data["templates"] = addTemplateArr
	}
	if len(allocateHostGroupList) != 0 {
		var addGroupArr []map[string]int64
		for _,groupId := range allocateHostGroupList {
			groupMaps := make(map[string]int64)
			groupMaps["groupid"] = groupId
			addGroupArr = append(addGroupArr,groupMaps)
		}
		data["groups"] = addGroupArr
	}
	//if(len(interfaceTypeList) != 0){
	//	var interfaceTypeArr []map[string]string
	//	for _,typeId := range interfaceTypeList {
	//		interfaceTypeMaps := make(map[string]string)
	//		interfaceTypeMaps["type"] = typeId
	//		interfaceTypeMaps["ip"] = "172.31.0.245"
	//		interfaceTypeMaps["useip"] = "1"
	//		interfaceTypeMaps["main"] = "1"
	//		interfaceTypeMaps["bulk"] = "1"
	//		interfaceTypeMaps["dns"] = ""
	//		interfaceTypeMaps["ensure"] = "present"
	//		if(typeId == "1"){
	//			interfaceTypeMaps["port"] = "10050"
	//		} else if(typeId == "2") {
	//			interfaceTypeMaps["port"] = "161"
	//		} else if(typeId == "3") {
	//			interfaceTypeMaps["port"] = "623"
	//		} else if(typeId == "4") {
	//			interfaceTypeMaps["port"] = "12345"
	//		}
	//		interfaceTypeArr = append(interfaceTypeArr,interfaceTypeMaps)
	//	}
	//	data["interfaces"] = interfaceTypeArr
	//}
	data["hostid"] = hostId
	if hostname != "" {
		data["host"] = hostname
	}
	if visibleName != "" {
		data["name"] = visibleName
	}
	result,_ := api.ZabbixRequest("host.update",data)
	if  result.Error.Code == 0 {
		return true,""
	} else {
		return false,result.Error.Error()
	}
}

func CheckHostGroupNameExist(groupName,environment string)(bool,error) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err
	}
	_, err = api.Login()
	if err != nil {
		return false, err
	}
	data := make(map[string]interface{})
	groupNameMap := make(map[string]interface{})
	groupNameMap["name"] = groupName
	data["output"] = "extend"
	data["filter"] = groupNameMap
	result,_ := api.ZabbixRequest("hostgroup.get",data)
	if len(result.Result.([]interface{})) != 0 {
		return false,errors.New("主机组名已存在！")
	} else {
		return true,nil
	}
}

func AddHostGroupName(groupName,environment string)(int64,bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return 0, false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return 0, false, err.Error()
	}
	data := make(map[string]interface{})
	data["name"] = groupName
	result,err := api.ZabbixRequest("hostgroup.create",data)
	if err == nil && result.Result != nil {
		res, _ := json.Marshal(result.Result)
		type ZabbixHost map[string][]string
		var ret ZabbixHost
		err = json.Unmarshal(res, &ret)
		groupId,_ := strconv.Atoi(ret["groupids"][0])
		return int64(groupId),true,""
	} else {
		return 0,false,result.Error.Error()
	}
}

func DeleteHostGroupName(groupId,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	var groupIdArr []string
	groupIdArr = append(groupIdArr, groupId)
	data := groupIdArr
	result,_ := api.ZabbixRequest("hostgroup.delete",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

func UpdateHostGroupName(groupId,groupName,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	data := make(map[string]interface{})
	data["groupid"] = groupId
	data["name"] = groupName
	result,_ := api.ZabbixRequest("hostgroup.update",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

func CheckTemplateNameExist(templateName,environment string)(bool,error) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err
	}
	_, err = api.Login()
	if err != nil {
		return false, err
	}
	data := make(map[string]interface{})
	templateNameMap := make(map[string]interface{})
	templateNameMap["host"] = templateName
	data["output"] = "extend"
	data["filter"] = templateNameMap
	result,_ := api.ZabbixRequest("template.get",data)
	if len(result.Result.([]interface{})) != 0 {
		return false,errors.New("模板名已存在！")
	} else {
		return true,nil
	}
}

func AddTemplate(templateName,environment string, groupIds, templateIds []string)(int64,bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return 0, false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return 0, false, err.Error()
	}
	data := make(map[string]interface{})
	var groupIdArr []map[string]interface{}
	for _,groupId := range groupIds {
		groupIdMap := make(map[string]interface{})
		groupIdMap["groupid"] = groupId
		groupIdArr = append(groupIdArr, groupIdMap)
	}
	if len(templateIds) != 0 {
		var templateIdArr []map[string]interface{}
		for _,templateId := range templateIds {
			templateIdMap := make(map[string]interface{})
			templateIdMap["templateid"] = templateId
			templateIdArr = append(templateIdArr, templateIdMap)
		}
		data["templates"] = templateIdArr
	}
	data["host"] = templateName
	data["groups"] = groupIdArr
	result,err := api.ZabbixRequest("template.create",data)
	if err == nil && result.Result != nil {
		res, _ := json.Marshal(result.Result)
		type ZabbixHost map[string][]string
		var ret ZabbixHost
		err = json.Unmarshal(res, &ret)
		templateId,_ := strconv.Atoi(ret["templateids"][0])
		return int64(templateId),true,""
	} else {
		return 0,false,result.Error.Error()
	}
}

func DeleteTemplate(templateId,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	var templateIdArr []string
	templateIdArr = append(templateIdArr, templateId)
	data := templateIdArr
	result,_ := api.ZabbixRequest("template.delete",data)
	if result.Error.Code != 0 {
		return false,result.Error.Error()
	} else {
		return true,""
	}
}

func UpdateTemplate(templateId,templateName,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	data := make(map[string]interface{})
	data["templateid"] = templateId
	data["name"] = templateName
	result,_ := api.ZabbixRequest("template.update",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

//func AddUser(name,environment,sendto,period string, groupIds []string)(int64,bool,string) {
//	var url string
//	var mediaTypeId string
//	if environment == "cn" {
//		url = beego.AppConfig.String("zabbix::cn_url")
//		mediaTypeId = beego.AppConfig.String("media::cn_media_type_id")
//	} else {
//		url = beego.AppConfig.String("zabbix::foreign_url")
//		mediaTypeId = beego.AppConfig.String("media::foreign_media_type_id")
//	}
//	api, err := zabbix.NewAPI(url, username, password)
//	if err != nil {
//		return 0, false, err.Error()
//	}
//	_, err = api.Login()
//	if err != nil {
//		return 0, false, err.Error()
//	}
//	data := make(map[string]interface{})
//	var groupArr []map[string]interface{}
//	data["alias"] = name
//	for _,groupId := range groupIds {
//		groupMap := make(map[string]interface{})
//		groupMap["usrgrpid"] = groupId
//		groupArr = append(groupArr,groupMap)
//	}
//	data["usrgrps"] = groupArr
//	if sendto != "" {
//		mediaMap := make(map[string]interface{})
//		var mediaArr []map[string]interface{}
//		mediaMap["mediatypeid"] = mediaTypeId
//		mediaMap["sendto"] = sendto
//		mediaMap["active"] = 0
//		mediaMap["severity"] = 63
//		mediaMap["period"] = period
//		mediaArr = append(mediaArr, mediaMap)
//		data["user_medias"] = mediaArr
//	}
//	result,err := api.ZabbixRequest("user.create",data)
//	if err == nil && result.Result != nil {
//		res, _ := json.Marshal(result.Result)
//		type ZabbixHost map[string][]string
//		var ret ZabbixHost
//		err = json.Unmarshal(res, &ret)
//		userId,_ := strconv.Atoi(ret["userids"][0])
//		return int64(userId),true,""
//	} else {
//		return 0,false,result.Error.Error()
//	}
//}

func DeleteUser(userId,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	var userIdArr []string
	userIdArr = append(userIdArr, userId)
	data := userIdArr
	result,_ := api.ZabbixRequest("user.delete",data)
	if result.Error.Code != 0 {
		return false,result.Error.Error()
	} else {
		return true,""
	}
}

func UpdateUserGroup(userId,environment,sendto,period string,userGroups []int64)(bool,string) {
	var url,mediaTypeId string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
		mediaTypeId = beego.AppConfig.String("media::cn_media_type_id")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
		mediaTypeId = beego.AppConfig.String("media::foreign_media_type_id")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	data := make(map[string]interface{})
	if len(userGroups) != 0 {
		var userGroupIdArr []map[string]interface{}
		for _,groupId := range userGroups {
			groupIdMap := make(map[string]interface{})
			groupIdMap["usrgrpid"] = groupId
			userGroupIdArr = append(userGroupIdArr, groupIdMap)
		}
		data["usrgrps"] = userGroupIdArr
	}
	data["userid"] = userId
	if sendto != "" {
		mediaMap := make(map[string]interface{})
		var mediaArr []map[string]interface{}
		mediaMap["mediatypeid"] = mediaTypeId
		mediaMap["sendto"] = sendto
		mediaMap["active"] = 0
		mediaMap["severity"] = 63
		mediaMap["period"] = period
		mediaArr = append(mediaArr, mediaMap)
		data["user_medias"] = mediaArr
	}
	result,_ := api.ZabbixRequest("user.update",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

func CheckUserGroupNameExist(groupName,environment string)(bool,error) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err
	}
	_, err = api.Login()
	if err != nil {
		return false, err
	}
	data := make(map[string]interface{})
	groupNameMap := make(map[string]interface{})
	groupNameMap["name"] = groupName
	data["output"] = "extend"
	data["filter"] = groupNameMap
	result,_ := api.ZabbixRequest("usergroup.get",data)
	if len(result.Result.([]interface{})) != 0 {
		return false,errors.New("用户组名已存在！")
	} else {
		return true,nil
	}
}

func AddUserGroup(groupName,environment string)(int64,bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return 0, false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return 0, false, err.Error()
	}
	data := make(map[string]interface{})
	data["name"] = groupName
	result,err := api.ZabbixRequest("usergroup.create",data)
	if err == nil && result.Result != nil {
		res, _ := json.Marshal(result.Result)
		type ZabbixUser map[string][]string
		var ret ZabbixUser
		err = json.Unmarshal(res, &ret)
		groupId,_ := strconv.Atoi(ret["usrgrpids"][0])
		return int64(groupId),true,""
	} else {
		return 0,false,result.Error.Error()
	}
}

func DeleteUserGroup(groupId,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	var groupIdArr []string
	groupIdArr = append(groupIdArr, groupId)
	data := groupIdArr
	result,_ := api.ZabbixRequest("usergroup.delete",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

func ChangeUserGroupStatus(groupId,environment,groupStatus string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	data := make(map[string]interface{})
	data["usrgrpid"] = groupId
	if groupStatus == "true" {
		data["users_status"] = "0"
	} else {
		data["users_status"] = "1"
	}
	result,_ := api.ZabbixRequest("usergroup.update",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}

func UpdateUserGroupName(groupId,groupName,environment string)(bool,string) {
	var url string
	if environment == "cn" {
		url = beego.AppConfig.String("zabbix::cn_url")
	} else {
		url = beego.AppConfig.String("zabbix::foreign_url")
	}
	api, err := zabbix.NewAPI(url, username, password)
	if err != nil {
		return false, err.Error()
	}
	_, err = api.Login()
	if err != nil {
		return false, err.Error()
	}
	data := make(map[string]interface{})
	data["usrgrpid"] = groupId
	data["name"] = groupName
	result,_ := api.ZabbixRequest("usergroup.update",data)
	if result.Error.Code != 0 {
		return false, result.Error.Error()
	} else {
		return true, ""
	}
}