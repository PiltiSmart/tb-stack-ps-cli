package cmd

import (
	"fmt"
	"os"

	"github.com/piltismart/tb-stack-ps-cli/pkg/s3"
	"github.com/piltismart/tb-stack-ps-cli/pkg/ui"
	"github.com/spf13/cobra"
)

var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "AWS S3-compatible cloud storage operations (powered by MinIO)",
	Long: `==================================================================
  PiltiSmart S3 Cloud Storage CLI (MinIO Engine)
  Target Endpoint: http://localhost:9000 | Alias: myminio
==================================================================

Full AWS S3-like CRUD command suite backed by high-performance MinIO.
Supports bucket operations, object transfers, synchronization, and streaming.`,
	Run: func(cmd *cobra.Command, args []string) {
		printS3Help()
	},
}

func printS3Help() {
	cfg := s3.GetDefaultConfig()
	fmt.Println("==================================================================")
	fmt.Println("  PiltiSmart S3 Cloud Storage CLI (MinIO Engine)")
	fmt.Printf("  Target Endpoint: %s | Alias: %s\n", cfg.Endpoint, cfg.Alias)
	fmt.Println("==================================================================")
	fmt.Println("\nUsage:")
	fmt.Println("  pilti s3 <command> [arguments] [flags]")
	fmt.Println("\nAvailable S3 Commands:")
	fmt.Printf("  %-8s %s\n", "ls", "List S3 buckets or objects inside a bucket")
	fmt.Printf("  %-8s %s\n", "mb", "Make (create) a new S3 bucket")
	fmt.Printf("  %-8s %s\n", "rb", "Remove (delete) an S3 bucket")
	fmt.Printf("  %-8s %s\n", "cp", "Copy files to S3, from S3, or between S3 paths")
	fmt.Printf("  %-8s %s\n", "rm", "Remove (delete) an object or prefix from S3")
	fmt.Printf("  %-8s %s\n", "cat", "Display contents of an S3 object to stdout")
	fmt.Printf("  %-8s %s\n", "stat", "Display metadata and attributes of an S3 object")
	fmt.Printf("  %-8s %s\n", "sync", "Synchronize / mirror local directory with S3 bucket")
	fmt.Printf("  %-8s %s\n", "setup", "Verify or install MinIO client ('mc') and configure alias")
	fmt.Printf("  %-8s %s\n", "config", "Display active MinIO endpoint and credential settings")
	fmt.Println("\nExamples:")
	fmt.Println("  pilti s3 ls")
	fmt.Println("  pilti s3 ls s3://mybucket")
	fmt.Println("  pilti s3 mb s3://mybucket")
	fmt.Println("  pilti s3 cp ./file.txt s3://mybucket/")
	fmt.Println("  pilti s3 cp s3://mybucket/file.txt ./")
	fmt.Println("  pilti s3 cat s3://mybucket/config.json")
	fmt.Println("  pilti s3 sync ./dist s3://mybucket/dist")
	fmt.Println("  pilti s3 rm s3://mybucket/file.txt")
	fmt.Println("  pilti s3 rb s3://mybucket --force")
	fmt.Println("==================================================================")
}

// executeMC prepares mc and executes the command with given arguments
func executeMC(mcCmd string, args ...string) error {
	cfg := s3.GetDefaultConfig()
	mcPath, err := s3.EnsureMC(cfg)
	if err != nil {
		ui.Error("%v", err)
		return err
	}

	fullArgs := append([]string{mcCmd}, args...)
	return s3.RunMC(mcPath, fullArgs...)
}

// pilti s3 ls [target]
var s3LsCmd = &cobra.Command{
	Use:   "ls [target]",
	Short: "List S3 buckets or objects inside a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		target := cfg.Alias
		if len(args) > 0 && args[0] != "" {
			target = s3.NormalizePath(args[0], cfg.Alias)
		}

		mcArgs := []string{target}
		if len(args) > 1 {
			for _, a := range args[1:] {
				mcArgs = append(mcArgs, s3.NormalizePath(a, cfg.Alias))
			}
		}

		if err := executeMC("ls", mcArgs...); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 mb <bucket>
var s3MbCmd = &cobra.Command{
	Use:   "mb <s3://bucket>",
	Short: "Make (create) a new S3 bucket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		bucket := s3.NormalizePath(args[0], cfg.Alias)

		if err := executeMC("mb", bucket); err != nil {
			os.Exit(1)
		}
		ui.Success("Bucket '%s' created successfully.", args[0])
	},
}

// pilti s3 rb <bucket> [--force]
var s3RbForce bool
var s3RbCmd = &cobra.Command{
	Use:   "rb <s3://bucket>",
	Short: "Remove (delete) an S3 bucket",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		bucket := s3.NormalizePath(args[0], cfg.Alias)

		mcArgs := []string{}
		if s3RbForce {
			mcArgs = append(mcArgs, "--force")
		}
		mcArgs = append(mcArgs, bucket)

		if err := executeMC("rb", mcArgs...); err != nil {
			os.Exit(1)
		}
		ui.Success("Bucket '%s' removed successfully.", args[0])
	},
}

