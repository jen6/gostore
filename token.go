package gostore

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const tokenUrl = "https://matlink.fr/token/email/gsfid"

type Token struct {
	TokenStr string
	GsfStr   string
}

type TokenDispenser struct {
	tok Token
	mtx sync.Mutex
}

func (td *TokenDispenser) GetToken() Token {
	td.mtx.Lock()
	defer td.mtx.Unlock()
	return td.tok
}

func (td *TokenDispenser) RefreshToken() {
	td.mtx.Lock()
	defer td.mtx.Unlock()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	var client = http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	req, _ := http.NewRequest("GET", tokenUrl, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	tokens := strings.Fields(string(body))

	td.tok = Token{TokenStr: tokens[0], GsfStr: tokens[1]}
}
