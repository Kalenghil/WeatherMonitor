package file_io

type DevicesInfoMap map[string]DevicesInfo

type DevicesInfo struct {
	DeviceName    string
	DeviceSerials []string
	DeviceSensors []string
	DeviceCounter int
}

type DeviceResponse struct {
	Date         string     `json:"Date" `
	DeviceName   string     `json:"uName"`
	SerialNumber string     `json:"serial"`
	Data         DeviceData `json:"data"`
}

type DeviceData interface {
	TableQuery() string
}

type ReferenceBarometerData struct {
	SystemSerial  string `json:"system_Serial"`
	SystemVersion string `json:"system_Version"`
	SystemRSSI    string `json:"system_RSSI"`
	SystemMAC     string `json:"system_MAC"`
	SystemIP      string `json:"system_IP"`

	RTCTime string `json:"RTC_time"`
	RTCDate string `json:"RTC_date"`

	BMP280_0_Temp      float32 `json:"BMP280_0_temp,string"`
	BMP280_0_Pressure  float32 `json:"BMP280_0_pressure,string"`
	BMP280_1_Temp      float32 `json:"BMP280_1_temp,string"`
	BMP280_1_Pressure  float32 `json:"BMP280_1_pressure,string"`
	BMP280_2_Temp      float32 `json:"BMP280_2_temp,string"`
	BMP280_2_Pressure  float32 `json:"BMP280_2_pressure,string"`
	BMP280_3_Temp      float32 `json:"BMP280_3_temp,string"`
	BMP280_3_Pressure  float32 `json:"BMP280_3_pressure,string"`
	BMP280_4_Temp      float32 `json:"BMP280_4_temp,string"`
	BMP280_4_Pressure  float32 `json:"BMP280_4_pressure,string"`
	BMP280_5_Temp      float32 `json:"BMP280_5_temp,string"`
	BMP280_5_Pressure  float32 `json:"BMP280_5_pressure,string"`
	BMP280_6_Temp      float32 `json:"BMP280_6_temp,string"`
	BMP280_6_Pressure  float32 `json:"BMP280_6_pressure,string"`
	BMP280_7_Temp      float32 `json:"BMP280_7_temp,string"`
	BMP280_7_Pressure  float32 `json:"BMP280_7_pressure,string"`
	BMP280_8_Temp      float32 `json:"BMP280_8_temp,string"`
	BMP280_8_Pressure  float32 `json:"BMP280_8_pressure,string"`
	AverageTemperature float32 `json:"weather_temp,string"`
	AveragePressure    float32 `json:"weather_pressure,string"`
	MedianTemperature  float32 `json:"weather_temperature_mediana,string"`
	MedianPressure     float32 `json:"weather_pressure_mediana,string"`
}

func (ReferenceBarometerData) TableQuery() string {
	return ReferenceBarometerTableQuery
}

func (ReferenceBarometerData) InsertQuery() string {
	return ReferenceBaremeterInsertQuery
}

type RosaK1Data struct {
	SystemVersion       string  `json:"system_Version"`
	SystemVoltage       string  `json:"system_Upit"`
	SystemRSSI          string  `json:"system_RSSI"`
	SystemMAC           string  `json:"system_MAC"`
	SoilHumidity        int     `json:"soil_soilH,string"`
	SoilTemperature     float32 `json:"soil_soilT,string"`
	ColorTemp           float32 `json:"color_tempCT,string"`
	ColorLux            float32 `json:"color_lux,string"`
	ColorClearComponent int     `json:"color_clearC,string"`
	ColorIRComponent    int     `json:"color_ir,string"`
	ColorRedComponent   int     `json:"color_redC,string"`
	ColorBlueComponent  int     `json:"color_blueC,string"`
	ColorGreenComponent int     `json:"color_greenC,string"`
	LightLux            float32 `json:"light_lux,string"`
	LightBlink          float32 `json:"light_blink,string"`
	WeatherTemperature  float32 `json:"weather_temp,string"`
	WeatherHumidity     float32 `json:"weather_humidity,string"`
	WeatherPressure     float32 `json:"weather_pressure,string"`
}

func (RosaK1Data) TableQuery() string {
	return RosaK1TableQuery
}

func (RosaK1Data) InsertQuery() string {
	return RosaK1InsertQuery
}

type RosaK2Data struct {
	SystemVersion       string  `json:"system_Version"`
	SystemVoltage       string  `json:"system_Upit"`
	SystemRSSI          string  `json:"system_RSSI"`
	SystemMAC           string  `json:"system_MAC"`
	SoilHumidity        int     `json:"soil_soilH,string"`
	SoilTemperature     float32 `json:"soil_soilT,string"`
	ColorTemp           float32 `json:"color_tempCT,string"`
	ColorLux            float32 `json:"color_lux,string"`
	ColorClearComponent int     `json:"color_clearC,string"`
	ColorIRComponent    int     `json:"color_ir,string"`
	ColorRedComponent   int     `json:"color_redC,string"`
	ColorBlueComponent  int     `json:"color_blueC,string"`
	ColorGreenComponent int     `json:"color_greenC,string"`
	LightLux            float32 `json:"light_lux,string"`
	LightBlink          float32 `json:"light_blink,string"`
	WeatherTemperature  float32 `json:"weather_temp,string"`
	WeatherHumidity     float32 `json:"weather_humidity,string"`
	WeatherPressure     float32 `json:"weather_pressure,string"`
}

func (RosaK2Data) TableQuery() string {
	return RosaK2TableQuery
}

