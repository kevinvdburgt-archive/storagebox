package storage

import (
	"fmt"
	"net/http"
	"strings"
)

func mkdirs(baseURL, username, password, path string) error {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		parts = parts[:len(parts)-1]
	}
	current := baseURL
	client := &http.Client{}

	for _, part := range parts {
		if part == "" {
			continue
		}

		current += "/" + part

		req, _ := http.NewRequest("MKCOL", current, nil)
		req.SetBasicAuth(username, password)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("MKCOL failed: %w", err)
		}

		resp.Body.Close()
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to create dir %s (status %d)", current, resp.StatusCode)
		}
	}
	return nil
}
