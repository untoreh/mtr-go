package factory

import (
	"github.com/untoreh/mtr-go/i"
	"github.com/untoreh/mtr-go/services"
	"github.com/untoreh/mtr-go"
)

func Bing() (BingSrv i.Ep) {
	Bing := services.Ep(services.GEp)
	Bing.InitBing()
	return Bing
}

func Mtr(options interface{}) mtr_go.Mtr {
	mtr := &new(mtr_go.Mtr)
	mtr.AssignVariables(options)
	mtr.MakeServices()
	mtr.LangMatrix()
	return mtr
}
