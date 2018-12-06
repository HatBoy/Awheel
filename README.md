# Awheel：分布式信息收集工具，golang编写的又一个轮子

## 工具介绍
+ 在平时对web的渗透测试过程中，基础信息收集是第一步，有些域名子域名太多，对每个子域名进行信息收集太麻烦，而对于子域名、子目录的爆破单机运行速度太慢，因此结合上述需求实现了一个新的轮子，分布式信息收集工具，输入域名后会自动调用接口并结合字典爆破子域名，然后对每一各子域名信息进行收集，同时对每一各子域名进行目录爆破。后期可能添加更多的功能。

## 功能介绍
+ 域名WHOIS信息收集：调用 http://whois.chaicp.com/ 接口实现
+ 域名解析IP并判断是否使用CDN：调用超级ping接口 https://ping.aizhan.com/ 实现
+ IP归属地信息查询：调用 http://ip.taobao.com 接口实现
+ IP开放端口信息查询：调用zoomeye接口实现
+ WEB服务信息收集，如title、server、poweredby：通过HTTP响应头获取
+ CMS识别：调用 http://whatweb.bugscaner.com/look/、https://whatcms.org/、http://www.iguoli.cn/cms.php 三个接口综合实现
+ WAF识别：使用 https://github.com/Ekultek/WhatWaf、https://github.com/EnableSecurity/wafw00f 开源程序的策略，通过正则实现
+ 子域名获取：调用 https://www.virustotal.com/ui/domains/ 接口结合字典爆破获取子域名
+ 子目录爆破：内置字典爆破子目录

## 基本架构
+ Redis作为分布式任务分法服务器，server端运行后监控Redis中是否有任务，有则执行并将结果发送回Redis中保存，没有则等待，client负责发送任务到Redis，可通过client端查看任务执行进度和执行结果
+ server端口执行任务的worker分为两种，一种是执行信息收集的worker，默认10个，一种是执行子域名爆破和子目录爆破的worker，默认为90个。

## 安装编译
+ 本程序用golang语言编写，可跨平台运行，测试编写版本为1.11.2，低版本可能会出现第三方库兼容问题
+ 安装好golang，并设置好GOROOT、GOPATH通过下面步骤编译程序
+ 第一步：安装第三方库
```
HTML解析包:go get github.com/PuerkitoBio/goquery
DNS解析包：go get github.com/miekg/dns
JSON处理:go get github.com/bitly/go-simplejson
Redis客户端：go get github.com/go-redis/redis
```
+ 第二步：编译server端和client端
```
go build -i -o server server.go
go build -i -o client client.go
```

## 参数使用
### 参数说明
+ server端参数说明：
    + -addr：Redis地址端口 (default "127.0.0.1:6379")
    + -passwd：Redis认证密码
    + -db：Redis数据库
    + -mworker：信息收集任务worker数量 (default 10)
    + -sworker：爆破任务worker数量 (default 90)
+ client端参数说明：
    + -addr：Redis地址端口 (default "127.0.0.1:6379")
    + -passwd：Redis认证密码
    + -db：Redis数据库
    + -send：发送任务到Redis
    + -show：显示运行结果 (default true)
    + -url：目标URL，例如:http://www.target.com
    + -subdomain：是否查询子域名 (default true)
    + -brudomain：是否进行子域名爆破
    + -firstdomain：一级域名，如：target.com
    + -dirtype：目标域名使用什么字典进行子目录爆破: php,asp,aspx,jsp,dir,mdb,auto (default "dir")
    + -subdirtype：目标子域名使用什么字典进行子目录爆破: php,asp,aspx,jsp,dir,mdb,auto (default "dir")
    + -subdir：是否进行子目录爆破 (default true)
    + -subdirrec：是否对目标子域名也进行子目录爆破
    + -clear：清空Redis全部数据

