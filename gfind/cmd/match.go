package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	utils "github.com/xuender/go-utils"
)

var matchCmd = &cobra.Command{
	Use:     "match [文件通配类型...]",
	Aliases: []string{"m"},
	Short:   "定义通配",
	Long: `
  设置需要索引的文件包含/排除关系`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置
		include := viper.GetStringSlice("include")
		exclude := viper.GetStringSlice("exclude")
		if len(args) > 0 {
			isDelete, _ := cmd.Flags().GetBool("delete")
			isExclude, _ := cmd.Flags().GetBool("exclude")
			// 删除通配
			if isDelete {
				if isExclude {
					del(&exclude, args)
				} else {
					del(&include, args)
				}
			} else {
				// 排除文件通配
				if isExclude {
					fmt.Println("排除通配类型", args)
					add(&exclude, &include, args)
				} else {
					fmt.Println("增加通配类型", args)
					add(&include, &exclude, args)
				}
			}
			// 保存配置
			viper.Set("include", include)
			viper.Set("exclude", exclude)
			viper.WriteConfig()
		}
		// 输出配置
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

// 增加配置
func add(target, other *[]string, ms []string) {
	m := map[string]bool{}
	for _, o := range *target {
		m[o] = true
	}
	for _, t := range ms {
		m[t] = true
	}

	*target = append([]string{})
	ss := utils.StringSlice(*other)
	for k := range m {
		*target = append(*target, k)
		ss.Delete(k)
	}
	*other = append(ss)
}

// 删除配置
func del(target *[]string, ms []string) {
	ss := utils.StringSlice(*target)
	ss.Delete(ms...)
	*target = append(ss)
}

func init() {
	rootCmd.AddCommand(matchCmd)
	matchCmd.Flags().BoolP("exclude", "e", false, "排除通配类型")
	matchCmd.Flags().BoolP("delete", "d", false, "删除通配类型")
}
