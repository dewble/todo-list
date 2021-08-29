package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

// 인터페이스를 implement , func(m *memoryHandler)
func (s *sqliteHandler) GetTodos() []*Todo {
	// 데이터를 읽어와서 반환해줘야한다. 그 반환값을 가지고 있을 list 를 가지고 있어야한다
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id,name,completed, createdAt FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// rows 는 행들이다 각 행마다 하나의 레코드를 나타내고 있다.
	// 각 행을 돌면서 레코드를 나타낸다
	// 다음행이 없으면 return false
	for rows.Next() {
		var todo Todo
		// todo 에 data 삽입
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		// todos 에 데이터 저장
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqliteHandler) AddTodo(name string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (name, completed, createdAt) VALUES (?,?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(name, false)
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()

	return &todo
}

func (s *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}

	// 영향을 받은 레코드의 개수를 알려준다
	cnt, _ := rst.RowsAffected()

	// 영향받은 행이 있을경우, true
	return cnt > 0
}

func (s *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	// 기존 레코드는 두고 completed 값만 변경

	stmt, err := s.db.Prepare("UPDATE todos SET completed = ? WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}

	// 영향을 받은 레코드의 개수를 알려준다
	cnt, _ := rst.RowsAffected()

	// 영향받은 행이 있을경우, true
	return cnt > 0

}

// 종료되기전에 db를 닫아준다
func (s *sqliteHandler) Close() {
	s.db.Close()
}

func newSqliteHandler(filepath string) DBHandler {
	database, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}

	// data를 저장할 테이블 생성
	// 한줄은 ""
	// 여러줄은 ``
	// todos table이 없으면 만들어라
	// AUTOINCREMENT 값이 자동으로 하나씩 증가한다
	statement, _ := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		)`)

	// 실행
	statement.Exec()

	// sqliteHandler 의 멤버변수 db 에 저장해준다
	return &sqliteHandler{db: database}
}
