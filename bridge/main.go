package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"

	"github.com/hugo-mods/discussions/bridge/pkg/client"
	"github.com/hugo-mods/discussions/bridge/pkg/config"
)

func main() {
	cfg := config.FromEnvironment()
	fmt.Println("got config:", cfg)
	if eventName := os.Getenv("GITHUB_EVENT_NAME"); eventName != "" {
		fmt.Println("triggered by:", eventName)
		fmt.Println("  event path:", os.Getenv("GITHUB_EVENT_PATH"))
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("REPO_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	client := client.New(httpClient, cfg.RepoOwner, cfg.RepoName)

	categories, err := client.Categories()
	if err != nil {
		fatal("could not retrieve categories: %v", err)
	}
	if len(categories) == 0 {
		fatal("could not find any categories. please ensure that discussions are enabled and there is at least one category.")
	}
	category := categories.ByName(cfg.CategoryName, 1).First()
	if category == nil {
		fatal("could not find discussion with name %q", cfg.CategoryName)
	}
	fmt.Println("got category ID:", category.ID)

	discussions, err := client.Discussions(category.ID)
	if err != nil {
		fatal("could not get discussions for category %q", category.ID)
	}
	data, err := json.Marshal(discussions)
	if err != nil {
		fatal("could not marshal JSON")
	}
	fmt.Printf("got %d discussions\n", len(discussions))

	if err := os.MkdirAll(filepath.Dir(cfg.OutputFile), 0777); err != nil {
		fatal("could not create directories to write discussions to JSON file: %v", err)
	}
	if err := os.WriteFile(cfg.OutputFile, data, 0666); err != nil {
		fatal("could not write discussions to JSON file: %v", err)
	}
	fmt.Println("successfully wrote discussions")

}

func fatal(msg string, arg ...interface{}) {
	panic(fmt.Sprintf(msg, arg...))
}
