package app

import (
	"encoding/json"
	"example/todolist2/model"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {

	getSessionID = func(r *http.Request) string {
		return "testsessionId"
	}

	// db를 지워주지 않고 테스트를 다시하면 db에 테스트값이 더해져 테스트 오류 발생
	os.Remove("./test.db")
	assert := assert.New(t)

	// app handler, db 닫아주기
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	/*
		POST 테스트
	*/
	// addTodoHandler 에서 FormValue 로 보내기때문에 PostForm으로 받는다
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 	rd.JSON(w, http.StatusCreated, todo) 읽어보기

	// refactoring 하며 todo -> model.Todo로 변경
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	// 서버가 저장한 ID
	id1 := todo.ID

	// addTodoHandler 에서 FormValue 로 보내기때문에 PostForm으로 받는다
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 	rd.JSON(w, http.StatusCreated, todo) 읽어오기
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2")
	id2 := todo.ID

	/*
	 GET 테스트, getTodoListHandle에서 JSON으로 todo list를 받는다
	*/
	/*
		func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
		list := []*Todo{}
		for _, v := range todoMap {
			list = append(list, v)
		}
		rd.JSON(w, http.StatusOK, list)
		}
	*/
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}

	// JSON 으로 읽어오기
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)

	// todos 개수 확인
	assert.Equal(len(todos), 2)

	// todo list 검증
	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.ID == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("TestID should be id1 or id2"))
		}
	}

	/*
		GET - complete-todo 테스트
	*/

	// formvalue 로 complete := r.FormValue("complete") == "true" 를 받는다

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
	// integer 와 string의 + 가 안되므로, strconv로 변경해준다
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// GET 테스트, complete 검증, 변경되는지 확인
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}

	// JSON 으로 읽어오기
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)

	// todos 개수 확인
	assert.Equal(len(todos), 2)

	// todo completed 검증
	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	/*
	 DELETE 테스트
	*/

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// DELETE 됐는지 확인
	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}

	// JSON 으로 읽어오기
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)

	// todos 개수 확인
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}

}
