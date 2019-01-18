#-*-coding:utf8-*-

from ldap3 import Server, Connection
import pymysql
import log

LDAP_HOST = '172.31.1.50:389'
BASE_DN = 'ou=users,dc=jollychic,dc=com'
USER = 'cn=Manager,dc=jollychic,dc=com'
PASSWORD = 'se12pa'

class myLdap:

    def __init__(self,ldap_host,base_dn,user,password):
        self.ldap_host = ldap_host
        self.base_dn = base_dn
        self.user = user
        self.password = password


    def sync_users(self):
        ldap_server = Server(self.ldap_host)
        conn = Connection(ldap_server, user = self.user, password = self.password, check_names=True, lazy=False,
                          raise_exceptions=False)
        conn.open()
        conn.bind()
        params = {'search_base': self.base_dn, 'search_filter': '(&(cn=*))', 'attributes': ['cn', 'givenName', 'email']}
        res = conn.search(**params)
        user_dict = {}
        if res:
            items = conn.entries
            for item in items:
                if item:
                    username = str(item.cn).strip()
                    if username == "Admin":
                        continue
                    user_dict[username] = {}
                    if str(item.givenName).strip() == '[]':
                        given_name = ""
                    else:
                        given_name = str(item.givenName).strip()
                    if str(item.email).strip() == '[]':
                        email = ""
                    else:
                        email = str(item.email).strip()
                    user_dict[username]["given_name"] = given_name
                    user_dict[username]["email"] = email
        return user_dict

class mysqlClient:

    def __init__(self,host,username,passwd,db):
        self.host = host
        self.username = username
        self.passwd = passwd
        self.db = db
        self.__conn = pymysql.Connect(host=self.host, port=3306, user=self.username, passwd=self.passwd, db=self.db,
                                      charset="utf8")
        self.__cursor = self.__conn.cursor()
        self.user_list = []

    def get_users(self):
        sql = "select username from users"
        self.__cursor.execute(sql)
        for row in self.__cursor.fetchall():
            username = row[0]
            self.user_list.append(username)
        return self.user_list

    def compare_users(self):
        user_list = self.get_users()
        ldap = myLdap(LDAP_HOST, BASE_DN, USER, PASSWORD)
        ldap_users = ldap.sync_users()
        ldap_user_set = set(ldap_users.keys())
        local_user_set = set(user_list)
        return ldap_user_set - local_user_set,local_user_set - ldap_user_set,ldap_users

    def handler_data(self,log_obj):
        add_user_set, delete_user_set,ldap_users = self.compare_users()
        if delete_user_set:
            log_obj.info("Deleted users:%s",",".join(delete_user_set))
            delete_users = "("
            for user in delete_user_set:
                delete_users += "'%s'," % user
            delete_users = delete_users.rstrip(",")
            delete_users += ");"
            sql = "delete from users where username in %s" % delete_users
            self.__cursor.execute(sql)
        if add_user_set:
            log_obj.info("Added users:%s",",".join(add_user_set))
            sql2 = "insert into users(`username`,`givenname`,`email`) values"
            for user in add_user_set:
                sql2 += "('%s','%s','%s')," % (user,ldap_users.get(user).get("given_name"),ldap_users.get(user).get("email"))
            sql2 = sql2.rstrip(",")
            sql2 += ";"
            self.__cursor.execute(sql2)
        self.__cursor.close()
        self.__conn.close()

if __name__ == "__main__":
    obj = mysqlClient("localhost","root","","cmdb")
    log_obj = log.init_logger("../logs/sync.log")
    obj.handler_data(log_obj)


