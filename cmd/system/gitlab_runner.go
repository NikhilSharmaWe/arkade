// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/spf13/cobra"
)

func MakeInstallGitLabRunner() *cobra.Command {
	command := &cobra.Command{
		Use:   "gitlab-runner",
		Short: "Install GitLab Runner",
		Long:  `Install GitLab Runner for self-hosted CI.`,
		Example: `  arkade system install gitlab-runner
  arkade system install gitlab-runner --version <version>`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/usr/local/bin/gitlab-runner", "Installation path, where a gitlab-runner subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing GitLab Runner to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "arm"
		}

		if version == "" {
			version = "latest"
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		dlURL := fmt.Sprintf("https://gitlab-runner-downloads.s3.amazonaws.com/%s/binaries/gitlab-runner-linux-%s", version, dlArch)
		dlCmd := exec.Command("sudo", "curl", "-L", "--progress-bar", "--output", installPath, dlURL)

		fmt.Printf("Downloading from: %s\n", dlURL)

		if err := dlCmd.Run(); err != nil {
			fmt.Println(installPath)
			fmt.Println(dlURL)
			return fmt.Errorf("failed to download GitLab Runner binary: %s", err.Error())
		}

		chmodCmd := exec.Command("sudo", "chmod", "+x", installPath)

		fmt.Printf("Downloaded to: %s\n", installPath)

		if err := chmodCmd.Run(); err != nil {
			return fmt.Errorf("failed to set execute permissions for GitLab Runner binary: %s", err.Error())
		}

		fmt.Println("GitLab Runner installation completed successfully!")

		return nil
	}

	return command
}
