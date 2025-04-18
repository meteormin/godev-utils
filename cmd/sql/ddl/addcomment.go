package ddl

import (
	"bufio"
	"fmt"
	"github.com/meteormin/godev-utils/pkg/ddl"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

var AddCommentCmd = &cobra.Command{
	Use:   "addComment",
	Short: "Generate SQL ALTER statements to add comments to columns",
	Long: `Generate SQL ALTER statements to add comments to columns. For example:
	- generate from ddl string: addComment "CREATE TABLE users (id BIGINT, name VARCHAR(255));"
	- generate from sql file: addComment tables.sql
	`,
	Args: cobra.ExactArgs(1),
	RunE: runAddComment,
}

// CREATE TABLE 블록을 분리
func splitCreateTableBlocks(ddl string) []string {
	// naive splitter: "CREATE TABLE" 기준으로 나누되 다시 붙임
	parts := strings.Split(ddl, "CREATE TABLE")
	var blocks []string
	for _, part := range parts[1:] {
		block := "CREATE TABLE" + part
		blocks = append(blocks, block)
	}
	return blocks
}

func generateAlterStatements(tableName string, columnMap map[string]string) []string {
	var alterStatements []string
	for colDef, comment := range columnMap {
		stmt := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s COMMENT '%s';",
			tableName, strings.TrimRight(colDef, ","), comment)
		alterStatements = append(alterStatements, stmt)
	}
	return alterStatements
}

func writeToFile(fileName string, lines []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}

var (
	outputPath string
)

func init() {
	AddCommentCmd.Flags().StringVarP(&outputPath, "output", "o", "./alter_statements.sql", "output directory for generated SQL files")
}

func validFlags() error {
	dir := path.Dir(outputPath)
	if _, err := os.Stat(dir); err != nil {
		if err := os.Mkdir(dir, 0o755); err != nil {
			return err
		}
	}

	fmt.Println("outputPath:", outputPath)

	return nil
}

func runAddComment(cmd *cobra.Command, args []string) error {
	if err := validFlags(); err != nil {
		return err
	}

	var sql string
	if _, err := os.Stat(args[0]); err != nil {
		fmt.Println("Reading SQL DDL string")
		sql = args[0]
	} else {
		fmt.Println("Reading SQL DDL file")
		file, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		sql = string(file)
	}

	var allAlters []string
	for _, block := range splitCreateTableBlocks(sql) {
		tableName, commentMap := ddl.ExtractColumnComments(block)
		if tableName == "" || len(commentMap) == 0 {
			continue
		}
		alterSQLs := generateAlterStatements(tableName, commentMap)
		allAlters = append(allAlters, alterSQLs...)
		allAlters = append(allAlters, "\n")
	}

	if err := writeToFile(outputPath, allAlters); err != nil {
		fmt.Println("파일 저장 실패:", err)
	} else {
		fmt.Println("alter_comments.sql 파일로 저장 완료.")
	}

	return nil
}
