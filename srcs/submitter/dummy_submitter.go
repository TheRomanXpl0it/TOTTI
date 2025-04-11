package submitter

import "sub/config"

type DummySubmitter struct {
	Submitter
}

func newDummySubmitter(c *config.Config) *DummySubmitter {
	return &DummySubmitter{Submitter: Submitter{
		conf:              c,
		subAccepted:       "accepted",
		subInvalid:        "invalid",
		subOld:            "too old",
		subYourOwn:        "your own",
		subStolen:         "already stolen",
		subNop:            "from NOP team",
		subNotAvailable:   "is not available",
		subServiceDown:    NO_SUB,
		subDistpatchError: NO_SUB,
		subNotActive:      NO_SUB,
		subCritical:       NO_SUB,
	}}
}

func (s *DummySubmitter) submitFlags(flags []string) ([]Response, error) {
	res := make([]Response, len(flags))
	for _, flag := range flags {
		res = append(res, Response{
			Msg:    s.subAccepted,
			Flag:   flag,
			Status: "",
		})
	}
	return res, nil
}
