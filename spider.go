package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/300brand/logger"
	"github.com/300brand/spider/rule"
	"github.com/gorilla/mux"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Spider struct {
}

type ruleSet struct {
	Host       string
	Rule       rule.Rule
	RuleStr    string
	LastUpdate time.Time
}

var SpiderConf = struct {
	DSN *string
}{
	DSN: flag.String("spider.dsn", "root:@tcp(localhost:49158)/spider", "MySQL DSN"),
}

var _ Handler = new(Spider)

func init() {
	Register("spider", new(Spider))
}

func (s *Spider) Router(r *mux.Router) {
	r.HandleFunc("/allrules", s.HandleAllRules)
	r.HandleFunc("/validate", s.HandleValidate)
}

func (s *Spider) HandleAllRules(w http.ResponseWriter, r *http.Request) {
	response := Response{Success: true}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if response.Response, response.Error = s.dbGetRules(); response.Error != nil {
		response.Success = false
	}
	if err := enc.Encode(response); err != nil {
		logger.Error.Printf("HandleAllRules: %s", err)
	}
}

func (s *Spider) HandleValidate(w http.ResponseWriter, r *http.Request) {
	response := Response{}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		logger.Error.Printf("HandleValidate: %s", err)
	}
}

func (s *Spider) db() (db *sql.DB, err error) {
	return sql.Open("mysql", *SpiderConf.DSN+"?parseTime=true")
}

func (s *Spider) dbGetRules() (rules []ruleSet, err error) {
	db, err := s.db()
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT host, json, updated FROM rules ORDER BY host")
	if err != nil {
		return
	}

	rules = make([]ruleSet, 0, 16)
	for rows.Next() {
		var set ruleSet
		var data []byte
		if err = rows.Scan(&set.Host, &data, &set.LastUpdate); err != nil {
			return
		}
		if err = json.Unmarshal(data, &set.Rule); err != nil {
			return
		}
		var buf bytes.Buffer
		if err = json.Indent(&buf, data, "", "\t"); err != nil {
			return
		}
		set.RuleStr = buf.String()
		rules = append(rules, set)
	}
	return
}

func (s *Spider) dbAddRule() {}
