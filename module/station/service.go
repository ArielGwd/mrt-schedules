package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ArielGwd/mrt-schedules.git/module/common/client"
)

type Service interface {
	getAllStation() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
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

func (srv *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiun/"

	// hit external API
	byteResponse, err := client.DoRequest(srv.client, url)
	if err != nil {
		return nil, err
	}

	var schedules []Schedule
	err = json.Unmarshal(byteResponse, &schedules)
	if err != nil {
		return nil, err
	}

	// schedule selected by id station
	var selectedSchedule Schedule
	for _, schedule := range schedules {
		if schedule.StationId == id {
			selectedSchedule = schedule
			break
		}
	}

	if selectedSchedule.StationId == "" {
		return nil, errors.New("Station not found")
	}

	response, err = ConvertDataToResponses(selectedSchedule)
	if err != nil {
		return nil, err
	}

	return
}

func ConvertDataToResponses(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHiTripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return nil, err
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return nil, err
	}

	// convert to response
	for _, schedule := range scheduleLebakBulusParsed {
		if schedule.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        schedule.Format("15:04"),
			})
		}
	}

	for _, schedule := range scheduleBundaranHIParsed {
		if schedule.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHiTripName,
				Time:        schedule.Format("15:04"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, schedule := range schedules {
		trimmedTime := strings.TrimSpace(schedule)
		if trimmedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			err = errors.New("Invalid time format " + trimmedTime)
			return
		}

		response = append(response, parsedTime)
	}

	return
}
