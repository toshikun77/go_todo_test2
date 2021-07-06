package controllers

import (
	"log"
	"net/http"
	"todo_app/app/models"
)

//131
//ハンドラーを記述します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//"/signup" でリクエストを受け取ると、signup用のHTMLファイルを生成して、
//クライアントにその生成したHTMLを返します
//server.go に記述した generateHTML() を使用します
//第1引数 http.ResponseWriter を指定します
//第2引数 ブラウザに表示したい data を指定します
//第3引数 使用するテンプレートファイルを指定します。可変長引数になります
//if r.Method == "GET" {} と記述して、メソッドリクエストが GET の場合の処理を記述します
//else if r.Method == "POST" {} と記述して、メソッドリクエストが POST の場合の処理を記述します
//err := r.ParseForm() と記述して、入力フォームの解析を行ないます
//user := models.User{} と記述して、User の構造体からインスタンスを生成します
//インスタンスの各フィールドに、入力フォームの解析結果を登録します
//r.PostFormValue("name") と記述することで、name 属性の値を取得できます
//if err := user.CreateUser(); err != nil {} と記述して、DBの users テーブルにレコードを登録します
//135
//_, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、session が存在せず、ログインしていないので、
//generateHTML(w, "Hello", "layout", "public_navbar", "signup") と記述して、
//誰でもアクセスできる signup ページを表示させます
//エラーが発生しない場合、session が存在するので、/todos のURLを指定して、
//ユーザ専用の index.html へリダイレクトします
func signup(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, "Hello", "layout", "public_navbar", "signup")
		} else {
			http.Redirect(w, r, "/todos", 302)
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}

		//ユーザの登録が完了した後は、TOPページにリダイレクトさせます
		//StatusCode は 302 とします
		//一時的に別のURLへ遷移させたい時に使用する status code です
		http.Redirect(w, r, "/", 302)
	}
}

//133
//ハンドラーを記述します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//"/login" でリクエストを受け取ると、login用のHTMLファイルを生成して、
//クライアントにその生成したHTMLを返します
//server.go に記述した generateHTML() を使用します
//第1引数 http.ResponseWriter を指定します
//第2引数 ブラウザに表示したい data を指定します。今回は nil とします
//第3引数 使用するテンプレートファイルを指定します。可変長引数になります
//135
//_, err := session(w,r ) と記述して、
//session 有無の判定結果を取得します
//エラーが発生する場合は、session が存在せず、ログインしていないので、
//generateHTML(w, "Hello", "layout", "public_navbar", "login") と記述して、
//誰でもアクセスできる login ページを表示させます
//エラーが発生しない場合、session が存在するので、/todos のURLを指定して、
//ユーザ専用の index.html へリダイレクトします
func login(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		http.Redirect(w, r, "/todos", 302)
	}

}


//133
//ユーザ認証用の authenticate ハンドラーを記述します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。
//"/authenticate" でリクエストを受け取ると、Webサーバーは、
//このハンドラーによって入力フォームの値を解析し、パスワード認証します。
//パスワード認証が成功した場合は、セッションを作成し、ログインした状態の画面へ遷移します
//パスワード認証が失敗した場合は、ログインしていない状態で、ログイン画面へ遷移します

func authenticate(w http.ResponseWriter, r *http.Request) {

	//err := r.ParseForm() と記述して、入力フォームの解析を行ないます
	err := r.ParseForm()
	//	if err != nil {
	//		log.Println(err)
	//	}

	//password 属性のある入力欄から、password の値を取得する
	//email 属性のある入力欄から、email の値を取得する
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")

	//user, err := models.GetUserByEmail(email) と記述して、
	//ログイン時に、ブラウザに入力された email アドレスから
	//DB内の users テーブルに存在する、ユーザのレコードを取得します
	//エラー発生時は、
	//http.Redirect(w, r, "/login", 302) と記述して、
	//ログイン画面にリダイレクトさせます
	user, err := models.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", 302)
	}

	//user.PassWord には、users テーブルから取得した、暗号化された値が登録されているので、
	//models.Encrypt(password) と記述することによって、
	//ブラウザに入力されたパスワードの値も暗号化してから、両者を比較します
	//もし同じ値なら、ログインできるので、session を生成します
	//後で、この session を使用して、クッキーを生成して、ログイン情報を保存します
	//値が異なる場合は、ログインできないので、/login にリダイレクトさせます
	if user.PassWord == models.Encrypt(password) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}

	//cookie を生成します
	//http.Cookie の構造体のフィールドに、
	//Name, Value, HttpOnly の値をそれぞれ登録します
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    session.UUID,
		HttpOnly: true,
	}

	//http.ResponseWriter のヘッダーに cookie のヘッダーを追加します
	//第1引数 http.ResponseWriter を指定する
	//第2引数 cookie のメモリアドレスを指定する
	http.SetCookie(w, &cookie)

	//cookie の生成後は、TOPページにリダイレクトさせます
	//本来はログインしたユーザだけが閲覧できるページにリダイレクトさせます
	//まだ、そのページは作成していませんので、TOPページにリダイレクトさせます
	//StatusCode は 302 とします
	http.Redirect(w, r, "/", 302)
	} else {

	//値が異なる場合は、ログインできないので、/login にリダイレクトさせます
	http.Redirect(w, r, "/login", 302)

	}
}

//136
//ログアウト用のハンドラーを記述します
//ハンドラーとは、Webサーバー上に存在し、
//クライアントから、あるURLのリクエストを受け取って、
//何らかの処理を実行し、クライアントへレスポンスを返す関数のことです
//つまり、APIのことです。

func logout(w http.ResponseWriter, r *http.Request) {

	//cookie, err := r.Cookie("_cookie") と記述して、
	//cookie を取得します
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	//if err != http.ErrNoCookie {} と記述して、
	//返ってきた err が http.ErrNoCookie でなければ、
	//session インスタンスを生成します
	//users.go で作成した、session.DeleteSessionByUUID() メソッドを実行して、
	//session を削除します
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}

	//login ページにリダイレクトさせます
	http.Redirect(w, r, "/login", 302)

}
