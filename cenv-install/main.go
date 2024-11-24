package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/echo-webkom/cenv/refs/heads/main/install.sh | bash")
	fmt.Println("Installing latest release...")

	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done")
}

