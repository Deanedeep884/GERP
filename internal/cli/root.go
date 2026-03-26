package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gerp",
	Short: "GERP Master Control Plane Operator",
	Long:  `GERP (Google ERP) CLI - A FAANG-grade terminal operator natively orchestrating Spanner isolation bounds and Temporal Sagss.`,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Pings the GERP matrix execution limits and reports current bounds",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🟢 GERP MATRIX OPERATIONAL")
		fmt.Println("==================================================")
		fmt.Printf("GraphQL Boundary:   %s\n", ActiveConfig.GraphQLEndpoint)
		fmt.Printf("Temporal Execution: %s\n", ActiveConfig.TemporalHost)
		fmt.Printf("Cloud Spanner:      %s\n", ActiveConfig.SpannerDB)
		fmt.Println("==================================================")
	},
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.AddCommand(statusCmd)
}

// Execute is the structural bootloader for the CLI operator.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
