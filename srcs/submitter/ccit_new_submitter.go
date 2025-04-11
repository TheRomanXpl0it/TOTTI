package submitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sub/config"
	"time"
)

type CCITNewSubmitter struct {
	Submitter
}

func newCCITNewSubmitter(c *config.Config) *CCITNewSubmitter {
	return &CCITNewSubmitter{Submitter: Submitter{
		conf:              c,
		subAccepted:       "flag claimed",
		subInvalid:        "invalid",
		subOld:            "too old",
		subYourOwn:        "your own",
		subStolen:         "already claimed",
		subNop:            "from nop team",
		subNotAvailable:   "not available",
		subServiceDown:    NO_SUB,
		subDistpatchError: "the check which dispatched this flag didn't terminate successfully",
		subNotActive:      "the flag is not active yet",
		subCritical:       "notify the organizers",
	}}
}

func (s *CCITNewSubmitter) submitFlags(flags []string) ([]Response, error) {
	flagsJSON, err := json.Marshal(flags)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", s.conf.SubUrl, bytes.NewBuffer(flagsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Team-Token", s.conf.TeamToken)

	client := http.Client{Timeout: time.Duration(s.conf.SubInterval) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		urlErr, ok := err.(*url.Error)
		if ok && urlErr.Timeout() {
			return nil, fmt.Errorf("timeout: %v", err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resd response body: %v: Response: %+v",
			err, resp)
	}

	var res []Response
	err = json.Unmarshal(body, &res)
	if err == nil {
		return res, nil
	}

	if resp.StatusCode == 500 {
		var res ResponseError
		err = json.Unmarshal(body, &res)
		if err == nil {
			if res.Code == "RATE_LIMIT" {
				return nil, fmt.Errorf("rate limit exceeded: %s", res.Message)
			}
			return nil, fmt.Errorf("server error: %s: %s", res.Code, res.Message)
		}
	}

	return nil, fmt.Errorf(
		"failed to unmarshal response: %v: Response: %+v Body: %v",
		err, resp, string(body))
}
