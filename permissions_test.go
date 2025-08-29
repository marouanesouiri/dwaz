/************************************************************************************
 *
 * dwaz (Discord Wrapper API for Zwafriya), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package dwaz

import "testing"

func TestPermissionsNames(t *testing.T) {
	perms := PermissionAddReactions | PermissionBanMembers

	names := perms.Names()

	if len(names) != 2 {
		t.Errorf("Expeced the names slice len to be '2' got %d", len(names))
	}

	for _, name := range names {
		if name != PermissionNameAddReactions && name != PermissionNameBanMembers {
			t.Errorf("Expeced the name to be 'AddReactions' or 'BanMembers' go %s", name)
		}
	}
}
