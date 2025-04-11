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

type PomoSubmitter struct {
	Submitter
}

func newPomoSubmitter(c *config.Config) *PomoSubmitter {
	return &PomoSubmitter{Submitter: Submitter{
		conf:              c,
		subAccepted:       "accepted",
		subInvalid:        "invalid or too old",
		subOld:            "too old",
		subYourOwn:        "your own",
		subStolen:         "already stolen",
		subNop:            NO_SUB,
		subNotAvailable:   "is not available",
		subServiceDown:    "service is down",
		subDistpatchError: NO_SUB,
		subNotActive:      NO_SUB,
		subCritical:       NO_SUB,
	}}
}

func (s *PomoSubmitter) submitFlags(flags []string) ([]Response, error) {
	flagsJSON, err := json.Marshal(flags)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", s.conf.SubUrl, bytes.NewBuffer(flagsJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/plain")
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

	if resp.Header.Get("Content-Type") != "application/json" {
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
