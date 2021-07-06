package models

import (
	"log"
	"time"
)

type Todo struct {
	ID        int
	Content   string
	UserID    int
	CreatedAt time.Time
}

//117
//Todos テーブルにレコードを追加するメソッドを記述します
//User の構造体のメソッドとして作成します
//引数には、content を渡します
func (u *User) CreateTodo(content string) (err error) {
	cmd := `INSERT INTO todos (
		content,
		user_id,
		created_at) VALUES ($1, $2, $3)`

	//Db は、models パッケージに含まれていますので、Db.Exec() を実行できます
	//values (?,?,?) に対応する値を、Db.Exec() に渡します
	//created_at は、time.Now() を呼び出して、その返り値を使用します
	_, err = Db.Exec(cmd,
		content,
		u.ID,
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

//118
//DBの todos テーブルからレコードを取得するコードを書いていきます

func GetTodo(id int) (todo Todo, err error) {

	//todo インスタンスを生成します
	todo = Todo{}

	//SELECTコマンドを作成します
	//検索する id は、後ほど渡すので、WHERE id = ? と記述します
	cmd := `select id, content, user_id, created_at
	from todos where id = $1`

	//Db.QueryRow(cmd, id) で一行のレコードを取得します
	//第1引数 コマンドを指定する
	//第2引数　検索 id の value を指定する
	//QueryRow() メソッドによって、 1つだけレコードを取得できます
	//.Scan(&todo.ID, .....) と記述して、user フィールドに、取得したレコードを登録します
	err = Db.QueryRow(cmd, id).Scan(
		&todo.ID,
		&todo.Content,
		&todo.UserID,
		&todo.CreatedAt,
	)

	return todo, err
}

//119
//DBの todos テーブルから複数のレコードを取得するコードを書いていきます

func GetTodos() (todos []Todo, err error) {

	//複数データを取得します
	cmd := "select id, content, user_id, created_at FROM todos"

	//コマンドの条件に合う全てのレコードについて、DBのテーブルからすべて取得します
	rows, err := Db.Query(cmd)

	if err != nil {
		log.Println(err)
	}

	//for rows.Next() {} と記述して、取得したレコードを順番に、todos のインスタンスのフィールドに
	//格納していきます
	for rows.Next() {
		var t Todo
		err = rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedAt)
		if err != nil {
			log.Println(err)
		}
		todos = append(todos, t)
	}

	//Db.Query() の場合は、クローズの解放処理が必要です
	rows.Close()

	return todos, err
}

//120
//DBの todos テーブルから、特定のユーザのレコードを取得するコードを書いていきます

func (u *User) GetTodosByUser() (todos []Todo, err error) {

	//複数データを取得します
	//where user_id = ? と記述することで、特定のユーザIDのカラムをもつレコードだけ取得します
	cmd := `select id, content, user_id, created_at FROM todos where user_id = $1`

	//コマンドの条件に合う全てのレコードについて、DBのテーブルからすべて取得します
	rows, err := Db.Query(cmd, u.ID)

	if err != nil {
		log.Println(err)
	}

	//for rows.Next() {} と記述して、取得したレコードを順番に、todos のインスタンスのフィールドに
	//格納していきます
	for rows.Next() {
		var t Todo
		err = rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedAt)
		if err != nil {
			log.Println(err)
		}
		todos = append(todos, t)
	}

	//Db.Query() の場合は、クローズの解放処理が必要です
	rows.Close()

	return todos, err
}

//121
//Todos テーブル内のレコードを更新するメソッドを書いていきます

func (t *Todo) UpdateTodo() error {
	//todos テーブルのレコードを更新します
	//cmd := `update todos set content = ?, user_id = ? where id =?` と記述することで
	//Db.Exec() で渡した値を UPDATE コマンドに設定できます
	//このような記述にすることで、SQLインジェクションによる、レコードの書き換えを防ぐことができます
	cmd := `update todos set content = $1, user_id = $2 where id =$3`

	//コマンドを実行します
	//結果を使用しないので、_, err と記述します
	//第1引数 コマンドを指定します
	//第2引数 content = ? に渡す値を指定します
	//第3引数 user_id = ? に渡す値を指定します
	//t.ID の値をもつレコードのカラム を t.Content と t.UserID に更新します
	_, err := Db.Exec(cmd, t.Content, t.UserID, t.ID)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}

//122
//DBの todos テーブルのレコードを削除するコードを書いていきます

func (t *Todo) DeleteTodo() (err error) {

	//データの削除
	//WHERE で id を指定します
	cmd := `delete from todos where id =$1`

	//コマンドを実行します
	//結果を使用しないので、_, err と記述します
	//第1引数 コマンドを指定します
	//第2引数 where id = ? に渡す値を指定します
	//この場合は、id == t.ID である、レコードを削除することになります
	_, err = Db.Exec(cmd, t.ID)

	if err != nil {
		log.Fatalln(err)
	}

	return err
}
