package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

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
		fmt.Println("Генерация проекта...")
		progress := 0
		for progress < 100 {
			progress += 10
			fmt.Printf("\rПрогресс: [%s%s] %d%%", strings.Repeat("=", progress/10), strings.Repeat(" ", 10-progress/10), progress)
			time.Sleep(100 * time.Millisecond)
		}
		fmt.Println()
		if err := generator.Generate(&config); err != nil {
			fmt.Printf("Ошибка генерации проекта: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✨ Проект %s успешно сгенерирован! ✨\n", config.Name)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
