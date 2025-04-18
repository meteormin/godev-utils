package ddl

import (
	"fmt"
	"regexp"
	"strings"
)

func ExtractTables(ddl string) []string {
	// "CREATE TABLE" 구문을 찾아내는 로직을 구현합니다.
	// 정규표현식을 사용하여 CREATE TABLE 구문을 찾습니다.
	re := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+` + "`(.*?)`" + `\s*\((.*?)\)\s*;`)
	matches := re.FindAllStringSubmatch(ddl, -1)

	tables := make([]string, len(matches))
	for i, match := range matches {
		tables[i] = strings.TrimSpace(match[0]) // Use the whole match, not just the submatches
	}

	return tables
}

func ExtractColumns(table string) (map[string][][]string, error) {
	// "CREATE TABLE" 구문을 분석하여 테이블 이름과 필드를 추출하는 로직을 구현합니다.
	// 정규표현식을 사용하여 테이블 이름과 필드를 찾습니다.
	re := regexp.MustCompile(`(?is)CREATE\s+TABLE\s+` + "`(.*?)`" + `\s*\((.*?)\)\s*;`)
	match := re.FindStringSubmatch(table)
	if len(match) < 3 {
		return nil, fmt.Errorf("failed to extract table name and fields")
	}

	tableName := match[1]
	fieldsStr := match[2]

	// 필드를 분석합니다.
	fields := [][]string{}

	for _, fieldStr := range strings.Split(fieldsStr, ",") {
		trimmedField := strings.TrimSpace(fieldStr)
		// Extract field name and type
		// Example: `analyte_id` varchar(10) NOT NULL DEFAULT '0'
		fieldParts := strings.Fields(trimmedField)

		if len(fieldParts) < 2 {
			continue
		}

		fieldName := strings.ReplaceAll(fieldParts[0], "`", "")
		fieldType := fieldParts[1]

		fields = append(fields, []string{fieldType, fieldName})
	}

	return map[string][][]string{tableName: fields}, nil
}

func ExtractColumnComments(block string) (string, map[string]string) {
	lines := strings.Split(block, "\n")
	var tableName string
	columnComments := make(map[string]string)

	columnRegex := regexp.MustCompile("(?i)`?(\\w+)`?\\s+[^,]*?COMMENT\\s+'([^']*)'")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(line), "CREATE TABLE") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				tableName = strings.Trim(parts[2], "`")
			}
			continue
		}
		if matches := columnRegex.FindStringSubmatch(line); len(matches) == 3 {
			colDef := strings.Split(line, "COMMENT")[0]
			comment := matches[2]
			columnComments[colDef] = comment
		}
	}

	return tableName, columnComments
}
