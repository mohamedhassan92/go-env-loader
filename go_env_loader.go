package goenvloader


import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"github.com/joho/godotenv"
)

//EnvConfig adapter implementation for loading config from the env
type EnvConfig struct {
}

//Init initializes the adapter
func (adp *EnvConfig) Init() error {
	if err := godotenv.Load(".env.local"); err != nil { // loads .env.local under the root folder if exists
		fmt.Println("No .env.local file found")
	}
	return nil
}

//Load : loads entity from the env variables based on the field tag "env"
func (adp *EnvConfig) Load(i interface{}) error {
	rv := reflect.ValueOf(i)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidLoadError{reflect.TypeOf(i)}
	}
	v := rv.Elem()
	t := v.Type()
	for j := 0; j < t.NumField(); j++ {
		tf := t.Field(j)
		tag := tf.Tag.Get("env")
		if tag == "" {
			return &InvalidLoadTagError{tf.Name}
		}
		vf := v.FieldByName(tf.Name)
		switch tf.Type.Kind() {
		case reflect.String:
			vf.SetString(getEnv(tag, ""))
		case reflect.Int:
			vf.SetInt(int64(getEnvAsInt(tag, 0)))
		case reflect.Bool:
			vf.SetBool(getEnvAsBool(tag, false))
		default:
			return &UnSupportedFieldTypeError{tf}
		}
	}
	return nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// An InvalidLoadError describes an invalid argument passed to Load.
// (The argument to Load must be a non-nil pointer.)
type InvalidLoadError struct {
	Type reflect.Type
}

func (e *InvalidLoadError) Error() string {
	if e.Type == nil {
		return "env: Load(nil)"
	}

	if e.Type.Kind() != reflect.Struct {
		return "env: Load(non-struct " + e.Type.String() + ")"
	}
	return "env: Load(nil " + e.Type.String() + ")"
}

// An InvalidLoadTagError describes a missing env tag on the struct passed to Load.
// (All the fields under the struct must have the env tag)
type InvalidLoadTagError struct {
	FN string
}

func (e *InvalidLoadTagError) Error() string {
	return "env: field " + e.FN + " is missing the tag 'env'"
}

// An UnSupportedFieldTypeError describes an unsupported field type on the struct passed to Load.
// (All the fields under the struct must have supported types)
type UnSupportedFieldTypeError struct {
	SF reflect.StructField
}

func (e *UnSupportedFieldTypeError) Error() string {
	return "env: field type " + e.SF.Type.String() + " is not supported for field " + e.SF.Name + ""
}
