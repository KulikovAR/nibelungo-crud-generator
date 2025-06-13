package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KulikovAR/go-crud-generator/internal/domain"
	"github.com/KulikovAR/go-crud-generator/internal/usecase"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "CRUD project generator",
	Long:  `A generator for creating Go projects with CRUD operations based on domain entities`,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate [json-file]",
	Short: "Generate CRUD project from JSON configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile := args[0]

		// Читаем конфигурационный файл
		data, err := os.ReadFile(jsonFile)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}

		// Парсим конфигурацию
		var config domain.ProjectConfig
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			os.Exit(1)
		}

		// Создаем генератор
		generator := usecase.NewGenerator()

		// Генерируем проект
		if err := generator.Generate(&config); err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project %s generated successfully!\n", config.Name)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
