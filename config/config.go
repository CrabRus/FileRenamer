package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// читає рядок з консолі та повертає його
func readLine(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)

	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return "", errors.New("порожній ввід")
	}

	return text, nil
}

// Читає директорію
func ReadDirectory() (string, error) {
	dir, err := readLine("Вкажіть директорію для перейменування: ")
	if err != nil {
		return "", err
	}

	// перевірка: чи існує директорія?
	info, err := os.Stat(dir)
	if err != nil {
		return "", errors.New("директорія не існує")
	}

	if !info.IsDir() {
		return "", errors.New("вказаний шлях не є директорією")
	}

	return dir, nil
}

func ReadPattern() (string, error) {
	pattern, err := readLine("Вкажіть шаблон файлів (наприклад, *.png): ")
	if err != nil {
		return "", err
	}

	if len(pattern) < 2 {
		return "", errors.New("шаблон занадто короткий")
	}

	return pattern, nil
}

func ReadAction() (string, error) {
	action, err := readLine("Оберіть дію (prefix, suffix, replace, extension, lowercase, uppercase): ")
	if err != nil {
		return "", err
	}

	allowed := []string{"prefix", "suffix", "replace", "extension", "lowercase", "uppercase"}

	for _, a := range allowed {
		if action == a {
			return action, nil
		}
	}

	return "", errors.New("невідома дія")
}

func ReadParameter(action string) (string, error) {

	switch action {

	case "prefix":
		return readLine("Введіть значення префікса: ")

	case "suffix":
		return readLine("Введіть значення суфікса: ")

	case "replace":
		old, err := readLine("Що замінити: ")
		if err != nil {
			return "", err
		}
		newVal, err := readLine("На що замінити: ")
		if err != nil {
			return "", err
		}
		return old + "|" + newVal, nil

	case "extension":
		ext, err := readLine("Нове розширення (без крапки): ")
		if err != nil {
			return "", err
		}
		if strings.Contains(ext, ".") {
			return "", errors.New("не потрібно писати крапку у розширенні")
		}
		return ext, nil

	case "lowercase", "uppercase":
		return "", nil
	}

	return "", errors.New("невідомий тип правила")
}

func ReadBool(prompt string) (bool, error) {
	var input string
	fmt.Print(prompt)
	_, err := fmt.Scanln(&input)
	if err != nil {
		return false, err
	}

	input = strings.ToLower(strings.TrimSpace(input))

	if input == "true" || input == "1" || input == "yes" {
		return true, nil
	}
	if input == "false" || input == "0" || input == "no" {
		return false, nil
	}

	return false, errors.New("введіть true або false")
}
