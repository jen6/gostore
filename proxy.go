package gostore

import (
	"bufio"
	"os"
	"sync"
)

type ProxyRotater struct {
	counter   int
	proxyLen  int
	proxyList []string
	mtx       sync.Mutex
}

var proxySetter *ProxyRotater

func GetProxyRotater() *ProxyRotater {
	if proxySetter == nil {
		proxySetter = &ProxyRotater{}
	}
	return proxySetter
}

func NewProxyRotater(proxyList string) (*ProxyRotater, error) {
	fp, err := os.Open(proxyList)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	ps := GetProxyRotater()
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		ps.proxyList = append(ps.proxyList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	ps.proxyLen = len(ps.proxyList)
	ps.counter = 0

	return ps, nil
}

func (ps *ProxyRotater) Next() (string, bool) {
	var result string
	rotated := false
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	if ps.counter == ps.proxyLen {
		rotated = true
		ps.counter = 0
	}
	result = ps.proxyList[ps.counter]
	ps.counter += 1
	return result, rotated
}
