package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

var _ = utf8.MaxRune

func main() {
	h := newAppHandler()

	http.HandleFunc("/startProfiling", h.makeHandler("startProfiling", h.handleStartProfiling))
	http.HandleFunc("/finishProfiling", h.makeHandler("finishProfiling", h.handleFinishProfiling))
	http.HandleFunc("/stats", h.makeHandler("stats", h.handleStats))
	http.HandleFunc("/formatDate", h.makeHandler("formatDate", h.handleFormatDate))
	http.HandleFunc("/jsonInspect", h.makeHandler("jsonInspect", h.handleJsonInspect))
	http.HandleFunc("/validateIdent", h.makeHandler("validateIdent", h.handleValidateIdent))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newAppHandler() *appHandler {
	return &appHandler{
		numRequests: make(map[string]uint64),
	}
}

type appHandler struct {
	mu              sync.RWMutex
	numRequests     map[string]uint64
	runningProfiler bool
}

type requestStat struct {
	name  string
	count uint64
}

func (h *appHandler) requestStart(name string) {
	log.Printf("start %s request", name)
	h.mu.Lock()
	defer h.mu.Unlock()
	h.numRequests[name]++
}

func (h *appHandler) collectRequestStats() []requestStat {
	items := make([]requestStat, 0, 32)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, v := range h.numRequests {
		items = append(items, requestStat{name: k, count: v})
	}
	return items
}

func (h *appHandler) makeHandler(name string, handle func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h.requestStart(name)
		err := handle(w, req)
		if err != nil {
			log.Printf("%s error: %v", name, err)
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
		log.Printf("finished %s request", name)
	}
}

func (h *appHandler) handleStartProfiling(w http.ResponseWriter, req *http.Request) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.runningProfiler {
		return errors.New("already running a profiler")
	}
	f, err := os.Create("cpu.out")
	if err != nil {
		return err
	}
	pprof.StartCPUProfile(f)
	h.runningProfiler = true

	w.Write([]byte("OK\n"))

	return nil
}

func (h *appHandler) handleFinishProfiling(w http.ResponseWriter, req *http.Request) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.runningProfiler = false

	pprof.StopCPUProfile()

	w.Write([]byte("OK\n"))

	return nil
}

func (h *appHandler) handleStats(w http.ResponseWriter, req *http.Request) error {
	stats := h.collectRequestStats()
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].name < stats[j].name
	})
	for _, kv := range stats {
		w.Write([]byte(fmt.Sprintf("%s: %d\n", kv.name, kv.count)))
	}
	return nil
}

func (h *appHandler) handleFormatDate(w http.ResponseWriter, req *http.Request) error {
	arg := req.URL.Query().Get("arg")
	isNumber, err := regexp.MatchString(`\d+`, arg)
	if err != nil {
		panic(err)
	}
	if isNumber {
		unixTime, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return err
		}
		t := time.Unix(unixTime, 0)
		s := strings.Join([]string{t.String(), "\n"}, "")
		_, err = io.WriteString(w, s)
		return err
	}
	return errors.New("missing or invalid ?arg parameter")
}

func (h *appHandler) handleJsonInspect(w http.ResponseWriter, req *http.Request) error {
	arg := req.URL.Query().Get("arg")
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(arg), &obj); err != nil {
		return err
	}
	for k, v := range obj {
		rv := reflect.ValueOf(v)
		fmt.Fprintf(w, "%s: %s\n", k, reflect.TypeOf(rv.Interface()).String())
	}
	if len(obj) == 0 {
		w.Write([]byte("empty object\n"))
	}
	return nil
}

func (h *appHandler) handleValidateIdent(w http.ResponseWriter, req *http.Request) error {
	arg := req.URL.Query().Get("arg")

	if len(arg) == 0 {
		return errors.New("empty arg")
	}
	firstChar := []rune(arg)[0]
	if !unicode.IsLetter(firstChar) {
		return errors.New("first char should be a letter")
	}
	for i, ch := range arg {
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			return fmt.Errorf("invalid char at offset %d", i)
		}
	}

	w.Write([]byte("OK\n"))

	return nil
}
