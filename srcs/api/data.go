package api

import (
	"strconv"
	"sub/db"
	"sub/log"
	"time"
)

type FlagData struct {
	Accepted uint64 `json:"accepted"`
	Queued   uint64 `json:"queued"`
	Expired  uint64 `json:"expired"`
	Error    uint64 `json:"error"`
}

type ChartsData struct {
	Flags    FlagData            `json:"flags"`
	Exploits map[string]FlagData `json:"exploits"`
	Teams    map[string]FlagData `json:"teams"`
	Rounds   map[int]FlagData    `json:"rounds"`
}

var FirstTick int64

func composeChartsData(allFlags []db.Flag, tick int, exploit string) ChartsData {
	var flags FlagData
	exploits := make(map[string]FlagData)
	teams := make(map[string]FlagData)
	rounds := make(map[int]FlagData)

	for i := FirstTick; i < time.Now().Unix(); i += int64(conf.RoundDuration) {
		rounds[int(i-FirstTick)/conf.RoundDuration] = FlagData{}
	}

	// TODO rounds data
	for _, flag := range allFlags {
		if exploit != "" && flag.ExploitName != exploit {
			continue
		}

		currentExploit := exploits[flag.ExploitName]
		team := teams[flag.TeamIP]
		flagTick := int(flag.Time.T-uint32(FirstTick)) / conf.RoundDuration
		if tick > 0 {
			//now := time.Now().Unix() - int64(tick*conf.RoundDuration)
			currTick := int(time.Now().Unix()-int64(FirstTick)) / conf.RoundDuration
			if flagTick <= currTick-tick {
				teams[flag.TeamIP] = team
				exploits[flag.ExploitName] = currentExploit
				continue
			}
		}
		round := rounds[flagTick]

		switch flag.ServerResponse {
		case db.DB_SUCC:
			team.Accepted++
			currentExploit.Accepted++
			flags.Accepted++
			round.Accepted++
		case db.DB_ERR:
			team.Error++
			currentExploit.Error++
			flags.Error++
			round.Error++
		case db.DB_NSUB:
			team.Queued++
			currentExploit.Queued++
			flags.Queued++
			round.Queued++
		case db.DB_EXPIRED:
			team.Expired++
			currentExploit.Expired++
			flags.Expired++
			round.Expired++
		}
		teams[flag.TeamIP] = team
		exploits[flag.ExploitName] = currentExploit
		rounds[flagTick] = round
	}

	return ChartsData{
		Flags:    flags,
		Exploits: exploits,
		Teams:    teams,
		Rounds:   rounds,
	}
}

func composeData(start string, exploit string) ChartsData {
	tick := 0
	if start != "" {
		var err error
		if start[0] == 'l' {
			tick, err = strconv.Atoi(start[1:])
			if err != nil {
				log.Errorf("Invalid tick length: %v, Atoi error: %v\n", start, err)
				tick = 0
			}
		} else {
			tick, err = strconv.Atoi(start)
			if err != nil {
				log.Errorf("Invalid tick length: %v, Atoi error: %v\n", start, err)
				tick = 0
			}
		}
	}

	allFlags, err := conf.DB.GetAllFlags()
	if err != nil {
		log.Errorf("DB error retrieving all flags: %v\n", err)
		return ChartsData{}
	}

	return composeChartsData(allFlags, tick, exploit)
}
