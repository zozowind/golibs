package util

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/stvp/go-toml-config"
)

//TOMLToStruct 解析Toml文件到结构体
func TOMLToStruct(file string, obj interface{}) error {
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
	return config.Parse(file)
}

//StructToURLValue 结构体到url.Values
func StructToURLValue(v interface{}, t string) (url.Values, error) {
	values := make(url.Values)
	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return values, nil
		}
		val = val.Elem()
	}

	if v == nil {
		return values, nil
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("query: Values() expects struct input. Got %v", val.Kind())
	}

	err := reflectValue(values, val, "", t)
	return values, err
}

var timeType = reflect.TypeOf(time.Time{})

var encoderType = reflect.TypeOf(new(Encoder)).Elem()

//Encoder interface
type Encoder interface {
	EncodeValues(key string, v *url.Values) error
}

// reflectValue populates the values parameter from the struct fields in val.
// Embedded structs are followed recursively (using the rules defined in the
// Values function documentation) breadth-first.
func reflectValue(values url.Values, val reflect.Value, scope string, t string) error {
	var embedded []reflect.Value

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}

		sv := val.Field(i)
		tag := sf.Tag.Get(t)
		if tag == "-" {
			continue
		}
		name, opts := parseTag(tag)
		if name == "" {
			if sf.Anonymous && sv.Kind() == reflect.Struct {
				// save embedded struct for later processing
				embedded = append(embedded, sv)
				continue
			}

			name = sf.Name
		}

		if scope != "" {
			name = scope + "[" + name + "]"
		}

		if opts.Contains("omitempty") && isEmptyValue(sv) {
			continue
		}

		if sv.Type().Implements(encoderType) {
			if !reflect.Indirect(sv).IsValid() {
				sv = reflect.New(sv.Type().Elem())
			}

			m := sv.Interface().(Encoder)
			if err := m.EncodeValues(name, &values); err != nil {
				return err
			}
			continue
		}

		if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
			var del byte
			if opts.Contains("comma") {
				del = ','
			} else if opts.Contains("space") {
				del = ' '
			} else if opts.Contains("semicolon") {
				del = ';'
			} else if opts.Contains("brackets") {
				name = name + "[]"
			}

			if del != 0 {
				s := new(bytes.Buffer)
				first := true
				for i := 0; i < sv.Len(); i++ {
					if first {
						first = false
					} else {
						s.WriteByte(del)
					}
					s.WriteString(valueString(sv.Index(i), opts))
				}
				vStr := s.String()
				if vStr != "" {
					values.Add(name, vStr)
				}
			} else {
				for i := 0; i < sv.Len(); i++ {
					k := name
					if opts.Contains("numbered") {
						k = fmt.Sprintf("%s%d", name, i)
					}
					vStr := valueString(sv.Index(i), opts)
					if vStr != "" {
						values.Add(k, vStr)
					}
				}
			}
			continue
		}

		for sv.Kind() == reflect.Ptr {
			if sv.IsNil() {
				break
			}
			sv = sv.Elem()
		}

		if sv.Type() == timeType {
			vStr := valueString(sv, opts)
			if vStr != "" {
				values.Add(name, vStr)
			}
			continue
		}

		if sv.Kind() == reflect.Struct {
			reflectValue(values, sv, name, t)
			continue
		}

		vStr := valueString(sv, opts)
		if vStr != "" {
			values.Add(name, vStr)
		}
	}
	for _, f := range embedded {

		if err := reflectValue(values, f, scope, t); err != nil {
			return err
		}
	}

	return nil
}

// isEmptyValue checks if a value should be considered empty for the purposes
// of omitting fields with the "omitempty" option.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	if v.Type() == timeType {
		return v.Interface().(time.Time).IsZero()
	}

	return false
}

// tagOptions is the string following a comma in a struct field's "url" tag, or
// the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	s := strings.Split(tag, ",")
	return s[0], s[1:]
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	for _, s := range o {
		if s == option {
			return true
		}
	}
	return false
}

func valueString(v reflect.Value, opts tagOptions) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Bool && opts.Contains("int") {
		if v.Bool() {
			return "1"
		}
		return "0"
	}

	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		if opts.Contains("unix") {
			return strconv.FormatInt(t.Unix(), 10)
		}
		return t.Format(time.RFC3339)
	}

	return fmt.Sprint(v.Interface())
}

//URLValueToStruct URLValueToStruct
func URLValueToStruct(form url.Values, obj interface{}, t string) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		return fmt.Errorf("%v must be a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)
		tags := strings.Split(fieldT.Tag.Get(t), ",")
		var tag string
		if len(tags) == 0 || len(tags[0]) == 0 {
			tag = fieldT.Name
		} else if tags[0] == "-" {
			continue
		} else {
			tag = tags[0]
		}

		value := form.Get(tag)
		if len(value) == 0 {
			continue
		}

		var required = len(tags) > 1 && tags[1] == "required"
		switch fieldT.Type.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			fieldV.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			fieldV.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}
			fieldV.SetFloat(x)
		case reflect.Interface:
			fieldV.Set(reflect.ValueOf(value))
		case reflect.String:
			if required && value == "" {
				err := fmt.Errorf("%s is null", tag)
				return err
			}
			fieldV.SetString(value)
		}
	}
	return nil
}
