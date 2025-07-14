package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: go run cmd/generate-module/main.go <nome_modulo>")
		os.Exit(1)
	}

	moduleName := strings.ToLower(os.Args[1])

	// Criar diretório do módulo
	moduleDir := filepath.Join("internal", moduleName)
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		fmt.Printf("Erro ao criar diretório: %v\n", err)
		os.Exit(1)
	}

	// Dados para os templates
	data := struct {
		ModuleName      string
		ModuleNameCap   string
		ModuleNameUpper string
	}{
		ModuleName:      moduleName,
		ModuleNameCap:   strings.Title(moduleName),
		ModuleNameUpper: strings.ToUpper(moduleName),
	}

	// Criar arquivos básicos
	files := map[string]string{
		"doc.go":        docTemplate,
		"model.go":      modelTemplate,
		"repository.go": repositoryTemplate,
		"service.go":    serviceTemplate,
		"graphql.go":    graphqlTemplate,
		"init.go":       initTemplate,
	}

	for filename, tmplContent := range files {
		filepath := filepath.Join(moduleDir, filename)

		tmpl, err := template.New(filename).Parse(tmplContent)
		if err != nil {
			fmt.Printf("Erro no template %s: %v\n", filename, err)
			continue
		}

		file, err := os.Create(filepath)
		if err != nil {
			fmt.Printf("Erro ao criar %s: %v\n", filename, err)
			continue
		}

		if err := tmpl.Execute(file, data); err != nil {
			fmt.Printf("Erro ao executar template %s: %v\n", filename, err)
		}
		file.Close()
	}

	fmt.Printf("Módulo %s criado em %s\n", moduleName, moduleDir)
}
