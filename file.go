/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile downloads a file from the given URL and saves it in the specified directory.
// The filename is derived from baseName and the Content-Type returned by the server.
// Returns the full path of the saved file.
func DownloadFile(url, baseName, dir string) (string, error) {
	if url == "" {
		return "", errors.New("URL is empty")
	}

	resp, err := http.Head(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch headers: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch headers: status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		exts = []string{filepath.Ext(baseName)}
	}
	ext := exts[0]

	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	finalName := name + ext
	fullPath := filepath.Join(dir, finalName)

	respGet, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file: status %d", respGet.StatusCode)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, respGet.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fullPath, nil
}

// Base64Image represents a base64-encoded image data URI string.
type Base64Image = string

// NewImageFile reads an image file and returns its base64 data URI string.
//
// Example output: "data:image/png;base64,<base64-encoded-bytes>"
func NewImageFile(path string) (Base64Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	mimeType := http.DetectContentType(data)
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("not an image file: detected MIME type %s", mimeType)
	}

	if _, _, err := image.DecodeConfig(bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("invalid image data: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
