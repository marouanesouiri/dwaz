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
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAttachment_Save(t *testing.T) {
	const base64Png = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR4nGNgYAAAAAMAAWgmWQ0AAAAASUVORK5CYII="
	imageData, _ := base64.StdEncoding.DecodeString(base64Png)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(imageData)
	}))
	defer server.Close()

	att := &Attachment{
		Filename:    "avatar.jpg",
		ContentType: "image/png",
		URL:         server.URL,
	}

	tempDir := t.TempDir()
	err := att.Save(tempDir)
	if err != nil {
		t.Fatalf("Attachment.Save() error: %v", err)
	}

	wantPath := filepath.Join(tempDir, "avatar.png")
	if _, err := os.Stat(wantPath); err != nil {
		t.Fatalf("Expected file %s to be created, but got error: %v", wantPath, err)
	}
}
