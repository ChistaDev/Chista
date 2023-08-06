package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Chista-Framework/Chista/logger"
)

// Function to calculate the Levenshtein distance between two strings
func LevenshteinDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := 1; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			dp[i][j] = min(min(dp[i-1][j]+1, dp[i][j-1]+1), dp[i-1][j-1]+cost)
		}
	}

	return dp[m][n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Return given domain's [subdomain, hostname, tld]
func ParseDomain(domain string) (subdomain, hostname, tld string, err error) {
	// Parse the domain using the net/url package
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return "", "", "", err
	}

	// Split the hostname by dots to extract subdomain and TLD parts
	parts := strings.Split(parsedURL.String(), ".")
	numParts := len(parts)

	// Handle cases with 0, 1, or 2 parts
	if numParts == 0 {
		return "", "", "", fmt.Errorf("invalid domain format: %s", domain)
	} else if numParts == 1 {
		return "", parts[0], "", nil
	} else if numParts == 2 {
		return "", parts[0], parts[1], nil
	}

	// For domains with more than 2 parts, extract subdomain, hostname, and TLD
	subdomain = parts[0]
	hostname = parts[1]
	tld = parts[numParts-1]

	return subdomain, hostname, tld, nil
}

// Function to find new strings similar to the given input string
func GenerateSimilarDomains(input string, threshold int, tld string) []string {
	logger.Log.Debugln("Generating similar domains...")

	var similarStrings []string

	// Generate new strings with Levenshtein distance less than or equal to the threshold
	for char := 'a'; char <= 'z'; char++ {
		for i := 0; i <= len(input); i++ {
			newStr := input[:i] + string(char) + input[i:]
			if LevenshteinDistance(input, newStr) <= threshold {
				similarStrings = append(similarStrings, (newStr + "." + tld))
			}
		}
		for i := 0; i < len(input); i++ {
			newStr := input[:i] + input[i+1:]
			if LevenshteinDistance(input, newStr) <= threshold {
				similarStrings = append(similarStrings, (newStr + "." + tld))
			}
		}
	}

	return similarStrings
}
