package station

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ArielGwd/mrt-schedules.git/module/common/client"
)

type Service interface {
	getAllStation() (response []StationResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (srv *service) getAllStation() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	// hit external API
	byteResponse, err := client.DoRequest(srv.client, url)
	if err != nil {
		return nil, err
	}

	var stations []Station
	err = json.Unmarshal(byteResponse, &stations)

	for _, station := range stations {
		response = append(response, StationResponse(station))
	}

	return
}
