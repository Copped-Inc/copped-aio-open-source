package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func (s *Settings) Save() error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(os.Getenv("APPDATA"), "Copped AIO", "payments.json"), b, 0644)
}
