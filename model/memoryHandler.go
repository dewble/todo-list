package model

import "time"

// var todoMap map[int]*Todo 를 옮긴다
// DB를 사용하는 핸들러 추가, dbHandler가 implement하고 있다
// var handler dbHandler 이것만 들고 있으면 된다
type memoryHandler struct {
	todoMap map[int]*Todo
}

// 인터페이스를 implement , func(m *memoryHandler)
func (m *memoryHandler) getTodos() []*Todo {
	list := []*Todo{}
	for _, v := range m.todoMap {
		list = append(list, v)
	}
	return list
}
func (m *memoryHandler) addTodo(name string) *Todo {
	id := len(m.todoMap) + 1
	todo := &Todo{id, name, false, time.Now()}
	m.todoMap[id] = todo
	return todo
}
func (m *memoryHandler) removeTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok {
		delete(m.todoMap, id)
		return true
	}
	return false
}

func (m *memoryHandler) completeTodo(id int, complete bool) bool {
	if todo, ok := m.todoMap[id]; ok {
		todo.Completed = complete
		return true
	}
	return false
}

func newMemoryHandler() dbHandler {
	m := &memoryHandler{}
	m.todoMap = make(map[int]*Todo) // map을 initialize
	return m                        // dbHandler interface type으로 반환, memoryHandler 가 dbHandler를 implement하고 있기 때문에
}
