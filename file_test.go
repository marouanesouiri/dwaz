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
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func decodeBase64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic("invalid base64 in test: " + err.Error())
	}
	return b
}

func TestNewImageFile(t *testing.T) {
	// small valid PNG image (1x1 transparent pixel) (thanks ChatGPT for that)
	// "iVBORw0K" is the Base64 representation of the PNG file signature (magic number)
	const base64Png = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR4nGNgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="

	imgPath := filepath.Join(t.TempDir(), "test.png")
	if err := os.WriteFile(imgPath, decodeBase64(base64Png), 0644); err != nil {
		t.Fatalf("failed to write temp image file: %v", err)
	}

	dataURI, err := NewImageFile(imgPath)
	if err != nil {
		t.Fatalf("unexpected error from NewImageFile: %v", err)
	}

	if !strings.HasPrefix(dataURI, "data:image/png;base64,") {
		t.Errorf("unexpected data URI prefix: got %q", dataURI[:30])
	}
}

func TestNewImageFile_NonImage(t *testing.T) {
	txtPath := filepath.Join(t.TempDir(), "not_image.txt")
	if err := os.WriteFile(txtPath, []byte("hello world"), 0644); err != nil {
		t.Fatalf("failed to write text file: %v", err)
	}

	_, err := NewImageFile(txtPath)
	if err == nil {
		t.Error("expected error for non-image file, got nil")
	}
}
