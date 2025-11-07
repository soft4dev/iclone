package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update clonei to the latest version",
	Long: `Update clonei to the latest version by downloading and installing 
the latest release from GitHub. This command will automatically detect 
your operating system and architecture, then update the binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating clonei: %v\n", err)
			os.Exit(1)
		}
	},
}

func runUpdate() error {
	var updateCmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows PowerShell update command
		script := "irm 'https://raw.githubusercontent.com/soft4dev/clonei/main/scripts/install.ps1' | iex"
		updateCmd = exec.Command("powershell", "-Command", script)
	case "darwin", "linux":
		// macOS/Linux bash update command
		script := "curl -fsSL https://raw.githubusercontent.com/soft4dev/clonei/main/scripts/install.sh | bash -s update"
		updateCmd = exec.Command("bash", "-c", script)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	updateCmd.Stdin = os.Stdin

	fmt.Println("Updating clonei to the latest version...")
	return updateCmd.Run()
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
