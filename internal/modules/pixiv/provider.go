package pixiv

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"marichan-api/internal/config"
	"net/http"
	"net/url"
	"strings"
)

type Provider struct {
	env    *config.Env
	client *http.Client
}

func NewProvider(env *config.Env, client *http.Client) *Provider {
	return &Provider{
		env:    env,
		client: client,
	}
}

func (p *Provider) newRequest(ctx context.Context, method, path, referer string) (*http.Request, error) {
	fullURL := strings.TrimRight(p.env.PixivBaseURL, "/") + path

	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "Application/json")
	req.Header.Set("User-Agent", p.env.PixivUserAgent)

	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	if p.env.PixivPHPSESSID != "" {
		req.Header.Set("Cookie", "PHPSESSID="+p.env.PixivPHPSESSID)
	}

	return req, nil
}

func doJSON[T any](p *Provider, req *http.Request) (*T, error) {
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(raw, &errResp); err == nil && errResp.Message != "" {
			return nil, fmt.Errorf("%s", errResp.Message)
		}
		return nil, fmt.Errorf("pixiv request failed: status=%d body=%s", resp.StatusCode, string(raw))
	}

	if len(raw) == 0 {
		return nil, fmt.Errorf("Pixiv returned empty body")
	}

	var parsed AjaxResponse[T]
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("decode pixiv response: %w", err)
	}

	if parsed.Error {
		return nil, fmt.Errorf("pixiv error: %s", parsed.Message)
	}

	return &parsed.Body, nil
}

func (p *Provider) GetIllustDetail(ctx context.Context, artworkID string) (*illusDetailBody, error) {
	path := fmt.Sprintf("/ajax/illust/%s", artworkID)
	referer := fmt.Sprintf("%s/artworks/%s", strings.TrimRight(p.env.PixivBaseURL, "/"), artworkID)

	req, err := p.newRequest(ctx, http.MethodGet, path, referer)
	if err != nil {
		return nil, err
	}

	return doJSON[illusDetailBody](p, req)
}

func (p *Provider) GetIllustPages(ctx context.Context, artworkID string) (*IllustPagesBody, error) {
	path := fmt.Sprintf("/ajax/illust/%s/pages", artworkID)
	referer := fmt.Sprintf("%s/artworks/%s", strings.TrimRight(p.env.PixivBaseURL, "/"), artworkID)

	req, err := p.newRequest(ctx, http.MethodGet, path, referer)
	if err != nil {
		return nil, err
	}

	return doJSON[IllustPagesBody](p, req)
}

func (p *Provider) searchArtworks(ctx context.Context, keyword, order, mode, sMode string, page int) (*SearchArtworksBody, error) {
	escaped := url.PathEscape(keyword)

	q := url.Values{}
	q.Set("word", keyword)
	q.Set("type", "all")
	q.Set("order", order)
	q.Set("mode", mode)
	q.Set("s_mode", sMode)

	if page > 0 {
		q.Set("p", fmt.Sprintf("%d", page))
	}

	path := fmt.Sprintf("/ajax/search/artworks/%s?%s", escaped, q.Encode())

	req, err := p.newRequest(ctx, http.MethodGet, path, p.env.PixivBaseURL+"/")
	if err != nil {
		return nil, err
	}

	return doJSON[SearchArtworksBody](p, req)
}

func (p *Provider) GetUserFull(ctx context.Context, userID string) (*UserFullBody, error) {
	path := fmt.Sprintf("/ajax/user/%s?full=1", userID)
	referer := fmt.Sprintf("%s/member.php?id=%s", strings.TrimRight(p.env.PixivBaseURL, "/"), userID)

	req, err := p.newRequest(ctx, http.MethodGet, path, referer)
	if err != nil {
		return nil, err
	}

	return doJSON[UserFullBody](p, req)
}

func (p *Provider) DownloadImage(ctx context.Context, imageURL string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", p.env.PixivUserAgent)
	req.Header.Set("Referer", "https://www.pixiv.net/")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("failed to download image, status: %d", resp.StatusCode)
	}

	return resp, nil
}
