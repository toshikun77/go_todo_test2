package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"todo_app/app/models"
	"todo_app/config"
)

//128
//ハンドラー関数で、テンプレートを渡して表示する部分を記述します
//テンプレートを渡して表示する部分は、これから作成していく複数のハンドラー関数で
//共通で利用できるようにします
//第1引数 http.ResponseWriterを指定します
//第2引数 data を指定します
//第3引数 可変長引数として文字列型を指定します
//var files []string と記述して、文字列型のスライスを作成します
//fmt.Sprintf("app/views/templates/%s.html",file) と記述して、
//ファイルパスを作成します
//たとえば、app/views/templates/top のようなファイルパスを生成します
//files = append() と記述して、作成したファイルパスを 文字列型のスライス files に追加します
func generateHTML(w http.ResponseWriter, data interface{}, filnames ...string) {
	var files []string

	for _, file := range filnames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}

	//gotrading 89 テンプレートのキャッシュとハンドラーを復習してね
	//templates := template.Must(template.ParseFiles(files...))と記述して、
	//テンプレートファイル xxxx.html を読み込んで、templates に代入してキャッシングしておくよ
	//キャッシングしておくことで、テンプレートを毎回つくらずに済むよ
	//template.ParseFiles(files...) の引数の箇所は、files... として、
	//可変長引数とすることで、files に格納されている複数のテンプレートファイルを引数に渡せるよ
	//templates には、template.Template の構造体のポインタ型が格納されるよ
	templates := template.Must(template.ParseFiles(files...))

	//テンプレートを実行するよ
	//第1引数 http.ResponseWriterを指定します
	//第2引数 実行するテンプレートを指定する
	//第3引数 data を指定する
	//layout.html ファイル内で、{{define "layout"}} と定義しているので、
	//実行するテンプレートを明示的に指定する必要があります
	//第2引数は、"layout.html" ではなく、"layout" のみの記述でOKだよ
	templates.ExecuteTemplate(w, "layout", data)
}

//135
//server.go に cookie を取得する関数をつくっていきます
//cookie, err := r.Cookie("_cookie") と記述して、
//http.Request から http.Cookie の構造体のポインタ型を取得します
//http.Cookie の構造体のフィールドは以下となっています
//cookie := http.Cookie{
//	Name:     "_cookie",
//	Value:    session.UUID,
//	HttpOnly: true,
//}
//sess = models.Session{UUID: cookie.Value} と記述して、
//session インスタンスを生成します
//if ok, _ := sess.CheckSession(); !ok {} と記述し、
//session が存在していなければ、
//err = fmt.Errorf() でエラーを生成します
//session が存在すれば、エラーは返ってきません
//この関数を使って、アクセス制限を実装していきます
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			err = fmt.Errorf("Invalid session")
		}
	}
	return sess, err
}

//143
//"^/todos/(edit|update)/([0-9]+)$" と記述すると、
//URLの先頭が、todos で、/の後に、edit または、update で、その次の / のあとに、
//0-9 の文字列が1個以上入っているものですね。
//q = validPath.FindStringSubmatch(r.URL.Path) と記述すると、
//[/todos/edit/1 edit 1]とか、[/todos/update/1 update 1]　のような
//Pathの配列を得ることができるよ
//このようにして、1 のような ID がhttp.Request の構造体から摘出できるようになります
//144
//URL の正規表現のパターンに、delete のURLも追加します
//FindStringSubmatch() を使ったときに、/todos/delete/1 のようなURLも取得できるようになります

var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

//143
//parseURL ハンドラーを作成します
//今回はリクエストがあれば、そのURLからIDを取得する処理が必要なので
//そのコードを記述していきます
//引数は、ResponseWriter型と *Request型 と int型を指定します
//int型の箇所はIDを指定します
//返り値は、HandlerFunc 型 を指定します
//HandlerFunc 型は、ResponseWriter型と *Request型を引数に持つ関数の型です
//StartMainServer() 内で、http.HandleFunc() を使って、URLと関数を登録します
//この http.HandleFunc() を呼び出すには、インターフェースは不要で、
//第2引数に、単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ
//渡せばOKです
//そのファンクションを生成するための関数が、parseURL() 関数です
//parseURL() 関数の中で、todoEdit() 等のハンドラーが実行されます

func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	// /todos/edit/1
	//validPath と URL がマッチした箇所を取得したいので、FindStringSubmatch() を
	//使用します
	//第1引数に検索対象の文字列を渡します
	//第2引数に取得する回数を指定します。-1 の場合、無制限となります
	//Compile 済みの正規表現のパターンに一致したすべての文字列をスライスにして q に返します
	q := validPath.FindStringSubmatch(r.URL.Path)

	//fmt.Println(q)
	//[/todos/edit/1 edit 1]

	//スライス q が空であれば、NotFoundを返します
	if q == nil {
		http.NotFound(w, r)
		return
	}

	//q[2]に格納されているIDを Intger 型に変換します
	qi, err := strconv.Atoi(q[2])

	if err != nil {
		http.NotFound(w, r)
		return
	}

	//引数で渡した、http.ResponseWriter, *http.Request と、取得したID qi を
	//func に渡して、関数を実行します
	//parseURL 関数内で、1 などの ID を取得して、最終的には todoEdit() 等のハンドラーを実行します
	fn(w, r, qi)
	}
}



