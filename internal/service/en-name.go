package service

import (
	"bytes"
	_ "embed"
	"encoding/csv"
)

//go:embed metadata/serviceNameEn.csv
var ServiceNameEnCSV []byte

type ServiceNameEn struct {
	InfID   string
	InfNMEn string
}

var serviceNameEnList []ServiceNameEn
var serviceNameEnMap map[string]string

func init() {
	reader := csv.NewReader(bytes.NewReader(ServiceNameEnCSV))
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		serviceNameEnList = append(serviceNameEnList, ServiceNameEn{
			InfID:   record[0],
			InfNMEn: record[2],
		})
	}

	serviceNameEnMap = make(map[string]string)
	for _, item := range serviceNameEnList {
		serviceNameEnMap[item.InfID] = item.InfNMEn
	}
}

func GetServiceNameEnList() []ServiceNameEn {
	return serviceNameEnList
}

func GetServiceNameEn(ID string) string {
	return serviceNameEnMap[ID]
}
