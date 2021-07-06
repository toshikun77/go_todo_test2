package config

import (
	"log"
	"todo_app/utils"

	"gopkg.in/go-ini/ini.v1"
)

//106
//構造体 ConfigList を作成します
//config.ini を読み込んで、値を登録するために、構造体を宣言します
//この構造体のフィールドに、config.ini で読み込んだ値を登録します
//126
//Static フィールドを追加します
type ConfigList struct {
	Port      string
	SQLDriver string
	DbName    string
	LogFile   string
	Static    string
}

//106
//構造体 ConfigList の型をもつ、Config を作成します
//var Config ConfigList と記述して、グローバルに宣言します
//外部のパッケージからこの Config インスタンスを呼び出すことができるようになります
var Config ConfigList

//106
//func init() {} と記述すると、main() 関数より先に関数を実行できるよ
//LoadConfig() を init() の中で実行することで、main() より先に実行できるよ
//108
//utils.LoggingSettings(Config.LogFile) と記述して、config.ini で指定したファイル名の
//ログファイルを作成して、ログを書き込み始めるよ
func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

//106
//func LoadConfig() {} で config.ini を読み込むよ
//cfg, _ := ini.Load("config.ini") と記述すると、
//config.ini ファイルを読み込むよ
//ConfigList{} 、つまり ConfigList の 構造体の中で、
//Port:      cfg.Section("web").Key("port").MustInt(8080), と記述すると、
//config.ini ファイルの "8080" の値を取ってきて、 ConfigList の Port フィールドに設定してくれるよ
//値が空の場合は、デフォルト値として 8080 と ConfigList に設定されるよ
//ConfigList{} 、つまり ConfigList の 構造体の中で、
//DbName:    cfg.Section("db").Key("name").String(), と記述すると、
//config.ini ファイルの "webapp.sql" の値を取ってきて、 ConfigList のDbName フィールドに
//設定してくれるよ
//値が空の場合は、デフォルト値として、空文字が ConfigList に設定されるよ
//ConfigList{} 、つまり ConfigList の 構造体の中で、
//SQLDriver: cfg.Section("db").Key("driver").String(), と記述すると、
//config.ini ファイルの "sqlite3" の値を取ってきて、 ConfigList の SQLDriver フィールドに
//設定してくれるよ
//値が空の場合は、デフォルト値として、空文字が ConfigList に設定されるよ
//LogFile: cfg.Section("web").Key("logfile").String(), と記述すると、
//config.ini ファイルの "webapp.log" の値を取ってきて、 ConfigList の LogFile フィールドに
//設定してくれるよ
//値が空の場合は、デフォルト値として、空文字が ConfigList に設定されるよ
//126
//Static:    cfg.Section("web").Key("static").String(), と記述すると、
//config.ini ファイルの "app/views" の値を取ってきて、 ConfigList の Static フィールドに
//設定してくれるよ
//値が空の場合は、デフォルト値として、空文字が ConfigList に設定されるよ
func LoadConfig() {

	cfg, err := ini.Load("config.ini")

	if err != nil {
		log.Fatalln(err)
	}

	Config = ConfigList{
		Port:      cfg.Section("web").Key("port").MustString("8080"),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		LogFile:   cfg.Section("web").Key("logfile").String(),
		Static:    cfg.Section("web").Key("static").String(),
	}

}
