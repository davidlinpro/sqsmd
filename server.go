package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

func handleRequests() {
	r := mux.NewRouter()
	r.HandleFunc("/stats", statsHandler)
	r.HandleFunc("/ws_stats", wsStatsHandler)

	// for pprof
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))

	// for frontend -- better would be to serve this elsewhere but this is just a simple app
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("app")))

	log.Fatal(http.ListenAndServe(":8888", r))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	stats, err := json.Marshal(CacheGet())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(stats)
	return
}
