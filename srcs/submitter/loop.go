package submitter

import (
	"fmt"
	"strings"
	"sub/config"
	"sub/db"
	"sub/log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Counters struct {
	Invalid  uint
	Yours    uint
	Nop      uint
	Old      uint
	Stolen   uint
	Accepted uint
}

var (
	subInvalid  string
	subYourOwn  string
	subNop      string
	subOld      string
	subStolen   string
	subAccepted string
)

func loadSubmitter(conf *config.Config) SubmitterInterface {
	switch conf.SubProtocol {
	case "dummy":
		return newDummySubmitter(conf)
	case "ccit":
		return newCCITSubmitter(conf)
	default:
		log.Fatalf("Unknown submitter protocol: %s", conf.SubProtocol)
	}
	return nil
}

func logFlagsCount(c Counters, totalCount int) {
	// TODO: add colors and precompute the strings
	sLog := fmt.Sprintf("Submitted %d flags: %d Accepted", totalCount, c.Accepted)
	if c.Invalid > 0 {
		sLog = fmt.Sprintf("%s %d Invalid", sLog, c.Invalid)
	}
	if c.Yours > 0 {
		sLog = fmt.Sprintf("%s %d Yours", sLog, c.Yours)
	}
	if c.Nop > 0 {
		sLog = fmt.Sprintf("%s %d Nop", sLog, c.Nop)
	}
	if c.Old > 0 {
		sLog = fmt.Sprintf("%s %d Old", sLog, c.Old)
	}
	if c.Stolen > 0 {
		sLog = fmt.Sprintf("%s %d Stolen", sLog, c.Stolen)
	}
	log.Notice(sLog)
}

func updateSubmittedFlags(conf *config.Config, responses []Response, totalCount int) {
	counters := Counters{}

	for _, response := range responses {
		msg := strings.ToLower(response.Msg)

		if strings.Contains(msg, subInvalid) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_ERR})
			counters.Invalid++
		} else if strings.Contains(msg, subYourOwn) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_ERR})
			counters.Yours++
		} else if strings.Contains(msg, subNop) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_ERR})
			counters.Nop++
		} else if strings.Contains(msg, subOld) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_EXPIRED})
			counters.Old++
		} else if strings.Contains(msg, subStolen) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_SUCC})
			counters.Stolen++
		} else if strings.Contains(msg, subAccepted) {
			conf.DB.UpdateFlag(db.Flag{Flag: response.Flag, Status: db.DB_SUB, ServerResponse: db.DB_SUCC})
			counters.Accepted++
		} else {
			log.Criticalf("Invalid response: %+v\n", response)
		}
	}

	logFlagsCount(counters, totalCount)
}

func Loop(conf *config.Config) {
	submitter := loadSubmitter(conf)
	subInvalid = strings.ToLower(submitter.SubInvalid())
	subYourOwn = strings.ToLower(submitter.SubYourOwn())
	subNop = strings.ToLower(submitter.SubNop())
	subOld = strings.ToLower(submitter.SubOld())
	subStolen = strings.ToLower(submitter.SubStolen())
	subAccepted = strings.ToLower(submitter.SubAccepted())

	log.Infof("Starting submitter loop with %v submitter\n", submitter.Conf().SubProtocol)
	queue := NewOrderedSet()
	for {
		start_time := uint32(time.Now().Unix())

		expirationTime := primitive.Timestamp{T: start_time - uint32(conf.FlagAlive)}

		flagsFromDB, err := conf.DB.GetFlags(expirationTime)
		if err != nil {
			log.Errorf("error fetching flags from db: %v\n", err)
			continue
		}
		log.Debugf("flags fetched from DB: %v\n", len(flagsFromDB))
		for _, flag := range flagsFromDB {
			queue.Add(flag.Flag)
		}
		queueSize := queue.Len()

		for i := 0; i < min(conf.SubLimit, queueSize); {
			flags := make([]string, 0, min(conf.SubLimit, queueSize))
			for range min(conf.SubMaxPayloadSize, queueSize) {
				front := queue.Pop(true)
				if front == nil {
					log.Criticalf("error: queue ended too early\n")
					break
				}
				flags = append(flags, front.(string))
			}
			log.Noticef("Submitting %d flags of %d queued\n", len(flags), queueSize)

			responses, err := submitter.submitFlags(flags)
			if err != nil {
				if strings.Contains(err.Error(), "rate limit exceeded") ||
					strings.Contains(err.Error(), "timeout: ") {
					log.Warningf("error submitting flags: %v\n", err)
				} else {
					log.Errorf("error submitting flags: %v\n", err)
				}
				break
			}
			log.Infof("Received %d responses\n", len(responses))
			i += len(responses)
			updateSubmittedFlags(conf, responses, len(flags))
		}

		count, err := conf.DB.UpdateExpiredFlags(expirationTime)
		if err != nil {
			log.Errorf("error updating expired flags: %v\n", err)
		} else {
			if count > 0 {
				log.Warningf("Updated expired %d flags\n", count)
			} else {
				log.Infof("Updated expired %d flags\n", count)
			}
		}

		end_time := uint32(time.Now().Unix())
		duration := int(end_time - start_time)

		if duration < submitter.Conf().SubInterval {
			time.Sleep(time.Duration(submitter.Conf().SubInterval-duration) * time.Second)
		}
	}
}
