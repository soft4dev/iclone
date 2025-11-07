package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/soft4dev/clonei/internal/color"
	"github.com/soft4dev/clonei/internal/projects"
	"github.com/spf13/cobra"
)

var (
	project string
	install bool
	cd      bool
)

var rootCmd = &cobra.Command{
	Use:   "clonei",
	Short: "clone and install deps of project",
	Long: `
		It clones provided repo using git and install dependencies according to project type. eg. npm, pnpm, go, rust....
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if r := checkGitInstalled(); r != nil {
			return fmt.Errorf("git is not installed")
		}
		repoUrl := args[0]

		// Extract the directory name from the repo URL
		// e.g., "https://github.com/user/repo.git" -> "repo"
		projectDirName := repoUrl
		if idx := strings.LastIndex(projectDirName, "/"); idx != -1 {
			projectDirName = projectDirName[idx+1:]
		}
		projectDirName = strings.TrimSuffix(projectDirName, ".git")

		// Check if project directory already exists
		if _, err := os.Stat(projectDirName); err == nil {
			return fmt.Errorf("project directory '%s' already exists in the current location", projectDirName)
		}

		color.PrintSuccess("🚀 Cloning repository: %s", repoUrl)
		gitCloneOutput := exec.Command("git", "clone", repoUrl)
		gitCloneOutput.Stdout = os.Stdout
		gitCloneOutput.Stderr = os.Stderr
		gitCloneOutput.Stdin = os.Stdin
		if err := gitCloneOutput.Run(); err != nil {
			return fmt.Errorf("error cloning repo: %w", err)
		}

		var projectHandler projects.ProjectHandler
		projectDetector := projects.DefaultDetector()
		if project == "AUTO" {
			var err error
			if projectHandler, err = projectDetector.FindProjectHandler(projectDirName); err != nil {
				return err
			}
		} else {
			projectHandler = projectDetector.ProjectHandlerFromName(project)
		}

		if projectHandler == nil {
			return fmt.Errorf("no handler found for project type '%s'\nAvailable project types: %s", project, projectDetector.GetAvailableProjectTypes())
		}
		color.PrintSuccess("\n📦 Installing dependencies for %s project...")
		if err := projectHandler.Install(projectDirName); err != nil {
			return err
		}
		color.PrintSuccess("✓ Dependencies installed successfully \n")

		if cd {
			if err := os.Chdir(projectDirName); err != nil {
				return fmt.Errorf("failed to change directory: %w", err)
			}
		}

		color.PrintSuccess("project: %s", project)
		color.PrintSuccess("url: %s", args[0])
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func Execute() {
	// Silence Cobra's default error printing
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	err := rootCmd.Execute()
	if err != nil {
		color.PrintError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&project, "project", "p", "AUTO", "Project type (npm, go, rust, etc.). Use AUTO for auto-detection")
	rootCmd.Flags().BoolVarP(&install, "install", "i", true, "controls whether to install dependencies after clone")
	rootCmd.Flags().BoolVarP(&cd, "cd", "c", true, "controls whether to change directory into the project folder after clone")
}

func checkGitInstalled() error {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git is not installed or not available in PATH")
	}
	return nil
}
