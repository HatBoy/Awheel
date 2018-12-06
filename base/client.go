package base

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"encoding/json"
	"awheel/data"
)

//var client *redis.Client

//func init() {
//	client = redis.NewClient(&redis.Options{
//		Addr:     "192.168.100.8:6379",
//		Password: "123456",
//		DB:       0,
//	})
//	_, err := client.Ping().Result()
//	if err != nil {
//		log.Fatal(err)
//		os.Exit(1)
//	}
//}

/*
Redis连接客户端，及各种操作封装
*/

func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.100.8:6379",
		Password: "123456",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return client
}

// 将Struct数据转为JSON字符串
func Struct2Json(target interface{}) string {
	jstring, _ := json.Marshal(target)
	return string(jstring)
}

//将Target JSON字符串转化为struct
func Target2Struct(jstring string) *data.Target {
	target := new(data.Target)
	err := json.Unmarshal([]byte(jstring), &target)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return target
}

// baseinfoResult json 转化为 struct
func BaseInfoResult2Struct(jstring string) *data.BaseInfoResult {
	result := new(data.BaseInfoResult)
	err := json.Unmarshal([]byte(jstring), &result)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return result
}

func DirResult2Struct(jstring string) *data.DirResult {
	result := new(data.DirResult)
	err := json.Unmarshal([]byte(jstring), &result)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return result
}


//发送值到redis列表中RPUSH
func Send2Redis(client *redis.Client, key string, value interface{}) {
	jstring := Struct2Json(value)
	_, err := client.RPush(key, jstring).Result()
	if err != nil {
		log.Fatalln(err)
	}
}

//从redis列表中取出值
func GetTarget(client *redis.Client, key string) *data.Target {
	value, err := client.LPop(key).Result()
	if err != nil {
		return nil
	}
	target := Target2Struct(value)
	return target
}

//将值扔进Redis SET集合中去重
func Add2Redis(client *redis.Client, key, value string) {
	_, err := client.SAdd(key, value).Result()
	if err != nil {
		log.Fatalln(err)
	}
}

func IsExist(client *redis.Client, key, value string) bool {
	ok, err := client.SIsMember(key, value).Result()
	if err != nil {
		log.Fatalln(err)
	}
	return ok
}
