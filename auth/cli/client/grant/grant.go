package grant

import (
	"github.com/spf13/cobra"
)

var ClientGrantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Add and remove client grants.",
}