func (RosaK2Data) InsertQuery() string {
	return RosaK2InsertQuery
}

type HydraLData struct {
	SystemSerial      string  `json:"system_Serial"`
	SystemVersion     string  `json:"system_Version"`
	SystemRSSI        string  `json:"system_RSSI"`
	SystemMAC         string  `json:"system_MAC"`
	SystemIP          string  `json:"system_IP"`
	BME280Temperature float32 `json:"BME280_temp,string"`
	BME280Humidity    float32 `json:"BME280_humidity,string"`
	BME280Pressure    float32 `json:"BME280_pressure,string"`
}

func (HydraLData) TableQuery() string {
	return HydraLTableQuery
}

func (HydraLData) InsertQuery() string {
	return HydraLInsertQuery
}

type HydraL1Data struct {
	SystemSerial      string  `json:"system_Serial"`
	SystemVersion     string  `json:"system_Version"`
	SystemRSSI        string  `json:"system_RSSI"`
	SystemMAC         string  `json:"system_MAC"`
	SystemIP          string  `json:"system_IP"`
	BME280Temperature float32 `json:"BME280_temp,string"`
	BME280Humidity    float32 `json:"BME280_humidity,string"`
	BME280Pressure    float32 `json:"BME280_pressure,string"`
}

func (HydraL1Data) TableQuery() string {
	return HydraL1TableQuery
}

func (HydraL1Data) InsertQuery() string {
	return HydraL1InsertQuery
}

type PascalData struct {
	SystemSerial       string  `json:"system_Serial"`
	SystemVersion      string  `json:"system_Version"`
	SystemRSSI         string  `json:"system_RSSI"`
	SystemMAC          string  `json:"system_MAC"`
	SystemIP           string  `json:"system_IP"`
	WeatherTemperature float32 `json:"weather_temp,string"`
	WeatherPressure    float32 `json:"weather_pressure,string"`
}

func (PascalData) TableQuery() string {
	return PascalTableQuery
}

func (PascalData) InsertQuery() string {
	return PascalInsertQuery
}

type TestStudioData struct {
	SystemVersion       string  `json:"system_Version"`
	SystemUpit          string  `json:"system_Upit"`
	SystemWorktime      string  `json:"system_Worktime"`
	SystemRSSI          string  `json:"system_RSSI"`
	SystemMAC           string  `json:"system_MAC"`
	RTCDate             string  `json:"RTC_date"`
	RTCTime             string  `json:"RTC_time"`
	BH1750Lux           float32 `json:"BH1750_lux,string"`
	BH1750Blinkmin      float32 `json:"BH1750_blinkmin,string"`
	BH1750Blinkmax      float32 `json:"BH1750_blinkmax,string"`
	BH1750Blink         float32 `json:"BH1750_blink,string"`
	TCS34725IsSaturated float32 `json:"TCS34725_isSaturated,string"`
	TCS34725Red         float32 `json:"TCS34725_red,string"`
	TCS34725Green       float32 `json:"TCS34725_green,string"`
	TCS34725Blue        float32 `json:"TCS34725_blue,string"`
	TCS34725Clear       float32 `json:"TCS34725_clear,string"`
	TCS34725RedC        float32 `json:"TCS34725_redC,string"`
	TCS34725GreenC      float32 `json:"TCS34725_greenC,string"`
	TCS34725BlueC       float32 `json:"TCS34725_blueC,string"`
	TCS34725ClearC      float32 `json:"TCS34725_clearC,string"`
	TCS34725Ir          float32 `json:"TCS34725_ir,string"`
	TCS34725ColorTempCT float32 `json:"TCS34725_colorTempCT,string"`
	TCS34725Lux         float32 `json:"TCS34725_lux,string"`
	TCS34725LuxDER      float32 `json:"TCS34725_luxDER,string"`
	TCS34725LuxMax      float32 `json:"TCS34725_luxMax,string"`
	TCS34725LuxCCT      float32 `json:"TCS34725_luxCCT,string"`
	CCS811ECO2          float32 `json:"CCS811_eCO2,string"`
	CCS811TVOC          float32 `json:"CCS811_TVOC,string"`
	CCS811ErrFlag       float32 `json:"CCS811_ErrFlag,string"`
	CCS811ErrCode       float32 `json:"CCS811_ErrCode,string"`
	BMP280Temp          float32 `json:"BMP280_temp,string"`
	BMP280Pressure      float32 `json:"BMP280_pressure,string"`
	BME280Temp          float32 `json:"BME280_temp,string"`
	BME280Humidity      float32 `json:"BME280_humidity,string"`
	BME280Pressure      float32 `json:"BME280_pressure,string"`
	DS18B20Temp         float32 `json:"DS18B20_temp,string"`
	AM2321Temp          float32 `json:"AM2321_temp,string"`
	AM2321Humidity      float32 `json:"AM2321_humidity,string"`
	SBM20Static         float32 `json:"SBM20_static,string"`
	SBM20Dynamic        float32 `json:"SBM20_dynamic,string"`
	SBM20Impulses       float32 `json:"SBM20_impulses,string"`
}

func (TestStudioData) TableQuery() string {
	return TestStudioTableQuery
}

func (TestStudioData) InsertQuery() string {
	return TestStudioInsertQuery
}

type DefaultData struct {
	Data interface{} `json:"data"`
}

func (DefaultData) TableQuery() string {
	return ""
}

func (DefaultData) InsertQuery() string {
	return ""
}
