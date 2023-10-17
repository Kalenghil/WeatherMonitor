package file_io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"io"
	"log"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

const (
	DefaultDevicesInfoPath = "./info/devices_info.json"
)

type ResponseMap map[string]*DeviceResponse

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getTableName(data *DeviceData) string {
	DataType := reflect.TypeOf(*data)
	TableName := toSnakeCase(DataType.Name()) + "s"
	return TableName
}

func getSQLData(elem DeviceResponse) []interface{} {
	var Data []interface{}
	Data = append(Data, elem.Date)
	Data = append(Data, structs.Values(elem.Data)...)

	return Data
}

func (r *DeviceResponse) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	var raw json.RawMessage
	err := dec.Decode(&raw)
	if err != nil {
		return err
	}
	var header struct {
		Date         string `json:"Date"`
		DeviceName   string `json:"uName"`
		SerialNumber string `json:"serial"`
	}
	err = json.Unmarshal(raw, &header)
	result := DeviceResponse{Date: header.Date, DeviceName: header.DeviceName, SerialNumber: header.SerialNumber, Data: nil}

	switch header.DeviceName {
	case "Тест Студии":
		tgt := &struct {
			Data TestStudioData `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "Роса-К-1":
		tgt := &struct {
			Data RosaK1Data `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "РОСА К-2":
		tgt := &struct {
			Data RosaK2Data `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "Hydra-L":
		tgt := &struct {
			Data HydraLData `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "Hydra-L1":
		tgt := &struct {
			Data HydraL1Data `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "Опорный барометр":
		tgt := &struct {
			Data ReferenceBarometerData `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	case "Паскаль":
		tgt := &struct {
			Data PascalData `json:"data"`
		}{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = tgt.Data
	default:
		tgt := &DefaultData{}
		err = json.Unmarshal(raw, &tgt)
		result.Data = *tgt
	}
	*r = result
	return err
}

func UnmarshalJSONIntoResponseMap(filename string) (*ResponseMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file")
	} else {
		log.Println("File opened successfully")
	}
	AllEntries := make(ResponseMap, 0)
	rawData, err := io.ReadAll(file)

	err = json.Unmarshal(rawData, &AllEntries)
	if err != nil {
		panic(err)
	}

	file.Close()
	return &AllEntries, err
}

func GetResponseMapsFromFiles() (*ResponseMap, error) {
	files := []string{"./raw_data/json/0030-31.json",
		"./raw_data/json/0031-01.json",
		"./raw_data/json/1-2.json",
		"./raw_data/json/2-3.json",
		"./raw_data/json/3-4.json",
		"./raw_data/json/4-5.json",
		"./raw_data/json/5-6.json",
		"./raw_data/json/6-7.json",
		"./raw_data/json/7-8.json",
		"./raw_data/json/8-9.json",
		"./raw_data/json/9-10.json",
		"./raw_data/json/10-11.json",
		"./raw_data/json/11-12.json",
		"./raw_data/json/12-13.json",
		"./raw_data/json/13-14.json",
		"./raw_data/json/14-15.json",
		"./raw_data/json/15-16.json",
		"./raw_data/json/16-17.json",
		"./raw_data/json/17-18.json",
		"./raw_data/json/18-19.json",
		"./raw_data/json/19-20.json",
		"./raw_data/json/20-21.json",
		"./raw_data/json/21-22.json",
		"./raw_data/json/22-23.json",
		"./raw_data/json/23-24.json",
		"./raw_data/json/24-25.json",
		"./raw_data/json/25-26.json",
		"./raw_data/json/26-27.json",
		"./raw_data/json/27-28.json"}

	slices.Reverse(files)
	responseMap := make(ResponseMap, 0)
	for _, fileName := range files {
		fmt.Println(fileName + " opened")
		localResponseMap, err := UnmarshalJSONIntoResponseMap(fileName)
		if err != nil {
			return nil, err
		}
		maps.Copy(responseMap, *localResponseMap)

	}

	return &responseMap, nil
}

func GetDeviceInfo(filename string) (DevicesInfoMap, error) {
	var InfoMap DevicesInfoMap
	f, err := os.Open(filename)
	if err != nil {
		return InfoMap, err
	}

	raw, err := io.ReadAll(f)
	if err != nil {
		return InfoMap, nil
	}
	err = json.Unmarshal(raw, &InfoMap)
	return InfoMap, err
}

func MergeJSONFromSource(filename string, responseMap interface{}) error {
	raw, err := json.Marshal(responseMap)
	fmt.Println("JSON Marshaled into raw")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, raw, 0644)
	fmt.Println("File Written")
	return err
}
