#-*-coding:utf8-*-

import redis
import pymysql
import time
import datetime
import json
import os
import log

class RemoteRedis:

    def __init__(self):
        self.__conn = redis.Redis(host='172.31.1.13', port=6379)
        self.chan_pub = "test"
        self.chan_sub = "test"

    def subscribe(self):
        pub = self.__conn.pubsub()
        pub.subscribe(self.chan_sub)
        pub.parse_response()
        return pub

class LocalRedis:

    key = "%s-test" % time.strftime("%Y%m%d", time.localtime())

    def __init__(self,log):
        self.__conn = redis.Redis(host="127.0.0.1",port=6379)
        self.log = log

    def _save(self):
        self.__conn.slaveof(host="172.31.1.13",port=6379)
        self.__conn.wait(1, 30000)
        if self.__conn.get(LocalRedis.key):
            return True
        else:
            return False

    def get(self):
        if self._save():
            self.log.info(u"服务端同步测试环境redis成功")
            value = self.__conn.get(LocalRedis.key)
            if value:
                json_value = json.loads(value)
                if json_value.get("unreachable"):
                    self.log.warning(u"本次连不上的机器有:%s" % str(json_value.get("unreachable")))
                del json_value["unreachable"]
                return json_value
        else:
            self.log.error(u"服务端同步测试环境redis失败")

