package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shopsmart/ssm2ssm"
	"github.com/shopsmart/ssm2ssm/pkg/service"
)

// New creates a new cobra command for the given version
func New(version string, svc service.Service) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ssm2ssm",
		Short: "Pulls SSM paramters into env format",
		Long: `SSM2SSM copies parameters from one SSM path to another

ssm2ssm /input/prefix /output/prefix`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var v bool
			v, _ = cmd.Flags().GetBool("verbose")
			if v {
				log.SetLevel(log.DebugLevel)
			}

			v, _ = cmd.Flags().GetBool("version")
			if v {
				fmt.Println(version)
				return
			}

			overwrite, _ := cmd.Flags().GetBool("overwrite")
			if v {
				fmt.Println(version)
				return
			}

			validArgs := []string{}
			for _, arg := range args {
				if arg != "" {
					validArgs = append(validArgs, arg)
				}
			}

			if len(validArgs) < 2 {
				if err := cmd.Help(); err != nil {
					log.Fatal(err)
				}
				return
			}

			err := ssm2ssm.Copy(svc, validArgs[0], validArgs[1], overwrite)
			if err != nil {
				log.Fatal(err)
				return
			}
		},
	}

	rootCmd.PersistentFlags().Bool("verbose", false, "enables verbose output")
	rootCmd.PersistentFlags().Bool("version", false, "prints the version and exits")
	rootCmd.PersistentFlags().Bool("overwrite", false, "overwrites the existing parameters under output prefix if already present")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, svc service.Service) {
	cmd := New(version, svc)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
