package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ZeaLoVe/alarm/g"
)

type ResponseHost struct {
	Ip       string `json:"ip,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

//endpoint-ip  kv
type SafeHostMap struct {
	sync.Mutex
	M map[string]string
}

var HostMap = &SafeHostMap{M: make(map[string]string)}

func (this *SafeHostMap) GetIP(endpoint string) string {
	this.Lock()
	defer this.Lock()
	return this.M[endpoint]
}

func (this *SafeHostMap) updateMap(hosts []ResponseHost) {
	this.Lock()
	defer this.Unlock()

	for _, host := range hosts {
		this.M[host.Endpoint] = host.Ip
	}
}

func (this *SafeHostMap) GetHosts(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	} else {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}
	var hosts_tmp []ResponseHost
	err = json.Unmarshal(body, &hosts_tmp)
	if err != nil {
		log.Println(err.Error())
		return
	}

	this.updateMap(hosts_tmp)
}

func (this *SafeHostMap) GetIPByEndpoint(endpoint string) string {
	this.Lock()
	defer this.Unlock()
	return this.M[endpoint]
}

func UpdateHostMap() {
	if g.Config().Hosts == nil {
		return
	}
	url := g.Config().Hosts.Api_url
	if url == "" {
		return
	}
	interval := g.Config().Hosts.Interval

	for {
		HostMap.GetHosts(url)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

//该函数可以搜索content，并从中提取 Endpoint: 后面跟的内容。判断依据是非数字字母下划线中划线
func GetEndpoint(content string) string {
	if !strings.Contains(content, "Endpoint:") {
		return ""
	}
	var index, end int
	index = strings.Index(content, "Endpoint:") + 9
	for i := index; i < len(content); i++ {
		if content[i] >= 'A' && content[i] <= 'Z' {
			continue
		} else if content[i] >= 'a' && content[i] <= 'z' {
			continue
		} else if content[i] >= '0' && content[i] <= '9' {
			continue
		} else if content[i] == '-' {
			continue
		} else if content[i] == '_' {
			continue
		} else {
			end = i
			break
		}
	}
	return content[index:end]
}

func TransContent(content string) string {
	ip := HostMap.GetIPByEndpoint(GetEndpoint(content))
	if strings.EqualFold(ip, "") {
		return fmt.Sprint(content)
	} else {
		return fmt.Sprintf("%s[ip:%s]", content, HostMap.GetIPByEndpoint(GetEndpoint(content)))
	}
}
