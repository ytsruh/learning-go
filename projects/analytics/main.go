package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mileusna/useragent"
)

var (
	forceIP         = ""
	events  *Events = &Events{ch: make(chan qdata)}
)

type TrackingData struct {
	Type          string `json:"type"`
	Identity      string `json:"identity"`
	UserAgent     string `json:"ua"`
	Event         string `json:"event"`
	Category      string `json:"category"`
	Referrer      string `json:"referrer"`
	ReferrerHost  string
	IsTouchDevice bool `json:"isTouchDevice"`
	OccuredAt     uint32
}

type Tracking struct {
	SiteID string       `json:"site_id"`
	Action TrackingData `json:"tracking"`
}

func main() {
	flag.StringVar(&forceIP, "ip", "", "force IP for request, useful in local")
	flag.Parse()

	if err := events.Open(); err != nil {
		log.Fatal(err)
	} else if err := events.EnsureTable(); err != nil {
		log.Fatal(err)
	}

	go events.Run()

	http.HandleFunc("/track", track)
	http.ListenAndServe(":9876", nil)
}

func track(w http.ResponseWriter, r *http.Request) {
	defer w.WriteHeader(http.StatusOK)

	data := r.URL.Query().Get("data")
	trk, err := decodeData(data)
	if err != nil {
		fmt.Print(err)
	}

	ua := useragent.Parse(trk.Action.UserAgent)

	headers := []string{"X-Forward-For", "X-Real-IP"}
	ip, err := ipFromRequest(headers, r)
	if err != nil {
		fmt.Println("error getting IP: ", err)
		return
	}

	geoInfo, err := getGeoInfo(ip.String())
	if err != nil {
		fmt.Println("error getting geo info: ", err)
		return
	}

	if len(trk.Action.Referrer) > 0 {
		u, err := url.Parse(trk.Action.Referrer)
		if err == nil {
			trk.Action.ReferrerHost = u.Host
		}
	}

	go events.Add(trk, ua, geoInfo)
}

func decodeData(s string) (data Tracking, err error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &data)
	return
}
