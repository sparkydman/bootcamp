package utils

import "encoding/json"

func EmbedStructFlat[T, U any](data U) (T, error) {
	var value T
	d, err := json.Marshal(data)
	if err != nil {
		return value, err
	}

	if err := json.Unmarshal(d, &value); err != nil {
		return value, err
	}
	return value, nil
}
