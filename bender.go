package main

import (
	"fmt"
	// "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	// "github.com/justinas/alice"
	"log"
	"net/http"
	"sync"
	"time"
)

type JobCenter struct {
	MUL    sync.Mutex
	JobMap map[string]chan string
}

var JC = JobCenter{
	JobMap: make(map[string]chan string, 5),
}

func (j JobCenter) AddJob(uid string, ch chan string) {
	j.MUL.Lock()
	time.Sleep(5 * time.Second)
	defer j.MUL.Unlock()
	j.JobMap[uid] = ch
}

func (j JobCenter) GetJob(uid string) chan string {
	j.MUL.Lock()
	defer j.MUL.Unlock()
	return j.JobMap[uid]
}

type TestCase struct {
	Name       string
	Passed     bool
	Skipped    bool
	Timeout    int
	IgnoreFail bool
	Log        string
}

type TestRun struct {
	UID       string
	Testcases []TestCase
	Passed    bool
	StartTime time.Time
}

// func fakeTestRun() []TestCase {
// 	t1 := TestCase{"t1", false, false, 300, false, ""}
// 	t2 := TestCase{"t2", false, false, 300, false, ""}
// 	t3 := TestCase{"t3", false, false, 300, false, ""}

// 	r := TestRun{t1, t2, t3}

// 	return r
// }

func lanuchTestRun(data chan string) {
	for d := range data {
		if d != "over" {
			log.Println(d)
		} else {
			log.Println("i am out")
		}
	}
	log.Println("run is over")
}

func SingleTestRun(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, ps.ByName("res"))
}

func LanuchTestPlan(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	JC.AddJob(ps.ByName("id"), make(chan string))
	go lanuchTestRun(JC.GetJob(ps.ByName("id")))
	fmt.Fprint(w, JC.JobMap)
}

func SingleTestCaseResult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ch := JC.GetJob(ps.ByName("id"))
	res := ps.ByName("res")
	if res == "quit" {
		close(ch)
	} else {
		ch <- res
	}
}

func main() {

	router := httprouter.New()
	router.GET("/id/:res", SingleTestRun)
	router.GET("/lanuch/:id", LanuchTestPlan)
	router.GET("/testcase/:id/:res", SingleTestCaseResult)

	log.Fatal(http.ListenAndServe(":8082", router))
}
