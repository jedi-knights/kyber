package commands

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jedi-knights/kyber/internal/domain/metrics"
)

func newListMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-metrics",
		Short: "List every registered metric and its default threshold.",
		RunE: func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tDEFAULT\tDIRECTION\tDESCRIPTION")
			for _, m := range metrics.DefaultRegistry().All() {
				dir := "higher is worse"
				if !m.HigherIsWorse() {
					dir = "lower is worse"
				}
				fmt.Fprintf(w, "%s\t%s\t%g\t%s\t%s\n",
					m.ID(), m.Name(), m.DefaultThreshold(), dir, m.Description())
			}
			return w.Flush()
		},
	}
}
