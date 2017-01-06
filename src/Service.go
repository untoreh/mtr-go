package mtr_go

type service interface {
	translate() []string
	genReq()
	preReq()
	getLangs() []string
}
