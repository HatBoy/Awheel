package main

import (
	"awheel/base"
	"awheel/info"
	"awheel/data"
	"strings"
	"time"
	"sync"
	"flag"
	"github.com/go-redis/redis"
	"log"
	"os"
)

/*
Server端，运行在服务器上，负责执行任务
*/

// 主要的worker，主要用于基础信息获取和分配爆破任务
func mainWorker(client *redis.Client) {
	for {
		target := base.GetTarget(client, "targets")
		if target == nil {
			time.Sleep(time.Second * 3)
			continue
		}

		// 正式运行时替换路径，下面的路径仅供测试使用
		dir := base.GetCurrentDirectory()
		path := dir + "/dicts/"
		//fmt.Println(path)

		//path := "D:/Code/GoCode/src/awheel/dicts/"

		targetId := target.TargetId
		domain := target.Domain
		proto := target.Proto
		domainId := target.DomainId
		master := target.Master
		isSubdmoain := target.IsSubdomian
		isDomainPan := target.IsDomainPan
		isBruteSubdomain := target.IsBruteSubdomain
		firstDomain := target.FirstDomain
		isSubdir := target.IsSubdir
		dirType := target.DirType
		subdirType := target.SubdirType
		isSubdirRec := target.IsSubdirRec

		base.Add2Redis(client, "subdomains", domain)
		//执行基本信息获取任务并将结果存入Redis
		if proto == "" {
			proto = base.HttpOrHttps(domain)
		}
		url := proto + "://" + domain
		isAlive := base.WebIsAlive(url)
		target.IsAlive = isAlive

		urlBaseInfo := info.BaseInfo(domain, url, isAlive)

		baseinfoResult := new(data.BaseInfoResult)
		baseinfoResult.Baseinfo = *urlBaseInfo
		baseinfoResult.Task = *target

		base.Send2Redis(client, "results", baseinfoResult)

		//通过接口获取子域名
		if isSubdmoain {
			subdomains := info.SubDomain(firstDomain)
			for _, sub := range subdomains {
				subtarget := new(data.Target)
				subtarget.TargetId = targetId
				subtarget.Proto = base.HttpOrHttps(sub)
				subtarget.Domain = sub
				subtarget.Master = false
				subtarget.DomainId = base.GetRandomString()
				subtarget.Types = "baseinfo"
				subtarget.IsSubdomian = false
				subtarget.IsDomainPan = isDomainPan
				subtarget.IsBruteSubdomain = false
				subtarget.FirstDomain = firstDomain
				subtarget.DirType = dirType
				subtarget.SubdirType = subdirType
				if !isSubdirRec {
					subtarget.IsSubdir = false
					subtarget.IsSubdirRec = false
				} else {
					subtarget.IsSubdir = isSubdir
					subtarget.IsSubdirRec = isSubdirRec
				}

				if !base.IsExist(client, "subdomains", sub) {
					base.Add2Redis(client, "subdomains", sub)
					base.Send2Redis(client, "targets", subtarget)
				}

			}
		}

		//判断是否进行子域名爆破，且不是域名泛解析
		if master && isSubdmoain && isBruteSubdomain && !isDomainPan {
			subFile := path + "subdomains.txt"
			lines := base.ReadLines(subFile)
			for _, line := range lines {
				line = strings.TrimSpace(line)
				subd := line + "." + firstDomain
				if base.IsExist(client, "subdomains", subd) {
					continue
				}

				subbruteTarget := new(data.Target)
				subbruteTarget.TargetId = targetId
				subbruteTarget.Proto = ""
				subbruteTarget.Domain = subd
				subbruteTarget.Master = false
				subbruteTarget.DomainId = ""
				subbruteTarget.Types = "subdomain"
				subbruteTarget.IsSubdomian = false
				subbruteTarget.IsDomainPan = false
				subbruteTarget.IsBruteSubdomain = false
				subbruteTarget.FirstDomain = firstDomain
				subbruteTarget.IsSubdir = isSubdir
				subbruteTarget.SubdirType = subdirType
				subbruteTarget.IsSubdirRec = isSubdirRec

				base.Send2Redis(client, "subdomain_targets", subbruteTarget)
			}
		}

		//子目录爆破
		if isSubdir && isAlive {
			bruteType := "auto"
			if master {
				bruteType = dirType
			} else {
				bruteType = subdirType
			}
			bruteTypeList := strings.Split(bruteType, ",")
			var dirs []string
			for _, t := range bruteTypeList {
				if t == "auto" {
					lines := base.ReadLines(path + "dir.txt")
					dirs = append(dirs, lines...)

					testt := base.TestWebType(proto + "://" + domain)
					if testt != "" {
						lines := base.ReadLines(path + testt + ".txt")
						dirs = append(dirs, lines...)
					}
					break
				} else {
					lines := base.ReadLines(path + t + ".txt")
					dirs = append(dirs, lines...)
				}

			}

			for _, dir := range dirs {
				subdirTarget := new(data.Target)
				subdirTarget.TargetId = targetId
				subdirTarget.Proto = proto
				subdirTarget.Domain = domain
				subdirTarget.Master = false
				subdirTarget.DomainId = domainId
				subdirTarget.Types = "subdir"
				subdirTarget.IsSubdomian = false
				subdirTarget.IsDomainPan = false
				subdirTarget.IsBruteSubdomain = false
				subdirTarget.FirstDomain = firstDomain
				subdirTarget.IsSubdir = isSubdir
				subdirTarget.Subdir = dir
				subdirTarget.SubdirType = subdirType
				subdirTarget.IsSubdirRec = isSubdirRec

				base.Send2Redis(client, "subdir_targets", subdirTarget)
			}

		}
	}

}

