package cpu

import (
	"strconv"
	"reflect"

	// project
	"github.com/shirou/gopsutil/cpu"
)

var cpuMap = map[string]string{
	"VendorID":       "vendor_id",
	"ModelName":      "model_name",
	"Cores":          "cpu_cores",
	"Mhz":            "mhz",
	"Family":         "family",
	"Model":          "model",
	"Stepping":       "stepping",
}

func getCpuInfo() (cpuInfo map[string]string, err error) {

	cpuInfo = make(map[string]string)
	
	cpuStat, err := cpu.Info()
	if err != nil {
		return
	}
	
	cpuInfo["cpu_logical_processors"] = strconv.Itoa(len(cpuStat));

	val := reflect.ValueOf(cpuStat[0])

	for i := 0; i < val.NumField(); i++ {

		typeField := val.Type().Field(i)
		
		if v, ok := cpuMap[typeField.Name]; ok {
			switch val.Field(i).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					cpuInfo[v] = strconv.FormatInt(val.Field(i).Int(), 10)
				case reflect.Float64:
					cpuInfo[v] = strconv.FormatFloat(val.Field(i).Float(), 'f', -1, 64)
 				case reflect.String:
					cpuInfo[v] = val.Field(i).String()
			}
	 	}
	}

	return
}
