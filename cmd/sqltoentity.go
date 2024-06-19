package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	pkgJava "github.com/meteormin/godev-utils/pkg/java"
	"github.com/spf13/cobra"
)

var sqlToEntityCmd = &cobra.Command{
	Use:   "sqlToEntityForJava",
	Short: "Generate entity classes from SQL DDL",
	Long: `Generate entity classes from SQL DDL. For example:
		- generate from ddl string: sqlToEntityForJava "CREATE TABLE users (id BIGINT, name VARCHAR(255));"
		- generate from sql file: sqlToEntityForJava tables.sql
	`,
	Args: cobra.ExactArgs(1),
	RunE: runSqlToEntity,
}

var (
	packageName string
	outputDir   string
	typeMap     string
)

func init() {
	rootCmd.AddCommand(sqlToEntityCmd)
	sqlToEntityCmd.Flags().StringVarP(&packageName, "package", "p", "com.example.entity", "package name for generated entity classes")
	sqlToEntityCmd.Flags().StringVarP(&outputDir, "out", "o", "./", "output directory for generated entity classes")
	sqlToEntityCmd.Flags().StringVarP(&typeMap, "typeMap", "t", "./typeMap.json", "type map for converting SQL types to Java types")
}

func validFlags() error {
	if packageName == "" {
		packageName = "com.example.entity"
	}

	if outputDir == "" {
		outputDir = "./"
	}

	if typeMap == "" {
		typeMap = "./typeMap.json"
	}

	if _, err := os.Stat(outputDir); err != nil {
		if err := os.Mkdir(outputDir, 0o755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(typeMap); err != nil {
		return err
	}

	fmt.Println("packageName:", packageName)
	fmt.Println("outputDir:", outputDir)
	fmt.Println("typeMap:", typeMap)

	return nil
}

// runSqlToEntity generates entity classes from SQL DDL
// and writes the generated classes to the output directory
// The generated entity classes are written in Java
// The generated entity classes are written in the package name specified by the user
// The generated entity classes are written in the output directory specified by the user
// The generated entity classes are written based on the type map specified by the user
// The type map is a JSON file that maps SQL types to Java types
// The generated entity classes are written in the format of a Java class
func runSqlToEntity(cmd *cobra.Command, args []string) error {
	// Implement the logic for generating entity classes from SQL DDL
	// 1. Read the input (SQL DDL string or file path)
	// 2. Parse the input to extract table names and fields
	// 3. Generate entity classes from the parsed data
	// 4. Write the generated entity classes to the output directory

	// Example: sqlToEntityForJava "CREATE TABLE users (id BIGINT, name VARCHAR(255));"
	// Example: sqlToEntityForJava tables.sql
	// Example: sqlToEntityForJava --package=com.example.entity --out=./output --typeMap=./typeMap.json "CREATE TABLE users (id BIGINT, name VARCHAR(255));"
	// Example: sqlToEntityForJava --package=com.example.entity --out=./output --typeMap=./typeMap.json tables.sql

	var sqlStr string
	if _, err := os.Stat(args[0]); err != nil {
		fmt.Println("Reading SQL DDL string")
		sqlStr = args[0]
	} else {
		fmt.Println("Reading SQL DDL file")
		file, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		sqlStr = string(file)
	}

	err := validFlags()
	if err != nil {
		return err
	}

	typeMapFile, err := os.ReadFile(typeMap)
	if err != nil {
		return err
	}

	var typeMap map[string]pkgJava.TypeMap
	err = json.Unmarshal(typeMapFile, &typeMap)
	if err != nil {
		return err
	}

	eg := pkgJava.NewEntityGenerator(
		pkgJava.Config{
			Pkg:     packageName,
			TypeMap: typeMap,
			ClassName: pkgJava.ClassName{
				Prefix: "",
				Suffix: "Entity",
			},
		})

	rs, err := eg.FromDDL(sqlStr)
	if err != nil {
		return err
	}

	fmt.Println("Generated entity classes:", len(rs))

	for key, value := range rs {
		filePath := outputDir + "/" + key + ".java"
		err := os.WriteFile(filePath, []byte(value), 0o644)
		if err != nil {
			fmt.Println("Error writing file:", filePath)
			return err
		}
	}

	return nil
}
