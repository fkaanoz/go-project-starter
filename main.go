package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var subDirs = []string{"/config", "/database", "/errors", "/handlers", "/models", "/services"}
var deps = []string{"github.com/gin-gonic/gin", "gorm.io/gorm", "gorm.io/driver/postgres", "github.com/spf13/viper"}
var files = []string{"/main.go", "/errors/errors.go", "/config/"}
var envTypes = []string{"json", "env"}

func main() {
	dir := flag.String("dir", "", "directory name for project (ex. test)")
	module := flag.String("module", "", "module name (ex. github.com/fkaanoz/test")
	env := flag.String("env", "env", "env type (.env or json)")
	flag.Parse()

	if (*env != "env") && (*env != "json") {
		log.Fatal("Wrong ENV type. It should be either env or json.")
	}

	files[2] += *env

	if err := CheckFlags([]string{*dir, *module, *env}); err != nil {
		log.Fatal("Check your flags! : ", err)
	}

	if err := CreateDirs(*dir); err != nil {
		log.Fatal("Error occurred when creating directories! : ", err)
	}

	if err := InitModule(*dir, *module); err != nil {
		RollBack(*dir)
		log.Fatal("Error occurred when initializing module! : ", err)
	}

	if err := InstallDeps(*dir); err != nil {
		RollBack(*dir)
		log.Fatal("Error occurred when installing dependencies! : ", err)
	}

	if err := CreateRequiredFiles(*dir); err != nil {
		RollBack(*dir)
		log.Fatal("Error occurred when creating files such main.go or .env! : ", err)
	}

	CreateReadme(*dir) // When all folders and files are created successfully, a readme file will be created. More than one run of the program will be prevented by looking at created readme file!

}

func CheckFlags(flags []string) error {
	for _, f := range flags {
		if len(f) == 0 {
			return errors.New("Not Enough Flag!")
		}

		if strings.Contains(f, " ") {
			return errors.New("Flags are wrong!")
		}
	}

	return nil
}

func CreateDirs(projectDir string) error {
	for _, dir := range subDirs {
		err := os.MkdirAll(projectDir+dir, 0777)

		if err != nil {
			return err
		}
	}

	return nil
}

func InitModule(projectDir string, module string) error {

	cmd := exec.Command("go", "mod", "init", module)

	cmd.Dir += "./" + projectDir

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func InstallDeps(projectDir string) error {

	for _, dep := range deps {
		cmd := exec.Command("go", "get", "-u", dep)

		cmd.Dir += "./" + projectDir

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func CreateRequiredFiles(projectDir string) error {

	for _, fileName := range files {
		_, err := os.Create("./" + projectDir + fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateReadme(projectDir string) {
	file, err := os.Create("./" + projectDir + "/README.md")

	if err != nil {
		log.Fatal("CREATING READ ME ERROR")
	}

	defer file.Close()

	file.WriteString("fkz")

	fmt.Println("Project created successfully!")
}

func RollBack(projectDir string) {

	_, err := os.ReadFile(projectDir + "/README.md")

	if err == nil {
		fmt.Println("Project is created successfully before!")
		return
	}

	err = os.RemoveAll(projectDir)

	if err != nil {
		log.Fatal("Cannot delete directory!", err)
	}
}
