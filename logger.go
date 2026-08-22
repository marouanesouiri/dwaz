/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

// Logger is the structured logging interface used throughout wumpkit.
//
// All methods accept alternating key-value pairs (slog style) for structured context.
// Standard library *slog.Logger natively implements this interface.
// If no custom logger is provided to the Client, it defaults to log/slog.Default().
type Logger interface {
	// Debug logs at debug level with structured key-value pairs.
	Debug(msg string, args ...any)
	// Info logs at info level with structured key-value pairs.
	Info(msg string, args ...any)
	// Warn logs at warn level with structured key-value pairs.
	Warn(msg string, args ...any)
	// Error logs at error level with structured key-value pairs.
	Error(msg string, args ...any)
}
