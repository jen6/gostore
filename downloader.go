package gostore

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	kdownLinkFront = "https://apkpure.com/"
	kdownLinkLast  = "/download"
)

func GetApk(path string, ai AppInfo) error {
	tr := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	}

	var client = http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	downLink := buildDownLink(ai)
	req, err := http.NewRequest("GET", downLink, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	link, ok := parseDownLink(resp.Body)
	if !ok {
		return errors.New("Cant find DownLink")
	}

	if !isApk(link) {
		return errors.New("Not apk")
	}

	err = downloadApk(link, path, ai.PkgName+".apk")
	if err != nil {
		return err
	}

	return nil
}

func buildDownLink(ai AppInfo) string {
	downLink := kdownLinkFront
	downLink += ai.AppName
	downLink += "/"
	downLink += ai.PkgName
	downLink += kdownLinkLast
	return downLink
}

func parseDownLink(body io.Reader) (string, bool) {
	doc, err := gq.NewDocumentFromReader(body)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	sel := doc.Find("#download_link")
	if sel.Length() != 1 {
		return "", false
	}

	downLink, ok := sel.Attr("href")
	return downLink, ok
}

func isApk(url string) bool {
	splited := strings.Split(url, "/")
	if splited[4] == "apk" {
		return true
	} else {
		return false
	}
}

func downloadApk(url, path, filename string) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Minute,
	)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"curl", "-L",
		"-o", filepath.Join(path, filename),
		url)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
