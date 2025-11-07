package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/soft4dev/clonei/internal/color"
)

/* PNPM */
type pnpmProjectHandler struct{}

func (n pnpmProjectHandler) Install(projectDir string) error {
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

type PnpmProject struct{}

func (pnpmProject *PnpmProject) Name() string {
	return "pnpm"
}

func (pnpmProject *PnpmProject) Detect(projectPath string) IProjectHandler {
	pnpmLockPath := filepath.Join(projectPath, "pnpm-lock.yaml")
	if _, err := os.Stat(pnpmLockPath); err == nil {
		return pnpmProjectHandler{}
	}
	return nil
}

func (pnpmProject *PnpmProject) ProjectHandler() IProjectHandler {
	return pnpmProjectHandler{}
}
