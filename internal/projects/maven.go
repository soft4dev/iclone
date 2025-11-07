package projects

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/soft4dev/clonei/internal/color"
)

type mavenProjectHandler struct{}

func (mavenProjectHandler *mavenProjectHandler) Install(projectDir string) error {
	if _, err := exec.LookPath("maven"); err != nil {
		return fmt.Errorf("maven not found; please install maven and ensure it's on your PATH")
	}

	color.PrintSuccess("  → Running maven ci...")
	init := exec.Command("mvn", "dependency:resolve")
	init.Dir = projectDir
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	init.Stdin = os.Stdin
	if err := init.Run(); err != nil {
		return fmt.Errorf("error initializing project (maven dependency:resolve): %w", err)
	}

	return nil
}

type MavenProject struct{}

func (mavenProject *MavenProject) Name() string {
	return "maven"
}

func (mavenProject *MavenProject) Detect(projectPath string) (IProjectHandler, error) {
	mavenProjectConfigPath := filepath.Join(projectPath, "pom.xml")
	if _, err := os.Stat(mavenProjectConfigPath); err == nil {
		return &mavenProjectHandler{}, nil
	}
	return nil, nil
}

func (mavenProject *MavenProject) ProjectHandler() IProjectHandler {
	return &mavenProjectHandler{}
}
