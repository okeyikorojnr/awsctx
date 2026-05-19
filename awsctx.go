package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "awsctx [profile]",
		Short: "Switch between AWS profiles (mimics kubectx)",
		Long:  `A fast, native Go tool to switch AWS profiles by spawning a sub-shell.`,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profiles := getProfiles()

			// Case 1: No args - List profiles
			if len(args) == 0 {
				current := os.Getenv("AWS_PROFILE")
				for _, p := range profiles {
					if p == current {
						fmt.Printf("\033[1;32m %s\033[0m\n", p)
					} else {
						fmt.Println(p)
					}
				}
				return
			}

			// Case 2: Profile switch
			target := args[0]
			if contains(profiles, target) {
				startSubShell(target)
			} else {
				fmt.Printf("error: profile '%s' not found\n", target)
				os.Exit(1)
			}
		},
		// This enables dynamic tab-completion for your profile names!
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getProfiles(), cobra.ShellCompDirectiveNoFileComp
		},
	}

	// Add the 'completion' command (built into Cobra)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please specify a shell (e.g., bash)")
				return
			}
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getProfiles() []string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".aws", "credentials")
	cfg, err := ini.Load(path)
	if err != nil {
		return []string{}
	}

	var profiles []string
	for _, section := range cfg.Sections() {
		name := section.Name()
		if name != "DEFAULT" && name != "default" {
			profiles = append(profiles, name)
		}
	}
	sort.Strings(profiles)
	return profiles
}

func startSubShell(profile string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), fmt.Sprintf("AWS_PROFILE=%s", profile))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	fmt.Printf("✔ Switched to AWS Profile: %s (exit to return)\n", profile)
	_ = cmd.Run()
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
