package utils

import (
	"io"
	"log"
	"os"
)

//func LoggingSettings(logFile string) {} と記述して、
//ログファイルを作成する関数をつくるよ(基礎編の28)
//logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//と記述して、ログファイルを作成します
//パーミッションは0666にして、すべてのユーザがこのログファイルに対して、
//読み取りと書き込みができるようにします。0666なので、このログファイルの実行はできません
//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) と記述して、
//日付、時間、ファイル名をログファイルに書き込みするように設定するよ
//log.Ldate = 日付| log.Ltime 　＝ 時間| log.Lshortfile ＝ ファイル名
//このようにして、ログファイルのフォーマットを設定するよ
//multiLogFile := io.MultiWriter(os.Stdout, logfile) と記述して、
//ログの書き込み先を、標準出力の os.Stdout と logfile の両方に指定します
//log.SetOutput(multiLogFile) と記述して、ログの出力先を設定するよ
//こうすると、画面でも結果が出力されて、ログファイルにも結果が書き込みされます
func LoggingSettings(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=logFile err=%s", err.Error())
	}
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}
