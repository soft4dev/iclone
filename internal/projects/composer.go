package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/soft4dev/clonei/internal/color"
)

type composerProjectHandler struct{}

func (composerProjectHandler *composerProjectHandler) Install(projectDir string) error {
	if _, err := exec.LookPath("composer"); err != nil {
		return fmt.Errorf("composer not found; please install composer and ensure it's on your PATH")
	}

	color.PrintSuccess("  → Running composer ci...")
	init := exec.Command("composer", "install")
	init.Dir = projectDir
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init.Stdin = os.Stdin
	if err := init.Run(); err != nil {
		return fmt.Errorf("error initializing project (composer install): %w", err)
	}

	return nil
}

type ComposerProject struct{}

func (composerProject *ComposerProject) Name() string {
	return "composer"
}

func (composerProject *ComposerProject) Detect(projectPath string) IProjectHandler {
	composerProjectConfigPath := filepath.Join(projectPath, "composer.json")
	if _, err := os.Stat(composerProjectConfigPath); err == nil {
		return &composerProjectHandler{}
	}
	return nil
}

func (composerProject *ComposerProject) ProjectHandler() IProjectHandler {
	return &composerProjectHandler{}
}
