package internal

import (
	"github.com/soft4dev/clonei/internal/projects"
)

type Project interface {
	Name() string
	Detect(projectPath string) (projects.IProjectHandler, error)
	ProjectHandler() projects.IProjectHandler
}

type ProjectDetector struct {
	projects []Project
}

// It tries to detect the project type automatically
// based on the given project path. It returns the corresponding handler
// if one is found, or nil if no type matches.
func (projectDetector *ProjectDetector) FindProjectHandlerAuto(projectPath string) (projects.IProjectHandler, error) {
	for _, projectType := range projectDetector.projects {
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

// It returns a default handler for the given project type name.
// If no matching project type is found, it returns nil.
func (projectDetector *ProjectDetector) FindProjectHandlerFromName(name string) projects.IProjectHandler {
	for _, project := range projectDetector.projects {
		if name == project.Name() {
			return project.ProjectHandler()
		}
	}
	return nil
}

func (projectDetector *ProjectDetector) RegisterProject(project Project) {
	projectDetector.projects = append(projectDetector.projects, project)
}

func (projectDetector *ProjectDetector) GetAvailableProjects() []string {
	var projectNames []string
	for _, project := range projectDetector.projects {
		projectNames = append(projectNames, project.Name())
	}
	return projectNames
}

// It initializes and returns a ProjectDetector instance
// that is used to detect project and find corresponding project handlers.
func GetProjectDetector() ProjectDetector {
	projectDetector := ProjectDetector{}
	projectDetector.RegisterProject(&projects.PnpmProject{})
	projectDetector.RegisterProject(&projects.NpmProject{})
	projectDetector.RegisterProject(&projects.CargoProject{})
	projectDetector.RegisterProject(&projects.MavenProject{})
	return projectDetector
}
