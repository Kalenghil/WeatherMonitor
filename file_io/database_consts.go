package file_io

const (

	// Here are string constants, from which the tables for storing weather devices data are creating

	ReferenceBarometerTableName  = `reference_barometer_datas`
	ReferenceBarometerTableQuery = `CREATE TABLE IF NOT EXISTS reference_barometer_datas (
	datetime TEXT,
	system_serial TEXT,
	system_version TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	system_ip TEXT,
	rtc_time TEXT,
	rtc_date TEXT,
	bmp280_0__temp REAL,
	bmp280_0__pressure REAL,
	bmp280_1__temp REAL,
	bmp280_1__pressure REAL,
	bmp280_2__temp REAL,
	bmp280_2__pressure REAL,
	bmp280_3__temp REAL,
	bmp280_3__pressure REAL,
	bmp280_4__temp REAL,
	bmp280_4__pressure REAL,
	bmp280_5__temp REAL,
	bmp280_5__pressure REAL,
	bmp280_6__temp REAL,
	bmp280_6__pressure REAL,
	bmp280_7__temp REAL,
	bmp280_7__pressure REAL,
	bmp280_8__temp REAL,
	bmp280_8__pressure REAL,
	average_temperature REAL,
	average_pressure REAL,
	median_temperature REAL,
	median_pressure REAL
	);`
	ReferenceBaremeterInsertQuery = `INSERT INTO reference_barometer_datas VALUES 
                                          (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	RosaK1TableName  = "rosa_k1_datas"
	RosaK1TableQuery = `CREATE TABLE IF NOT EXISTS rosa_k1_datas (
	datetime TEXT,
	system_version TEXT,
	system_voltage TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	soil_humidity INTEGER,
	soil_temperature REAL,
	color_temp REAL,
	color_lux REAL,
	color_clear_component INTEGER,
	color_i_r_component INTEGER,
	color_red_component INTEGER,
	color_blue_component INTEGER,
	color_green_component INTEGER,
	light_lux REAL,
	light_blink REAL,
	weather_temperature REAL,
	weather_humidity REAL,
	weather_pressure REAL
	);`
	RosaK1InsertQuery = `INSERT INTO rosa_k1_datas VALUES 
                              (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	RosaK2TableName  = "rosa_k2_datas"
	RosaK2TableQuery = `CREATE TABLE IF NOT EXISTS rosa_k2_datas (
    datetime TEXT,
	system_version TEXT,
	system_voltage TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	soil_humidity INTEGER,
	soil_temperature REAL,
	color_temp REAL,
	color_lux REAL,
	color_clear_component INTEGER,
	color_ir_component INTEGER,
	color_red_component INTEGER,
	color_blue_component INTEGER,
	color_green_component INTEGER,
	light_lux REAL,
	light_blink REAL,
	weather_temperature REAL,
	weather_humidity REAL,
	weather_pressure REAL
	);`
	RosaK2InsertQuery = `INSERT INTO rosa_k2_datas VALUES 
                              (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	HydraLTableName  = "hydra_l_datas"
	HydraLTableQuery = `CREATE TABLE IF NOT EXISTS hydra_l_datas (
	datetime TEXT,
	system_serial TEXT,
	system_version TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	system_ip TEXT,
	bme280_temperature REAL,
	bme280_humidity REAL,
	bme280_pressure REAL
	);`
	HydraLInsertQuery = `INSERT INTO hydra_l_datas VALUES 
                              (?,?,?,?,?,?,?,?)`

	HydraL1TableName  = "hydra_l1_datas"
	HydraL1TableQuery = `CREATE TABLE IF NOT EXISTS hydra_l1_datas (
	datetime TEXT,
	system_serial TEXT,
	system_version TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	system_ip TEXT,
	bme280_temperature REAL,
	bme280_humidity REAL,
	bme280_pressure REAL
	);`
	HydraL1InsertQuery = `INSERT INTO hydra_l1_datas VALUES 
                              (?,?,?,?,?,?,?,?)`

	PascalTableName  = "pascal_datas"
	PascalTableQuery = `CREATE TABLE IF NOT EXISTS pascal_datas (
	datetime TEXT,
	system_serial TEXT,
	system_version TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	device_name TEXT,
	weather_temperature REAL,
	weather_pressure REAL
	);`
	PascalInsertQuery = `INSERT INTO pascal_datas VALUES
                             (?,?,?,?,?,?,?,?)`

	TestStudioTableName  = "test_studio_datas"
	TestStudioTableQuery = `CREATE TABLE IF NOT EXISTS test_studio_datas (
    datetime TEXT,
	system_version TEXT,
	system_upit TEXT,
	system_worktime TEXT,
	system_rssi TEXT,
	system_mac TEXT,
	rtc_date TEXT,
	rtc_time TEXT,
	bh1750_lux REAL,
	bh1750_blinkmin REAL,
	bh1750_blinkmax REAL,
	bh1750_blink REAL,
	tcs34725_is_saturated REAL,
	tcs34725_red REAL,
	tcs34725_green REAL,
	tcs34725_blue REAL,
	tcs34725_clear REAL,
	tcs34725_red_c REAL,
	tcs34725_green_c REAL,
	tcs34725_blue_c REAL,
	tcs34725_clear_c REAL,
	tcs34725_ir REAL,
	tcs34725_color_temp_ct REAL,
	tcs34725_lux REAL,
	tcs34725_lux_der REAL,
	tcs34725_lux_max REAL,
	tcs34725_lux_cct REAL,
	ccs811_eco2 REAL,
	ccs811_tvoc REAL,
	ccs811_err_flag REAL,
	ccs811_err_code REAL,
	bmp280_temp REAL,
	bmp280_pressure REAL,
	bme280_temp REAL,
	bme280_humidity REAL,
	bme280_pressure REAL,
	ds18_b20_temp REAL,
	am2321_temp REAL,
	am2321_humidity REAL,
	sbm20_static REAL,
	sbm20_dynamic REAL,
	sbm20_impulses REAL
	);`
	TestStudioInsertQuery = `INSERT INTO test_studio_datas VALUES 
                                  (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	FuncPer3HoursSelectQuery = `
SELECT datetime, sensor_data
FROM
(SELECT datetime, CAST(_FUNC(%s) AS REAL) as sensor_data,
       (cast(strftime('%%H', datetime) AS INTEGER) / 3) AS H,
       (cast(strftime('%%d', %s) AS INTEGER) -
        cast(strftime('%%d', datetime) AS INTEGER) - 1) * 24 / 3 AS T
FROM (%s)
WHERE datetime BETWEEN (%s) AND (%s) _SERIAL_CHECK
GROUP BY T, H)
ORDER BY datetime;`

	FuncPerHourSelectQuery = `
SELECT datetime, sensor_data
FROM
    (SELECT datetime, CAST(_FUNC(%s) AS REAL) as sensor_data, strftime('%%Y-%%m-%%d %%H', datetime) AS H
    FROM (%s)
    WHERE datetime BETWEEN (%s) AND (%s) _SERIAL_CHECK
    GROUP BY H)
ORDER BY datetime;`

	FuncPerDaySelectQuery = `
SELECT datetime, sensor_data
FROM
    (SELECT datetime, CAST(_FUNC(%s) AS REAL) as sensor_data, date(datetime) AS D
    FROM (%s)
    WHERE datetime BETWEEN (%s) AND (%s) _SERIAL_CHECK
    GROUP BY D)
ORDER BY datetime;`

	DefaultDataBaseName   = "weather_data.db"
	DefaultDataBaseFlavor = "sqlite3"
	DefaultBatchSize      = 700
)
