package captcha

import (
	"sync"
)

var runningSolver = (&runningSolverStruct{}).new()
var queue = (&queueStruct{}).new()

type runningSolverStruct struct {
	mu sync.Mutex
	m  map[string][]string
}

func (r *runningSolverStruct) new() *runningSolverStruct {
	r.m = make(map[string][]string)
	return r
}

func (r *runningSolverStruct) len(site string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.m[site])
}

func (r *runningSolverStruct) add(site, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[site] = append(r.m[site], id)
}

func (r *runningSolverStruct) remove(site, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, v := range r.m[site] {
		if v == id {
			r.m[site] = append(r.m[site][:i], r.m[site][i+1:]...)
			return
		}
	}
}

type queueStruct struct {
	mu sync.Mutex
	m  map[string][]Captcha
}

func (q *queueStruct) new() *queueStruct {
	q.m = make(map[string][]Captcha)
	return q
}

func (q *queueStruct) get() map[string][]Captcha {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.m
}

func (q *queueStruct) getSite(site string) []Captcha {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.m[site]
}

func (q *queueStruct) set(site string, c []Captcha) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.m[site] = c
}
