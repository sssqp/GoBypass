package build

import (
	"github.com/EmYiQing/GoBypass/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func Build(code string, ldflags bool, hide bool, race bool) {
	log.Info("build...")
	if ldflags || hide || race {
		advanceBuild(code, ldflags, hide, race)
		return
	}
	cmd := []string{
		"go",
		"build",
		"-o",
		"output.exe",
		"output/main.go",
	}
	privateBuild(code, cmd)
}

func Garble(code string, ldflags bool, hide bool, race bool) {
	log.Info("garble build...")
	if ldflags || hide || race {
		advanceGarble(code, ldflags, hide, race)
		return
	}
	cmd := []string{
		"go",
		"build",
		"-o",
		"output.exe",
		"output/main.go",
	}
	privateGrable(code, cmd)
}

func advanceBuild(code string, ldflags bool, hide bool, race bool) {
	cmd := []string{
		"go",
		"build",
		"-o",
		"output.exe",
		"-ldflags",
		"",
		"output/main.go",
	}
	if ldflags && hide {
		cmd[5] = "-s -w -H windowsgui"
	}
	if ldflags && !hide {
		cmd[5] = "-s -w"
	}
	if !ldflags && hide {
		cmd[5] = "-H windowsgui"
	}
	if race {
		cmd[5] = "-s -w"
		cmd = append(cmd, "output/main.go")
		cmd[6] = "-race"
	}
	privateBuild(code, cmd)
}

func advanceGarble(code string, ldflags bool, hide bool, race bool) {
	cmd := []string{
		"go",
		"build",
		"-o",
		"output.exe",
		"-ldflags",
		"",
		"output/main.go",
	}
	if ldflags && hide {
		cmd[5] = "-s -w -H windowsgui"
	}
	if ldflags && !hide {
		cmd[5] = "-s -w"
	}
	if !ldflags && hide {
		cmd[5] = "-H windowsgui"
	}
	if race {
		log.Error("can not use race in garble")
		return
	}
	privateGrable(code, cmd)
}

func privateGrable(code string, command []string) {
	_ = os.RemoveAll(filepath.Join(".", "output.exe"))
	newPath := filepath.Join(".", "output")
	_ = os.MkdirAll(newPath, os.ModePerm)
	_ = ioutil.WriteFile("output/main.go", []byte(code), 0777)
	cmd := exec.Command("./garble/garble.exe", command...)
	err := cmd.Run()
	if err == nil {
		log.Info("build success")
		log.Info("file: output.exe")
	} else {
		log.Error("error")
	}
	_ = os.RemoveAll(newPath)
}

func privateBuild(code string, command []string) {
	_ = os.RemoveAll(filepath.Join(".", "output.exe"))
	newPath := filepath.Join(".", "output")
	_ = os.MkdirAll(newPath, os.ModePerm)
	_ = ioutil.WriteFile("output/main.go", []byte(code), 0777)
	cmdPart := checkOS()
	for i, s := range command {
		if i == 5 {
			cmdPart += "='" + s + "'"
			continue
		}
		cmdPart += " "
		cmdPart += s
	}
	
	var err error
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C" , cmdPart)
		_ , err = cmd.Output()
	}else {
		cmd := exec.Command("sh", "-c" , cmdPart)
		_ , err = cmd.Output()
	}

	if err == nil {
		log.Info("build success")
		log.Info("file: output.exe")
	} else {
		log.Error("error")
	}
	_ = os.RemoveAll(newPath)
}

func checkOS() string {
	osName := runtime.GOOS
	if osName == "windows"{
		return ""
	}else {
		return "env CGO_ENABLED=0 GOOS=windows GOARCH=amd64"
	}
}