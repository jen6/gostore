package gostore

import (
	"crypto/tls"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const DefaultFetchSize = 60

func GetNewAppsReader(start, sz int) (io.ReadCloser, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var client = http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	reqForm := url.Values{}
	reqForm.Add("start", strconv.Itoa(start))
	reqForm.Add("num", strconv.Itoa(DefaultFetchSize))
	reqForm.Add("numChildren", "0")
	reqForm.Add("cctcss", "square-cover")
	reqForm.Add("cllayout", "NORMAL")
	reqForm.Add("ipf", "1")
	reqForm.Add("xhr", "1")

	req, err := http.NewRequest("POST", "https://play.google.com/store/apps/collection/topselling_new_free?auth_user=0", strings.NewReader(reqForm.Encode()))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	//	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	fmt.Println(reqForm.Encode())
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

func GetNewAppList(body io.Reader, sz int) []AppInfo {
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

	if sz < len(result) {
		result = result[:sz]
	}

	return result
}
