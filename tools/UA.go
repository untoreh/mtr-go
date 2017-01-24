package tools

import (
	"io/ioutil"
	"github.com/levigross/grequests"
	"log"
	"encoding/csv"
	"os"
	"net/http/cookiejar"
	"math/rand"
)

type UA struct {
	domain string
	remote string
	local  string
	agents []string
}

func (ua *UA) downloadStrings() ([]byte) {
	jar, _ := cookiejar.New(nil)
	item := grequests.RequestOptions{
		Headers: map[string]string{
			"Host": "ua.theafh.net",
			"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0",
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.5",
			"Accept-Encoding": "gzip, deflate",
			"Referer": "http://ua.theafh.net/list.php?s=windows+chrome&include=yes&class=all&do=desc",
			"Connection": "keep-alive",
			"Upgrade-Insecure-Requests": "1",
		},
		UseCookieJar: true,
		CookieJar: jar,
	}
	//gen cookies
	resp, _ := grequests.Get(ua.domain, &item)
	//open url for download
	resp, _ = grequests.Get(ua.remote, &item)

	if err := resp.DownloadToFile(ua.local); err != nil {
		log.Println("Unable to download file: ", err)
	}
	defer resp.ClearInternalBuffer()
	defer resp.RawResponse.Body.Close()

	return resp.Bytes()
}

func (ua *UA) Get() string {
	n := rand.Int() % len(ua.agents)
	return ua.agents[n]
}

func (ua *UA) parseStrings(data []byte) []string {
	f, _ := os.Open(ua.local)
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	agents := make([]string, len(lines))
	for k := range lines {
		agents[k] = lines[k][0]
	}
	return agents
}

func (ua *UA) New() *UA {
	var data []byte
	var err error

	ua.domain = "http://ua.theafh.net/"
	ua.remote = "http://ua.theafh.net/list.php?s=windows+chrome&include=yes&class=all&do=desc&download=true"
	ua.local = "UAstringss.csv"

	if data, err = ioutil.ReadFile(ua.local); err != nil {
		//data = ua.downloadStrings()
		log.Fatal("Strings not fetched, make sure the file is present")
	}
	ua.agents = ua.parseStrings(data)
	return ua
}