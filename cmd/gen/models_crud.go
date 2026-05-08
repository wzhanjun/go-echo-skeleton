package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

var genCrudCmd = &cobra.Command{
	Use:   "crud <ModelName>...",
	Short: "为指定模型生成 repo + service 脚手架",
	Long: `根据模型名生成 internal/repo/impl/{model}_repo.go 和 internal/services/{model}_service.go。
模型名采用 PascalCase，例如: gen crud User Order Ticket`,
	RunE: func(cmd *cobra.Command, args []string) error {
		models, _ := cmd.Flags().GetStringSlice("models")
		models = append(models, args...)
		if len(models) == 0 {
			return cmd.Help()
		}

		modPath := modulePath()
		for _, name := range models {
			m := modelInfo{
				Name:     name,
				Plural:   plural(name),
				VarName:  lowerFirst(name),
				ModPath:  modPath,
			}
			genFile("internal/repo/impl/"+snake(name)+"_repo.go", repoTPL, m)
			genFile("internal/services/"+snake(name)+"_service.go", serviceTPL, m)
		}
		return nil
	},
}

type modelInfo struct {
	Name    string
	Plural  string
	VarName string
	ModPath string
}

func modulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "github.com/wzhanjun/go-echo-skeleton"
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "github.com/wzhanjun/go-echo-skeleton"
}

func genFile(path, tpl string, m modelInfo) {
	code := execTpl(tpl, m)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		exit("mkdir %s: %v", dir, err)
	}
	formatted, err := format.Source([]byte(code))
	if err != nil {
		formatted = []byte(code)
	}
	if err := os.WriteFile(path, formatted, 0644); err != nil {
		exit("write %s: %v", path, err)
	}
	fmt.Printf("  created %s\n", path)
}

func execTpl(tpl string, m modelInfo) string {
	var buf bytes.Buffer
	template.Must(template.New("").Parse(tpl)).Execute(&buf, m)
	return buf.String()
}

func snake(s string) string {
	var out []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(r))
	}
	return string(out)
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(string(s[0])) + s[1:]
}

func plural(s string) string { return strings.ToLower(s) + "s" }

func exit(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

const repoTPL = `package impl

import (
	"{{.ModPath}}/internal/models"
	"{{.ModPath}}/internal/repo"
)

type {{.Name}}Repo struct {
	repo.BaseRepo[models.{{.Name}}]
}

func New{{.Name}}Repo() *{{.Name}}Repo {
	return &{{.Name}}Repo{NewBaseRepoImpl[models.{{.Name}}]()}
}
`

const serviceTPL = `package services

import (
	"{{.ModPath}}/internal/dto"
	"{{.ModPath}}/internal/models"
	"{{.ModPath}}/internal/repo/impl"
	"{{.ModPath}}/pkg/db"
	"xorm.io/xorm"
)

type {{.Name}}Service struct {
	repo *impl.{{.Name}}Repo
}

func New{{.Name}}Service() *{{.Name}}Service {
	return &{{.Name}}Service{repo: impl.New{{.Name}}Repo()}
}

func (s *{{.Name}}Service) Get(id interface{}) (*models.{{.Name}}, error) {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.GetById(session, id)
}

func (s *{{.Name}}Service) List(params dto.PageParams) ([]*models.{{.Name}}, int64, error) {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Paginate(session, params, func(sess *xorm.Session) *xorm.Session {
		return sess.OrderBy("id DESC")
	})
}

func (s *{{.Name}}Service) Create(m *models.{{.Name}}) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Create(session, m)
}

func (s *{{.Name}}Service) Update(id interface{}, m *models.{{.Name}}) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Update(session, id, m)
}

func (s *{{.Name}}Service) Delete(id interface{}) error {
	session := db.GetEngine().NewSession()
	defer session.Close()
	return s.repo.Delete(session, id)
}
`

func init() {
	genCrudCmd.Flags().StringSliceP("models", "m", nil, "模型名列表")
	rootCmd.AddCommand(genCrudCmd)
}
