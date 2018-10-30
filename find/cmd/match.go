package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xuender/go-utils"
)

var matchCmd = &cobra.Command{
	Use:     "match 文件通配类型...",
	Aliases: []string{"m"},
	Short:   "定义匹配文件",
	Long: `
  设置需要索引的文件匹配,忽略关系`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("缺少通配类型")
		}
		isExclude, _ := cmd.Flags().GetBool("exclude")
		include := viper.GetStringSlice("include")
		exclude := viper.GetStringSlice("exclude")
		// 包含的文件通配
		if isExclude {
			fmt.Println("增加排除通配类型", args)
			m := map[string]bool{}
			for _, o := range exclude {
				m[o] = true
			}
			for _, t := range args {
				m[t] = true
			}
			exclude = []string{}
			ss := utils.StringSlice(include)
			for k := range m {
				exclude = append(exclude, k)
				ss.Delete(k)
			}
			include = ss
		} else {
			fmt.Println("增加目录文件通配类型", args)
			m := map[string]bool{}
			for _, o := range include {
				m[o] = true
			}
			for _, t := range args {
				m[t] = true
			}
			include = []string{}
			ss := utils.StringSlice(exclude)
			for k := range m {
				include = append(include, k)
				ss.Delete(k)
			}
			exclude = ss
		}
		viper.Set("include", include)
		viper.Set("exclude", exclude)
		viper.WriteConfig()
		fmt.Println("Include:")
		for _, m := range include {
			fmt.Printf("\t- %s\n", m)
		}
		fmt.Println("Exclude:")
		for _, m := range exclude {
			fmt.Printf("\t- %s\n", m)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
	matchCmd.Flags().BoolP("exclude", "e", false, "排除的通配类型")
}
