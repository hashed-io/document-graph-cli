package cmd

import (
	"github.com/hashed-io/document-graph-cli/e"
	"github.com/spf13/cobra"
	"github.com/streamingfast/logging"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var cfgFile string

var zlog *zap.Logger

func init() {
	logging.Register("github.com/hashed-io/document-graph-cli/cmd", &zlog)
}

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dgctl",
	Short: "CLI for Document Graph",
	Long:  "CLI for Document Graph",
}

func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".dgctl.yaml", "configuration file")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetDefault("blockpause", "0s")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".dgctl")
	}

	viper.SetEnvPrefix("DGCTL")
	viper.AutomaticEnv() // read in environment variables that match prefix

	err := viper.ReadInConfig()
	if err != nil {
		zlog.Sugar().Errorf("fatal error config file: %s \n", err)
	}

	zlog.Debug("configuration file used", zap.String("filename", viper.ConfigFileUsed()))

	e.E()
}
