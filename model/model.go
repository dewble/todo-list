package model

import (
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// 인터페이스 추가
type dbHandler interface {
	// 외부로 공개되는 인터페이스가 아니므로 소문자로 작성
	getTodos() []*Todo
	addTodo(name string) *Todo
	removeTodo(id int) bool
	completeTodo(id int, complete bool) bool
}

var handler dbHandler

// 초기화, 이 패키지가 처음으로 initialize 될때 실행, 한번만 호출됨
func init() {
	// handler = newMemoryHandler()
	// sqllite 사용시 아래 핸들러가 추가
	handler = newSqliteHandler()

}

// app.go Todo → model.Todo 로 변경
// app.go todoMap -> function 을 추가해서 사용
/*
func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	rd.JSON(w, http.StatusOK, list)
}
*/

func GetTodos() []*Todo {
	return handler.getTodos()
}

/*
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	id := len(todoMap) + 1
	todo := &Todo{id, name, false, time.Now()}
	todoMap[id] = todo
	rd.JSON(w, http.StatusCreated, todo)
}
*/
func AddTodo(name string) *Todo {
	return handler.addTodo(name)
}

/*
func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// if _, ok := todoMap[id]; ok {
	// 	delete(todoMap, id)
	// 	rd.JSON(w, http.StatusOK, Success{true})
	// } else {
	// 	rd.JSON(w, http.StatusOK, Success{false})
	// }
}
*/
func RemoveTodo(id int) bool {
	return handler.removeTodo(id)
}

/*
func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	if todo, ok := todoMap[id]; ok {
		todo.Completed = complete
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}
*/

func CompleteTodo(id int, complete bool) bool {
	return handler.completeTodo(id, complete)
}
