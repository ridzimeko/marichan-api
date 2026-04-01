package traceanime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"marichan-api/internal/pkg/httpclient"
	"mime/multipart"
	"net/http"
)

type Service struct {
	client *http.Client
}

type TraceAnimeResponse struct {
	FrameCount int    `json:"frameCount"`
	Error      string `json:"error"`
	Result     []struct {
		Anilist    int     `json:"anilist"`
		Filename   string  `json:"filename"`
		Episode    any     `json:"episode"`
		From       float64 `json:"from"`
		To         float64 `json:"to"`
		Similarity float64 `json:"similarity"`
		Video      string  `json:"video"`
		Image      string  `json:"image"`
	} `json:"result"`
}

func NewService() *Service {
	return &Service{client: httpclient.New()}
}

func (s *Service) SearchByImage(ctx context.Context, file *multipart.FileHeader) (*TraceAnimeResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, src); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.trace.moe/search", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("trace anime request failed: %s", string(raw))
	}

	var result TraceAnimeResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
