package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/ihaiker/gokit/files"
	"flag"
	"os/exec"
	"io"
	"bufio"
	"time"
)

const (
	HTTP_80_CONF = `upstream ${HTTP_DOMAIN_KEY}_80_host {
		server ${HTTP_DOMAIN_PROXIES}
	}
	server {
		listen  80;
		server_name ${HTTP_DOMAIN};
		location / {
			proxy_pass  http://${HTTP_DOMAIN_KEY}_80_host;
			proxy_redirect default;
			proxy_set_header Host $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_set_header REMOTE-HOST $remote_addr;
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		}
	}`
	HTTP_FILE_CONF = `server {
		listen  80;
		server_name ${HTTP_DOMAIN};
		location / {
			root ${HTTP_DOMAIN_PROXIES}
			index index.html index.htm;
			autoindex on;
		}
	}`

	HTTP_443_CONF = `
	upstream ${HTTP_DOMAIN_KEY}_443_host {
		server ${HTTP_DOMAIN_PROXIES}
	}
	server {
		listen  443;
		server_name ${HTTP_DOMAIN};
		location / {
			proxy_pass  http://${HTTP_DOMAIN_KEY}_443_host;
			proxy_redirect default;
			proxy_set_header Host $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_set_header REMOTE-HOST $remote_addr;
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		}
		ssl on;
		ssl_session_timeout 5m;
	    ssl_protocols SSLv2 SSLv3 TLSv1;
	    ssl_ciphers ALL:!ADH:!EXPORT56:RC4+RSA:+HIGH:+MEDIUM:+LOW:+SSLv2:+EXP;
	    ssl_prefer_server_ciphers on;
		ssl_certificate  vhosts/https/server.crt;
	    ssl_certificate_key  vhosts/https/server.key;
	}`

	MYSELF_CONF = "/usr/local/etc/nginx/vhosts/myself.conf"

	VHOSTS_JSON = "/usr/local/etc/nginx/vhosts/vhosts.json"
)

func getConf(httpPort, domain string, proxies ...interface{}) string {
	domainKey := strings.Replace(domain, ".", "_", -1)

	template := HTTP_80_CONF
	if httpPort == "443" {
		template = HTTP_443_CONF
	} else if fileKit.IsDir(proxies[0].(string)) {
		template = HTTP_FILE_CONF
	}

	out := strings.Replace(template, "${HTTP_DOMAIN}", domain, -1)
	out = strings.Replace(out, "${HTTP_DOMAIN_KEY}", domainKey, -1)
	proxyAddress := ""
	for c, proxy := range proxies {
		proxyAddress = proxyAddress + proxy.(string) + ";"
		if c < len(proxies) - 1 {
			proxyAddress = proxyAddress + "\n"
		}
	}
	out = strings.Replace(out, "${HTTP_DOMAIN_PROXIES}", proxyAddress, -1)
	return "\n#-------------------- " + domain + " --------------\n" + out
}

var conf = flag.String("c", VHOSTS_JSON, "the hosts config file")

func rewriteConf() {

	fmt.Println(*conf)

	if f, err := os.Open(*conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if body, err := ioutil.ReadAll(f); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		jsonObj := new(map[string]interface{})
		if err := json.Unmarshal(body, jsonObj); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			nginxConf := ""
			for httpPort, val := range *jsonObj {
				httpConf := val.(map[string]interface{})
				for domain, serverNames := range httpConf {
					switch serverNames.(type) {
					case string:
						nginxConf += getConf(httpPort, domain, serverNames)
					case []interface{}:
						nginxConf += getConf(httpPort, domain, serverNames.([]interface{})...)
					}
				}
			}
			os.Remove(MYSELF_CONF)
			if wf, err := os.OpenFile(MYSELF_CONF, (os.O_CREATE | os.O_RDWR), os.ModePerm & (^os.ModeAppend)); err != nil {
				fmt.Println(err)
			} else {
				wf.WriteString(nginxConf)
				wf.Close()
			}
		}
	}
}

func readOut(stop chan int, reader io.Reader) {
	r := bufio.NewReader(reader)
	for {
		line, _, err := r.ReadLine()
		fmt.Println(string(line))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func reload() {
	stop := make(chan int)
	reader, writer, _ := os.Pipe()
	go readOut(stop, reader)
	writer.WriteString("start\n")

	cmd := exec.Command("nginx", "-t")
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); err != nil {
		time.Sleep(time.Second)
		fmt.Println(err)
		os.Exit(1)
	}

	cmd = exec.Command("nginx", "-s", "reload")
	cmd.Stdout = writer
	cmd.Stderr = writer
	if err := cmd.Run(); err != nil {
		time.Sleep(time.Second)
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("OK")
	}
}

func main() {
	rewriteConf()
	reload()
}
