package utils

import (
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego"
	"fmt"
	"errors"
)

func Connect()(redis.Conn,error) {
	host := beego.AppConfig.String("cache::redis_host")
	password := beego.AppConfig.String("cache::redis_password")
	c, err := redis.Dial("tcp", host)
	_,err2 := c.Do("AUTH", password)
	if err != nil {
		return nil,errors.New(fmt.Sprintf("Connect to redis error:%v",err))
	}
	if err2 != nil {
		return nil,errors.New("Wrong redis password")
	}
	return c,nil
}

// SetCache
//func SetCache(key string, value interface{}, timeout int) error {
//	data, err := Encode(value)
//	if err != nil {
//		return err
//	}
//	if cc == nil {
//		return errors.New("cc is nil")
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println(r)
//			cc = nil
//		}
//	}()
//	timeouts := time.Duration(timeout) * time.Second
//	err = cc.Put(key, data, timeouts)
//	if err != nil {
//		fmt.Println(err)
//		fmt.Println("SetCache失败，key:" + key)
//		return err
//	} else {
//		return nil
//	}
//}

func GetKey(conn redis.Conn,key string) (bool) {
	value,err := redis.String(conn.Do("GET",key))
	fmt.Printf("%v\n",value)
	defer conn.Close()
	if err == nil {
		return true
	} else {
		return false
	}
}

// DelCache
//func DelCache(key string) error {
//	if cc == nil {
//		return errors.New("cc is nil")
//	}
//	defer func() {
//		if r := recover(); r != nil {
//			//fmt.Println("get cache error caught: %v\n", r)
//			cc = nil
//		}
//	}()
//	err := cc.Delete(key)
//	if err != nil {
//		return errors.New("Cache删除失败")
//	} else {
//		return nil
//	}
//}

// Encode
// 用gob进行数据编码
//
//func Encode(data interface{}) ([]byte, error) {
//	buf := bytes.NewBuffer(nil)
//	enc := gob.NewEncoder(buf)
//	err := enc.Encode(data)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}

// Decode
// 用gob进行数据解码
//
//func Decode(data []byte, to interface{}) error {
//	buf := bytes.NewBuffer(data)
//	dec := gob.NewDecoder(buf)
//	return dec.Decode(to)
//}