// 主要用于各种爆破，如子域名、子目录爆破
func subWorker(client *redis.Client) {
	for {
		target := base.GetTarget(client, "subdomain_targets")
		if target == nil {
			target = base.GetTarget(client, "subdir_targets")
			if target == nil {
				time.Sleep(time.Second * 3)
				continue
			}
		}

		targetId := target.TargetId
		domain := target.Domain
		proto := target.Proto
		domainId := target.DomainId
		types := target.Types
		master := target.Master
		firstDomain := target.FirstDomain
		isSubdir := target.IsSubdir
		dirType := target.DirType
		subDir := target.Subdir
		subdirType := target.SubdirType
		isSubdirRec := target.IsSubdirRec

		if types == "subdomain" {
			ok, err := info.Nslookup(domain)
			if err != nil {

			}
			//成功，存在该子域名
			if ok && err == nil {
				subbruteTarget := new(data.Target)
				subbruteTarget.TargetId = targetId
				subbruteTarget.Proto = ""
				subbruteTarget.Domain = domain
				subbruteTarget.Master = false
				subbruteTarget.DomainId = base.GetRandomString()
				subbruteTarget.Types = "baseinfo"
				subbruteTarget.IsSubdomian = false
				subbruteTarget.IsDomainPan = false
				subbruteTarget.IsBruteSubdomain = false
				subbruteTarget.FirstDomain = firstDomain
				subbruteTarget.IsSubdir = isSubdir
				subbruteTarget.SubdirType = subdirType
				subbruteTarget.IsSubdirRec = isSubdirRec

				base.Send2Redis(client, "targets", subbruteTarget)
			}
		} else if types == "subdir" {
			url := proto + "://" + domain + subDir
			ok, statusCode := info.SubDirCheck(url)
			if ok {
				dirResult := new(data.DirResult)
				dirResult.TargetId = targetId
				dirResult.FirstDomain = firstDomain
				dirResult.Proto = proto
				dirResult.Domain = domain
				dirResult.DomainId = domainId
				dirResult.Subdir = subDir
				dirResult.Master = master
				dirResult.DirType = dirType
				dirResult.SubdirType = subdirType
				dirResult.IsSubdirRec = isSubdirRec
				dirResult.StatusCode = statusCode

				base.Send2Redis(client, "subdir_results", dirResult)
			}
		}
	}
}

func main() {

	addr := flag.String("addr", "127.0.0.1:6379", "Redis ip:port")
	passwd := flag.String("passwd", "", "Redis password")
	db := flag.Int("db", 0, "Redis DB")
	mworkers := flag.Int("mworker", 10, "Main workers nums")
	sworkers := flag.Int("sworker", 90, "Sub workers nums")

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

	target := new(data.Target)
	target.Proto = "https"
	target.Domain = "www.freebuf.com"
	target.TargetId = base.GetRandomString()
	target.Master = true
	target.DomainId = base.GetRandomString()
	target.Types = "baseinfo"
	target.IsSubdomian = true
	target.IsBruteSubdomain = true
	target.IsDomainPan = false
	target.FirstDomain = "freebuf.com"
	target.IsSubdir = true
	target.Subdir = ""
	target.DirType = "php,dir"
	target.SubdirType = "auto"
	target.IsSubdirRec = true

	//base.Send2Redis(client,"targets", target)

	var wg sync.WaitGroup
	wg.Add(*mworkers + *sworkers)

	for i := 0; i < *mworkers; i++ {
		go mainWorker(client)
	}

	for j := 0; j < *sworkers; j++ {
		go subWorker(client)
	}

	wg.Wait()

}
