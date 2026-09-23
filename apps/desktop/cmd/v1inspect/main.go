// v1inspect reports whether a V1 SQLite database is ready for a V2 migration
// dry run. It opens the source database read-only and writes JSON to stdout.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/zbango/gym-saas/apps/desktop/internal/v1inspect"
)

func main() {
	var sourcePath string
	flag.StringVar(&sourcePath, "db", "", "path to the V1 SQLite database")
	flag.Parse()
	if sourcePath == "" {
		fmt.Fprintln(os.Stderr, "usage: v1inspect -db /path/to/v1.sqlite")
		os.Exit(2)
	}

	report, err := v1inspect.Inspect(context.Background(), sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "inspect V1 database: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "encode inspection report: %v\n", err)
		os.Exit(1)
	}
}
