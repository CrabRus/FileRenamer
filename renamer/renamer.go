package renamer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"renamer/model"
	"strings"
)

// Пошук файлів за шаблоном
func FindFiles(dir string, pattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		match, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if match {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func applyRule(filePath string, rule model.Rule) (string, error) {
	dir := filepath.Dir(filePath)
	name := filepath.Base(filePath)
	ext := filepath.Ext(name)
	nameWithoutExt := strings.TrimSuffix(name, ext)

	switch rule.Action {

	case "prefix":
		newName := rule.Parameter + nameWithoutExt + ext
		return filepath.Join(dir, newName), nil

	case "suffix":
		newName := nameWithoutExt + rule.Parameter + ext
		return filepath.Join(dir, newName), nil

	case "replace":
		parts := strings.Split(rule.Parameter, "|")
		if len(parts) != 2 {
			return "", errors.New("некоректний параметр replace (має бути: старе|нове)")
		}
		old := parts[0]
		newStr := parts[1]

		newName := strings.ReplaceAll(nameWithoutExt, old, newStr) + ext
		return filepath.Join(dir, newName), nil

	case "extension":
		newName := nameWithoutExt + "." + rule.Parameter
		return filepath.Join(dir, newName), nil

	case "lowercase":
		newName := strings.ToLower(nameWithoutExt) + ext
		return filepath.Join(dir, newName), nil

	case "uppercase":
		newName := strings.ToUpper(nameWithoutExt) + ext
		return filepath.Join(dir, newName), nil
	}

	return "", errors.New("невідоме правило")
}

func RenameFiles(files []string, rule model.Rule) []model.RenameResult {
	var results []model.RenameResult

	for _, oldPath := range files {
		newPath, err := applyRule(oldPath, rule)

		if err != nil {
			results = append(results, model.RenameResult{
				OldName: oldPath,
				NewName: "",
				Success: false,
				Error:   err,
			})
			continue
		}

		err = os.Rename(oldPath, newPath)

		if err != nil {
			results = append(results, model.RenameResult{
				OldName: oldPath,
				NewName: newPath,
				Success: false,
				Error:   err,
			})
			continue
		}

		results = append(results, model.RenameResult{
			OldName: oldPath,
			NewName: newPath,
			Success: true,
			Error:   nil,
		})
	}

	return results
}

// DryRun
func DryRun(files []string, rule model.Rule) []model.RenameResult {
	var results []model.RenameResult

	for _, oldPath := range files {
		newPath, err := applyRule(oldPath, rule)

		if err != nil {
			results = append(results, model.RenameResult{
				OldName: oldPath,
				NewName: "",
				Success: false,
				Error:   err,
			})
			continue
		}

		results = append(results, model.RenameResult{
			OldName: oldPath,
			NewName: newPath,
			Success: true,
			Error:   nil,
		})
	}

	return results
}

// SaveBackup зберігає результати перейменування у JSON
func SaveBackup(results []model.RenameResult, filename string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// UndoBackup відміняє перейменування за backup файлом
func UndoBackup(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var results []model.RenameResult
	err = json.Unmarshal(data, &results)
	if err != nil {
		return err
	}

	for _, r := range results {
		if _, err := os.Stat(r.NewName); err == nil {
			if err := os.Rename(r.NewName, r.OldName); err != nil {
				fmt.Printf("❌ Не вдалося відкотити %s → %s: %v\n", r.NewName, r.OldName, err)
			} else {
				fmt.Printf("✅ Відкотили: %s → %s\n", r.NewName, r.OldName)
			}
		}
	}
	return os.Remove(filename)
}
