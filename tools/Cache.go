package tools

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/patrickmn/go-cache"
)

type Ca struct {
	*cache.Cache
}

var NoExpiration = cache.NoExpiration

var Cache = Ca{}
var fLock sync.Mutex
var cache_path = "mtr_cache.raw"

func (c Ca) Init() {
	if _, err := os.Stat(cache_path); os.IsNotExist(err) {
		log.Print("Cache file not found, creating new one...")
		Cache = Ca{cache.New(-1, Seconds(30))}
		c.Save()
	} else {
		decodedCache := map[string]cache.Item{}
		b := new(bytes.Buffer)
		d := gob.NewDecoder(b)
		err = d.Decode(&decodedCache)
		if err != nil && err != io.EOF {
			log.Print("Corrupted cache file, creating new one...")
			log.Print(err)
			Cache = Ca{cache.New(-1, Seconds(30))}
			c.Save()
		}
		Cache = Ca{cache.NewFrom(-1, Seconds(30), decodedCache)}
		log.Print("Successfully loaded cache from file")
	}
	return
}

func (c Ca) Save() {
	fLock.Lock()
	defer fLock.Unlock()
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	e.Encode(Cache.Cache)
	err := ioutil.WriteFile(cache_path, b.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func (c Ca) ToBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	gob.Register(value)
	err := enc.Encode(&value) // dereference ifc for the encoder to encode type ifc
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c Ca) ToTypes(bts []byte) interface{} {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	var ho interface{}
	dec.Decode(&ho)
	return ho
}

func (c Ca) SetBytes(key string, value interface{}) error {
	bts, err := c.ToBytes(value)
	if err != nil {
		log.Print(err)
	}
	Cache.Set(key, bts, -1)
	return nil
}

func (c Ca) GetBytes(key string) interface{} {
	value, found := Cache.Get(key)
	if !found {
		log.Print("Key not found in cache")
	}
	return c.ToTypes(value.([]byte))
}
