package executor

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"io/ioutil"
	"bufio"
	"syscall"
	"encoding/csv"
)

func isRequestFile(path string) bool {
	return strings.HasSuffix(strings.ToUpper(path), ".REQ")
}

func requestToMe(machineName string) bool {
	if machineName == "*" {
		return true
	}

	mine, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(mine) != strings.ToUpper(machineName)
}

// リクエストファイルを読み込みコマンドを実行をする
func Process(path string) (flag bool, stdout, stderr string, exitCode int, err2 error) {
	if !isRequestFile(path) {
		// リクエストファイルではない
		return
	}

	machineName, command := getContent(path)

	if len(machineName) > 0 && len(command) > 0 {
		if requestToMe(machineName) {
			// ホスト名が一致しない
			return
		}

		stdout, stderr, exitCode, err2 = execute(command)
		fmt.Println("stdout: ", stdout)
		fmt.Println("stderr: ", stderr)
		fmt.Println("exitCode: ", exitCode)
		fmt.Println("err: ", err2)
		flag = true
	}

	return
}

func getContent(path string) (machineName, command string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	// マシン名
	if !s.Scan() {
		println("ホスト名が読み取れません: ", path)
		return
	}
	machineName = s.Text()

	// コマンド
	if !s.Scan() {
		println("コマンドが読み取れません: ", path)
		return
	}
	command = s.Text()

	return
}

// [execcommandexample/main\.go at master · hnakamur/execcommandexample](https://github.com/hnakamur/execcommandexample/blob/master/main.go)
func execute(command string) (stdout, stderr string, exitCode int, err error) {
	fmt.Println("実行します: ", command)
	// '"'で囲まれた部分は1要素として' 'で区切る
	r := csv.NewReader(strings.NewReader(command))
    r.Comma = ' '
    fields, err := r.Read()
    if err != nil {
        panic(err)
    }
	commands := fields

	cmd := exec.Command(commands[0], commands[1:]...)
	stdout, stderr, exitCode, err = runCommand(cmd)
	return
}

func runCommand(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err = cmd.Start(); err != nil {
		panic(err)
	}

	stdout_, err2 := ioutil.ReadAll(outReader)
	if err2 != nil {
		panic(err2)
	}
	stdout = string(stdout_)
	stderr_, err2 := ioutil.ReadAll(errReader)
	if err2 != nil {
		panic(err2)
	}
	stderr = string(stderr_)

	if err = cmd.Wait(); err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}
