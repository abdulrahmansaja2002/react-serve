package binder

import (
	"echo-react-serve/server/models/dto"
	"encoding/json"
	"fmt"

	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func BindMultipartForm(c echo.Context, out interface{}) error {
	return bindMultipartForm(c, out, "")
}

func bindMultipartForm(c echo.Context, out interface{}, parentKey string) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	v := reflect.ValueOf(out).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// Generate the form-data key
		key := parentKey
		if key != "" {
			key += "-"
		}
		key += formTag

		// Handle nested fields
		if field.Type.Kind() == reflect.Struct {
			// Initialize the nested struct if it's nil
			// if fieldValue.IsNil() {
			// 	fieldValue.Set(reflect.New(field.Type).Elem())
			// }
			nested := fieldValue.Addr().Interface()
			if err := bindMultipartForm(c, nested, key); err != nil {
				return err
			}
			continue
		}

		// Handle slices of files
		// if field.Type.Kind() == reflect.Slice && field.Type.Elem() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
		if strings.HasPrefix(key, "files") {
			files := form.File[key]
			slice := reflect.MakeSlice(field.Type, len(files), len(files))
			for i, file := range files {
				slice.Index(i).Set(reflect.ValueOf(file))
			}
			fieldValue.Set(slice)
			continue
		}
		if strings.HasPrefix(key, "target") {
			targetStrs := form.Value[key] // form.Value[key] should be form[key]
			var targets []dto.MemberTarget
			for _, targetStr := range targetStrs {
				var target dto.MemberTarget
				if err := json.Unmarshal([]byte(targetStr), &target); err != nil {
					return fmt.Errorf("Error unmarshaling target: %v")
				}
				targets = append(targets, target)
			}
			fieldValue.Set(reflect.ValueOf(targets))
			continue
		}

		// Handle single files
		if field.Type == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			if files, ok := form.File[key]; ok && len(files) > 0 {
				fieldValue.Set(reflect.ValueOf(files[0]))
			}
			continue
		}

		// Handle regular form values
		switch field.Type.Kind() {
		case reflect.String:
			if values, ok := form.Value[key]; ok && len(values) > 0 {
				fieldValue.SetString(values[0])
			}
		case reflect.Int, reflect.Int64:
			if values, ok := form.Value[key]; ok && len(values) > 0 {
				if intVal, err := strconv.Atoi(values[0]); err == nil {
					fieldValue.SetInt(int64(intVal))
				}
			}
		}
	}

	return nil
}
