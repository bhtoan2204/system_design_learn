package util

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

func NormalizeType(typeName string) string {
	primitiveTypes := map[string]string{
		"double":        "float64",
		"float":         "float64",
		"int32":         "int32",
		"int64":         "int64",
		"uint32":        "uint32",
		"uint64":        "uint64",
		"sint32":        "int32",
		"sint64":        "int64",
		"fixed32":       "uint32",
		"fixed64":       "uint64",
		"sfixed32":      "int32",
		"sfixed64":      "int64",
		"bool":          "bool",
		"string":        "string",
		"bytes":         "[]byte",
		"optional bool": "*bool",
	}

	if normalized, ok := primitiveTypes[typeName]; ok {
		return normalized
	}
	return typeName
}

func GetAllProtoFiles(path string) []string {
	var protoPaths []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".proto") {
			protoPaths = append(protoPaths, path)
		}

		return nil
	})

	if err != nil {
		fmt.Println("failed to get all proto files, return whatever found, error:", err)
	}

	return protoPaths
}

func ReadFileAsString(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func CreateFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	// folder has already existed
	return nil
}

func LowerFirstChar(input string) string {
	if len(input) <= 1 {
		return input
	}

	return strings.ToLower(input[:1]) + input[1:]
}

func UpperFirstChar(input string) string {
	if len(input) <= 1 {
		return input
	}

	return strings.ToUpper(input[:1]) + input[1:]
}

func ToCamelCase(input string) string {
	return strcase.ToLowerCamel(input)
}

func ToPascalCase(input string) string {
	return strcase.ToCamel(input)
}

func ToLowerCase(input string) string {
	return strings.ToLower(input)
}

func ConvertSliceStrToUCWord(input []string) string {
	mapString := ""
	for _, str := range input {
		mapString += UpperFirstChar(str)
	}
	return mapString
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// TODO: remove this function
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 {
			parts[i] = strings.Title(parts[i])
		}
	}

	return strings.Join(parts, "")
}

func WriteFile(data []byte, filePath string) error {
	dir, _ := filepath.Split(filePath)

	if _, err := os.Stat(dir); err == nil {
		os.Remove(filePath)
	} else {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}
	return nil
}

func WriteFileWithFormat(content []byte, path string) error {
	var err error
	content, err = format.Source(content)
	if err != nil {
		return fmt.Errorf("FORMATTING_ERROR: %v", err)
	}
	return WriteFile(content, path)
}

func Find[T any](list []T, f func(T) bool) T {
	var newValue T
	for _, v := range list {
		if f(v) {
			newValue = v
			break
		}
	}
	return newValue
}
