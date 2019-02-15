package cmd

import (
	"fmt"
	"github.com/freshautomations/stoml/defaults"
	"github.com/freshautomations/stoml/exit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math"
	"os"
)

const help = `Usage:
  stoml filename key [flags]

Flags:
  -h, --help      help for stoml
      --version   version for stoml
`

const epsilon = 1e-9 // Margin of error

func CheckArgs(cmd *cobra.Command, args []string) error {
	validateArgs := cobra.ExactArgs(2)
	if err := validateArgs(cmd, args); err != nil {
		return err
	}

	fileName := args[0]
	_, err := os.Stat(fileName)
	return err
}

func RunRoot(cmd *cobra.Command, args []string) (output string, err error) {
	fileName := args[0]
	key := args[1]

	viper.SetConfigFile(fileName)
	err = viper.ReadInConfig()
	if err != nil {
		if _, IsUnsupportedExtension := err.(viper.UnsupportedConfigError); IsUnsupportedExtension {
			viper.SetConfigType("toml")
			err = viper.ReadInConfig()
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	result := viper.Get(key)
	if result == nil {
		result = ""
	}
	if f32, IsFloat32 := result.(float32); IsFloat32 {
		if _, frac := math.Modf(math.Abs(float64(f32))); frac < epsilon || frac > 1.0-epsilon {
			result = int32(f32)
		}
	}
	if f64, IsFloat64 := result.(float64); IsFloat64 {
		if _, frac := math.Modf(math.Abs(f64)); frac < epsilon || frac > 1.0-epsilon {
			result = int64(f64)
		}
	}

	output = fmt.Sprintf("%v", result)
	return
}

func runRootWrapper(cmd *cobra.Command, args []string) {
	if result, err := RunRoot(cmd, args); err != nil {
		exit.Fail(err)
	} else {
		exit.Succeed(result)
	}
}

func Execute() error {
	var rootCmd = &cobra.Command{
		Version: defaults.Version,
		Use:     "stoml",
		Short:   "STOML - simple toml parser for Shell",
		Long: `A simplified TOML (also known as a more formal INI) parser for the Linux Shell.
Source and documentation is available at https://github.com/freshautomations/stoml`,
		Args: CheckArgs,
		Run:  runRootWrapper,
	}
	rootCmd.SetUsageTemplate(help)

	return rootCmd.Execute()
}