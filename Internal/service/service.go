package service

import (
	"database/sql"
	"regexp"
	"strings"
)

type Service struct {
	DB *sql.DB
}



// ParseMarkdownLink extracts text and URL from markdown link format [text](url)
func (s *Service) ParseMarkdownLink(content string) (text, url string, isValid bool) {
	// Regex pattern to match [text](url) format
	pattern := `\[([^\]]+)\]\(([^)]+)\)`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(content)
	if len(matches) == 3 {
		text = strings.TrimSpace(matches[1])
		url = strings.TrimSpace(matches[2])

		// Basic URL validation
		if text != "" && url != "" && (strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "www.")) {
			// Add https:// if URL starts with www.
			if strings.HasPrefix(url, "www.") {
				url = "https://" + url
			}
			return text, url, true
		}
	}

	return "", "", false
}

// DetectMarkdownLink checks if content contains a markdown link pattern
func (s *Service) DetectMarkdownLink(content string) bool {
	pattern := `\[([^\]]+)\]\(([^)]+)\)`
	re := regexp.MustCompile(pattern)
	return re.MatchString(content)
}

func (s *Service) GetLinkPlaceholder() string {
	return "[test](google.com)"
}
