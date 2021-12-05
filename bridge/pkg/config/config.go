package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	RepoOwner    string
	RepoName     string
	CategoryName string
	OutputFile   string
}

func FromEnvironment() *Config {
	getOrDefault := func(key, fallback string) string {
		val := os.Getenv(key)
		if val == "" {
			return fallback
		}
		return val
	}
	repo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	if len(repo) != 2 {
		panic(fmt.Sprintf("invalid repo: %s", repo))
	}
	return &Config{
		RepoOwner:    repo[0],
		RepoName:     repo[1],
		CategoryName: getOrDefault("CATEGORY_NAME", "Blog"),
		OutputFile:   getOrDefault("OUTPUT_FILE", "data/discussions.json"),
	}
}

func (c *Config) String() string {
	return fmt.Sprintf(`
Repo:
	Owner: %s
	Name: %s
Discussions:
	CategoryName: %s
OutputFile: %s
`, c.RepoOwner, c.RepoName, c.CategoryName, c.OutputFile)
}
