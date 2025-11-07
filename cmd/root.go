package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/soft4dev/clonei/internal"
	"github.com/soft4dev/clonei/internal/color"
	customErrors "github.com/soft4dev/clonei/internal/errors"
	projectHandler "github.com/soft4dev/clonei/internal/projects"
	"github.com/soft4dev/clonei/internal/utils"
	"github.com/spf13/cobra"
)

var (
	projectType string
	install     bool
	cd          bool
)

var rootCmd = &cobra.Command{
	Use:   "clonei",
	Short: "clone and install deps of project",
	Long:  `It clones provided repo using git and install dependencies according to project type. eg. npm, pnpm, go, rust....`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		projectType, _ := cmd.Flags().GetString("project")
		projectDetector := internal.GetProjectDetector()
		availableProjectTypes := projectDetector.GetAvailableProjects()
		if !utils.ContainsStringInStringSlice(availableProjectTypes, projectType) {
			return customErrors.NewCustomError(fmt.Sprintf("unsupported project type '%s'\nAvailable project types: \n %s", projectType, availableProjectTypes), customErrors.ErrorTypeError, false)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if git is installed
		if err := utils.CheckGitInstalled(); err != nil {
			return customErrors.NewCustomError("Error: git is not installed", customErrors.ErrorTypeError, false)
		}
		repoUrl := args[0]

		// Extract the directory name from the repo URL e.g., "https://github.com/user/repo.git" -> "repo"
		var projectDirName string
		if idx := strings.LastIndex(repoUrl, "/"); idx != -1 {
			projectDirName = repoUrl[idx+1:]
		}
		projectDirName = strings.TrimSuffix(projectDirName, ".git")

		// Check if project directory already exists
		if _, err := os.Stat(projectDirName); err == nil {
			return customErrors.NewCustomError(fmt.Sprintf("project directory '%s' already exists in the current location", projectDirName), customErrors.ErrorTypeError, false)
		}

		// clone the repository
		color.PrintInfo("Step 1: Cloning repository: %s", repoUrl)
		gitCloneOutput := exec.Command("git", "clone", repoUrl)
		gitCloneOutput.Stdout = os.Stdout
		gitCloneOutput.Stderr = os.Stderr
		gitCloneOutput.Stdin = os.Stdin
		if err := gitCloneOutput.Run(); err != nil {
			return customErrors.NewCustomError("", customErrors.ErrorTypeError, false)
		}

		// project handler handles the dependencies installation using install() function
		var projectHandler projectHandler.IProjectHandler

		// detects project type and set the projectHandler accordingly
		projectDetector := internal.GetProjectDetector()
		if projectType == "AUTO" {
			if projectHandler = projectDetector.FindProjectHandlerAuto(projectDirName); projectHandler == nil {
				return customErrors.NewCustomError(fmt.Sprintf("no handler found for project type '%s'\nAvailable project types: %s", projectType, projectDetector.GetAvailableProjects()), customErrors.ErrorTypeInfo, false)
			}
		} else {
			projectHandler = projectDetector.FindProjectHandlerFromName(projectType)
		}

		// it will be nil if user specified project type is not found by the function FindProjectHandlerFromName()
		if projectHandler == nil {
			return customErrors.NewCustomError(fmt.Sprintf("no handler found for project type '%s'\nAvailable project types: %s", projectType, projectDetector.GetAvailableProjects()), customErrors.ErrorTypeWarning, false)
		}

		color.PrintSuccess("\n📦 Installing dependencies for %s project...")
		if err := projectHandler.Install(projectDirName); err != nil {
			return customErrors.NewCustomError(fmt.Sprintf("failed to install dependencies: %s", err), customErrors.ErrorTypeWarning, false)
		}

		color.PrintSuccess("✓ Dependencies installed successfully \n")
		if cd {
			if err := os.Chdir(projectDirName); err != nil {
				return customErrors.NewCustomError(fmt.Sprintf("failed to change directory: %s", err), customErrors.ErrorTypeWarning, false)
			}
		}

		color.PrintSuccess("project: %s", projectType)
		color.PrintSuccess("url: %s", args[0])
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		if cmdErr, ok := err.(*customErrors.CustomError); ok {
			if cmdErr.MessageType == customErrors.ErrorTypeError {
				color.PrintError(err.Error())
			}
			if cmdErr.MessageType == customErrors.ErrorTypeWarning {
				color.PrintWarning(err.Error())
			}
			if cmdErr.MessageType == customErrors.ErrorTypeInfo {
				color.PrintInfo(err.Error())
			}
			if cmdErr.ShowUsage {
				rootCmd.Usage()
			}
		} else {
			color.PrintError(err.Error())
			rootCmd.Usage()
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.Flags().StringVarP(&projectType, "project", "p", "AUTO", "Project type (npm, go, rust, etc.). Use AUTO for auto-detection")
	rootCmd.Flags().BoolVarP(&install, "install", "i", true, "controls whether to install dependencies after clone")
	rootCmd.Flags().BoolVarP(&cd, "cd", "c", true, "controls whether to change directory into the project folder after clone")
}
