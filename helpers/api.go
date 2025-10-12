package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func PostJSONRequest(
	url string,
	payload interface{},
	respObj interface{},
	maxRetries int,
	backoff time.Duration,
) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			lastErr = err
		} else {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = err
			} else if err := json.Unmarshal(body, respObj); err != nil {
				lastErr = err
			} else {
				// Success, return nil error
				return nil
			}
		}
		if attempt < maxRetries {
			time.Sleep(backoff * time.Duration(attempt)) // exponential backoff
		}
	}
	return errors.New("PostJSONRequest failed after retries: " + lastErr.Error())
}
