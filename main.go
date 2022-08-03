package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

var (
	unsupportedCriteria      = errors.New("unsupported criteria")
	emptyDepartureStation    = errors.New("empty departure station")
	emptyArrivalStation      = errors.New("empty arrival station")
	badArrivalStationInput   = errors.New("bad arrival station input")
	badDepartureStationInput = errors.New("bad departure station input")
)

func (t *Train) UnmarshalJSON(b []byte) (err error) {
	type TrainDuplicate Train

	tt := struct {
		ArrivalTime   string
		DepartureTime string
		*TrainDuplicate
	}{
		TrainDuplicate: (*TrainDuplicate)(t),
	}
	err = json.Unmarshal(b, &tt)
	if err != nil {
		return err
	}

	t.ArrivalTime, err = time.Parse("15:04:05", tt.ArrivalTime)
	if err != nil {
		return err
	}

	t.DepartureTime, err = time.Parse("15:04:05", tt.DepartureTime)
	if err != nil {
		return err
	}
	return nil
}

func (t Trains) sortByCriteria(criteria string) {
	switch criteria {
	case "price":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].Price < t[j].Price
		})
	case "arrival-time":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].ArrivalTime.Before(t[j].ArrivalTime)
		})
	case "departure-time":
		sort.SliceStable(t, func(i, j int) bool {
			return t[i].DepartureTime.Before(t[j].DepartureTime)
		})
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	var stationLimit = 1

	filename, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer filename.Close()

	data, err := ioutil.ReadAll(filename)

	if err != nil {
		log.Fatal(err)
	}

	var t Trains

	jsonErr := json.Unmarshal(data, &t)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if departureStation == "" {
		return nil, emptyDepartureStation
	}

	depStation, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, badDepartureStationInput
	}

	if depStation < stationLimit {
		return nil, badDepartureStationInput
	}

	if arrivalStation == "" {
		return nil, emptyArrivalStation
	}

	arrStation, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, badArrivalStationInput
	}

	if arrStation < stationLimit {
		return nil, badArrivalStationInput
	}

	if criteria != "price" && criteria != "arrival-time" && criteria != "departure-time" {
		return nil, unsupportedCriteria
	}
	var filtered Trains
	for _, tr := range t {
		if tr.ArrivalStationID == arrStation && tr.DepartureStationID == depStation {
			filtered = append(filtered, tr)
		}
	}
	
	if len(filtered) == 0 {
		return nil, nil
	}

	if len(filtered) > 2 {
		filtered = filtered[:3]
	}
	
	filtered.sortByCriteria(criteria)
	return filtered[:3], nil
}

func PrintTrains(trains Trains) {
	for _, train := range trains {
		fmt.Printf("{TrainID: %d, "+
			"DepartureStationID: %d, "+
			"ArrivalStationID: %d, "+
			"Price: %.2f, "+
			"ArrivalTime: time.Date(%d, time.%s, %d, %d, %d, %d, %d time.%s), "+
			"DepartureTime: time.Date(%d, time.%s, %d, %d, %d, %d, %d time.%s)}\n",
			train.TrainID,
			train.DepartureStationID,
			train.ArrivalStationID,
			train.Price,
			train.ArrivalTime.Year(),
			train.ArrivalTime.Month().String(),
			train.ArrivalTime.Day(),
			train.ArrivalTime.Hour(),
			train.ArrivalTime.Minute(),
			train.ArrivalTime.Second(),
			train.ArrivalTime.Nanosecond(),
			train.ArrivalTime.Location(),
			train.DepartureTime.Year(),
			train.DepartureTime.Month().String(),
			train.DepartureTime.Day(),
			train.DepartureTime.Hour(),
			train.DepartureTime.Minute(),
			train.DepartureTime.Second(),
			train.DepartureTime.Nanosecond(),
			train.DepartureTime.Location())
	}
}

func main() {
	var departureStation, arrivalStation, criteria string

	fmt.Print("Departure station: ")
	fmt.Scanln(&departureStation)

	fmt.Print("Arrival station: ")
	fmt.Scanln(&arrivalStation)

	fmt.Print("Criteria: ")
	fmt.Scanln(&criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		log.Fatal(err)
	}

	PrintTrains(result)
}