// pilti s3 cp <src> <dst> [--recursive]
var s3CpRecursive bool
var s3CpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy files to S3, from S3, or between S3 paths",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		src := s3.NormalizePath(args[0], cfg.Alias)
		dst := s3.NormalizePath(args[1], cfg.Alias)

		mcArgs := []string{}
		if s3CpRecursive {
			mcArgs = append(mcArgs, "--recursive")
		}
		mcArgs = append(mcArgs, src, dst)

		if err := executeMC("cp", mcArgs...); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 rm <target> [--recursive] [--force]
var s3RmRecursive bool
var s3RmForce bool
var s3RmCmd = &cobra.Command{
	Use:   "rm <s3://bucket/key>",
	Short: "Remove (delete) an object or prefix from S3",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		target := s3.NormalizePath(args[0], cfg.Alias)

		mcArgs := []string{}
		if s3RmRecursive {
			mcArgs = append(mcArgs, "--recursive")
		}
		if s3RmForce {
			mcArgs = append(mcArgs, "--force")
		}
		mcArgs = append(mcArgs, target)

		if err := executeMC("rm", mcArgs...); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 cat <target>
var s3CatCmd = &cobra.Command{
	Use:   "cat <s3://bucket/key>",
	Short: "Display contents of an S3 object to stdout",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		target := s3.NormalizePath(args[0], cfg.Alias)

		if err := executeMC("cat", target); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 stat <target>
var s3StatCmd = &cobra.Command{
	Use:   "stat <s3://bucket/key>",
	Short: "Display metadata and attributes of an S3 object",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		target := s3.NormalizePath(args[0], cfg.Alias)

		if err := executeMC("stat", target); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 sync <src> <dst> [--delete] [--overwrite]
var s3SyncDelete bool
var s3SyncCmd = &cobra.Command{
	Use:   "sync <source> <destination>",
	Short: "Synchronize / mirror local directory with S3 bucket",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		src := s3.NormalizePath(args[0], cfg.Alias)
		dst := s3.NormalizePath(args[1], cfg.Alias)

		mcArgs := []string{}
		if s3SyncDelete {
			mcArgs = append(mcArgs, "--remove")
		}
		mcArgs = append(mcArgs, "--overwrite", src, dst)

		// mc mirror handles directory syncing
		if err := executeMC("mirror", mcArgs...); err != nil {
			os.Exit(1)
		}
	},
}

// pilti s3 setup
var s3SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Verify or install MinIO client ('mc') and configure alias",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		ui.PrintBanner("PiltiSmart S3 / MinIO Environment Setup")
		ui.Info("Target Endpoint : %s", cfg.Endpoint)
		ui.Info("MinIO Alias     : %s", cfg.Alias)
		ui.Info("Access Key      : %s", cfg.AccessKey)

		mcPath := s3.FindMCExecutable()
		if mcPath == "" {
			ui.Info("Installing MinIO client ('mc')...")
			var err error
			mcPath, err = s3.InstallMC()
			if err != nil {
				ui.Error("Failed to install mc: %v", err)
				os.Exit(1)
			}
		} else {
			ui.Success("Found existing MinIO client: %s", mcPath)
		}

		if err := s3.EnsureAlias(mcPath, cfg); err != nil {
			ui.Error("Failed to set MinIO alias: %v", err)
			os.Exit(1)
		}

		ui.Success("PiltiSmart S3 client is fully configured and ready to use!")
		fmt.Printf("\nTry running: pilti s3 ls\n\n")
	},
}

// pilti s3 config
var s3ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Display active MinIO endpoint and credential settings",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := s3.GetDefaultConfig()
		fmt.Println("==================================================================")
		fmt.Println("  PiltiSmart S3 / MinIO Configuration")
		fmt.Println("==================================================================")
		fmt.Printf("  MinIO Endpoint : %s\n", cfg.Endpoint)
		fmt.Printf("  Default Alias  : %s\n", cfg.Alias)
		fmt.Printf("  Access Key     : %s\n", cfg.AccessKey)
		fmt.Printf("  Secret Key     : %s\n", "******** (configured)")
		mcPath := s3.FindMCExecutable()
		if mcPath != "" {
			fmt.Printf("  mc Binary Path : %s\n", mcPath)
		} else {
			fmt.Println("  mc Binary Path : Not installed (will auto-install on first run)")
		}
		fmt.Println("==================================================================")
	},
}

func init() {
	s3RbCmd.Flags().BoolVarP(&s3RbForce, "force", "f", false, "Force delete bucket even if not empty")
	s3CpCmd.Flags().BoolVarP(&s3CpRecursive, "recursive", "r", false, "Copy recursively")
	s3RmCmd.Flags().BoolVarP(&s3RmRecursive, "recursive", "r", false, "Remove recursively")
	s3RmCmd.Flags().BoolVarP(&s3RmForce, "force", "f", false, "Force removal without prompt")
	s3SyncCmd.Flags().BoolVar(&s3SyncDelete, "delete", false, "Delete files in destination that do not exist in source")

	s3Cmd.AddCommand(s3LsCmd)
	s3Cmd.AddCommand(s3MbCmd)
	s3Cmd.AddCommand(s3RbCmd)
	s3Cmd.AddCommand(s3CpCmd)
	s3Cmd.AddCommand(s3RmCmd)
	s3Cmd.AddCommand(s3CatCmd)
	s3Cmd.AddCommand(s3StatCmd)
	s3Cmd.AddCommand(s3SyncCmd)
	s3Cmd.AddCommand(s3SetupCmd)
	s3Cmd.AddCommand(s3ConfigCmd)

	s3Cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printS3Help()
	})
}
