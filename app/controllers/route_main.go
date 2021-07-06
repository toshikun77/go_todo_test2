package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

//124
//ハンドラーを記述します

//t, err := template.ParseFiles("top.html") と記述して、
//引数に渡したテンプレートファイル top.html を解析して、テンプレートの構造体を生成します
//構造体のポインタ型は、t に格納されます
//t.Execute(w, nil) と記述して、テンプレートを実行します
//第1引数 http.ResponseWriter を指定します
//第2引数 データを指定します
//こうして、解析済みのテンプレート t を実行します
//第2引数は、今回は、nil としますが、"XXXX" と記述してデータを指定すれば、
//tmpl.html の中の {{ . }} に第2引数のデータが組み込まれます
//そして、http.ResponseWriter にわたす最終的な html を生成します
//その html をブラウザに表示します
//"Hello" というデータを第2引数に渡すと、
//top.html の中の {{ . }} に第2引数のデータが組み込まれます
//そして、http.ResponseWriter にわたす最終的な html を生成します
//128
//server.go に記述した generateHTML() を使用します
//第1引数 http.ResponseWriter を指定します
//第2引数 ブラウザに表示したい data を指定します
//第3引数 使用するテンプレートファイルを指定します。可変長引数になります
//top ハンドラーは、layout.html と top.html を使うので、layout, top を string型で指定します
//131
//ハンドラーの内部の genarateHTML() の引数に、"public_navbar" のテンプレートファイルを追加します
//135
//_, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、session が存在せず、ログインしていないので、
//generateHTML(w, "Hello", "layout", "public_navbar", "top") と記述して、
//誰でもアクセスできる top ページを表示させます
//エラーが発生しない場合、session が存在するので、/todos のURLを指定して、
//ユーザ専用の index.html へリダイレクトします
func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, "Hello", "layout", "public_navbar", "top")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}
}

//135
//index ハンドラーを作成します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//sess, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、ログインしていないので、
//http.Redirect(w, r, "/", 302) と記述して、topページにリダイレクトさせます
//session が存在する場合は、
//generateHTML(w, nil, "layout", "public_navbar", "index") と記述して
//index.html を表示します
//第1引数 http.ResponseWriter を指定します
//第2引数 ブラウザに表示したい data を指定します。今回は、nil とします
//第3引数 使用するテンプレートファイルを指定します。可変長引数になります
//ハンドラーの内部の genarateHTML() の引数に、"private_navbar" のテンプレートファイルを追加します
//139
//sess, err := session(w, r) と修正します
//user, err := sess.GetUserBySession() と記述して
//user.go に記述したメソッドを使用してユーザのレコードを取得します
func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}

		//139
		//特定の user の todo を todos テーブルから取得し、todos に格納します
		todos, _ := user.GetTodosByUser()

		//user の Todos フィールドに、todos を格納します
		user.Todos = todos

		//第2引数の値を、user に修正しデータを渡します
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

//141
//todoNew ハンドラーを作成します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//sess, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、session が存在せず、ログインしていないので、
//http.Redirect(w, r, "/login", 302) と記述して、
//誰でもアクセスできる login ページを表示させます
//session が存在する場合は、
//generateHTML(w, nil, "layout", "private_navbar", "todo_new") と記述して、
//ユーザ専用の todo_new の画面を表示します

func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

//141
// todoSave ハンドラーを作成します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//sess, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、session が存在せず、ログインしていないので、
//http.Redirect(w, r, "/login", 302) と記述して、
//誰でもアクセスできる login ページを表示させます
func todoSave(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {

	//session が存在する場合は、
	//err = r.ParseForm() と記述して、ブラウザに入力されたデータを解析します
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//user, err := sess.GetUserBySession() と記述して
	//user.go に記述したメソッドを使用してユーザのレコードを取得します
	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
	}

	//解析したデータのうち、content 属性をもつ値を取得して、content に格納します
	content := r.PostFormValue("content")

	//if err := user.CreateTodo(content); err !=nil {}
	//と記述して、content を DB の todos テーブルに書き込みます
	if err := user.CreateTodo(content); err != nil {
		log.Println(err)
	}

	//DB の todos テーブルに書き込みが完了したら、/todos のユーザ専用ページにリダイレクトさせます
	http.Redirect(w, r, "/todos", 302)
	}
}

