package admin

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Create, update, and delete admin users.",
}
