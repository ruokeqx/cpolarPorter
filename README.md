
`cpolarPorter` fetch port map information from `cpolar` local server and update it to dns txt record periodically.

### Usage

```bash
DOMAIN_DOMAIN=xxx.com DOMAIN_RR=xxx ALICLOUD_ACCESS_KEY=xxx ALICLOUD_SECRET_KEY=xxx CPOLAR_URL=http://127.0.0.1:9200 CPOLAR_USERNAME=xxx CPOLAR_PASSWORD=xxx go run main.go port.go
```

```bash
~> dig txt xxx.xxx.com

; <<>> DiG 9.18.28-0ubuntu0.24.04.1-Ubuntu <<>> txt xxx.xxx.com
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 14309
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;xxx.xxx.com.              IN      TXT

;; ANSWER SECTION:
xxx.xxx.com.       600     IN      TXT     "{\"22\":\"xxxxx\",\"80\":\"xxxxx\"}"

;; Query time: 66 msec
;; SERVER: xxx.xxx.xxx.xxx#53(xxx.xxx.xxx.xxx) (UDP)
;; WHEN: Thu Oct 24 00:00:00 CST 2024
;; MSG SIZE  rcvd: 74
```

```bash
~> dig +short txt xxx.xxx.com
"{\"22\":\"xxxxx\",\"80\":\"xxxxx\"}"
```


### Tech Detail


```bash
curl 'http://127.0.0.1:9200/api/v1/user/login' \
  --data-raw '{"email":"xxx","password":"xxx"}' \
  --insecure
```

```bash
curl 'http://127.0.0.1:9200/api/v1/tunnels' \
  -H 'Authorization: Bearer xxx' \
  --insecure
```

### Credits

[go-acme/lego](https://github.com/go-acme/lego)
