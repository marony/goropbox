package main

// go get github.com/mitchellh/go-homedir

import (
	"./monitor"
	"./executor"
	"fmt"
	"os"
	"os/signal"
	"github.com/mitchellh/go-homedir"
	_ "reflect"
)

// 監視するディレクトリ
const dropbox_dir = "~/Dropbox/goropbox"
// 

// まいんちゃん
func main() {
	
	// [【Go言語】Ctrl\+cなどによるSIGINTの捕捉とdeferの実行 \- DRYな備忘録](http://otiai10.hatenablog.com/entry/2018/02/19/165228)
	defer teardown()

	// シグナル用チャネル
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

	// 終了検知用チャネル
    done := make(chan error, 1)
    go do(done)

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
func do(done chan<- error) {
	// ディレクトリパスの"~"を展開する
	dir, err := homedir.Expand(dropbox_dir)
	if err != nil {
		panic(err)
	}

	fmt.Println("監視ディレクトリ: " + dir)

	monitor.Execute(dir, executor.Process)

	// 終了
	done <- nil
    close(done)
}

func teardown() {
	fmt.Println("データのあとかたづけ")
}
