package tools

import (
	"github.com/patrickmn/go-cache"
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type Ca struct{
	*cache.Cache
}
var Cache = Ca{cache.New(-1, 30 * time.Second)}
var NoExpiration = cache.NoExpiration

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

func (c Ca) ToTypes(bts []byte) (interface{}) {
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
	if !found  {
		log.Print("Key not found in cache")
	}
	return c.ToTypes(value.([]byte))
}
