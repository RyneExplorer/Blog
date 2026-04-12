package utils

// TruncateRunes 按 Unicode 字符截取前 max 个字符（用于摘要）
func TruncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
