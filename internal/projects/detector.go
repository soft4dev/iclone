package projects

import (
	"os"
	"path/filepath"
)

type ProjectType interface {
	Name() string
	Detect(projectPath string) (ProjectHandler, error)
	DefaultProjectHandler() ProjectHandler
}

type ProjectDetector struct {
	projectTypes []ProjectType
}

func (self *ProjectDetector) FindProjectHandler(projectPath string) (ProjectHandler, error) {
	for _, projectType := range self.projectTypes {
		projectHandler, err := projectType.Detect(projectPath)
		if err != nil {
			return nil, err
		}
		if projectHandler != nil {
			return projectHandler, nil
		}
	}
	return nil, nil
}

func (self *ProjectDetector) ProjectHandlerFromName(name string) ProjectHandler {
	for _, projectType := range self.projectTypes {
		if name == projectType.Name() {
			return projectType.DefaultProjectHandler()
		}
	}
	return nil
}

func (self *ProjectDetector) RegisterDetector(projectType ProjectType) {
	self.projectTypes = append(self.projectTypes, projectType)
}

func (self *ProjectDetector) GetAvailableProjectTypes() []string {
	var projectTypes = []string{}
	for _, projectType := range self.projectTypes {
		projectTypes = append(projectTypes, projectType.Name())
	}
	return projectTypes
}

func DefaultDetector() ProjectDetector {
	projectDetector := ProjectDetector{}
	projectDetector.RegisterDetector(&PnpmProjectType{})
	projectDetector.RegisterDetector(&NpmProjectType{})
	return projectDetector
}

/* PNPM */
type PnpmProjectType struct{}

func (self *PnpmProjectType) Name() string {
	return "pnpm"
}

func (self *PnpmProjectType) Detect(projectPath string) (ProjectHandler, error) {
	pnpmLockPath := filepath.Join(projectPath, "pnpm-lock.yaml")
	if _, err := os.Stat(pnpmLockPath); err == nil {
		return pnpmHandler{}, nil
	}
	return nil, nil
}

func (self *PnpmProjectType) DefaultProjectHandler() ProjectHandler {
	return pnpmHandler{}
}

/* NPM */
type NpmProjectType struct{}

func (self *NpmProjectType) Name() string {
	return "npm"
}

func (self *NpmProjectType) Detect(projectPath string) (ProjectHandler, error) {
	npmLockPath := filepath.Join(projectPath, "package-lock.json")
	if _, err := os.Stat(npmLockPath); err == nil {
		return npmHandler{}, nil
	}
	return nil, nil
}

func (self *NpmProjectType) DefaultProjectHandler() ProjectHandler {
	return npmHandler{}
}
