package mtr_go

import (
	"github.com/patrickmn/go-cache"
	"net/http"
)

/*

 */
type mtr struct {
	c *cache;
	gz http.Client;
	in string;
	arr bool;
	services []string;
	defGlue string;
	defSplitglue string;
	splitGlue string;
	glue string;
	ep ep;
	merge bool;
	matrix map[string]string;
        txtrq languageCode;
	srv []service;
	httpOpts []string;
	options map[string]string;
	target string;
	source string;
}


