package helper

import (
	"encoding/json"
	"fmt"
)

func DebugJSON(data any) error {
	json, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}
