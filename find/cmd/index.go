package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	search "../.."
)

var indexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"i"},
	Short:   "索引文件",
	Long: `
  根据通配符索引文件或目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("索引数据库:", _dbPath)
		engine, _ := search.NewLevelEngine(_dbPath)
		defer engine.Close()

		include := viper.GetStringSlice("include")
		exclude := viper.GetStringSlice("exclude")
		files := map[string]bool{}
		for _, indexPath := range args {
			fmt.Println("索引:", indexPath)
			err := filepath.Walk(indexPath, func(walkPath string, f os.FileInfo, err error) error {
				if f == nil {
					return err
				}
				if f.IsDir() {
					return nil
				}
				// 排除
				for _, e := range exclude {
					if match(e, walkPath) {
						return nil
					}
				}
				// 匹配
				isMatch := walkPath == indexPath
				if !isMatch {
					for _, i := range include {
						if match(i, walkPath) {
							isMatch = true
							break
						}
					}
				}
				if isMatch {
					fullName, _ := filepath.Abs(walkPath)
					key := []byte(fullName)
					has, err := engine.Has(key)
					if err != nil {
						return err
					}
					if has {
						doc, err := engine.Get(key)
						if err != nil {
							return err
						}
						if f.ModTime().Before(doc.Modified) {
							return nil
						}
					}
					files[fullName] = true
				}
				return nil
			})
			if err != nil {
				fmt.Println("目录扫描错误", err.Error())
			}
		}
		for k := range files {
			fmt.Println("索引文件:", k)
			b, err := ioutil.ReadFile(k)
			if err != nil {
				fmt.Println(err.Error())
			}
			doc := search.Document{
				Key:      []byte(k),
				Title:    k,
				Content:  string(b),
				Modified: time.Now(),
			}
			engine.Put(&doc)
		}

		if n, _ := cmd.Flags().GetBool("number"); n {
			fmt.Println("索引文档:", engine.IndexNum())
		}
		if l, _ := cmd.Flags().GetBool("list"); l {
			fmt.Println("索引文档:")
			engine.IndexKeys(func(key []byte) {
				fmt.Println("\t -", string(key[2:]))
			})
		}
		return nil
	},
}

// 匹配
func match(pattern, name string) bool {
	// 目录匹配
	if m, err := path.Match(pattern, path.Dir(name)); err == nil && m {
		return true
	}
	// 文件匹配
	if m, err := path.Match(pattern, path.Base(name)); err == nil && m {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(indexCmd)
	flags := indexCmd.Flags()
	flags.StringVarP(&_dbPath, _db, "d", "db", "数据库目录")
	flags.BoolP("number", "n", false, "索引文件数量")
	flags.BoolP("list", "l", false, "显示索引文件")
}
