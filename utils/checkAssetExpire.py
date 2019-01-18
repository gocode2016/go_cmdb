#-*-coding:utf8-*-

import pymysql
from datetime import datetime,timedelta
from mail import SendMail

class CheckAssetExpire:

    def __init__(self):
        self.__conn = pymysql.Connect(host="localhost",port=3306,user="root",passwd="",db="cmdb",charset="utf8")
        self.__cursor = self.__conn.cursor()
        self.expire_asset_list = []

    def _select(self):
        sql = "select inner_ip,outer_ip,expire_time from asset where expire_time is NOT NULL or expire_time != ''"
        self.__cursor.execute(sql)
        for row in self.__cursor.fetchall():
            inner_ip = row[0]
            outer_ip = row[1]
            expire_time = row[2]
            if self._check(expire_time):
                if outer_ip:
                    self.expire_asset_list.append(outer_ip)
                else:
                    self.expire_asset_list.append(inner_ip)

    def _check(self,expire_time):
        now = datetime.now().date()
        tommorow = now + timedelta(days=1)
        if str(tommorow) == expire_time:
            return True
        else:
            return False

    def send(self):
        self._select()
        expire_assets = ",".join(self.expire_asset_list)
        expire_assets = expire_assets.rstrip(",")
        if expire_assets:
            content = u"Hi:\n    IP为 %s 的服务器明天就要过期，请提前做好相关准备！（系统邮件，请勿回复！）" % expire_assets
            SendMail.send(content)

if __name__ == "__main__":
    c = CheckAssetExpire()
    c.send()
