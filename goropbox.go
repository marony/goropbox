package main

// go get github.com/mitchellh/go-homedir

import (
	"./monitor"
	"./executor"
	"fmt"
	"flag"
	"os"
	"os/signal"
	"github.com/mitchellh/go-homedir"
	"time"
	"strconv"
	_ "reflect"
)

// まいんちゃん
func main() {
	// コマンド引数
	var (
		dir = flag.String("dir", "~/Dropbox/goropbox", "監視するディレクトリ")
		interval = flag.Int("interval", 60, "監視する間隔(秒)")
		count = flag.Int("count", 0, "監視する回数(0 = 無限)")
	)
	flag.Parse()

	// [【Go言語】Ctrl\+cなどによるSIGINTの捕捉とdeferの実行 \- DRYな備忘録](http://otiai10.hatenablog.com/entry/2018/02/19/165228)
	defer teardown()

	// シグナル用チャネル
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

	// 終了検知用チャネル
    done := make(chan error, 1)
    go do(done, *dir, *interval, *count)

    select {
    case sig := <-c:
        fmt.Println("シグナル来た:", sig)
        /*
         teardown中に再度SIGINTが来る場合を考慮し、
         send on closed channelのpanicを避ける。
       */
        // close(c)
        return
    case err := <-done:
        fmt.Println("loopの終了:", err)
    }

	fmt.Println("終了")
}

// 実際の処理
func do(done chan<- error, dir string, interval, count int) {
	// ディレクトリパスの"~"を展開する
	dir, err := homedir.Expand(dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("監視ディレクトリ: " + dir)

	monitor.Execute(dir, interval, count, executor.Process, complete)

	// 終了
	done <- nil
    close(done)
}

// リクエストファイルをリネームし、実行結果をファイルとして出力する
func complete(path, stdout, stderr string, exitCode int, err error) {
	mine, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// 結果出力
	{
		f, err := os.Create(path + ".res")
		if err != nil {
			panic(err)
		}

		defer f.Close()
		f.Write(([]byte)(mine))
		f.Write(([]byte)("\n"))
		f.Write(([]byte)(time.Now().String()))
		f.Write(([]byte)("\n"))
		f.Write(([]byte)(strconv.Itoa(exitCode)))
		f.Write(([]byte)("\n"))
	}

	if len(stdout) > 0 {
		f, err := os.Create(path + ".out")
		if err != nil {
			panic(err)
		}

		defer f.Close()
		f.Write(([]byte)(stdout))
	}

	if len(stderr) > 0 {
		f, err := os.Create(path + ".err")
		if err != nil {
			panic(err)
		}

		defer f.Close()
		f.Write(([]byte)(stderr))
	}

	// リクエストファイルのファイル名変更
	if err = os.Rename(path, path + ".done"); err != nil {
		panic(err)
	}
}

func teardown() {
	fmt.Println("データのあとかたづけ")
}
