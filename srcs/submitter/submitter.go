package submitter

import "sub/config"

type Response struct {
	Msg    string `json:"msg"`
	Flag   string `json:"flag"`
	Status bool   `json:"status"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SubmitterInterface interface {
	submitFlags(flags []string) ([]Response, error)
	Conf() *config.Config
	SubAccepted() string
	SubInvalid() string
	SubOld() string
	SubYourOwn() string
	SubStolen() string
	SubNop() string
	SubNotAvailable() string
	SubServiceDown() string
}

type Submitter struct {
	conf            *config.Config
	subAccepted     string
	subInvalid      string
	subOld          string
	subYourOwn      string
	subStolen       string
	subNop          string
	subNotAvailable string
	subServiceDown  string
}

const NO_SUB = "##########"

func (s *Submitter) Conf() *config.Config {
	return s.conf
}

func (s *Submitter) SubAccepted() string {
	return s.subAccepted
}

func (s *Submitter) SubInvalid() string {
	return s.subInvalid
}

func (s *Submitter) SubOld() string {
	return s.subOld
}

func (s *Submitter) SubYourOwn() string {
	return s.subYourOwn
}

func (s *Submitter) SubStolen() string {
	return s.subStolen
}

func (s *Submitter) SubNop() string {
	return s.subNop
}

func (s *Submitter) SubNotAvailable() string {
	return s.subNotAvailable
}

func (s *Submitter) SubServiceDown() string {
	return s.subServiceDown
}
