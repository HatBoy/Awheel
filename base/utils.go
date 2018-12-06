package base

import (
	"path/filepath"
	"os"
	"log"
	"strings"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	mrand "math/rand"
	crand "crypto/rand"
	"math/big"
)

/*
一些用具函数
*/

// 获取当前执行文件目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
		panic(err)
		os.Exit(0)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 读取文件并返回文件行数组，widnows和Linux分隔符会有不一样的
func ReadLines(path string) []string {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		panic(err)
		os.Exit(0)
	}
	data := string(raw)
	lines := strings.Split(data, "\r\n")
	return lines
}

// URL处理，将URL分离为域名和协议
func UrlSplit(url string) (string, string, bool) {
	url = strings.TrimSpace(url)
	var proto, domain string
	if strings.HasPrefix(url, "http://") {
		proto = "http"
		domain = strings.Replace(url, "http://", "", -1)
		domain = strings.Replace(domain, "/", "", -1)
		return proto, domain, true
	} else if strings.HasPrefix(url, "https://") {
		proto = "https"
		domain = strings.Replace(domain, "/", "", -1)
		return proto, domain, true
	} else {
		return "", "", false
	}

}

func GetRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	seed, _ := crand.Int(crand.Reader, big.NewInt(1000000))
	r := mrand.New(mrand.NewSource(seed.Int64()))
	for i := 0; i < 16; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	ctx := md5.New()
	ctx.Write([]byte(string(result)))
	return hex.EncodeToString(ctx.Sum(nil))
}

