package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/h2non/filetype.v1"

	search "../.."
)

var indexCmd = &cobra.Command{
	Use:     "index",
	Aliases: []string{"i"},
	Short:   "索引文件",
	Long: `
  根据通配符索引文件或目录`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := GetString(cmd, _db)
		fmt.Println("索引数据库:", db)
		engine, _ := search.NewLevelEngine(db)
		defer engine.Close()

		include := viper.GetStringSlice("include")
		exclude := viper.GetStringSlice("exclude")
		files := map[string]bool{}
		for _, indexPath := range args {
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
				if isMatch && isTxt(walkPath) {
					fullName, _ := filepath.Abs(walkPath)
					fmt.Println("选中文件:", fullName)
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
				return err
			}
		}
		for k := range files {
			fmt.Println("更新索引:", k)
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

func isTxt(f string) bool {
	t, err := filetype.MatchFile(f)
	return err == nil && t.MIME.Type == ""
}

// 匹配
func match(pattern, name string) bool {
	// 目录匹配
	for _, p := range strings.Split(name, string(os.PathSeparator)) {
		if m, err := path.Match(pattern, p); err == nil && m {
			return true
		}
	}
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
	flags.StringP(_db, "d", "db", "数据库目录")
	flags.BoolP("number", "n", false, "索引文件数量")
	flags.BoolP("list", "l", false, "显示索引文件")
}
