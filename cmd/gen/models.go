package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "从数据库反向生成 model 结构体",
	Long: `读取 config/config.yaml 中的 mysql 配置，生成临时 reverse 配置文件并调用 reverse 命令。

需要先安装 reverse: go install xorm.io/reverse@latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, _ := cmd.Flags().GetString("config-dir")
		tables, _ := cmd.Flags().GetStringSlice("tables")
		tmplPath, _ := cmd.Flags().GetString("template")

		dsn, err := buildDSN(configDir)
		if err != nil {
			return err
		}

		// 生成临时 reverse 配置文件
		tmpDir, err := os.MkdirTemp("", "reverse-*")
		if err != nil {
			return fmt.Errorf("创建临时目录失败: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		cfgPath := filepath.Join(tmpDir, "reverse.yml")
		if err := writeReverseConfig(cfgPath, dsn, tables, tmplPath); err != nil {
			return fmt.Errorf("生成 reverse 配置文件失败: %w", err)
		}

		reverseBin, err := exec.LookPath("reverse")
		if err != nil {
			return fmt.Errorf("未找到 reverse 命令，请先安装: go install xorm.io/reverse@latest")
		}

		run := exec.Command(reverseBin, "-f", cfgPath)
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		return run.Run()
	},
}

func writeReverseConfig(path, dsn string, tables []string, tmplPath string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, _ = fmt.Fprintf(f, `kind: reverse
name: mydb
source:
  database: mysql
  conn_str: "%s"
targets:
- type: codes
  language: golang
  output_dir: internal/models
  multiple_files: true
  table_mapper: snake
  column_mapper: snake
`, dsn)

	if tmplPath != "" {
		abs, _ := filepath.Abs(tmplPath)
		_, _ = fmt.Fprintf(f, "  template_path: %s\n", abs)
	}

	if len(tables) > 0 {
		_, _ = fmt.Fprintln(f, "  include_tables:")
		for _, t := range tables {
			_, _ = fmt.Fprintf(f, "    - %s\n", t)
		}
	}

	return nil
}

func buildDSN(dir string) (string, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	if err := v.ReadInConfig(); err != nil {
		return "", fmt.Errorf("读取配置失败: %w\n请确认 %s/config.yaml 存在 (可复制 config.yaml.example 并填入 mysql 配置)", err, dir)
	}

	host := v.GetString("mysql.host")
	port := v.GetInt("mysql.port")
	user := v.GetString("mysql.username")
	pass := v.GetString("mysql.password")
	dbname := v.GetString("mysql.dbname")

	if host == "" || dbname == "" {
		return "", fmt.Errorf("mysql 配置不完整，请检查 %s/config.yaml 中的 mysql 配置", dir)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", user, pass, host, port, dbname), nil
}

func init() {
	genModelsCmd.Flags().StringP("config-dir", "c", "config", "配置文件目录")
	genModelsCmd.Flags().StringSliceP("tables", "t", nil, "指定表名（不指定则生成全部）")
	genModelsCmd.Flags().StringP("template", "", "", "自定义模板文件路径")
	rootCmd.AddCommand(genModelsCmd)
}
