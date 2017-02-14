package tools

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/levigross/grequests"
)

type UA struct {
	domain  string
	remote  string
	gh      string
	local   string
	agents  []string
	agentsL int
}

//func (ua *UA) downloadStrings() ([]byte) {
//	jar, _ := cookiejar.New(nil)
//	item := grequests.RequestOptions{
//		Headers: map[string]string{
//			"Host": "ua.theafh.net",
//			"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0",
//			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
//			"Accept-Language": "en-US,en;q=0.5",
//			"Accept-Encoding": "gzip, deflate",
//			"Referer": "http://ua.theafh.net/list.php?s=windows+chrome&include=yes&class=all&do=desc",
//			"Connection": "keep-alive",
//			"Upgrade-Insecure-Requests": "1",
//		},
//		UseCookieJar: true,
//		CookieJar: jar,
//	}
//	gen cookies
//	resp, _ := grequests.Get(ua.domain, &item)
//open url for download
//resp, _ = grequests.Get(ua.remote, &item)
//
//if err := resp.DownloadToFile(ua.local); err != nil {
//	log.Println("Unable to download file: ", err)
//}
//defer resp.RawResponse.Body.Close()
//}

func (ua *UA) downloadFromGithub() {
	var resp *grequests.Response
	var err error
get:
	for ret := 0; ret < 3; ret++ {
		resp, err = grequests.Get(ua.gh, nil)
		if err == nil {
			break get
		}

	}
	if err := resp.DownloadToFile(ua.local); err != nil {
		log.Println("Unable to download file: ", err)
	}
	defer resp.RawResponse.Body.Close()
}

func (ua *UA) Get() string {
	return ua.agents[rand.Int()%ua.agentsL]
}

func (ua *UA) parseStrings() []string {
	f, _ := os.Open(ua.local)
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		log.Print(err)
	}
	agents := make([]string, len(lines))
	for k := range lines {
		agents[k] = lines[k][0]
	}
	return agents
}

func (ua *UA) New() *UA {
	var err error

	ua.domain = "http://ua.theafh.net/"
	ua.remote = "http://ua.theafh.net/list.php?s=windows+chrome&include=yes&class=all&do=desc&download=true"
	ua.gh = "https://cdn.rawgit.com/untoreh/mtr-go/master/UAstrings.csv"
	ua.local = "UAstrings.csv"

	if _, err = ioutil.ReadFile(ua.local); err != nil {
		// Direct download gives forbidden ...
		//ua.downloadStrings()
		// use a list saved on github
		ua.downloadFromGithub()
	}
	ua.agents = ua.parseStrings()
	ua.agentsL = len(ua.agents)
	return ua
}
