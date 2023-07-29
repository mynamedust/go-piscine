package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type LocationType struct {
	lat, lon float64
}

type RestaurantsType struct {
	name, address, phone string
	id                   int
	location             LocationType
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
	reustarant.id, err = strconv.Atoi(line[0])
	if err != nil {
		fmt.Println("error: atoi failed: ", err)
		return nil, err
	}
	reustarant.name = line[1]
	reustarant.address = line[2]
	reustarant.phone = line[3]
	reustarant.location.lon, err = strconv.ParseFloat(line[4], 64)
	if err != nil {
		fmt.Println("error: float64 parsing failed: ", err)
		return nil, err
	}
	reustarant.location.lat, err = strconv.ParseFloat(line[5], 64)
	if err != nil {
		fmt.Println("error: float64 parsing failed: ", err)
		return nil, err
	}
	return reustarant, nil
}
