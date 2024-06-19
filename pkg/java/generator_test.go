package java

import (
	"testing"
)

func TestEntityGenerator_FromDDL(t *testing.T) {
	type fields struct {
		pkg       string
		typeMap   map[string]TypeMap
		templates []EntityGeneratorTemplate
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		// Add test cases here
		{
			name: "Test case 1",
			fields: fields{
				pkg: "com.example.entity",
				typeMap: map[string]TypeMap{
					"VARCHAR": {
						SqlType:     "VARCHAR",
						JavaType:    "String",
						JavaPackage: "",
					},
					"BIGINT": {
						SqlType:     "BIGINT",
						JavaType:    "Long",
						JavaPackage: "",
					},
				},
				templates: []EntityGeneratorTemplate{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eg := &EntityGenerator{
				cfg: Config{
					Pkg:     tt.fields.pkg,
					TypeMap: tt.fields.typeMap,
				},
				templates: tt.fields.templates,
			}

			got, err := eg.FromDDL("CREATE TABLE users (id BIGINT, name VARCHAR(255));")
			if (err != nil) != tt.wantErr {
				t.Errorf("EntityGenerator.FromDDL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for key, value := range got {
				t.Logf("%s:\n%s", key, value)
			}
		})
	}
}
