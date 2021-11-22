package cmd

import (
	"strings"

	"github.com/hashed-io/document-graph-cli/e"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	RootCmd.PersistentFlags().StringP("vault-file", "", "./eosc-vault.json", "Wallet file that contains encrypted key material")
	RootCmd.PersistentFlags().IntP("expiration", "", 30, "Set time before transaction expires, in seconds. Defaults to 30 seconds.")

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
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	recurseViperCommands(RootCmd, nil)

	err := viper.ReadInConfig()
	if err != nil {
		zlog.Sugar().Errorf("fatal error config file: %s \n", err)
	}

	zlog.Debug("configuration file used", zap.String("filename", viper.ConfigFileUsed()))

	e.E()
}

func recurseViperCommands(root *cobra.Command, segments []string) {
	// Stolen from: github.com/abourget/viperbind
	var segmentPrefix string
	if len(segments) > 0 {
		segmentPrefix = strings.Join(segments, "-") + "-"
	}

	root.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "global-" + f.Name
		viper.BindPFlag(newVar, f)
	})
	root.Flags().VisitAll(func(f *pflag.Flag) {
		newVar := segmentPrefix + "cmd-" + f.Name
		viper.BindPFlag(newVar, f)
	})

	for _, cmd := range root.Commands() {
		recurseViperCommands(cmd, append(segments, cmd.Name()))
	}
}