class MySQLClient:

    def __init__(self,value,log):
        self.__conn = pymysql.Connect(host="localhost",port=3306,user="root",passwd="",db="cmdb",charset="utf8")
        self.__cursor = self.__conn.cursor()
        self.asset_list = []
        self.asset_dict = {}
        self.new_asset_list = []
        self.update_asset_dict = {}
        self.update_asset_list = []
        self.value = value
        self.log = log

    def insert(self):
        if self.new_asset_list:
            for asset in self.new_asset_list:
                v = self.value.get(asset)
                if isinstance(v.get("Disk"), dict):
                    disk = ""
                    for k1,v1 in v.get("Disk").items():
                        disk += "%s:%s" % (k1,v1)
                insert_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                sql = "insert into asset(`hostname`,`inner_ip`,`outer_ip`,`os`,`os_release`,`cpu`,`memory`,`disk`,`sn`,`asset_type`,`asset_environment`,`manufacturer`,`model`,`asset_status`,`idc_id`,`create_time`,`update_time`) values('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s');"
                data = (v.get("Hostname"),v.get("InnerIp"),v.get("OuterIp"),v.get("Os"),v.get("OsRelease"),v.get("Cpu"),v.get("Memory"),disk,v.get("Sn"),v.get("AssetType"),v.get("AssetEnvironment"),v.get("Manufacturer"),v.get("Model"),v.get("AssetStatus"),v.get("Idc"),insert_time,insert_time)
                self.__cursor.execute(sql % data)
                self.__conn.commit()
        else:
            pass

    def update(self):
        if self.update_asset_dict:
            for k,v in self.update_asset_dict.items():
                sql = "update asset set "
                if len(self.update_asset_dict[k]) == 1:
                    for k1,v1 in v.items():
                        sql += "%s = '%s' " % (k1,v1.get("newVal"))
                else:
                    for k1,v1 in v.items():
                        sql += "%s = '%s'," % (k1,v1.get("newVal"))
                update_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                sql += ",update_time='%s' where inner_ip='%s';" % (update_time,k)
                if sql.find(",,update_time") > 0:
                    sql = sql.replace(",,update_time",",update_time")
                self.__cursor.execute(sql)
                self.__conn.commit()
        else:
            pass

    def select(self):
        sql = "select * from asset where asset_environment = 'test'"
        self.__cursor.execute(sql)
        for row in self.__cursor.fetchall():
            InnerIp = row[2]
            Hostname = row[1]
            OuterIp = row[3]
            Os = row[4]
            OsRelease = row[5]
            Cpu = row[6]
            Memory = row[7]
            Disk = row[8]
            Sn = row[9]
            AssetType = row[10]
            AssetEnvironment = row[11]
            Manufacturer = row[12]
            Model = row[13]
            AssetStatus = row[18]
            self.asset_list.append(InnerIp)
            self.asset_dict[InnerIp] = {"Hostname":Hostname,"OuterIp":OuterIp,"Os":Os,"OsRelease":OsRelease,"Cpu":Cpu,"Memory":Memory,"Disk":Disk,
                                        "Sn":Sn,"AssetType":AssetType,"AssetEnvironment":AssetEnvironment,"Manufacturer":Manufacturer,
                                        "Model":Model,"AssetStatus":AssetStatus}

    def compare(self,value):
        for k,v in value.items():
            if k in self.asset_list:
                self.update_asset_dict[k] = {}
                disk = ""
                if isinstance(v.get("Disk"), dict):
                    for k1, v1 in v.get("Disk").items():
                        disk += "%s:%s " % (k1, v1)
                disk = disk.strip()
                if v.get("Hostname") != self.asset_dict[k].get("Hostname"):
                    self.update_asset_dict[k]["hostname"] = {"newVal":v.get("Hostname"),"oldVal":self.asset_dict[k].get("Hostname")}
                if v.get("OuterIp") != self.asset_dict[k].get("OuterIp"):
                    self.update_asset_dict[k]["outer_ip"] = {"newVal": v.get("OuterIp"),"oldVal": self.asset_dict[k].get("OuterIp")}
                if v.get("Os") != self.asset_dict[k].get("Os"):
                    self.update_asset_dict[k]["os"] = {"newVal":v.get("Os"),"oldVal":self.asset_dict[k].get("Os")}
                if v.get("OsRelease") != self.asset_dict[k].get("OsRelease"):
                    self.update_asset_dict[k]["os_release"] = {"newVal": v.get("OsRelease"),"oldVal": self.asset_dict[k].get("OsRelease")}
                if v.get("Sn") != self.asset_dict[k].get("Sn"):
                    self.update_asset_dict[k]["sn"] = {"newVal":v.get("Sn"),"oldVal":self.asset_dict[k].get("Sn")}
                if v.get("Cpu") != self.asset_dict[k].get("Cpu"):
                    self.update_asset_dict[k]["cpu"] = {"newVal": v.get("Cpu"),"oldVal": self.asset_dict[k].get("Cpu")}
                if v.get("Memory") != self.asset_dict[k].get("Memory"):
                    self.update_asset_dict[k]["memory"] = {"newVal":v.get("Memory"),"oldVal":self.asset_dict[k].get("Memory")}
                if disk != self.asset_dict[k].get("Disk").strip():
                    self.update_asset_dict[k]["disk"] = {"newVal": disk,"oldVal": self.asset_dict[k].get("Disk").strip()}
            else:
                self.new_asset_list.append(k)
        if self.new_asset_list:
            self.log.info(u"新增加了%d条资产，IP是%s" % (len(self.new_asset_list),self.new_asset_list))
        for k in list(self.update_asset_dict.keys()):
            if not self.update_asset_dict[k]:
                del self.update_asset_dict[k]
        if self.update_asset_dict:
            self.log.info(u"本次资产变更信息:%s",str(self.update_asset_dict))

    def insert_log(self):
        if self.update_asset_dict:
            for k,v in self.update_asset_dict.items():
                name = "robot"
                ip = k
                content = "Type:update,"
                for k1,v1 in v.items():
                    content += "Field:%s,oldValue:%s,newValue:%s," % (k1,v1.get("oldVal"),v1.get("newVal"))
                create_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                sql = "insert into log(`name`,`ip`,`content`,`create_time`) values('%s','%s','%s','%s')" % (
                name, ip, content, create_time)
                self.__cursor.execute(sql)
                self.__conn.commit()
        elif self.new_asset_list:
            for ip in self.new_asset_list:
                asset = self.value.get(ip)
                name = "robot"
                ip = ip
                content = "Type:add,Hostname:%s,InnerIp:%s" % (asset.get("Hostname"),asset.get("InnerIp"))
                create_time = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                sql = "insert into log(`name`,`ip`,`content`,`create_time`) values('%s','%s','%s','%s')" % (name, ip, content, create_time)
                self.__cursor.execute(sql)
                self.__conn.commit()
        else:
            pass

def main():
    log_file = os.path.dirname(os.getcwd())+"/logs/test.log"
    log_obj = log.init_logger(log_file)
    obj = RemoteRedis()
    redis_sub = obj.subscribe()
    while True:
        msg = redis_sub.parse_response()
        if msg[2] == b'test finish':
            log_obj.info(u"服务端已经收到测试环境发布的消息")
            obj = LocalRedis(log_obj)
            value = obj.get()
            obj2 = MySQLClient(value, log_obj)
            obj2.select()
            obj2.compare(value)
            obj2.insert()
            obj2.update()
            obj2.insert_log()

if __name__ == "__main__":
    main()


