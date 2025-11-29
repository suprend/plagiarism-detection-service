package config

import "os"

func ServerPort() string {
	if v := os.Getenv("PORT"); v != "" {
		return v
	}
	return "8082"
}

func FilestorageURL() string {
	if v := os.Getenv("FILESTORAGE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func PlagiarismURL() string {
	if v := os.Getenv("PLAGIARISM_URL"); v != "" {
		return v
	}
	return "http://localhost:8081"
}

func WordcloudURL() string {
	if v := os.Getenv("WORDCLOUD_URL"); v != "" {
		return v
	}
	return "https://quickchart.io/wordcloud"
}

func WordcloudDir() string {
	if v := os.Getenv("WORDCLOUD_DIR"); v != "" {
		return v
	}
	return "tmp-files/wordclouds"
}
