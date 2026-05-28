package dvhcvn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *Service) fetchAndParse(ctx context.Context) ([]Province, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTP, err)
	}
	if s.cfg.userAgent != "" {
		req.Header.Set("User-Agent", s.cfg.userAgent)
	}

	resp, err := s.cfg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTP, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: HTTP %d", ErrHTTP, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTP, err)
	}

	var provinces []Province
	if err := json.Unmarshal(body, &provinces); err != nil {
		return nil, fmt.Errorf("%w: JSON decode: %w", ErrInvalidSchema, err)
	}

	return provinces, nil
}
