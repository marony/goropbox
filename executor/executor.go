package executor

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"io/ioutil"
	"bufio"
	"syscall"
)

func isRequestFile(path string) bool {
	return strings.Contains(strings.ToUpper(path), ".REQ")
}

// リクエストファイルを読み込みコマンドを実行をする
// リクエストファイルをリネームし、l実行結果をファイルとして出力する
func Process(path string) {
	if !isRequestFile(path) {
		println("リクエストファイルではありません: ", path)
		return
	}

	fmt.Println("処理します: ", path)
	machineName, command := getContent(path)

	mine, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	if strings.ToUpper(mine) != strings.ToUpper(machineName) {
		println("ホスト名が一致しません: ", mine, " != ", machineName)
		return
	}

	stdout, stderr, exitCode, err := execute(command)
	fmt.Println("stdout: ", stdout)
	fmt.Println("stderr: ", stderr)
	fmt.Println("exitCode: ", exitCode)
	fmt.Println("err: ", err)

	// TODO: リネーム & 結果ファイル作成

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
		println("ホスト名が読み取れません")
		return
	}
	machineName = s.Text()

	// コマンド
	if !s.Scan() {
		println("コマンドが読み取れません")
		return
	}
	command = s.Text()

	return
}

// [execcommandexample/main\.go at master · hnakamur/execcommandexample](https://github.com/hnakamur/execcommandexample/blob/master/main.go)
func execute(command string) (stdout, stderr string, exitCode int, err error) {
	fmt.Println("実行します: ", command)
	commands := strings.Split(command, " ")
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

	err = cmd.Start()
	if err != nil {
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

	err = cmd.Wait()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}
