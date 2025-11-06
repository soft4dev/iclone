package projects

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/soft4dev/clonei/internal/color"
)

type ProjectHandler interface {
	Install(projectDir string) error
}

type npmHandler struct{}
type pnpmHandler struct{}
type cargoHandler struct{}

func (n npmHandler) Install(projectDir string) error {
	// check npm command exists
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("npm not found; please install npm and ensure it's on your PATH")
	}

	color.PrintSuccess("  → Running npm ci...")
	init := exec.Command("npm", "ci")
	init.Dir = projectDir
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init.Stdin = os.Stdin
	if err := init.Run(); err != nil {
		return fmt.Errorf("error initializing project (npm ci): %w", err)
	}

	return nil
}

func (n pnpmHandler) Install(projectDir string) error {
	// check pnpm command exists
	if _, err := exec.LookPath("pnpm"); err != nil {
		return fmt.Errorf("pnpm not found; please install pnpm and ensure it's on your PATH")
	}

	color.PrintSuccess("  → Running pnpm install --frozen-lockfile...")
	init := exec.Command("pnpm", "install", "--frozen-lockfile")
	init.Dir = projectDir
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init.Stdin = os.Stdin
	if err := init.Run(); err != nil {
		return fmt.Errorf("error initializing project (pnpm install): %w", err)
	}

	return nil
}

func (self cargoHandler) Install(prjectDir string) error {
	return nil
}
