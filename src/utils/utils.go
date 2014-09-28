package utils

import "strings"
import "strconv"
import "path/filepath"
import "regexp"

func NormalizeBoolean(val string, def bool) bool {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "true", "yes", "on", "1":
		return true
	case "false", "no", "off", "0":
		return false
	default:
		return def
	}
}

func NormalizeInt(val string, def int) int {
	if res, err := strconv.Atoi(val); err != nil {
		return res
	} else {
		return def
	}
}

func SanitizeFilename(name string) string {
	name = filepath.Clean(name)
	rex, _ := regexp.Compile(`[^\w]+`)
	name = rex.ReplaceAllString(name, "-")
	return strings.TrimSpace(name)
}
