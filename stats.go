package main

import (
	"github.com/zenoss/glog"

	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// Send the list of stats to the TSDB.
func Post(destination string, stats []Sample) error {
	payload := map[string][]Sample{"metrics": stats}
	data, err := json.Marshal(payload)
	if err != nil {
		glog.Warningf("Couldn't marshal stats: ", err)
		return err
	}
	statsreq, err := http.NewRequest("POST", destination, bytes.NewBuffer(data))
	if err != nil {
		glog.Warningf("Couldn't create stats request: ", err)
		return err
	}
	statsreq.Header["User-Agent"] = []string{"Zenoss Metric Publisher"}
	statsreq.Header["Content-Type"] = []string{"application/json"}
	resp, reqerr := http.DefaultClient.Do(statsreq)
	if reqerr != nil {
		glog.Warningf("Couldn't post stats: ", reqerr)
		return reqerr
	}
	if !strings.Contains(resp.Status, "200 OK") {
		glog.Warningf("couldn't post stats: ", resp.Status)
		return nil
	}
	resp.Body.Close()
	return nil
}

// Sample is a single metric measurement
type Sample struct {
	Metric    string            `json:"metric"`
	Value     string            `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Tags      map[string]string `json:"tags,omitempty"`
}
