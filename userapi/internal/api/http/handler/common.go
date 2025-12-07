package handler

import "strings"

func extractWorkID(path, suffix string) (string, bool) {
	if !strings.HasPrefix(path, "/works/") || !strings.HasSuffix(path, suffix) {
		return "", false
	}
	work := strings.TrimPrefix(path, "/works/")
	work = strings.TrimSuffix(work, suffix)
	work = strings.TrimSuffix(work, "/")
	if work == "" {
		return "", false
	}
	return work, true
}
