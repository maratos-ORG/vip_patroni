package checker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	log "vip_patroni/internal/logging"
)

type Response struct {
	State         string `json:"state"`
	Role          string `json:"role"`
	ServerVersion int    `json:"server_version"`
}

func GetPatroniStatus(url string) (Response, error) {
	var result Response
	resp, err := http.Get(url)
	if err != nil {
		log.Error("no response from request: %s", err)
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		log.Error("can't read from response body: %s", err)
		return result, err
	}

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Error("an not unmarshal JSON: %s", err)
		return result, err
	}
	return result, nil
}
