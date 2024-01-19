package cmd

import (
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/routes"
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
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		routes.Route(r)
		return r.Run(":8080")
	},
}
