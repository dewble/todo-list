package model

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type pqHandler struct {
	db *sql.DB
}

// 인터페이스를 implement , func(m *memoryHandler)
func (s *pqHandler) GetTodos(sessionId string) []*Todo {
	// 데이터를 읽어와서 반환해줘야한다. 그 반환값을 가지고 있을 list 를 가지고 있어야한다
	todos := []*Todo{}

	// 세션에 해당하는것만 가쟈온다
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos WHERE sessionId=?", sessionId)
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

func (s *pqHandler) AddTodo(name string, sessionId string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (sessionId, name, completed, createdAt) VALUES (?, ?, ?, datetime('now'))")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(sessionId, name, false)
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

func (s *pqHandler) RemoveTodo(id int) bool {
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

func (s *pqHandler) CompleteTodo(id int, complete bool) bool {
	// 기존 레코드는 두고 completed 값만 변경

	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
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
func (s *pqHandler) Close() {
	s.db.Close()
}

func newPQHandler(dbConn string) DBHandler {
	database, err := sql.Open("postgre", dbConn)
	if err != nil {
		panic(err)
	}

	// data를 저장할 테이블 생성
	// 한줄은 ""
	// 여러줄은 ``
	// todos table이 없으면 만들어라
	// AUTOINCREMENT 값이 자동으로 하나씩 증가한다
	statement, err := database.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			sessionId STRING,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		);
		CREATE INDEX IF NOT EXISTS sessionIdIndexOnTodos ON todos (
			sessionId ASC
		);`)

	if err != nil {
		panic(err)
	}

	// 실행
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	// sqliteHandler 의 멤버변수 db 에 저장해준다
	return &pqHandler{db: database}
}
