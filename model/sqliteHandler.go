package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type sqlliteHandler struct {
	db *sql.DB
}

// 인터페이스를 implement , func(m *memoryHandler)
func (s *sqlliteHandler) getTodos() []*Todo {
	return nil
}
func (s *sqlliteHandler) addTodo(name string) *Todo {
	return nil
}
func (s *sqlliteHandler) removeTodo(id int) bool {
	return false
}

func (s *sqlliteHandler) completeTodo(id int, complete bool) bool {
	return false

}

// 종료되기전에 db를 닫아준다
func (s *sqlliteHandler) close() {
	s.db.Close()

}

func newSqliteHandler() dbHandler {
	database, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	// data를 저장할 테이블 생성
	statement, _ := database.Prepare(
		// 한줄은 ""
		// 여러줄은 ``
		// todos table이 없으면 만들어라
		// AUTOINCREMENT 값이 자동으로 하나씩 증가한다
		`CREATE TABLE IF NOT EXISTS todos ( 
			id        INTEGER  PRIMARY KEY AUTOINCREMENT,
			name      TEXT,
			completed BOOLEAN,
			createdAt DATETIME
		)`)

	// 실행
	statement.Exec()
	return &sqlliteHandler{}
}
