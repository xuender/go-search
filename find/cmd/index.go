package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"i"},
	Short:   "索引文件",
	Long: `
  根据通配符索引文件或目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("索引", _dbPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	flags := indexCmd.Flags()
	flags.StringVarP(&_dbPath, _db, "d", "db", "数据库目录")
}
