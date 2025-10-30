package cmd

import (
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ultrageopro/wpath/config"
	"github.com/ultrageopro/wpath/internal/out"
	"github.com/ultrageopro/wpath/internal/watcher"
)

var Args config.Args = config.Args{}

var rootCmd = &cobra.Command{
	Use:     "dirwatch --path <DIR> [--filter-name REGEX] [--since RFC3339|UNIX] [--no-color]",
	Short:   "CLI для наблюдения за изменениями в директории",
	Version: "0.1.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		if pathErr := validatePath(Args.FlagPath); pathErr != nil {
			return pathErr
		}
		if filterErr := validateFilter(Args.FlagFilterName); filterErr != nil {
			return filterErr
		}
		if sinceErr := validateSince(Args.FlagSinceStr); sinceErr != nil {
			return sinceErr
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&Args.FlagPath, "path", "", "Путь к директории для наблюдения (обязательный)")
	rootCmd.Flags().StringVar(&Args.FlagFilterName, "filter-name", "", "Регулярное выражение для фильтрации имен файлов")
	rootCmd.Flags().StringVar(&Args.FlagSinceStr, "since", "", "Порог времени (RFC3339, напр. 2025-10-30T12:34:56Z, или UNIX seconds)")
	rootCmd.Flags().BoolVar(&Args.FlagNoColor, "no-color", false, "Отключить цветной вывод")

	rootCmd.MarkFlagRequired("path")
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	printer := out.NewPathPrinter(Args.FlagNoColor)
	defer printer.Stop()

	mu := &sync.Mutex{}
	processor, err := watcher.NewProcessor(printer, Args)
	if err != nil {
		os.Exit(1)
	}
	if err := processor.Watch(Args.FlagPath, mu); err != nil {
		os.Exit(1)
	}
}
