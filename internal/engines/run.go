package engines

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

type RunOptions struct {
	PatchMarker      string   `json:"patchMarker"`
	UseGoBuildTag    bool     `json:"useGoBuildTag"`
	GoBuildTag       string   `json:"goBuildTag"`
	SkipEmpty        bool     `json:"skipEmpty"`
	TargetPath       string   `json:"targetPath"`
	TargetDir        string   `json:"targetDir"`
	IgnoreExtensions bool     `json:"ignoreExtensions"`
	TargetExtensions []string `json:"targetExtensions"`
}

type RunArgs struct {
	Engine     string     `json:"engine"`
	Args       EngineArgs `json:"args"`
	RunOptions RunOptions `json:"runOptions"`
}

// NewRunOptionsFromCLI creates a new RunOptions struct based on the CLI flags
// args should be a map of key-value pairs from CLI flags (e.g. from the -r flag)
func NewRunOptionsFromCLI(args map[string]string) RunOptions {
	// Start with the base options
	options := RunOptions{
		// init pointers fields
		TargetExtensions: make([]string, 0),
	}

	// Create a map of field names to field indices (lowercased for case-insensitive matching)
	optionsValue := reflect.ValueOf(&options).Elem()
	optionsType := optionsValue.Type()

	fieldMap := make(map[string]int)
	for i := 0; i < optionsType.NumField(); i++ {
		field := optionsType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Strip any tag options (like omitempty)
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
			fieldMap[strings.ToLower(jsonTag)] = i
		}
		// Also map by field name for convenience
		fieldMap[strings.ToLower(field.Name)] = i
	}

	for key, value := range args {
		// Look up field by lowercase name
		fieldIdx, exists := fieldMap[strings.ToLower(key)]
		if !exists {
			log.Printf("unknown run option: %s", key)
			continue
		}

		field := optionsValue.Field(fieldIdx)
		fieldType := optionsType.Field(fieldIdx)

		// Set the value according to the field's type
		switch field.Kind() {
		case reflect.String:
			field.SetString(value)

		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				log.Printf("invalid boolean value for %s: %s", key, value)
				continue
			}
			field.SetBool(boolValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Printf("invalid integer value for %s: %s", key, value)
				continue
			}
			field.SetInt(intValue)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintValue, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				log.Printf("invalid unsigned integer value for %s: %s", key, value)
				continue
			}
			field.SetUint(uintValue)

		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Printf("invalid float value for %s: %s", key, value)
				continue
			}
			field.SetFloat(floatValue)

		case reflect.Slice:
			// Handle slice types - currently only supporting []string
			if field.Type().Elem().Kind() == reflect.String {
				elements := strings.Split(value, ",")
				stringSlice := make([]string, 0, len(elements))
				for _, el := range elements {
					el = strings.TrimSpace(el)
					if el != "" {
						stringSlice = append(stringSlice, el)
					}
				}
				field.Set(reflect.ValueOf(stringSlice))
			} else {
				log.Printf("unsupported slice type for field %s", fieldType.Name)
			}

		default:
			log.Printf("unsupported type for field %s: %s", fieldType.Name, field.Type().String())
		}
	}

	return options
}
