package database

import (
	"WeatherMonitor/file_io"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"log"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type TypeOfFunc int
type TypeOfAveraging int

const (
	IsRaw   TypeOfFunc = 0
	IsMax   TypeOfFunc = 1
	IsMin   TypeOfFunc = 2
	IsAvg   TypeOfFunc = 3
	IsCount TypeOfFunc = 4

	NoGrouping TypeOfAveraging = 0
	PerHour    TypeOfAveraging = 1
	Per3Hours  TypeOfAveraging = 2
	PerDay     TypeOfAveraging = 3
)

type SensorDataRequest struct {
	TableName          string
	SensorName         string
	SensorSerial       string
	BeginDateTime      string
	EndDateTime        string
	TypeOfDataFunc     TypeOfFunc
	TypeOfDataGrouping TypeOfAveraging
}

type PlotDataElem struct {
	Datetime string
	Value    float32
}

type PlotDataArray []PlotDataElem

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getTableName(data *file_io.DeviceData) string {
	DataType := reflect.TypeOf(*data)
	TableName := ToSnakeCase(DataType.Name()) + "s"
	return TableName
}

func getSQLData(elem *file_io.DeviceResponse) []interface{} {
	var Data []interface{}
	Data = append(Data, elem.Date)
	Data = append(Data, structs.Values(elem.Data)...)
	return Data
}

func CreateDataBase() (*sql.DB, error) {
	file, err := os.Create(file_io.DefaultDataBaseName)
	if err != nil {
		return nil, err
	}
	err = file.Close()

	db, err := sql.Open(file_io.DefaultDataBaseFlavor, file_io.DefaultDataBaseName)
	log.Println("DB created successfully")

	return db, err
}

func OpenDataBaseFromFile(filePath string) (*sql.DB, error) {
	db, err := sql.Open(file_io.DefaultDataBaseFlavor, filePath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()

	return db, err
}

func PrepareDataBase(db *sql.DB) error {
	createReferenceBarometerData, err := db.Prepare(file_io.ReferenceBarometerTableQuery)
	if err != nil {
		return err
	}
	createRosaK1Data, err := db.Prepare(file_io.RosaK1TableQuery)
	if err != nil {
		return err
	}
	createRosaK2Data, err := db.Prepare(file_io.RosaK2TableQuery)
	if err != nil {
		return err
	}
	createHydraLData, err := db.Prepare(file_io.HydraLTableQuery)
	if err != nil {
		return err
	}
	createHydraL1Data, err := db.Prepare(file_io.HydraL1TableQuery)
	if err != nil {
		return err
	}
	createPascalData, err := db.Prepare(file_io.PascalTableQuery)
	if err != nil {
		return err
	}
	createTestStudioData, err := db.Prepare(file_io.TestStudioTableQuery)
	if err != nil {
		return err
	}

	_, err = createReferenceBarometerData.Exec()
	_, err = createRosaK1Data.Exec()
	_, err = createRosaK2Data.Exec()
	_, err = createHydraLData.Exec()
	_, err = createHydraL1Data.Exec()
	_, err = createPascalData.Exec()
	_, err = createTestStudioData.Exec()

	return err
}

func InsertDataIntoDB(db *sql.DB, entryMap *file_io.ResponseMap) (file_io.DevicesInfoMap, error) {
	Responses := make(map[string][]file_io.DeviceResponse, 0)
	devicesInfo := make(file_io.DevicesInfoMap, 0)
	AppendIfMissing := func(slice []string, i string) []string {
		for _, ele := range slice {
			if ele == i {
				return slice
			}
		}
		return append(slice, i)
	}
	for _, Entry := range *entryMap {
		if reflect.TypeOf(Entry.Data) == reflect.TypeOf(file_io.DefaultData{}) {
			continue
		}

		if device, exists := devicesInfo[Entry.DeviceName]; exists {
			device.DeviceSerials = AppendIfMissing(device.DeviceSerials, Entry.SerialNumber)
			device.DeviceCounter++
			devicesInfo[Entry.DeviceName] = device
		} else {
			devicesInfo[Entry.DeviceName] = file_io.DevicesInfo{DeviceName: reflect.TypeOf(Entry.Data).Name(),
				DeviceSensors: structs.Names(Entry.Data),
				DeviceSerials: []string{Entry.SerialNumber},
				DeviceCounter: 1,
			}
		}

		Responses[Entry.DeviceName] = append(Responses[Entry.DeviceName], *Entry)
		if len(Responses[Entry.DeviceName]) >= file_io.DefaultBatchSize {
			err := BatchDataToDB(db, Entry.DeviceName, Responses[Entry.DeviceName])
			if err != nil {
				return devicesInfo, err
			}
			Responses[Entry.DeviceName] = nil
		}
	}

	for DeviceName, DeviceResponseSlice := range Responses {
		err := BatchDataToDB(db, DeviceName, DeviceResponseSlice)
		if err != nil {
			return devicesInfo, err
		}
	}

	for _, elem := range devicesInfo {
		slices.Sort(elem.DeviceSerials)
	}

	return devicesInfo, nil
}

func BatchDataToDB(db *sql.DB, deviceName string, responses []file_io.DeviceResponse) error {
	BatchInsertQuery, _, _ :=
		sq.Insert(getTableName(&responses[0].Data)).
			Values(getSQLData(&responses[0])...).
			ToSql()

	ValueTemplateSlice := make(
		[]string,
		reflect.TypeOf(responses[0].Data).NumField()+1)

	for i := range ValueTemplateSlice {
		ValueTemplateSlice[i] = "?"
	}
	ValueTemplate := "(" + strings.Join(ValueTemplateSlice, ",") + ")"

	for i := 0; i < len(responses)-1; i++ {
		BatchInsertQuery += ", " + ValueTemplate
	}

	SQLDataSlice := make([]interface{}, 0)
	for _, DataEntry := range responses {
		SQLDataSlice = append(SQLDataSlice, getSQLData(&DataEntry)...)
	}

	BatchInsertStatement, err := db.Prepare(BatchInsertQuery)
	if err != nil {
		return err
	}

	_, err = BatchInsertStatement.Exec(SQLDataSlice...)
	log.Println(deviceName + ": Batched Successfully")
	return err
}

func CreateDB() (*sql.DB, error) {
	db, err := CreateDataBase()
	err = PrepareDataBase(db)

	return db, err
}

func trimSQLString(str string) string {
	return "'" + str + "'"
}

func prepareSelectRequest(request SensorDataRequest) string {
	var SelectRequest string
	if request.TypeOfDataGrouping == NoGrouping {
		SelectRequest, _, _ =
			sq.Select("datetime", "CAST((%s) AS REAL) AS sensor_data").
				From("%s").
				Where("datetime BETWEEN %s AND %s _SERIAL_CHECK").
				OrderBy("datetime").ToSql()
	} else {
		var BaseSelectRequest string
		switch request.TypeOfDataGrouping {
		case PerHour:
			BaseSelectRequest = file_io.FuncPerHourSelectQuery
		case Per3Hours:
			BaseSelectRequest = file_io.FuncPer3HoursSelectQuery
		case PerDay:
			BaseSelectRequest = file_io.FuncPerDaySelectQuery
		}

		RegexpExpr, _ := regexp.Compile("_FUNC")
		switch request.TypeOfDataFunc {
		case IsRaw:
			BaseSelectRequest = RegexpExpr.ReplaceAllString(BaseSelectRequest, "")
		case IsMax:
			BaseSelectRequest = RegexpExpr.ReplaceAllString(BaseSelectRequest, "MAX")
		case IsMin:
			BaseSelectRequest = RegexpExpr.ReplaceAllString(BaseSelectRequest, "MIN")
		case IsAvg:
			BaseSelectRequest = RegexpExpr.ReplaceAllString(BaseSelectRequest, "AVG")
		case IsCount:

		}

		SelectRequest = BaseSelectRequest
	}

	RegexpExpr, _ := regexp.Compile("_SERIAL_CHECK")
	if request.SensorSerial != "" {
		SelectRequest = RegexpExpr.ReplaceAllString(SelectRequest,
			"AND system_serial = "+trimSQLString(request.SensorSerial))
	} else {

		SelectRequest = RegexpExpr.ReplaceAllString(SelectRequest, "")
	}

	if request.TypeOfDataGrouping == Per3Hours {
		SelectRequest = fmt.Sprintf(SelectRequest,
			request.SensorName,
			trimSQLString(request.EndDateTime),
			request.TableName,
			trimSQLString(request.BeginDateTime),
			trimSQLString(request.EndDateTime))
	} else {
		SelectRequest = fmt.Sprintf(SelectRequest,
			request.SensorName,
			request.TableName,
			trimSQLString(request.BeginDateTime),
			trimSQLString(request.EndDateTime))
	}
	return SelectRequest
}

func GetDataFromDB(db *sql.DB, request SensorDataRequest) (*PlotDataArray, error) {
	PlotData := make(PlotDataArray, 0)

	SelectRequest := prepareSelectRequest(request)
	SelectStatement, err := db.Prepare(SelectRequest)
	if err != nil {
		panic(err)
	}
	rows, err := SelectStatement.Query()
	for rows.Next() {
		var elem PlotDataElem
		err = rows.Scan(&elem.Datetime, &elem.Value)
		if err != nil {
			return &PlotData, err
		}
		PlotData = append(PlotData, elem)
	}

	return &PlotData, err
}
