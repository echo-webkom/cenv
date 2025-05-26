package cmd

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"slices"

	"github.com/fatih/color"
)

var Version string

func checkIfLatestVersion() {
	if slices.Contains(os.Args, "--skip-version") {
		os.Args = slices.Delete(os.Args, slices.Index(os.Args, "--skip-version"), slices.Index(os.Args, "--skip-version")+1)
		return
	}

	if Version == "" {
		return
	}

	if !shouldCheck() {
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

// cenv creates the file .cenv-version in your home dir, which contains
// the last time for version check. If the file does not exist or the the time
// is past expiration, cenv checks again. This is to minimize delay when running
// cenv as checking for an update may be slow and annoying on slower internet.
func shouldCheck() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	filepath := path.Join(home, ".cenv-version")

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	expiration := time.Hour * 24
	checked, err := time.Parse(time.RFC3339, string(b))
	if err != nil || time.Now().Sub(checked) > expiration {
		file.Seek(0, io.SeekStart)
		file.WriteString(time.Now().Format(time.RFC3339))
		return true
	}

	return false
}
