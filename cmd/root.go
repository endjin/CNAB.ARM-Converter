package main

import (
	"fmt"
	"os"

	"github.com/endjin/CNAB.ARM-Converter/pkg/generator"
	"github.com/spf13/cobra"
)

// Version is set as part of build
var Version string

var bundleloc string
var outputloc string
var overwrite bool
var indent bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the atfcnab version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("atfcnab version: %v \n", Version)
	},
}

var rootCmd = &cobra.Command{
	Use:   "atfcnab",
	Short: "atfcnab generates an ARM template for executing a CNAB package using Azure ACI",
	Long:  `atfcnab generates an ARM template which can be used to execute Duffle in a container using ACI to perform actions on a CNAB Package, which in turn executes the CNAB Actions using the Duffle ACI Driver   `,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return generator.GenerateTemplate(bundleloc, outputloc, overwrite, indent)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&bundleloc, "bundle", "b", "bundle.json", "name of bundle file to generate template for , default is bundle.json")
	rootCmd.Flags().StringVarP(&outputloc, "file", "f", "azuredeploy.json", "file name for generated template,default is azuredeploy.json")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "specifies if to overwrite the output file if it already exists, default is false")
	rootCmd.Flags().BoolVarP(&indent, "indent", "i", false, "specifies if the json output should be indented")

	rootCmd.AddCommand(versionCmd)
}

// Execute runs the template generator
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
