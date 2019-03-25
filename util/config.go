package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/stvp/go-toml-config"
)

//ParseTOMLStruct 解析TOML文件到结构体
func ParseTOMLStruct(file string, obj interface{}) error {
	registerTOMLConfig("", obj)
	return config.Parse(file)
}

func registerTOMLConfig(prefix string, obj interface{}) {
	objT := reflect.TypeOf(obj)
	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		panic(fmt.Sprintf("%v must be a struct pointer", obj))
	}

	objT = objT.Elem()
	objV := reflect.ValueOf(obj).Elem()
	for i := 0; i < objT.NumField(); i++ {
		fieldT := objT.Field(i)
		configName := strings.TrimSpace(fieldT.Tag.Get("config"))
		if len(configName) == 0 {
			continue
		}
		defaultDefine := strings.TrimSpace(fieldT.Tag.Get("default"))
		if len(prefix) > 0 {
			configName = prefix + "." + configName
		}
		fieldV := objV.Field(i)
		switch fieldT.Type.Kind() {
		case reflect.Bool:
			defaultDefine = strings.ToUpper(defaultDefine)
			config.BoolVar(fieldV.Addr().Interface().(*bool), configName, defaultDefine == "1" || defaultDefine == "TRUE")
		case reflect.Int:
			x, err := strconv.ParseInt(defaultDefine, 10, 64)
			if err != nil && len(defaultDefine) > 0 {
				panic(err)
			}
			config.IntVar(fieldV.Addr().Interface().(*int), configName, int(x))
		case reflect.Int64:
			if fieldT.Type.String() == "time.Duration" {
				x, err := time.ParseDuration(defaultDefine)
				if err != nil && len(defaultDefine) > 0 {
					panic(err)
				}
				config.DurationVar(fieldV.Addr().Interface().(*time.Duration), configName, x)
			} else {
				x, err := strconv.ParseInt(defaultDefine, 10, 64)
				if err != nil && len(defaultDefine) > 0 {
					panic(err)
				}
				config.Int64Var(fieldV.Addr().Interface().(*int64), configName, x)
			}
		case reflect.Uint:
			x, err := strconv.ParseUint(defaultDefine, 10, 64)
			if err != nil && len(defaultDefine) > 0 {
				panic(err)
			}
			config.UintVar(fieldV.Addr().Interface().(*uint), configName, uint(x))
		case reflect.Uint64:
			x, err := strconv.ParseUint(defaultDefine, 10, 64)
			if err != nil && len(defaultDefine) > 0 {
				panic(err)
			}
			config.Uint64Var(fieldV.Addr().Interface().(*uint64), configName, x)
		case reflect.Float64:
			x, err := strconv.ParseFloat(defaultDefine, 64)
			if err != nil && len(defaultDefine) > 0 {
				panic(err)
			}
			config.Float64Var(fieldV.Addr().Interface().(*float64), configName, x)
		case reflect.String:
			config.StringVar(fieldV.Addr().Interface().(*string), configName, defaultDefine)
		case reflect.Struct:
			registerTOMLConfig(configName, fieldV.Addr().Interface())
		case reflect.Ptr:
			v := reflect.New(fieldV.Type().Elem())
			fieldV.Set(v)
			registerTOMLConfig(configName, v.Interface())
		case reflect.Map:
			if len(prefix) > 0 {
				panic("config name can not be more than two fields")
			}
			fieldV.Set(reflect.MakeMap(fieldT.Type))
			tags := strings.Split(configName, ",")
			for _, tag := range tags {
				var v reflect.Value
				switch fieldV.Type().Elem().Kind() {
				case reflect.Ptr:
					v = reflect.New(fieldV.Type().Elem().Elem())
				default:
					panic("value part of map type must be a struct pointer")
				}
				fieldV.SetMapIndex(reflect.ValueOf(tag), v)
				registerTOMLConfig(tag, v.Interface())
			}
		}
	}
}
