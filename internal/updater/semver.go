package updater

import (
	"strconv"
	"strings"
)

// NormalizeVersion trims leading 'v' or 'V' and whitespace.
func NormalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

// CompareVersions compares two semver strings (e.g. "1.2.0" and "1.3.0").
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func CompareVersions(v1, v2 string) int {
	norm1 := NormalizeVersion(v1)
	norm2 := NormalizeVersion(v2)

	if norm1 == norm2 {
		return 0
	}

	// Handle dev / non-semver
	if norm1 == "dev" {
		return -1
	}
	if norm2 == "dev" {
		return 1
	}

	// Split prerelease tag (e.g. 1.0.0-rc1)
	parts1 := strings.SplitN(norm1, "-", 2)
	parts2 := strings.SplitN(norm2, "-", 2)

	nums1 := strings.Split(parts1[0], ".")
	nums2 := strings.Split(parts2[0], ".")

	maxLen := len(nums1)
	if len(nums2) > maxLen {
		maxLen = len(nums2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(nums1) {
			n1, _ = strconv.Atoi(nums1[i])
		}
		if i < len(nums2) {
			n2, _ = strconv.Atoi(nums2[i])
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	// Check prereleases: normal release has higher precedence than prerelease (1.0.0 > 1.0.0-rc1)
	if len(parts1) > 1 && len(parts2) == 1 {
		return -1
	}
	if len(parts1) == 1 && len(parts2) > 1 {
		return 1
	}
	if len(parts1) > 1 && len(parts2) > 1 {
		if parts1[1] < parts2[1] {
			return -1
		} else if parts1[1] > parts2[1] {
			return 1
		}
	}

	return 0
}
