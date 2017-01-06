package mtr_go

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type stuff struct {
	misc string;
}

type ep struct {
	urls   map[string]string;
	params map[string]string;
	misc   map[string]stuff;
	active bool;
	mtr    mtr;
	s      service;
}

func (ep *ep) epInit() {
	ep.mtr = cache.New(false, 30 * time.Second);
}

func (ep *ep) epDefaults() {
	ua, found := ep.mtr.c.Get("mtr_ua_rnd");
	if (found != true) {
		ua = "user agent string";
		ep.mtr.c.Set("mtr_ua_rnd", ua, cache.NoExpiration);
	}
	ep.params["default"] = map[string]int{
		"headers": map[string]string{
			"User-Agent" : *ua.(string),
		},
		"query": []string{},
	};
	ep.misc["glue"] = ep.mtr.glue;
	ep.misc["splitGlue"] = ep.mtr.splitGlue;

	ep.active = true;
}

func translate() {

}

