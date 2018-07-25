package main

import (
	"flag"
	"fmt"
	gs "oss.navercorp.com/gungun-son/gostore"
	"sync/atomic"
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

	cnt := 0
	var appDownloaded int32
	appDownloaded = 0

	for ; conf.CrawlSize > 0; conf.CrawlSize -= gs.DefaultFetchSize {
		body, err := gs.GetNewAppsReader(cnt*gs.DefaultFetchSize, conf.CrawlSize)
		if err != nil {
			fmt.Println(err)
			return
		}

		apps := gs.GetNewAppList(body, conf.CrawlSize)
		fmt.Printf("app data crawl fin! : %d\n", len(apps))

		workChan := make(chan interface{}, 10)
		for _, app := range apps {
			go func(ap gs.AppInfo) {
				fmt.Println(ap.AppName, " is now downloading")
				err = gs.GetApk(conf.SavePath, ap)
				if err != nil {
					fmt.Println(ap.AppName, " : ", err)
				} else {
					atomic.AddInt32(&appDownloaded, 1)
				}
				<-workChan
			}(app)
			workChan <- nil
		}
		body.Close()
		cnt += 1
	}

	fmt.Printf("Total %d apps Crawled\n", appDownloaded)
}
