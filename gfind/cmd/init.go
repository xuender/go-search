package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// 数据库目录
	_db = "db-path"
)

// GetString 读取配置
func GetString(cmd *cobra.Command, name string) string {
	f := cmd.Flag(name)
	// 命令行优先
	if f.Changed {
		return f.Value.String()
	}
	ret := viper.GetString(name)
	if ret == "" {
		return f.Value.String()
	}
	return ret
}

// GetBool 读取配置
func GetBool(cmd *cobra.Command, name string) bool {
	f := cmd.Flag(name)
	b, _ := strconv.ParseBool(f.Value.String())
	// 命令行优先
	if f.Changed {
		return b
	}
	ret := viper.GetString(name)
	if ret == "" {
		return b
	}
	return viper.GetBool(name)
}
