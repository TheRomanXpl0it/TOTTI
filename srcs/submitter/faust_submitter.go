package submitter

import (
	"net"
	"strings"
	"sub/config"
)

type FaustSubmitter struct {
	Submitter
}

func newFaustSubmitter(c *config.Config) *FaustSubmitter {
	return &FaustSubmitter{Submitter: Submitter{
		conf:            c,
		subAccepted:     "OK",
		subInvalid:      "INV",
		subOld:          "OLD",
		subYourOwn:      "OWN",
		subStolen:       "DUP",
		subNop:          "from NOP team",
		subNotAvailable: "ERR",
	}}
}

func (s *FaustSubmitter) submitFlags(flags []string) ([]Response, error) {
	conn, err := net.Dial("tcp6", s.conf.SubUrl)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	responses := make([]Response, 0, len(flags))
	for _, flag := range flags {
		_, err = conn.Write([]byte(flag + "\n"))
		if err != nil {
			return nil, err
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		out := strings.SplitN(string(buf[:n-1]), " ", 2)
		msg := out[1]

		resp := Response{Msg: msg, Flag: flag}
		responses = append(responses, resp)
	}

	return responses, nil
}
