# goropbox

Dropboxのファイル同期機能を利用してリモートPCでコマンドを実行します。

リクエストをファイルとしてDropbox上のディレクトリに配置すると、リモートPCに同期されるとgoropboxによりコマンドが実行され実行結果がファイルとして作成されます。

![goropbox](goropbox.png)

## 使い方

リモートPCでgoropboxを実行する
```
$ go run goropbox.go -dir 監視ディレクトリ -interval 監視する間隔 -count 監視する回数
```

ローカルPCのDropbox上の監視ディレクトリにリクエストファイルを配置する

vi 監視するディレクトリ/request1.req
```
マシン名
sh -c "ls -l /"
```

DropboxでリモートPCにリクエストファイルが同期され、goropboxによりコマンドが実行され結果がファイルとして出力されます。

### 実行結果

cat 監視するディレクトリ/request1.req.res
```
マシン名
2019-06-18 19:40:33.3194002 +0900 DST m=+0.002673801
2019-06-18 19:40:33.3341825 +0900 DST m=+0.017456101
0
```

### 標準出力

標準出力がある場合のみ出力される。

cat 監視するディレクトリ/request1.req.out
```
total 112
drwxr-xr-x  1 root root    512 Jun  4 15:09 bin
drwxr-xr-x  1 root root    512 Jul 26  2018 boot
drwxr-xr-x  1 root root    512 Jun 20 14:01 dev
drwxr-xr-x  1 root root    512 Jun  7 11:13 etc
drwxr-xr-x  1 root root    512 Apr 18 15:29 home
-rwxr-xr-x  1 root root 112600 Jan  1  1970 init
drwxr-xr-x  1 root root    512 May 23 17:39 lib
drwxr-xr-x  1 root root    512 May 23 17:23 lib32
drwxr-xr-x  1 root root    512 Jul 26  2018 lib64
drwxr-xr-x  1 root root    512 Jul 26  2018 media
drwxr-xr-x  1 root root    512 Apr 18 15:28 mnt
drwxr-xr-x  1 root root    512 Jul 26  2018 opt
dr-xr-xr-x 24 root root      0 Jun 18 12:23 proc
drwx------  1 root root    512 Jul 26  2018 root
drwxr-xr-x  1 root root    512 Jun 18 12:23 run
drwxr-xr-x  1 root root    512 May 23 17:20 sbin
drwxr-xr-x  1 root root    512 Jul 19  2018 snap
drwxr-xr-x  1 root root    512 Jul 26  2018 srv
dr-xr-xr-x 12 root root      0 Jun 18 12:23 sys
drwxrwxrwt  1 root root    512 Jun 20 17:36 tmp
drwxr-xr-x  1 root root    512 May 23 17:23 usr
drwxr-xr-x  1 root root    512 Jul 26  2018 var
```

### 標準出力

エラー出力がある場合のみ出力される。

cat 監視するディレクトリ/request1.req.err
```
```

### ヘルプ
```
$ go run goropbox.go --help
```

## テキストファイル仕様

リクエストファイル(*.req)
```
マシン名
コマンド [引数]
```

結果ファイル(*.res)
```
マシン名
実行開始日時
実行完了日時
戻り値
```

