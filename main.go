package main

import (
	"fmt"
	"renamer/config"
	"renamer/model"
	"renamer/renamer"
)

func main() {
	fmt.Println("=== File Renamer v1.0 ===")

	for {
		fmt.Println("\nЩо хочете зробити?")
		fmt.Println("1. Перейменувати файли")
		fmt.Println("2. Відкотити останні зміни (undo)")
		fmt.Println("3. Вихід")

		var choice int
		fmt.Print("Введіть номер дії: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			renameFiles()
		case 2:
			undoChanges()
		case 3:
			fmt.Println("Вихід з програми.")
			return
		default:
			fmt.Println("❌ Невірний вибір, спробуйте ще раз.")
		}
	}
}

// функція для перейменування файлів
func renameFiles() {
	dir, err := config.ReadDirectory()
	if err != nil {
		fmt.Println("❌ Помилка:", err)
		return
	}

	pattern, err := config.ReadPattern()
	if err != nil {
		fmt.Println("❌ Помилка:", err)
		return
	}

	action, err := config.ReadAction()
	if err != nil {
		fmt.Println("❌ Помилка:", err)
		return
	}

	parameter, err := config.ReadParameter(action)
	if err != nil {
		fmt.Println("❌ Помилка:", err)
		return
	}

	rule := model.Rule{Action: action, Parameter: parameter}

	files, err := renamer.FindFiles(dir, pattern)
	if err != nil {
		fmt.Println("❌ Помилка пошуку:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Немає файлів, які відповідають шаблону")
		return
	}

	// dry run
	fmt.Println("\n=== Попередній перегляд змін (DRY RUN) ===")
	dry := renamer.DryRun(files, rule)
	for _, r := range dry {
		oldShort := shortName(r.OldName)
		if r.Success {
			newShort := shortName(r.NewName)
			fmt.Printf("• %s → %s\n", oldShort, newShort)
		} else {
			fmt.Printf("• %s → ❌ Помилка: %v\n", oldShort, r.Error)
		}
	}

	var answer string
	fmt.Print("\nВиконати перейменування? (y/n): ")
	fmt.Scanln(&answer)
	if answer != "y" && answer != "Y" {
		fmt.Println("❎ Операцію скасовано.")
		return
	}

	results := renamer.RenameFiles(files, rule)
	fmt.Println("\n=== Результат перейменування ===")
	success := 0
	for _, r := range results {
		oldShort := shortName(r.OldName)
		if r.Success {
			newShort := shortName(r.NewName)
			fmt.Printf("• %s → %s\n", oldShort, newShort)
			success++
		} else {
			fmt.Printf("• %s → ❌ Помилка: %v\n", oldShort, r.Error)
		}
	}
	fmt.Printf("\n✅ Успішно перейменовано: %d файлів з %d\n", success, len(results))

	backupFile := "backup.json"
	if err := renamer.SaveBackup(results, backupFile); err != nil {
		fmt.Println("❌ Не вдалося зберегти backup:", err)
	} else {
		fmt.Printf("💾 Backup збережено у %s\n", backupFile)
	}
}

func undoChanges() {
	backupFile := "backup.json"
	err := renamer.UndoBackup(backupFile)
	if err != nil {
		fmt.Println("❌ Помилка при відкаті:", err)
	} else {
		fmt.Println("✅ Всі зміни успішно відкотені.")
	}
}

func shortName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}

//
