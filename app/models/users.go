package models

import (
	"log"
	"time"
)

//User の構造体を作成します
//139
//todos フィールドを User の構造体に追加して、データをまとめます
type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo
}

//133
//構造体 Session の型を宣言します
type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatedAt time.Time
}

//DB の users テーブルに、User レコードを作成するメソッドを記述します
func (u *User) CreateUser() (err error) {
	cmd := `insert into users (
		uuid,
		name,
		email,
		password,
		created_at) values ($1,$2,$3,$4,$5)`

	//Db は、models パッケージ内の base.go で宣言していますので、Db.Exec() を実行できます
	//values (?,?,?,?,?) に対応する値を、Db.Exec() に渡します
	//UUDI は、base.go に記述した、createUUID() を呼び出して生成します
	//Encrypt(u.PassWord) は、平文のパスワードから、ハッシュ値を生成して、そのハッシュ値を16進数の文字列にした値となります
	//以下のような文字列になります
	//cf23df2207d99a74fbe169e3eba035e633b65d94
	//created_at は、time.Now() を呼び出して、その返り値を使用します
	_, err = Db.Exec(cmd,
		createUUID(),
		u.Name,
		u.Email,
		Encrypt(u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

//112
//DBの users テーブルから、ユーザのレコードを取得するコードを書いていきます

func GetUser(id int) (user User, err error) {

	//user インスタンスを生成します
	user = User{}

	//SELECTコマンドを作成します
	//検索する id は、後ほど渡すので、WHERE id = ? と記述します
	cmd := `select id, uuid, name, email, password, created_at from users where id = $1`

	//Db.QueryRow(cmd, id) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 id の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&user.ID, .....) と記述して、user フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)

	return user, err
}

//113
//users テーブルのレコードを更新するコードを書いていきます

func (u *User) UpdateUser() (err error) {

	//Update コマンドを作成します
	//Update する レコードの id は、後ほど渡すので、WHERE id = ? と記述します
	//このような記述にすることで、SQLインジェクションによる、レコードの書き換えを防ぐことができます
	cmd := `update users set name = $1, email = $2 where id = $3`

	//コマンドを実行します
	//結果を使用しないので、_, err と記述します
	//第1引数 コマンドを指定します
	//第2引数 name = ? に渡す値を指定します
	//第3引数 email = ? に渡す値を指定します
	//第4引数 where id = ? に渡す値を指定します
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

//114
//DBのusers テーブルのレコードを削除するコードを書いていきます

func (u *User) DeleteUser() (err error) {

	//データの削除
	//WHERE で id を指定します
	cmd := "delete from users where id =$1"

	//コマンドを実行します
	//結果を使用しないので、_, err と記述します
	//第1引数 コマンドを指定します
	//第2引数 where id = ? に渡す値を指定します
	//この場合は、id == u.ID である、レコードを削除することになります
	_, err = Db.Exec(cmd, u.ID)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

//133
//ログイン時に、ブラウザに入力された email アドレスから
//DB内の users テーブルに存在する、ユーザのレコードを取得するコードを書いていきます

func GetUserByEmail(email string) (user User, err error) {

	//user インスタンスを生成します
	user = User{}

	//SELECTコマンドを作成します
	//検索する email は、後ほど渡すので、WHERE email = ? と記述します
	cmd := `select id, uuid, name, email, password, created_at
	from users where email = $1`

	//Db.QueryRow(cmd, email) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 email の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&user.ID, .....) と記述して、user フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)

	return user, err
}

//133
//sessions テーブルに、
//セッションのレコードを追加するメソッドを書いていきます
//構造体 User のメソッドとします

func (u *User) CreateSession() (session Session, err error) {

	//session インスタンスを生成します
	session = Session{}

	//1つ目のコマンドで、セッション情報のレコードを DB の sessions テーブルに書き込みます

	cmd1 := `INSERT INTO sessions (
		uuid,
		email,
		user_id,
		created_at) VALUES ($1, $2, $3, $4)`

	//Db は、models パッケージに含まれていますので、Db.Exec() を実行できます
	//values (?,?,?,?) に対応する値を、Db.Exec() に渡します
	//UUDI は、createUUID() を呼び出して生成します
	//email と user_id は、u.Email と u.ID を渡します
	//created_at は、time.Now() を呼び出して、その返り値を使用します
	_, err = Db.Exec(cmd1,
		createUUID(),
		u.Email,
		u.ID,
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	//2つ目のコマンドで、セッション情報のレコードを DB の sessions テーブルから読み込みます
	//1つ目のコマンドで、DBの sessions テーブルへ書き込んだレコードの情報を、
	//DBの sessions テーブルから再度取得することになります

	//where user_id = ? and email = ? として、ユニークなレコードの情報を取得します
	cmd2 := `select id, uuid, email, user_id, created_at
	from sessions where user_id = $1 and email = $2`

	//Db.QueryRow(cmd, u.ID, u.Email) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 user_id の value を指定する
	//第3引数　検索 email の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&session.ID, .....) と記述して、session フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserID,
		&session.CreatedAt,
	)

	return session, err
}

//133
//セッションのレコードが、
//DB の sessions テーブルに存在するか否かを確認するメソッドを作成します
//session インスタンスの UUID フィールドの値で、sessions テーブルのレコードを検索するよ
//UUIDは以下のような値となっているよ
//c50bbe62-a4b3-11eb-88cd-7c04d0c2f23c

func (sess *Session) CheckSession() (valid bool, err error) {

	//where uuid = ? として、コマンドを作成します
	cmd := `select id, uuid, email, user_id, created_at
	from sessions where uuid = $1`

	//Db.QueryRow(cmd, sess.UUID) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 uuid の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&sess.ID, .....) と記述して、session フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserID,
		&sess.CreatedAt,
	)

	//error が発生する場合は、DBの sessions テーブルからのレコードの取得が失敗していることになります
	//セッションがDB の sessions テーブルに存在しないので、false とします
	if err != nil {
		valid = false
		return valid, err
	}

	//セッションのIDの値は１以上の整数となるため、
	//セッションのIDが初期値の 0 でなければ、DBの sessions テーブルからのレコードの取得が成功していることになります
	//セッション情報のレコードが、DB の sessions テーブルに存在するので、true とします
	if sess.ID != 0 {
		valid = true
	}
	return valid, err

}

//136
//DB の sessions テーブルから session の情報のレコードを削除するメソッドを作成します
//cmd := `delete from sessions where uuid = ?` と記述してコマンドを作成します
//_, err = Db.Exec(cmd, sess.UUID) と記述して、コマンドを実行します
//if err != nil {} と記述して、エラーが発生している場合は、エラーを出力します

func (sess *Session) DeleteSessionByUUID() (err error) {

	cmd := `delete from sessions where uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

//139
//session の構造体のフィールドの UserID を使って、
//users テーブルから、ユーザのレコードを取得するメソッドをつくります

func (sess *Session) GetUserBySession() (user User, err error) {

	//user インスタンスを生成します
	user = User{}

	cmd := `select id, uuid, name, email, created_at
	from users where id = $1`

	//Db.QueryRow(cmd, sess.UserID) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 session の id の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&user.ID, .....) と記述して、user フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd, sess.UserID).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	return user, err
}
