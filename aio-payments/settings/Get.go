package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Get() (*Settings, error) {

	file, err := os.OpenFile(filepath.Join(os.Getenv("APPDATA"), "Copped AIO", "payments.json"), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = json.NewDecoder(file).Decode(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil

}
