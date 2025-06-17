package db

import (
	"echo-react-serve/server/models/entity"
	"reflect"
)

// MergeMemberFiles merges the source MemberFiles into the target MemberFiles
func MergeMemberFiles(target interface{}, source interface{}) error {
	targetValue := reflect.ValueOf(target).Elem()
	sourceValue := reflect.ValueOf(source).Elem()

	// Iterate over all fields of the MemberFiles struct
	for i := 0; i < targetValue.NumField(); i++ {
		fieldName := targetValue.Type().Field(i).Name
		targetField := targetValue.Field(i)
		sourceField := sourceValue.FieldByName(fieldName)

		// Check if the field exists in the source
		if !sourceField.IsValid() {
			continue
		}

		// Handle nested structs (e.g., MemberPrograms)
		if targetField.Kind() == reflect.Struct {
			nestedTarget := targetField.Addr().Interface()
			nestedSource := sourceField.Addr().Interface()
			err := MergeMemberFiles(nestedTarget, nestedSource)
			if err != nil {
				return err
			}
			continue
		}

		// Check if the field is a slice of File
		if targetField.Kind() == reflect.Slice || targetField.Type().Elem() == reflect.TypeOf(entity.File{}) {
			// Append source files to the target field
			targetField.Set(reflect.AppendSlice(targetField, sourceField))
		}
	}

	return nil
}

// func main() {
// 	// Example usage
// 	target := &MemberFiles{
// 		Loa: []File{{Name: "loa1.pdf"}},
// 		Program: MemberPrograms{
// 			Pip: []File{{Name: "pip1.pdf"}},
// 		},
// 	}

// 	source := &MemberFiles{
// 		Loa: []File{{Name: "loa2.pdf"}},
// 		Program: MemberPrograms{
// 			Pip: []File{{Name: "pip2.pdf"}},
// 		},
// 	}

// 	err := MergeMemberFiles(target, source)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Println("Updated Loa:", target.Loa)
// 		fmt.Println("Updated Program.Pip:", target.Program.Pip)
// 	}
// }
