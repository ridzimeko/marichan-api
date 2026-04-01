package pixiv

import (
	"context"
	"marichan-api/internal/config"
	"marichan-api/internal/pkg/httpclient"
	"net/http"
)

type Service struct {
	provider *Provider
}

func NewService(env *config.Env) *Service {
	client := httpclient.New()
	return &Service{
		provider: NewProvider(env, client),
	}
}

type ArtworkDetailResult struct {
	Detail *illusDetailBody `json:"detail"`
	Pages  *IllustPagesBody `json:"pages"`
}

func (s *Service) GetArtworkDetail(ctx context.Context, artworkID string) (*ArtworkDetailResult, error) {
	detail, err := s.provider.GetIllustDetail(ctx, artworkID)
	if err != nil {
		return nil, err
	}

	pages, err := s.provider.GetIllustPages(ctx, artworkID)
	if err != nil {
		return nil, err
	}

	return &ArtworkDetailResult{
		Detail: detail,
		Pages:  pages,
	}, nil
}

func (s *Service) searchArtworks(ctx context.Context, keyword, order, mode, sMode string, page int) (*SearchArtworksBody, error) {
	if order == "" {
		order = "date_d"
	}
	if mode == "" {
		mode = "all"
	}
	if sMode == "" {
		sMode = "s_tag"
	}
	if page <= 0 {
		page = 1
	}

	return s.provider.searchArtworks(ctx, keyword, order, mode, sMode, page)
}

type ArtistDetailResult struct {
	User *UserFullBody `json:"user"`
}

func (s *Service) GetArtistDetail(ctx context.Context, userID string) (*ArtistDetailResult, error) {
	user, err := s.provider.GetUserFull(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &ArtistDetailResult{
		User: user,
	}, nil
}

func (s *Service) DownloadImage(ctx context.Context, imageURL string) (*http.Response, error) {
	return s.provider.DownloadImage(ctx, imageURL)
}
