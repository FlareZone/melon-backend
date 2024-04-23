package cmd

import (
	"fmt"
	"github.com/FlareZone/melon-backend/common/logger"
	"github.com/inconshreveable/log15"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	env string
	log = log15.New("m", "cmd")
)

var appCmd = &cobra.Command{
	Use:   "melon",
	Short: `melon service`,
	Long:  `melon service`,
}

func init() {
	appCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "Select the deployment environment (dev, test, prod)")
	cobra.OnInitialize(initConfig)
	appCmd.AddCommand(backendCmd)
}

func initConfig() {
	viper.SetEnvPrefix("ENV")
	viper.AutomaticEnv()
	// 根据环境变量确定加载的配置文件
	if env != "" {
		viper.SetConfigFile(fmt.Sprintf("./config.%s.yaml", env))
	} else {
		viper.SetConfigFile("/etc/melon/config.prod.yaml")
	}
	log15.Root().SetHandler(logger.InitLogHandle(log15.Root().GetHandler()))
	if err := viper.ReadInConfig(); err != nil {
		log.Error("viper read in config failed", "err", err)
		os.Exit(1)
	} else {
		url := viper.GetString("app_url")
		log.Debug("app_url", "url", url)
	}
	//migrate.Schema(viper.GetString("database.melon.dsn"))
}

func Execute() {
	if err := appCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
