package app

import (
	"example/todolist2/model"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// 쿠키스토어를 만든다.
// SESSION_KEY 라는 환경 변수를 만들어주고 가져와서 사용한다.
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

var rd *render.Render = render.New()

// 핸들러를 추가해서 close에 대한 책임을 준다
type AppHandler struct {
	http.Handler // 핸들러가  http.Handler 인터페이스를 포함
	db           model.DBHandler
}

// function 들을 Apphandler의 메서드로 변경
func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	// list := []*model.Todo{}
	// for _, v := range todoMap {
	// 	list = append(list, v)
	// }

	list := a.db.GetTodos()
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	todo := a.db.AddTodo(name)
	// id := len(todoMap) + 1
	// todo := &Todo{id, name, false, time.Now()}
	// todoMap[id] = todo
	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

// 쿠키에서 세션을 읽어오는 세션을 만든다. signin.go 에 작성한것과 같이 작성
func getSessionID(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	// 비어있는지 체크 -> 로그인을 안했다는 경우
	if val == nil {
		return ""
	}
	// 비어있지 않을 경우 string으로 변경해서 return
	return val.(string)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})

	}

	// if _, ok := todoMap[id]; ok {
	// 	delete(todoMap, id)
	// 	rd.JSON(w, http.StatusOK, Success{true})
	// } else {
	// 	rd.JSON(w, http.StatusOK, Success{false})
	// }
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"
	ok := a.db.CompleteTodo(id, complete)

	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})

	}

	// if todo, ok := todoMap[id]; ok {
	// 	todo.Completed = complete
	// 	rd.JSON(w, http.StatusOK, Success{true})
	// } else {
	// 	rd.JSON(w, http.StatusOK, Success{false})
	// }
}

// func addTestTodos() {
// 	todoMap[1] = &Todo{1, "Buy a milk", false, time.Now()}
// 	todoMap[2] = &Todo{2, "Exercise", true, time.Now()}
// 	todoMap[3] = &Todo{3, "Home work", false, time.Now()}
// }

// 새로운 인스턴스를 만들어서 밖에서 호출하도록 function 추가
func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 유저가 요청한 url이 signin.html 일경우 next()로 넘겨줘야한다. 그렇지 않을경우 무한루프
	if strings.Contains(r.URL.Path, "/signin.html") || strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	// 세션 ID가 쿠키에 있느냐로 signed in 여부 판단
	// if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	// if not user signed in -> redirect singin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)

}

func MakeHandler(filepath string) *AppHandler {
	// todoMap = make(map[int]*Todo)
	// todoMap[1] = &Todo{1, "num1", false, time.Now()}
	// todoMap[2] = &Todo{2, "num2", true, time.Now()}
	// todoMap[3] = &Todo{3, "num3", false, time.Now()}

	// addTestTodos()

	r := mux.NewRouter()

	// main.go 에 있던것을 가져와서 사용
	// 세션아이디를 체크하고 있으면 sign 없으면 login화면으로 넘긴다
	// 순서대로 확인한다. chain으로 되어있다
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(CheckSignin),
		negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	a := &AppHandler{
		Handler: n,
		// db 는 NewDBHandler()를 호출해서 결과값 저장

		db: model.NewDBHandler(filepath),
	}

	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	r.HandleFunc("/", a.indexHandler)

	// login page
	// 핸들러 생성, google에 로그인 요청
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	// 핸들러 생성,
	r.HandleFunc("/auth/google/callback", googleAuthCallback)

	return a
}
