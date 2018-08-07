package gostore

import (
	//"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	kdownLinkFront = "https://apkpure.com/"
	kdownLinkLast  = "/download"
)

func init() {
	//	refreshToken()
}

func TestToken(token Token) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Minute,
	)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gplaycli",
		"-s", "naver",
		"-ts", token.TokenStr,
		"-g", token.GsfStr,
	)
	_ = cmd.Run()
	time.Sleep(30 * time.Second)
	return
}

func GetApk(proxy, path string, ai AppInfo, token Token) error {
	err := downloadApk(proxy, path, ai.PkgName, token)
	if err != nil {
		return err
	}
	return nil

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

func downloadApk(proxy, path, pkgName string, token Token) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Minute,
	)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gplaycli",
		"-d", pkgName,
		"-f", path,
		"-dc", "gemini",
		"-ts", token.TokenStr,
		"-g", token.GsfStr,
		"-v",
	)
	//	var stderr bytes.Buffer
	//	cmd.Stderr = &stderr

	env := os.Environ()

	proxy_env := fmt.Sprintf("https_proxy=%s", proxy)
	fmt.Println(proxy_env)
	env = append(env, proxy_env)
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(pkgName, " : ", string(out))
		return err
	} else {
		fmt.Println(pkgName, " : ", string(out))
	}

	return nil
}