func StartMainServer() error {

	//126
	//URLへのアクセスに対して静的ページを返す場合は、http.FileServer を利用します。
	//config.Config.Static には、app/views が格納されているので
	//files := http.FileServer(http.Dir(config.Config.Static)) と記述すると、
	//以下のようなハンドラーのポインタ型が files に格納されます
	//&{app/views}
	//config.Config.Static には、app/views が格納されているので、
	//ファイルサーバーが実際に探すファイルの Path は、
	//以下のような、"app/views/static/......" のPath となります
	//"/app/views/static/file.txt"
	files := http.FileServer(http.Dir(config.Config.Static))

	//fmt.Println(files)
	//&{app/views}

	//126
	//ファイルサーバーが実際に探すファイルの Path には、/static/ は含まれないので、
	//"/static/" は削除したいよね
	//このときに、http.StripPrefix() 関数を使うよ
	//第1引数 ファイルサーバーが探すファイルの Path に含まれている、削除したい prefix を指定する
	//第2引数 ファイルサーバーのハンドラーのポインタ型を指定する
	//返り値は、第1引数で指定した、prefix が削除された状態のハンドラーのポインタ型を返すよ
	//以下のようなハンドラーのポインタ型が prefixHandler に格納されます
	//0x432e090
	//ファイルサーバーが実際に探すファイルの Path は、
	//"app/views/file.txt" となるよ
	prefixHandler := http.StripPrefix("/static/", files)

	//fmt.Println(prefixHandler)
	//0x432e090

	//126
	//http.Handle() を使って、URL とハンドラーを登録します
	//第1引数 URLを指定する
	//第2引数 ハンドラーのポインタ型を指定する
	//"static/" のURLでリクエストが届いたら、prefixHandler のハンドラーを実行するよ
	//このようにして、app/views 配下の静的ファイルを読み込みします
	//ちなみに、http.Handle() を呼び出すには、
	//http.Handler型(ServeHttpメソッドを持つインターフェース型)が必要です
	http.Handle("/static/", prefixHandler)

	//124
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行する関数を指定
	//"/"のアクセスがきた時は、top 関数を実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//top 関数は、app/controllers/route_main.go に記述します
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/", top)

	//131
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/"のアクセスがきた時は、signup ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//sign 関数は、app/controllers/route_auth.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/signup", signup)

	//133
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/"のアクセスがきた時は、login ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//login 関数は、app/controllers/route_auth.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/login", login)

	//133
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/"のアクセスがきた時は、login ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//login 関数は、app/controllers/route_auth.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/authenticate", authenticate)

	//135
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/todos"のアクセスがきた時は、index ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//index 関数は、app/controllers/route_main.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/todos", index)

	//136
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/logout"のアクセスがきた時は、logout ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//logout 関数は、app/controllers/route_auth.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/logout", logout)

	//141
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/todos/new"のアクセスがきた時は、todoNew ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//login 関数は、app/controllers/route_auth.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/todos/new", todoNew)

	//141
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//第2引数 実行するハンドラーを指定
	//"/todos/save"のアクセスがきた時は、todoSave ハンドラーを実行してください。というふうにできるんですね
	//リダイレクトやね
	//このように、http.HandleFunc() を複数記述して、URLごとに異なる処理を実行させることができます
	//todoSave 関数は、app/controllers/route_main.go に記述しています
	//ちなみに、http.HandleFunc() を呼び出すには、インターフェースは不要で、
	//単に http.ResponseWriter, *http.Request を引数に持つファンクションさえ渡せばOKです
	http.HandleFunc("/todos/save", todoSave)

	//143
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//"/todos/edit/" というように、末尾に "/" をつけています
	//"/" を付与しない場合は、URLが完全一致することを求められます
	//"/" を付与した場合は、要求された URLの先頭が、登録されたURLと一致するかを調べるんですね
	//"/todos/edit/1" とか、"/todos/edit/2" のような、末尾に数字が入ったURLの場合でも、
	//このハンドラーを実行することができるようになります
	//第2引数 実行するハンドラーを指定
	//parseURL(todoEdit) と記述します
	//ハンドラー関数をチェーンさせて実行させています
	//http.HandleFunc() の第2引数は、ResponseWriter型と *Request型を引数に持つ関数の型を
	//指定できます
	//そして、parseURL の返り値は、ResponseWriter型と *Request型を引数に持つ関数の型である、
	//HandlerFunc 型 となっていますので、parseURL() をhttp.HandleFunc() の第2引数に
	//渡すことができるのです
	//処理の流れとしては、/todos/edit/1 のURL でリクエストを受信した場合、
	//まず parseURL() が実行されます
	//そして、parseURL() 関数の処理の中で、リクエストに含まれる ID を取得して、
	//そのIDを第3引数に指定して、todoEdit ハンドラーが実行されます
	//todoEdit ハンドラーは、ID をもとにユーザの todo を DBの todos テーブルから取得して、
	//テンプレートに todo を渡して、HTML を生成して、ブラウザに出力します
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))

	//143
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//"/todos/update/" というように、"/" をつけています
	//"/" を付与しない場合は、URLが完全一致することを求められます
	//"/" を付与した場合は、要求された URLの先頭が、登録されたURLと一致するかを調べるんですね
	//"/todos/update/1" とか、"/todos/update/2" のような、末尾に数字が入ったURLの場合でも、
	//このハンドラーを実行することができるようになります
	//第2引数 実行するハンドラーを指定
	//parseURL(todoUpdate)) と記述します
	//ハンドラー関数をチェーンさせて実行させています
	//http.HandleFunc() の第2引数は、ResponseWriter型と *Request型を引数に持つ関数の型を
	//指定できます
	//そして、parseURL の返り値は、ResponseWriter型と *Request型を引数に持つ関数の型である、
	//HandlerFunc 型 となっていますので、parseURL() をhttp.HandleFunc() の第2引数に
	//渡すことができるのです
	//処理の流れとしては、/todos/update/1 のURL でリクエストを受信した場合、
	//まず parseURL() が実行されます
	//そして、parseURL() 関数の処理の中で、リクエストに含まれる ID を取得して、
	//todoUpdate ハンドラーが実行されます
	//todoUpdate ハンドラーは、ID をもとにユーザの todo を DBの todos テーブルに書き込みます、
	//その後、HTML を生成して、ブラウザに出力します
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))

	//144
	//http.HandleFunc() を使って、URLと関数を登録します
	//第1引数 URLを指定
	//"/todos/delete/" というように、"/" をつけています
	//"/" を付与しない場合は、URLが完全一致することを求められます
	//"/" を付与した場合は、要求された URLの先頭が、登録されたURLと一致するかを調べるんですね
	//"/todos/delete/1" とか、"/todos/delete/2" のような、末尾に数字が入ったURLの場合でも、
	//このハンドラーを実行することができるようになります
	//第2引数 実行するハンドラーを指定
	//parseURL(todoDelete)) と記述します
	//ハンドラー関数をチェーンさせて実行させています
	//http.HandleFunc() の第2引数は、ResponseWriter型と *Request型を引数に持つ関数の型を
	//指定できます
	//そして、parseURL の返り値は、ResponseWriter型と *Request型を引数に持つ関数の型である、
	//HandlerFunc 型 となっていますので、parseURL() をhttp.HandleFunc() の第2引数に
	//渡すことができるのです
	//処理の流れとしては、/todos/delete/1 のURL でリクエストを受信した場合、
	//まず parseURL() が実行されます
	//そして、parseURL() 関数の処理の中で、リクエストに含まれる ID を取得して、
	//todoDelete ハンドラーが実行されます
	//todoDelete ハンドラーは、ID をもとにユーザの todo を DBの todos テーブルに書き込みます、
	//その後、HTML を生成して、ブラウザに出力します
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))

	//124
	//第１引数は、string 型を指定します。
	//第1引数 ネットワークアドレス。デフォルト値は 80番ポートです。今回は 8080 にしています
	//Webサーバーを立ち上げるときに使うポートを記述してやるんですね
	//":"の前に何も書かなければ、ローカルホストとなります
	//今回の場合は、"localhost:8080" を引数に渡していることになります
	//第2引数 ハンドラーを指定します。 nil の場合は、デフォルトのマルチプレクサーを使います
	//nilを入れるとデフォルトのハンドラーでアクセスしたときに、"page not found" を返してくれます
	//今度は、&MyHandler{} を渡します
	//デフォルトのマルチプレクサーを使用しない場合は、すべての処理が、MyHandler によって実行されてしまいます
	//通常は、URLによって、処理を分けます
	//そのため、http.ListenAndServe() では、nil を渡してデフォルトのマルチプレクサーを使用します
	//なので、nil に戻しました

	//147
	//port := os.Getenv("PORT") と記述して、heroku 環境のポート番号を取得します
	//port には文字列型の値が格納されます
	port := os.Getenv("PORT")

	//return http.ListenAndServe(":"+port, nil) と記述して、サーバーを起動します
	return http.ListenAndServe(":"+port, nil)
}
