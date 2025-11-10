package loader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/osamikoyo/music-and-marks/logger"
	"go.uber.org/zap"
)

type Loader struct {
	logger *logger.Logger
}

func NewLoader(logger *logger.Logger) *Loader {
	return &Loader{
		logger: logger,
	}
}

func (l *Loader) SearchArtists(query string, limit int) (*SearchResult, error) {
	l.logger.Info("setuping search artists request",
		zap.String("query", query),
		zap.Int("limit", limit))

	baseURL := "https://musicbrainz.org/ws/2/artist"
	url := fmt.Sprintf("%s?query=%s&fmt=json&limit=%d", baseURL, query, limit)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		l.logger.Error("failed to create request",
			zap.String("url", url),
			zap.Error(err))

		return nil, fmt.Errorf("failed create request: %w", err)
	}

	req.Header.Set("User-Agent", "Music-service/1.0")

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		l.logger.Error("failed send request",
			zap.Error(err))

		return nil, fmt.Errorf("failed send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		l.logger.Error("http error",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)))

		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.logger.Error("failed read request body",
			zap.Error(err))

		return nil, fmt.Errorf("failed send response: %w", err)
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		l.logger.Error("failed parse result",
			zap.Error(err))

		return nil, fmt.Errorf("failed parse result: %w", err)
	}

	return &result, nil
}

func (l *Loader) SearchRelease(query string, limit, offset int) (*ReleaseSearchResult, error) {
	baseURL := "https://musicbrainz.org/ws/2/release"
	params := url.Values{}
	params.Add("query", query)
	params.Add("fmt", "json")
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("offset", fmt.Sprintf("%d", offset))

	reqURL := baseURL + "?" + params.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MyGoMusicApp/1.0 (your-email@example.com)")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ReleaseSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return &result, nil
}
