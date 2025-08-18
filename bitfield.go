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

// BitField is a type constraint that matches any integer type.
// It represents a value that can be used as a bitfield to store
// multiple boolean flags using bitwise operations.
type BitField interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// BitFieldAdd returns a new bitfield with the specified bitmasks set.
// Each bitmask corresponds to a flag value that will be added (ORed)
// into the bitfield.
//
// Example:
//
//	var flags uint8 = 0
//	flags = BitFieldAdd(flags, 1, 4) // sets bit 0 and 2 → flags = 5
func BitFieldAdd[T BitField](bitfield T, bitmasks ...T) T {
	for _, bitmask := range bitmasks {
		bitfield |= bitmask
	}
	return bitfield
}

// BitFieldRemove returns a new bitfield with the specified bitmasks cleared.
// Each bitmask corresponds to a flag value that will be removed (AND NOTed)
// from the bitfield.
//
// Example:
//
//	var flags uint8 = 5 // 0101
//	flags = BitFieldRemove(flags, 1) // clears bit 0 → flags = 4
func BitFieldRemove[T BitField](bitfield T, bitmasks ...T) T {
	for _, bitmask := range bitmasks {
		bitfield &^= bitmask
	}
	return bitfield
}

// BitFieldHas reports whether the given bitfield contains all of the specified
// bitmasks. It returns true if every bitmask is fully present in the bitfield.
//
// Example:
//
//	const (
//		A uint8 = 1 << iota // 0001
//		B                   // 0010
//		C                   // 0100
//	)
//
//	var flags = A | C // 0101
//
//	BitFieldHas(flags, A)    // true
//	BitFieldHas(flags, B)    // false
//	BitFieldHas(flags, A, C) // true
func BitFieldHas[T BitField](bitfield T, bitmasks ...T) bool {
	for _, bitmask := range bitmasks {
		if bitfield&bitmask != bitmask {
			return false
		}
	}
	return true
}

// BitFieldMissing returns a bitfield containing the subset of bitmasks
// that are not present in the given bitfield. If all specified bitmasks
// are already set, it returns zero.
//
// Example:
//
//	const (
//		A uint8 = 1 << iota // 0001
//		B                   // 0010
//		C                   // 0100
//	)
//
//	var flags = A | C // 0101
//
//	BitFieldMissing(flags, A, B, C) // 0010 (B is missing)
//	BitFieldMissing(flags, A)       // 0000 (A is present)
//	BitFieldMissing(flags, B, C)    // 0010 (only B is missing)
func BitFieldMissing[T BitField](bitfield T, bitmasks ...T) T {
	var missing T
	for _, bitmask := range bitmasks {
		if bitfield&bitmask == 0 {
			missing |= bitmask
		}
	}
	return missing
}
