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

type CCITSubmitter struct {
	Submitter
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newCCITSubmitter(c *config.Config) *CCITSubmitter {
	return &CCITSubmitter{Submitter: Submitter{
		conf:            c,
		subAccepted:     "accepted",
		subInvalid:      "invalid",
		subOld:          "too old",
		subYourOwn:      "your own",
		subStolen:       "already claimed",
		subNop:          "from NOP team",
		subNotAvailable: "is not available",
	}}
}

func (s *CCITSubmitter) submitFlags(flags []string) ([]Response, error) {
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

	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		return nil, fmt.Errorf("invalid content type: %+v", resp)
	}

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