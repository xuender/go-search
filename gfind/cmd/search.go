package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"

	search "../.."
)

var searchCmd = &cobra.Command{
	Use:     "search 关键词...",
	Aliases: []string{"s"},
	Short:   "搜索关键词",
	Long: `
  根据关键词搜索已索引的文件爱你`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, _ := filepath.Abs(GetString(cmd, _db))
		fmt.Println("索引数据库:", db)

		engine, _ := search.NewLevelEngine(db)
		defer engine.Close()

		if len(args) == 0 {
			return errors.New("请输入搜索内容")
		}
		str := strings.Join(args, " ")
		color.Println("搜索:", color.Red(str))

		start := time.Now()
		docs := engine.Search(str)
		end := time.Now()
		for _, d := range docs {
			fmt.Println("\t -", d.Title)
		}
		color.Println("运行时间:", color.Green(fmt.Sprintf("%v", end.Sub(start))))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	flags := searchCmd.Flags()
	flags.StringP(_db, "d", "db", "数据库目录")
}
