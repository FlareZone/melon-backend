package cmd

import (
	"github.com/FlareZone/melon-backend/common/logger"
	"github.com/FlareZone/melon-backend/common/migrate"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	cfgFile string
	log     = log15.New("m", "cmd")
)

var appCmd = &cobra.Command{
	Use:   "melon",
	Short: `melon service`,
	Long:  `melon service`,
}

func init() {
	appCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
	cobra.OnInitialize(initConfig)
	appCmd.AddCommand(backendCmd)
}

func initConfig() {
	viper.SetEnvPrefix("ENV")
	viper.AutomaticEnv()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("/etc/melon/config.yaml")
	}
	log15.Root().SetHandler(logger.InitLogHandle(log15.Root().GetHandler()))
	if err := viper.ReadInConfig(); err != nil {
		log.Error("viper read in config failed", "err", err)
		os.Exit(1)
	}
	migrate.Schema(viper.GetString("database.melon.dsn"))
}

func Execute() {

	if err := appCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