### 运行方式
+ server可在任何地方运行多个：
```
server -addr 192.168.100.8:6379 -passwd 123456 -db 0 -mworker 10 -sworker 100
```
+ 注意：server端必须和dicts目录放在同一个目录运行，否则无法找到爆破字典会报错。server目前不太稳定，容易崩溃，可通过supervisor等工具监控，实现崩溃重启
+ client分为发送任何和显示运行结果两个功能
+ client发送任务：
```
client -addr 192.168.100.8:6379 -passwd 123456 -db 0 -url http://www.target.com -subdomain=true -brudomain=true -brudomain=true -firstdomain target.com -dirtype php,asp,jsp -subdirtype auto -subdir=true -subdirrec=true -send=true
```
+ client显示运行结果：
```
client -addr 192.168.100.8:6379 -passwd 123456 -db 0 -show=true
```
+ 默认为显示结果，可不加-show参数
+ 运行结果展示显示界面如下：
```
[*] 运行进度: 子域名爆破运行完成，子目录爆破还剩请求 377371 个！
-------------------------------------
[*] 域名:  www.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://www.freebuf.com
[*] URL Title:  FreeBuf互联网安全新媒体平台
[*] URL HTTP响应Server:  Apache/2.2.21
[*] URL Robots文件:  User-agent: *
Disallow: /*?*
Disallow: /trackback
Disallow: /wp-*/
Disallow: */comment-page-*
Disallow: /*?replytocom=*
Disallow: */trackback
Disallow: /?random
Disallow: */feed
Disallow: /*.css$
Disallow: /*.js$
Sitemap: http://www.freebuf.com/sitemap.txt
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  static.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  http://static.freebuf.com
[*] URL HTTP响应Server:  Apache/2.2.21
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  my.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  http://my.freebuf.com
[*] URL Title:  登录中心
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  nginx
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  zhuanlan.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  http://zhuanlan.freebuf.com
[*] URL Title:  专栏 - FreeBuf 互联网安全新媒体平台 | 关注黑客与极客
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  wit.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://wit.freebuf.com
[*] URL Title:  WitAwards 2018互联网安全年度评选
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
[*] 域名子目录:
[+]  https://wit.freebuf.com/index.html 200
-------------------------------------
[*] 域名:  open.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  ://open.freebuf.com
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  job.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  http://job.freebuf.com
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  bar.freebuf.com
[*] 疑似使用CDN
[*] 域名对应IP列表:  []
[*] URL:  ://bar.freebuf.com
[*] IP端口信息:
-------------------------------------
[*] 域名:  shop.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://shop.freebuf.com
[*] URL Title:  商城 | FreeBuf.COM | 极客生活，创意人生
[*] URL HTTP响应Server:  Apache/2.2.21
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  search.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://search.freebuf.com
[*] URL Title:  搜索 - FreeBuf 互联网安全新媒体平台 | 关注黑客与极客
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  live.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://live.freebuf.com
[*] URL Title:  公开课 - FreeBuf公开课
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  fit.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://fit.freebuf.com
[*] URL Title:  FIT 2019 互联网安全创新大会
[*] URL HTTP响应Server:  Apache/2.2.21
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  company.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://company.freebuf.com
[*] URL Title:  企业空间 - FreeBuf 互联网安全新媒体平台 | 关注黑客与极客
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  freetalk.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  http://freetalk.freebuf.com
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  push.freebuf.com
[*] 疑似使用CDN
[*] 域名对应IP列表:  []
[*] URL:  ://push.freebuf.com
[*] IP端口信息:
-------------------------------------
[*] 域名:  ai.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://ai.freebuf.com
[*] URL Title:  专栏 - FreeBuf 互联网安全新媒体平台 | 关注黑客与极客
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  prize.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://prize.freebuf.com
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  Tengine
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  hackpwn.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  ://hackpwn.freebuf.com
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  geekpwn.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  ://geekpwn.freebuf.com
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  api.freebuf.com
[*] 域名IP:  123.57.248.123
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  ://api.freebuf.com
[*] IP端口信息:
[+]  ssh  5.3   2222 OpenSSH protocol 2.0
[+]  http  2.2.21   80 Apache httpd
-------------------------------------
[*] 域名:  pay.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  ://pay.freebuf.com
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
[*] 域名:  user.freebuf.com
[*] 域名IP:  101.200.172.135
[*] IP所在地址:  中国 北京 北京
[*] 域名注册商:  GoDaddy.com, LLC
[*] 域名注册人:  GoDaddy.com, LLC
[*] 域名联系邮箱:  abuse@godaddy.com
[*] 域名联系电话:  480-624-2505
[*] 域名创创建时间:  2010-08-21T15:24:37Z
[*] 域名过期时间:  2019-08-21T15:24:37Z
[*] 域名服务器:  whois.godaddy.com
[*] 域名DNS列表:
[+]  F1G1NS1.DNSPOD.NET
[+]  F1G1NS2.DNSPOD.NET
[*] 域名状态:
[+]  clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited
[+]  clientRenewProhibited https://icann.org/epp#clientRenewProhibited
[+]  clientTransferProhibited https://icann.org/epp#clientTransferProhibited
[+]  clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited
[*] URL:  https://user.freebuf.com
[*] URL Title:  登录中心
[*] URL HTTP响应PoweredBy:  PHP/5.6.22
[*] URL HTTP响应Server:  nginx
[*] IP端口信息:
[+]  http     443 Tengine httpd
[+]  http     80 Tengine httpd
-------------------------------------
```