package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func validate(value any, rules string) error {
	rulesList := strings.Split(rules, "|")
	for _, rule := range rulesList {
		if !isRuleValid(value, rule) {
			return fmt.Errorf("%s", rule)
		}
	}
	return nil
}

func isRuleValid(value any, rule string) bool {
	if rule == "" {
		return true
	}

	rulePairs := strings.Split(rule, ":")

	if rulePairs[0] == "required" {
		switch value.(type) {
		case nil:
			return false
		case string:
			return value.(string) != ""
		case []any:
			return len(value.([]any)) > 0
		default:
			return true
		}
	}

	if rulePairs[0] == "email" {
		_, err := mail.ParseAddress(value.(string))
		return err == nil
	}

	if rulePairs[0] == "numeric" {
		_, ok := value.(float64)
		return ok
	}

	if rulePairs[0] == "min" {
		min, err := strconv.Atoi(rulePairs[1])
		return err == nil && len(value.(string)) >= min
	}

	if rulePairs[0] == "in" {
		valid := strings.Split(rulePairs[1], ",")
		return in(value.(string), valid...)
	}

	if rulePairs[0] == "date" {
		_, err := time.Parse(time.DateOnly, value.(string))
		return err == nil
	}

	if rulePairs[0] == "array" {
		switch value.(type) {
		case []any:
			return true
		default:
			return false
		}
	}

	if rulePairs[0] == "boolean" {
		switch value.(type) {
		case bool:
			return true
		default:
			return false
		}
	}

	fmt.Printf("Unknown rule %v\n", rule)

	return true
}

func in(value string, options ...string) bool {
	for _, option := range options {
		if option == value {
			return true
		}
	}
	return false
}

func Valid[T any](data map[string]any) (T, map[string]string) {
	var v T
	vType := reflect.TypeOf(v)
	vIsPtr := vType.Kind() == reflect.Ptr

	if vIsPtr { // If T is a pointer, get the underlying element type
		vType = vType.Elem()
	}

	if vType.Kind() != reflect.Struct {
		return v, map[string]string{"error": "T must be a struct or a pointer to a struct"}
	}

	val := reflect.New(vType).Elem()

	errors := make(map[string]string)

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)

		fieldValue := val.Field(i)

		name := field.Tag.Get("json")
		rules := field.Tag.Get("validate")
		err := validate(data[name], rules)
		if err != nil {
			errors[name] = err.Error()
		} else {
			if data[name] == nil {
				continue
			}
			if err := setField(fieldValue, data[name]); err != nil {
				errors[name] = err.Error()
			}
		}
	}

	if len(errors) > 0 {
		return v, errors
	}

	if vIsPtr { // if T is pointer
		v = val.Addr().Interface().(T)
	} else { // If T is struct
		v = val.Interface().(T)
	}

	return v, nil
}

func DecodeValid[T any](r *http.Request) (T, map[string]string) {
	data := make(map[string]any)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		// NOTE: does this error needs to halt.
		// the error here can happen if the body is not valid json, or empty.
	}

	return Valid[T](data)
}

func setField(fieldValue reflect.Value, value interface{}) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(fmt.Sprint(value))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, ok := value.(float64); !ok {
			return fmt.Errorf("expected numeric value")
		} else {
			fieldValue.SetInt(int64(intValue))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, ok := value.(float64); !ok {
			return fmt.Errorf("expected numeric value")
		} else {
			fieldValue.SetUint(uint64(uintValue))
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, ok := value.(float64); !ok {
			return fmt.Errorf("expected float value")
		} else {
			fieldValue.SetFloat(floatValue)
		}
	case reflect.Bool:
		if boolValue, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean value")
		} else {
			fieldValue.SetBool(boolValue)
		}
	case reflect.Slice:
		if value, ok := value.([]any); !ok {
			return fmt.Errorf("expected slice of strings")
		} else {
			if v, err := convertAnyToStringSlice(value); err != nil {
				return fmt.Errorf("expected slice of strings")
			} else {
				fieldValue.Set(reflect.ValueOf(v))
			}
		}
	default:
		return fmt.Errorf("unsupported type: %v", fieldValue)
	}
	return nil
}

func convertAnyToStringSlice(anySlice []any) ([]string, error) {
	result := make([]string, len(anySlice))

	for i, v := range anySlice {
		switch value := v.(type) {
		case string:
			result[i] = value
		case fmt.Stringer:
			result[i] = value.String()
		default:
			result[i] = fmt.Sprintf("%v", value)
		}
	}

	return result, nil
}
