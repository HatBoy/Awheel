package main

import (
	"flag"
	"github.com/go-redis/redis"
	"log"
	"os"
	"awheel/data"
	"awheel/base"
	"strings"
	"fmt"
	"awheel/info"
	"strconv"
)

/*
client端，负责提交任务，并查看任务执行结果
*/

func printString(value string, info string) {
	if value != "" {
		fmt.Println("[*] URL "+info, strings.TrimSpace(value))
	}
}

func main() {

	addr := flag.String("addr", "127.0.0.1:6379", "Redis ip:port")
	passwd := flag.String("passwd", "", "Redis password")
	db := flag.Int("db", 0, "Redis DB")

	url := flag.String("url", "", "Target url, eg:http://www.target.com")
	isSubdomian := flag.Bool("subdomain", true, "is search subdomain")
	isBruteSubdomain := flag.Bool("brudomain", false, "is brute subdomain")
	firstDomain := flag.String("firstdomain", "", "First Domain")
	isSubdir := flag.Bool("subdir", true, "is scan subdir")
	dirType := flag.String("dirtype", "dir", "dir type can use: php,asp,aspx,jsp,dir,mdb,auto")
	subdirType := flag.String("subdirtype", "dir", "dir type can use: php,asp,aspx,jsp,dir,mdb,auto")
	subdirRec := flag.Bool("subdirrec", false, "Subdomain iterative brute subdirectory")
	isSend := flag.Bool("send", false, "send targets")
	isshow := flag.Bool("show", true, "show result")

	clear := flag.Bool("clear", false, "clear the redis")

	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr:     *addr,
		Password: *passwd,
		DB:       *db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// 清空数据库
	if *clear {
		client.FlushAll()
	}

	// 提交任务
	if *isSend {

		var proto, domain string
		if strings.HasPrefix(*url, "http://") || strings.HasPrefix(*url, "https://") {
			urls := strings.Split(*url, "://")
			proto = urls[0]
			domain = strings.Replace(urls[1], "/", "", -1)
		} else {
			fmt.Println("URL error! url must have proto, eg:http://www.target.com")
			os.Exit(1)
		}

		if *firstDomain == "" {
			fmt.Println("must have firstdomain!")
			os.Exit(1)
		}

		target := new(data.Target)
		target.Proto = proto
		target.Domain = domain
		target.TargetId = base.GetRandomString()
		target.Master = true
		target.DomainId = base.GetRandomString()
		target.Types = "baseinfo"
		target.IsSubdomian = *isSubdomian
		target.IsBruteSubdomain = *isBruteSubdomain
		target.FirstDomain = *firstDomain
		target.IsDomainPan = info.DomainIsPan(*firstDomain)
		target.IsSubdir = *isSubdir
		target.Subdir = ""
		target.DirType = *dirType
		target.SubdirType = *subdirType
		target.IsSubdirRec = *subdirRec

		base.Send2Redis(client, "targets", target)
		fmt.Println("[*] 任务发送成功！")
	}

	// 显示结果
	if *isshow {
		subdomainTargets, _ := client.LLen("subdomain_targets").Result()
		subdirTargets, _ := client.LLen("subdir_targets").Result()

		if subdomainTargets == 0 && subdirTargets == 0{
			fmt.Println("[*] 运行进度: 子域名爆破和子目录爆破均运行完成!")
		}else if subdomainTargets == 0 && subdirTargets != 0{
			fmt.Println("[*] 运行进度: 子域名爆破运行完成，子目录爆破还剩请求", subdirTargets, "个！")
		}else if subdomainTargets != 0 && subdirTargets == 0 {
			fmt.Println("[*] 运行进度: 子域名爆破还剩请求", subdomainTargets,"个，子目录爆破运行完成！")
		}else if subdomainTargets != 0 && subdirTargets == 0{
			fmt.Println("[*] 运行进度: 子域名爆破还剩请求", subdomainTargets,"个，子目录爆破还剩请求", subdirTargets, "个！")
		}

		resultNum, _ := client.LLen("results").Result()
		resultDatas, _ := client.LRange("results", 0, resultNum).Result()

		subdirNum, _ := client.LLen("subdir_results").Result()
		subdirResults, _ := client.LRange("subdir_results", 0, subdirNum).Result()

		domainDir := make(map[string]string)

		for _, s := range subdirResults {
			result := base.DirResult2Struct(s)
			domainId := result.DomainId
			subdir := result.Subdir
			statusCode := result.StatusCode

			if v, ok := domainDir[domainId]; ok {
				domainDir[domainId] = v + ";" + subdir + " " + strconv.Itoa(statusCode)
			} else {
				domainDir[domainId] = subdir + " " + strconv.Itoa(statusCode)
			}
		}

		for _, r := range resultDatas {
			fmt.Println("-------------------------------------")
			result := base.BaseInfoResult2Struct(r)
			baseInfo := result.Baseinfo
			target := result.Task
			domainId := target.DomainId
			domain := baseInfo.Domain
			fmt.Println("[*] 域名: ", domain)
			isCDN := baseInfo.IsCDN
			if isCDN {
				ipList := baseInfo.IpList
				fmt.Println("[*] 疑似使用CDN")
				fmt.Println("[*] 域名对应IP列表: ", ipList)
			} else {
				ip := baseInfo.Ip
				fmt.Println("[*] 域名IP: ", ip)
				country := baseInfo.Country
				region := baseInfo.Region
				city := baseInfo.City
				fmt.Println("[*] IP所在地址: ", country, region, city)
				fmt.Println("[*] 域名注册商: ", baseInfo.Registrar)
				fmt.Println("[*] 域名注册人: ", baseInfo.Registrant)
				fmt.Println("[*] 域名联系邮箱: ", baseInfo.Emali)
				fmt.Println("[*] 域名联系电话: ", baseInfo.Phone)
				fmt.Println("[*] 域名创创建时间: ", baseInfo.CreateDate)
				fmt.Println("[*] 域名过期时间: ", baseInfo.ExpireDate)
				fmt.Println("[*] 域名服务器: ", baseInfo.DomainServer)
				dns := baseInfo.DNS
				fmt.Println("[*] 域名DNS列表: ")
				for _, d := range dns{
					fmt.Println("[+] ", d)
				}
				status := baseInfo.Status
				fmt.Println("[*] 域名状态: ")
				for _, s := range status{
					fmt.Println("[+] ", s)
				}
			}
			url := baseInfo.Url
			fmt.Println("[*] URL: ", url)
			cms := baseInfo.Cms
			waf := baseInfo.Waf
			printString(cms, "CMS类型: ")
			printString(waf, "WAF类型: ")
			title := baseInfo.Title
			printString(title, "Title: ")
			poweredby := baseInfo.PoweredBy
			printString(poweredby, "HTTP响应PoweredBy: ")
			server := baseInfo.Server
			printString(server, "HTTP响应Server: ")
			via := baseInfo.Via
			printString(via, "HTTP响应Via: ")
			robots := baseInfo.Robots
			printString(robots, "Robots文件: ")
			portinfos := baseInfo.PortInfos
			fmt.Println("[*] IP端口信息: ")
			for _, port := range portinfos {
				fmt.Println("[+] ", port.Service, port.Hostname, port.Version, port.Device, port.Os, port.Port, port.App, port.Extrainfo)
			}

			if v, ok := domainDir[domainId]; ok {
				fmt.Println("[*] 域名子目录: ")
				dirs := strings.Split(v, ";")
				for _, d := range dirs {
					fmt.Println("[+] ", url + d)
				}
			}
		}
	}

}