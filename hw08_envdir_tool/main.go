package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func checkParams(inParams []string) error {
	if len(inParams) < 2 {
		return errors.New("Not enough arguments: " + strconv.Itoa(len(inParams)) + " need (at least): 2")
	}

	envDir := inParams[0]
	dirInfo, err := os.Stat(envDir)
	if os.IsNotExist(err) {
		return errors.New(envDir + " does not exist.")
	}
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return errors.New(envDir + " is not directory.")
	}

	childCommand := inParams[1]
	_, err = exec.LookPath(childCommand)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Println(`Эта утилита позволяет запускать программы, получая переменные окружения из определенной директории:
		- если директория содержит файл с именем 'S', первой строкой которого является 'T', то
		'envdir' удаляет переменную среды с именем 'S', если таковая существует, а затем добавляет
		переменную среды с именем 'S' и значением 'T';
		- имя 'S' не должно содержать '='; пробелы и табуляция в конце 'T' удаляются;
		 терминальные нули ('0x00') заменяются на перевод строки ('\n');
		- если файл полностью пустой (длина - 0 байт), то 'envdir' удаляет переменную окружения с именем 'S'.
		`)
	}

	// checking parameters and files
	flag.Parse()
	inParams := flag.Args()
	err := checkParams(inParams)
	if err != nil {
		log.Panic(err)
	}
	envDir := inParams[0]
	envs, err := ReadDir(envDir)
	if err != nil {
		log.Panic(err)
	}

	for key, env := range envs {
		if env.NeedRemove {
			os.Unsetenv(key)
		}
		os.Setenv(key, env.Value)
	}

	code := RunCmd(inParams[1:], envs)
	os.Exit(code)
}
