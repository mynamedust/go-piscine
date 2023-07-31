package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type LocationType struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RestaurantsType struct {
	Name     string       `json:"name,omitempty"`
	Address  string       `json:"address,omitempty"`
	Phone    string       `json:"phone,omitempty"`
	Id       int          `json:"id,omitempty"`
	Location LocationType `json:"location,omitempty"`
}

func ReadLine(r *csv.Reader) (*RestaurantsType, error) {
	reustarant := new(RestaurantsType)
	line, err := r.Read()
	if err != nil && err != io.EOF {
		fmt.Println("error: file reading failed: ", err)
		return nil, err
	} else if err == io.EOF {
		return nil, err
	}
	reustarant.Id, err = strconv.Atoi(line[0])
	if err != nil {
		fmt.Println("error: atoi failed: ", err)
		return nil, err
	}
	reustarant.Name = line[1]
	reustarant.Address = line[2]
	reustarant.Phone = line[3]
	reustarant.Location.Lon, err = strconv.ParseFloat(line[4], 64)
	if err != nil {
		fmt.Println("error: float64 parsing failed: ", err)
		return nil, err
	}
	reustarant.Location.Lat, err = strconv.ParseFloat(line[5], 64)
	if err != nil {
		fmt.Println("error: float64 parsing failed: ", err)
		return nil, err
	}
	return reustarant, nil
}
