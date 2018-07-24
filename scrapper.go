package gostore

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetNewAppsReader(sz int) (io.ReadCloser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var client = http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	reqForm := map[string]interface{}{
		"start":       0,
		"num":         sz,
		"numChildren": 0,
		"cctcss":      "sqare-cover",
		"cllayout":    "NORMAL",
		"ipf":         1,
		"xhr":         1,
	}

	jsonData, _ := json.Marshal(reqForm)
	buffer := bytes.NewBuffer(jsonData)

	req, err := http.NewRequest("POST", "https://play.google.com/store/apps/collection/topselling_new_free?hl=ko&authuser=0", buffer)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return resp.Body, nil
}

type AppInfo struct {
	AppName string
	PkgName string
}

func GetNewAppList(body io.Reader) []AppInfo {
	doc, err := gq.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
		return []AppInfo{}
	}

	var result []AppInfo
	sel := doc.Find("div .details > .title")

	sel.Each(func(i int, s *gq.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			return
		}
		pkgName := strings.Split(url, "=")[1]
		appName, ok := s.Attr("title")
		if !ok {
			return
		}

		result = append(result, AppInfo{AppName: appName, PkgName: pkgName})
	})

	return result
}