//143
//todoEdit ハンドラーを作成します
//このハンドラーは、parseURL() 関数内で実行されます

func todoEdit(w http.ResponseWriter, r *http.Request, id int) {

	//sess, err := session(w,r ) と記述して、
	//session 有無の判定結果を取得します
	//エラーが発生する場合は、session が存在せず、ログインしていないので、
	//http.Redirect(w, r, "/login", 302) と記述して、
	//誰でもアクセスできる login ページを表示させます
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {

	//_, err := sess.GetUserBySession() と記述して
	//user.go に記述したメソッドを使用してユーザのレコードを取得します
	//レコードを取得できない場合は、エラーが発生しますので、エラー処理を記述します
	_, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
	}

	//引数に渡した ID から todo を取得します
	t, err := models.GetTodo(id)
	if err != nil {
		log.Println(err)
	}

	//第2引数に、取得した todo が格納されている t を渡します
	//第5引数には、先程作成した、todo_edit を渡します
	generateHTML(w, t, "layout", "private_navbar", "todo_edit")
	}
}

//143
//todoUpdate ハンドラーを作成します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//このハンドラーは、parseURL() 関数内で実行されます

func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {

	//sess, err := session(w,r ) と記述して、
	//session 有無の判定結果を取得します
	//エラーが発生する場合は、session が存在せず、ログインしていないので、
	//http.Redirect(w, r, "/login", 302) と記述して、
	//誰でもアクセスできる login ページを表示させます
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {

	//session が存在する場合は、
	//err = r.ParseForm() と記述して、ブラウザに入力されたデータを解析します
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	//user, err := sess.GetUserBySession() と記述して
	//user.go に記述したメソッドを使用してユーザのレコードを取得します
	//レコードを取得できない場合は、エラーが発生しますので、エラー処理を記述します
	user, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
	}

	//解析したデータのうち、content 属性をもつ値を取得して、content に格納します
	content := r.PostFormValue("content")

	//構造体 Todo に紐づく、UpdateTodo() メソッドを使用したいので、
	//todo インスタンスを生成します
	t := &models.Todo{ID: id, Content: content, UserID: user.ID}

	//if err := t.UpdateTodo(); err !=nil {}
	//と記述して、content を DB の todos テーブルに書き込みます
	if err := t.UpdateTodo(); err != nil {
		log.Println(err)
	}

	//DB の todos テーブルに書き込みが完了したら、/todos のユーザ専用ページにリダイレクトさせます
	http.Redirect(w, r, "/todos", 302)
	}
}

//144
//todoDelete ハンドラーを作成します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//このハンドラーは、parseURL() 関数内で実行されます

func todoDelete(w http.ResponseWriter, r *http.Request, id int) {

	//sess, err := session(w,r ) と記述して、
	//session 有無の判定結果を取得します
	//エラーが発生する場合は、session が存在せず、ログインしていないので、
	//http.Redirect(w, r, "/login", 302) と記述して、
	//誰でもアクセスできる login ページを表示させます
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {

	//_, err := sess.GetUserBySession() と記述して
	//user.go に記述したメソッドを使用してユーザのレコードを取得します
	//レコードを取得できない場合は、エラーが発生しますので、エラー処理を記述します
	_, err := sess.GetUserBySession()
	if err != nil {
		log.Println(err)
	}

	//構造体 Todo に紐づく、DeleteTodo() メソッドを使用したいので、
	//削除対象のレコードを取得して、todo インスタンス t を生成します
	t, err := models.GetTodo(id)
	if err != nil {
		log.Println(err)
	}

	//if err := t.DeleteTodo(); err !=nil {}
	//と記述して、todo を DB の todos テーブルから削除します
	if err := t.DeleteTodo(); err != nil {
		log.Println(err)
	}

	//DB の todos テーブルからレコードを削除するのが完了したら、
	// /todos のユーザ専用ページにリダイレクトさせます
	http.Redirect(w, r, "/todos", 302)
	}
}