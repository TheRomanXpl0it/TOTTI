package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sub/config"
	"sub/db"
	"sub/log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var conf *config.Config

func configHandler(w http.ResponseWriter, r *http.Request) {
	teams := make(map[string]string)
	for i, ip := range conf.Teams {
		teams[fmt.Sprintf("team%d", i+1)] = ip
	}

	response := map[string]interface{}{
		"flag_format":    conf.FlagFormat,
		"round_duration": conf.RoundDuration,
		"teams":          teams,
		"flag_id_url":    conf.FlagIdUrl,
		"flag_lifetime":  conf.FlagAlive,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Errorf("json marshal error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func flagsHandler(w http.ResponseWriter, r *http.Request) {
	var flags []db.Flag
	err := json.NewDecoder(r.Body).Decode(&flags)
	if err != nil {
		log.Errorf("json decoder error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range len(flags) {
		flags[i].Time = primitive.Timestamp{T: uint32(time.Now().Unix())}
		flags[i].Status = db.DB_NSUB
		flags[i].ServerResponse = db.DB_NSUB
	}

	err = conf.DB.InsertFlags(flags)
	if err != nil {
		log.Errorf("DB insert error: %v\n", err)
	}

	w.WriteHeader(http.StatusOK)
}

func manualHandler(w http.ResponseWriter, r *http.Request) {
	urlsValues := r.URL.Query()
	flags, ok := urlsValues["flag"]
	if !ok {
		log.Warningf("No manual flag in request\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flag := flags[0]
	log.Infof("Manual Flag: %s\n", flag)
	if !conf.FlagRegex.Match([]byte(flag)) {
		log.Warningf("Invalid manual flag\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	insertFlags := []db.Flag{
		{
			Flag:           flag,
			Username:       "manual",
			ExploitName:    "manual",
			TeamIP:         conf.Teams[0],
			Time:           primitive.Timestamp{T: uint32(time.Now().Unix())},
			Status:         db.DB_NSUB,
			ServerResponse: db.DB_NSUB,
		},
	}

	err := conf.DB.InsertFlags(insertFlags)
	if err != nil {
		log.Errorf("DB insert error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var startTick string
	startTickCookie, err := r.Cookie("startTick")
	if err == nil {
		startTick = startTickCookie.Value
	}
	var exploit string
	exploitCookie, err := r.Cookie("exploit")
	if err == nil {
		exploit = exploitCookie.Value
	}

	json.NewEncoder(w).Encode(composeData(startTick, exploit))
}

func ServeAPI(c *config.Config) {
	conf = c

	http.HandleFunc("GET /api/config", configHandler)
	http.HandleFunc("POST /api/flags", flagsHandler)

	http.HandleFunc("POST /manual", manualHandler)
	http.HandleFunc("GET /data", dataHandler)
	http.Handle("GET /", http.FileServer(http.Dir("./static")))

	if err := http.ListenAndServe("0.0.0.0:5000", nil); err != nil {
		log.Criticalf("error starting the api server: %v\n", err)
	}
}
