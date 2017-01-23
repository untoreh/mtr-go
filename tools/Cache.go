package tools

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = cache.New(-1, 30 * time.Second)
var NoExpiration = cache.NoExpiration
