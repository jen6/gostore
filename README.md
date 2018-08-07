# gostore
playstore apk crawler

[Release](https://oss.navercorp.com/gungun-son/gostore/releases)에서 다운로드

## install
1. install gplaycli with pip3   
`pip3 install gplaycli`  

2. set config for gplaycli
```
[irteamsu@dev-gun-chome2-ncl ~]$ cat ~/.config/gplaycli/gplaycli.conf 
[Credentials]
gmail_address=aa
gmail_password=aa

#keyring_service=gplaycli
token=True
token_url=https://matlink.fr/token/email/gsfid

[Cache]
token=~/.cache/gplaycli/token

[Locale]
#locale=en_GB
locale=ko_KR
timezone=CEST
```

3. 프록시서버들 구축
[tiny proxy](https://tinyproxy.github.io/)
구글의 ip밴을 피하기 위해서 여러개의 유동적인 ip사용

4. config.json 만들기
config.json.bak을 참고해서  config.json을 제작 
proxy-list.txt도 함께

5. 실행
`./gostore -conf='./config.json'`
