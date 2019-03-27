package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/demands"
	"github.com/learnToCrypto/lakoposlati/internal/messages"
	"github.com/learnToCrypto/lakoposlati/internal/messages/msglist"
	"github.com/learnToCrypto/lakoposlati/internal/pagination"
)

func MyShipments(writer http.ResponseWriter, request *http.Request) {

	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {

		userI, err := sess.User()
		if err != nil {
			error_message(writer, request, "Cannot find user")
		}
		username := userI.Name
		var limit int = 8

		x, err := demands.DemandsNumUser(strconv.Itoa(userI.Id))
		//fmt.Println("number of demands from ", username, "is ", x)
		if err != nil {
			error_message(writer, request, "Cannot get demands")
		}

		dmnds, err := demands.DemandsByUserId(strconv.Itoa(userI.Id), limit, 0)
		if err != nil {
			error_message(writer, request, "Cannot get demands")
		}

		// m contains number of offers
		m := make(map[string]int)
		for _, v := range dmnds {
			m[v.Uuid] = v.NumReplies()
		}

		//paginator : total number of pages calculated using PageNum
		p := pagination.NewPaginator(limit, 1, pagination.PageNum(x, limit))
		// generate a slice of int containing page links
		p.GeneratePageLink()

		// funcMap defines a function called FormatTime that trans
		funcMap := template.FuncMap{
			"FormatTime": func(t time.Time) string {
				tn := time.Now().Local()
				th := t.Local()
				diff := tn.Sub(th)
				return humanizeDuration(diff)
			},
		}

		d := struct {
			Demands   []demands.Demand // Must be exported! is it?
			MsgNum    map[string]int
			Paginator pagination.Paginator
			Username  string
		}{
			Demands:   dmnds,
			MsgNum:    m,
			Paginator: *p,
			Username:  username,
		}

		generateHTMLwithFunc(writer, d, funcMap, "layout/base", "private/navbar", "demandList")
	}
}

//Inbox shows messages published by user
func Inbox(writer http.ResponseWriter, request *http.Request) {

	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {

		userI, err := sess.User()
		if err != nil {
			error_message(writer, request, "Cannot find user")
		}
		username := userI.Name
		var limit int = 8

		sm, err := msglist.Inbox(&userI)
		if err != nil {
			error_message(writer, request, "Cannot read messages")
		}

		x := len(sm)
		fmt.Println("messages: ", sm)
		fmt.Println("number: ", x)
		// m contains number of offers
		//paginator : total number of pages calculated using PageNum
		p := pagination.NewPaginator(limit, 1, pagination.PageNum(x, limit))
		// generate a slice of int containing page links
		p.GeneratePageLink()

		d := struct {
			Username  string
			Paginator pagination.Paginator
			Messages  []messages.Message
		}{
			Username:  username,
			Paginator: *p,
			Messages:  sm,
		}

		// funcMap defines a function called FormatTime that trans
		funcMap := template.FuncMap{
			"FormatTime": func(t time.Time) string {
				tn := time.Now().Local()
				th := t.Local()
				diff := tn.Sub(th)
				return humanizeDuration(diff)
			},
		}

		generateHTMLwithFunc(writer, d, funcMap, "layout/base", "private/navbar", "private/inbox")
	}
}
