package sticker

import (
	"context"
	"fmt"
	"image"
	"marichan-api/internal/config"
	"marichan-api/internal/pkg/fileutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type Service struct {
	env *config.Env
}

type ConvertResult struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func NewService(env *config.Env) *Service {
	return &Service{env: env}
}

func (s *Service) ConvertToSticker(ctx context.Context, file *multipart.FileHeader) (*ConvertResult, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowed[ext] {
		return nil, fmt.Errorf("unsupported file format")
	}

	tempPath, _, err := fileutil.SaveUploadedFile(file, s.env.TempDir)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempPath)

	img, err := imaging.Open(tempPath)
	if err != nil {
		return nil, err
	}

	resultImg := fitStickerCanvas(img)

	if err := fileutil.EnsureDir(s.env.UploadDir); err != nil {
		return nil, err
	}

	outName := strings.TrimSuffix(filepath.Base(tempPath), filepath.Ext(tempPath)+".webp")
	outPath := filepath.Join(s.env.UploadDir, outName)

	outFile, err := os.Create(outPath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	if err := webp.Encode(outFile, resultImg, &webp.Options{Lossless: true}); err != nil {
		return nil, err
	}

	bounds := resultImg.Bounds()

	return &ConvertResult{
		FileName: outName,
		FileURL:  fileutil.BuildFileUrl(s.env.BaseURL, "files", outName),
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
	}, nil
}

func fitStickerCanvas(img image.Image) image.Image {
	const size = 512

	resized := imaging.Fit(img, size, size, imaging.Lanczos)
	canvas := imaging.New(size, size, image.Transparent)
	return imaging.PasteCenter(canvas, resized)
}
