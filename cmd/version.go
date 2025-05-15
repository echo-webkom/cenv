package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"slices"

	"github.com/fatih/color"
)

var Version string

func checkIfLatestVersion() {
	if Version == "" {
		return
	}

	if slices.Contains(os.Args, "--skip-version") {
		return
	}

	resp, err := http.Get("https://api.github.com/repos/echo-webkom/cenv/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	var result struct {
		TagName string `json:"tag_name"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}

	latestVersion := strings.TrimSpace(result.TagName)
	isLatest := Version == latestVersion

	if !isLatest {
		color.Yellow("A new version of cenv is available. Run 'cenv upgrade' to upgrade to the latest version.")
	}
}
