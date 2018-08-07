package main

import (
	"flag"
	"fmt"
	gs "oss.navercorp.com/gungun-son/gostore"
	"time"
)

func main() {
	flagPath := flag.String("conf", "", "config file path")
	flag.Parse()

	if *flagPath == "" {
		flag.PrintDefaults()
		return
	}

	conf, err := gs.GetConfig(*flagPath)
	if err != nil {
		fmt.Println("Error in GetConfig : ", err)
		return
	}

	var tokenDispenser gs.TokenDispenser
	tokenDispenser.RefreshToken()

	proxyRotate, err := gs.NewProxyRotater(conf.ProxyListPath)
	if err != nil {
		fmt.Println("Error in NewProxyRotater : ", err)
		return
	}

	apkManager, err := gs.NewApkManager(conf.SavePath)
	if err != nil {
		fmt.Println("error in apkmanager : ", err)
		return
	}
	defer apkManager.Close()

	appDownloaded := 0
	for _, link := range conf.DownloadLinks {
		cnt := 0
		crawlSize := conf.CrawlSize
		for ; conf.CrawlSize > 0; crawlSize -= gs.DefaultFetchSize {
			body, err := gs.GetNewAppsReader(link, cnt*gs.DefaultFetchSize, conf.CrawlSize)
			if err != nil {
				fmt.Println(err)
				return
			}

			apps := gs.GetNewAppList(body, conf.CrawlSize)
			fmt.Printf("app data crawl fin! : %d\n", len(apps))
			body.Close()
			if len(apps) == 0 {
				break
			}

			for _, app := range apps {
				time.Sleep(15 * time.Second)
				fmt.Println(app.PkgName, " is now downloading")
				proxy, isRotate := proxyRotate.Next()
				if isRotate {
					tokenDispenser.RefreshToken()
					gs.TestToken(tokenDispenser.GetToken())
				}

				newApkinfo, err := gs.SearchApkInfo(proxy, app.PkgName, tokenDispenser.GetToken())
				if err != nil {
					fmt.Println("err in Search Apk : ", err)
					continue
				}
				apkinfo, ok := apkManager.ApkChecker[app.PkgName]
				if ok {
					isNew := apkinfo.LastUpdate.Before(newApkinfo.LastUpdate)
					if !isNew {
						fmt.Println(app.PkgName, " is already outdate")
						continue
					}
					fmt.Println(app.PkgName, " is going update")
				}

				time.Sleep(15 * time.Second)
				err = gs.GetApk(proxy, conf.SavePath, app, tokenDispenser.GetToken())
				if err != nil {
					fmt.Println(app.AppName, " : ", err)
				} else {
					appDownloaded += 1
					if ok {
						apkManager.UpdateApk(newApkinfo)
					} else {
						apkManager.NewApk(newApkinfo)
					}
				}
			}
			cnt += 1
		}
	}

	fmt.Printf("Total %d apps Crawled\n", appDownloaded)
}
