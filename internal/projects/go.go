package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/soft4dev/clonei/internal/color"
)

type goProjectHandler struct{}

func (goProjectHandler *goProjectHandler) Install(projectDir string) error {
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go not found; please install go and ensure it's on your PATH")
	}

	color.PrintSuccess("  → Running go ci...")
	init := exec.Command("go", "mod", "tidy")
	init.Dir = projectDir
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init.Stdin = os.Stdin
	if err := init.Run(); err != nil {
		return fmt.Errorf("error initializing project (go mod tidy): %w", err)
	}

	return nil
}

type GoProject struct{}

func (goProject *GoProject) Name() string {
	return "go"
}

func (goProject *GoProject) Detect(projectPath string) IProjectHandler {
	goProjectConfigPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goProjectConfigPath); err == nil {
		return &goProjectHandler{}
	}
	return nil
}

func (goProject *GoProject) ProjectHandler() IProjectHandler {
	return &goProjectHandler{}
}
