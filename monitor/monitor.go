package monitor

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"
)

type FileInfos []os.FileInfo
type ByModTime struct{ FileInfos }

func (fi ByModTime) Len() int {
    return len(fi.FileInfos)
}

func (fi ByModTime) Swap(i, j int) {
    fi.FileInfos[i], fi.FileInfos[j] = fi.FileInfos[j], fi.FileInfos[i]
}

// 古いもの順にソート
func (fi ByModTime) Less(i, j int) bool {
    return fi.FileInfos[i].ModTime().Unix() < fi.FileInfos[j].ModTime().Unix()
}

// ディレクトリからリクエストファイルを検索してfを呼び出す
func Execute(dir string, interval, count int, process func(string) (bool, string, string, int, error), complete func(string, time.Time, string, string, int, error)) {
	f := func() {
		// ファイル一覧の取得
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		sort.Sort(ByModTime{files})
		for _, fi := range files {
			path := filepath.Join(dir, fi.Name())
			startTime := time.Now()
			flag, stdout, stderr, exitCode, err := process(path)
			if flag {
				complete(path, startTime, stdout, stderr, exitCode, err)
			}
		}
	}
	
	if count == 0 {
		for {
			f()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	} else {
		for i := 0; i < count; i++ {
			f()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}
}
