package java

import (
	"fmt"
	"strings"
	"unicode"

	pkgDDL "github.com/meteormin/godev-utils/pkg/ddl"
	"github.com/meteormin/godev-utils/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EntityGeneratorTemplate struct {
	Package     string
	TableName   string
	ClassName   string
	Annotations []string
	Imports     []string
	Fields      []Field
}

type Field struct {
	Annotations []string
	Type        string
	Name        string
}

type Annotation struct {
	Name    string
	Package string
	Values  []string
}

type AnnotationMap struct {
	Sql        string
	Annotation Annotation
}
type TypeMap struct {
	SqlType     string `json:"sqlType"`
	JavaType    string `json:"javaType"`
	JavaPackage string `json:"javaPackage"`
}

func (t EntityGeneratorTemplate) TemplateFilename() string {
	return "java/entity.template"
}

type Config struct {
	Pkg         string          `json:"-"`
	TypeMap     []TypeMap       `json:"typeMap"`
	Annotations []AnnotationMap `json:"annotations"`
	ClassName   ClassName       `json:"-"`
}

type ClassName struct {
	Prefix string
	Suffix string
}

type EntityGenerator struct {
	cfg       Config
	templates []EntityGeneratorTemplate
}

func NewEntityGenerator(cfg Config) *EntityGenerator {
	return &EntityGenerator{
		cfg:       cfg,
		templates: make([]EntityGeneratorTemplate, 0),
	}
}

func (eg *EntityGenerator) SetPackage(pkg string) {
	eg.cfg.Pkg = pkg
}

func (eg *EntityGenerator) FromDDL(ddl string) (map[string]string, error) {
	tables := pkgDDL.ExtractTables(ddl)

	fmt.Println("parsed tables: ", len(tables))
	for _, table := range tables {
		fields, err := pkgDDL.ExtractColumns(table)
		if err != nil {
			return nil, err
		}

		addImports := make([]string, 0)
		templateFields := make([]Field, 0)
		for tableName, field := range fields {
			for _, f := range field {
				tm, err := eg.mappingType(f[0])
				if err != nil {
					fmt.Printf("error mapping type: %s\n", err.Error())
					continue
				}

				if tm.JavaPackage != "" {
					addImports = append(addImports, tm.JavaPackage)
				}

				templateFields = append(templateFields, Field{
					Type: tm.JavaType,
					Name: eg.snakeToCamelLowerFirst(f[1]),
				})
			}

			template := EntityGeneratorTemplate{
				Package:   eg.cfg.Pkg,
				TableName: tableName,
				ClassName: eg.cfg.ClassName.Prefix + eg.snakeToCamel(tableName) + eg.cfg.ClassName.Suffix,
				Fields:    templateFields,
				Imports:   addImports,
			}

			eg.templates = append(eg.templates, template)
		}
	}

	return eg.Parse()
}

func (eg *EntityGenerator) Parse() (map[string]string, error) {
	parsed := make(map[string]string, 0)
	for _, t := range eg.templates {
		rs, err := templates.NewParser(t).Parse()
		if err != nil {
			return nil, err
		}

		parsed[t.ClassName] = rs
	}

	return parsed, nil
}

func (eg *EntityGenerator) mappingType(typeStr string) (TypeMap, error) {
	typeStr = strings.Split(typeStr, "(")[0]
	for _, v := range eg.cfg.TypeMap {
		caser := cases.Lower(language.English)
		lowerCaseKey := caser.String(v.SqlType)
		lowerTypeStr := caser.String(typeStr)
		if strings.Contains(lowerCaseKey, lowerTypeStr) {
			return v, nil
		}
	}

	return TypeMap{}, fmt.Errorf("type %s not found", typeStr)
}

func (eg *EntityGenerator) mappingAnnotation(typeStr string) (Annotation, error) {
	for _, v := range eg.cfg.Annotations {
		if v.Sql == typeStr {
			return v.Annotation, nil
		}
	}

	return Annotation{}, fmt.Errorf("annotation %s not found", typeStr)
}

func (eg *EntityGenerator) snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := 0; i < len(parts); i++ {
		caser := cases.Title(language.English)
		parts[i] = caser.String(parts[i])
	}

	return strings.Join(parts, "")
}

func (eg *EntityGenerator) lowerFirstCharacter(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return ""
}

func (eg *EntityGenerator) snakeToCamelLowerFirst(s string) string {
	return eg.lowerFirstCharacter(eg.snakeToCamel(s))
}
