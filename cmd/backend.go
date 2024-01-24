package cmd

import (
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/routes"
	"github.com/FlareZone/melon-backend/internal/validator"
	"github.com/FlareZone/melon-backend/internal/worker"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var backendCmd = &cobra.Command{
	Use:   "api",
	Short: `Run melon api to start service`,
	Long:  `Run melon api to start service`,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		components.InitComponents()
		// gin 增加验证器
		validator.ValidatorRegister()
		// 执行 worker
		worker.InitWorker()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		routes.Route(r)
		return r.Run(":8080")
	},
}
