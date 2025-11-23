package imagekit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type ImageKitClient struct {
	client *imagekit.ImageKit
}

type UploadOptions struct {
	Folder            string
	UseUniqueFileName bool
	Tags              []string
	CustomCoordinates string
	ResponseFields    []string
	IsPrivateFile     bool
	Transformation    *TransformationOptions
}

type TransformationOptions struct {
	Pre  string
	Post string
}

func NewImageKitClient() *ImageKitClient {
	if config.AppConfig.ImageKitPrivateKey == "" || config.AppConfig.ImageKitPublicKey == "" || config.AppConfig.ImageKitEndpoint == "" {
		return nil
	}

	ik := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  config.AppConfig.ImageKitPrivateKey,
		PublicKey:   config.AppConfig.ImageKitPublicKey,
		UrlEndpoint: config.AppConfig.ImageKitEndpoint,
	})

	return &ImageKitClient{
		client: ik,
	}
}

// UploadImage uploads an image file with validation. This is the recommended method for image uploads.
func (c *ImageKitClient) UploadImage(
	ctx context.Context,
	file io.Reader,
	fileName string,
	opts *UploadOptions,
) (*uploader.UploadResult, error) {
	// Validate image extension
	if !ValidateImageExtension(fileName) {
		return nil, fmt.Errorf("invalid image extension: %s. Supported formats: jpg, jpeg, png, gif, webp, svg, bmp, ico", GetFileExtension(fileName))
	}

	return c.UploadFileWithOptions(ctx, file, fileName, opts)
}

// UploadImageFromBytes uploads an image from byte data (common for web uploads)
func (c *ImageKitClient) UploadImageFromBytes(
	ctx context.Context,
	fileData []byte,
	fileName string,
	opts *UploadOptions,
) (*uploader.UploadResult, error) {
	if len(fileData) == 0 {
		return nil, fmt.Errorf("file data cannot be empty")
	}

	// Validate image extension
	if !ValidateImageExtension(fileName) {
		return nil, fmt.Errorf("invalid image extension: %s. Supported formats: jpg, jpeg, png, gif, webp, svg, bmp, ico", GetFileExtension(fileName))
	}

	file := bytes.NewReader(fileData)
	return c.UploadFileWithOptions(ctx, file, fileName, opts)
}

// UploadFile uploads a file to ImageKit with basic options (returns URL only)
func (c *ImageKitClient) UploadFile(file io.Reader, fileName string) (string, error) {
	result, err := c.UploadFileWithOptions(context.Background(), file, fileName, nil)
	if err != nil {
		return "", err
	}
	return result.Url, nil
}

// UploadFileWithOptions uploads a file with advanced options
func (c *ImageKitClient) UploadFileWithOptions(
	ctx context.Context,
	file io.Reader,
	fileName string,
	opts *UploadOptions,
) (*uploader.UploadResult, error) {
	if file == nil {
		return nil, fmt.Errorf("file reader cannot be nil")
	}

	if fileName == "" {
		fileName = fmt.Sprintf("upload_%d%s", time.Now().Unix(), ".jpg")
	}

	uploadParams := c.buildUploadParams(fileName, opts)

	uploadResponse, err := c.client.Uploader.Upload(ctx, file, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &uploadResponse.Data, nil
}

// UploadBase64 uploads a base64 encoded image
func (c *ImageKitClient) UploadBase64(
	ctx context.Context,
	base64Data string,
	fileName string,
	opts *UploadOptions,
) (*uploader.UploadResult, error) {
	if base64Data == "" {
		return nil, fmt.Errorf("base64 data cannot be empty")
	}

	if fileName == "" {
		fileName = fmt.Sprintf("upload_%d.jpg", time.Now().Unix())
	}

	// Validate image extension if fileName is provided
	if fileName != "" && !ValidateImageExtension(fileName) {
		return nil, fmt.Errorf("invalid image extension: %s. Supported formats: jpg, jpeg, png, gif, webp, svg, bmp, ico", GetFileExtension(fileName))
	}

	uploadParams := c.buildUploadParams(fileName, opts)

	uploadResponse, err := c.client.Uploader.Upload(ctx, base64Data, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload base64 image: %w", err)
	}

	return &uploadResponse.Data, nil
}

// UploadFromURL uploads an image from a URL
func (c *ImageKitClient) UploadFromURL(
	ctx context.Context,
	url string,
	fileName string,
	opts *UploadOptions,
) (*uploader.UploadResult, error) {
	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	if fileName == "" {
		// Extract filename from URL or generate one
		fileName = extractFileNameFromURL(url)
		if fileName == "" {
			fileName = fmt.Sprintf("upload_%d.jpg", time.Now().Unix())
		}
	}

	// Validate image extension if we can determine it
	if fileName != "" && !ValidateImageExtension(fileName) {
		return nil, fmt.Errorf("invalid image extension: %s. Supported formats: jpg, jpeg, png, gif, webp, svg, bmp, ico", GetFileExtension(fileName))
	}

	uploadParams := c.buildUploadParams(fileName, opts)

	uploadResponse, err := c.client.Uploader.Upload(ctx, url, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload from URL: %w", err)
	}

	return &uploadResponse.Data, nil
}

// DeleteFile deletes a file from ImageKit by file ID
func (c *ImageKitClient) DeleteFile(ctx context.Context, fileID string) error {
	if fileID == "" {
		return fmt.Errorf("file ID cannot be empty")
	}

	_, err := c.client.Media.DeleteFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// buildUploadParams constructs upload parameters
func (c *ImageKitClient) buildUploadParams(fileName string, opts *UploadOptions) uploader.UploadParam {
	params := uploader.UploadParam{
		FileName: fileName,
	}

	// Set default UseUniqueFileName to true if opts is nil
	if opts == nil {
		useUnique := true
		params.UseUniqueFileName = &useUnique
		return params
	}

	// Handle options
	if opts.Folder != "" {
		params.Folder = opts.Folder
	}

	params.UseUniqueFileName = &opts.UseUniqueFileName

	if len(opts.Tags) > 0 {
		params.Tags = strings.Join(opts.Tags, ",")
	}

	if opts.CustomCoordinates != "" {
		params.CustomCoordinates = opts.CustomCoordinates
	}

	if len(opts.ResponseFields) > 0 {
		params.ResponseFields = strings.Join(opts.ResponseFields, ",")
	}

	if opts.IsPrivateFile {
		val := true
		params.IsPrivateFile = &val
	}

	// Note: Transformation options are not supported in UploadParam
	// Transformations should be applied to the URL after upload or configured
	// in ImageKit dashboard for URL transformations

	return params
}

// extractFileNameFromURL extracts filename from URL
func extractFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		fileName := parts[len(parts)-1]
		// Remove query parameters if any
		if idx := strings.Index(fileName, "?"); idx != -1 {
			fileName = fileName[:idx]
		}
		return fileName
	}
	return ""
}

// GetFileExtension returns the file extension from a filename
func GetFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext)
}

// ValidateImageExtension checks if the file extension is a valid image format
func ValidateImageExtension(fileName string) bool {
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".bmp":  true,
		".ico":  true,
	}

	ext := GetFileExtension(fileName)
	return validExtensions[ext]
}
