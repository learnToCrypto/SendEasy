package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/demands"
	"github.com/learnToCrypto/lakoposlati/internal/logger"
	"github.com/learnToCrypto/lakoposlati/internal/messages"
	"github.com/learnToCrypto/lakoposlati/internal/messages/msglist"
	"github.com/learnToCrypto/lakoposlati/internal/pagination"
	"github.com/learnToCrypto/lakoposlati/internal/sessions"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

// GET /demand/new
// Show the new demand form page
func NewDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("NewDemand")
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {

		username, err := user.UsernamebySession(sess.UserId)
		if err != nil {
			username = "My Account"
		}
		d := struct {
			Username string
		}{
			Username: username,
		}
		generateHTML(writer, d, "layout/base", "private/navbar", "new.demand")
	}
}

// POST / demand/create/1
// Create the demand (insert in database) / demand.go
func CreateDemand1(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("CreateDemand")

	err := request.ParseForm()
	if err != nil {
		logger.Danger(err, "Cannot parse form")
	}

	object := request.PostFormValue("obj")
	collection := request.PostFormValue("src")
	delivery := request.PostFormValue("dest")
	//fmt.Println("object: ", object)
	//fmt.Println("collection: ", collection)
	//fmt.Println("delivery: ", delivery)

	sess, err := session(writer, request)
	//fmt.Println("session in demand/create/1: ", sess)
	//	fmt.Println("err in demand/create/1: ", err)
	if err != nil {
		fmt.Println("public")
		d := struct {
			Object      string
			Location    string
			Destination string
		}{
			Object:      object,
			Location:    collection,
			Destination: delivery,
		}
		generateHTML(writer, d, "layout/base", "public/navbar", "public/new.demand")
	} else {
		fmt.Println("private")
		username, err := user.UsernamebySession(sess.UserId)
		if err != nil {
			username = "My Account"
		}
		d := struct {
			Username    string
			Object      string
			Location    string
			Destination string
		}{
			Username:    username,
			Object:      object,
			Location:    collection,
			Destination: delivery,
		}

		generateHTML(writer, d, "layout/base", "private/navbar", "private/new.demand")
	}
}

// POST / demand/create/2
func CreateDemandPriv(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("CreateDemand")
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	}

	userD, err := sess.User()
	if err != nil {
		logger.Danger(err, "Cannot get user from session")
	}

	err = request.ParseForm()
	if err != nil {
		logger.Danger(err, "Cannot parse form")
	}

	object := request.PostFormValue("obj")
	collection := request.PostFormValue("src")
	delivery := request.PostFormValue("dest")
	timeframe := request.PostFormValue("timeframe")
	status := 0

	//fmt.Println("object2: ", object)
	//fmt.Println("collection2: ", collection)
	//fmt.Println("delivery2: ", delivery)
	//fmt.Println("timeframe: ", timeframe)

	if object == "" || collection == "" || delivery == "" || timeframe == "" {
		error_message(writer, request, "Cannot create demand - not all fields filled")
	} else if _, err := demands.CreateDemand(&userD, object, collection, delivery, timeframe, status); err != nil {
		logger.Danger(err, "Cannot create demand")
	}
	http.Redirect(writer, request, "/demand/list/1", 302) // change /1 with proper pagination, last demands first
}

// POST / demand/create/3
func CreateDemandPub(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("CreateDemand")
	err := request.ParseForm()
	if err != nil {
		logger.Danger(err, "Cannot parse form")
	}
	userI := user.User{
		FirstName: request.PostFormValue("first_name"),
		LastName:  request.PostFormValue("last_name"),
		Name:      request.PostFormValue("first_name") + " " + request.PostFormValue("last_name"),
		Email:     request.PostFormValue("email"),
		Password:  request.PostFormValue("password"),
	}
	if err := userI.Create(); err != nil {
		logger.Danger(err, "Cannot create user")
	}

	//Here session should be made
	fmt.Println("cookie creating in demand/create/3")
	sess, err := sessions.CreateSession(&userI)
	if err != nil {
		logger.Danger(err, "Cannot create session")
	}

	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    sess.Uuid,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600 * 8,
	}
	http.SetCookie(writer, &cookie)

	object := request.PostFormValue("obj")
	collection := request.PostFormValue("src")
	delivery := request.PostFormValue("dest")
	timeframe := request.PostFormValue("timeframe")
	status := 0

	//fmt.Println("object2: ", object)
	//fmt.Println("collection2: ", collection)
	//fmt.Println("delivery2: ", delivery)
	//fmt.Println("timeframe: ", timeframe)

	if object == "" || collection == "" || delivery == "" || timeframe == "" {
		error_message(writer, request, "Cannot create demand - not all fields filled")
	} else if _, err := demands.CreateDemand(&userI, object, collection, delivery, timeframe, status); err != nil {
		logger.Danger(err, "Cannot create demand")
	}
	//sess1, err1 := session(writer, request)
	//fmt.Println("session in public route: ", sess1)
	//fmt.Println("err in public route: ", err1)

	http.Redirect(writer, request, "/demand/list/1", 302) // change /1 with proper pagination, last demands first
}

// GET /demand/read
// Show the details of the thread, including the posts and the form to write a post
func ReadDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("ReadDemand")
	vals := request.URL.Query()
	uuid := vals.Get("id")
	demand, err := demands.DemandByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot read demand")
	} else {
		sess, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &demand, "layout/base", "public/navbar", "public/demand")
		} else {
			username, err := user.UsernamebySession(sess.UserId)
			if err != nil {
				username = "My Account"
			}
			sm, err := msglist.Messages(&demand)
			if err != nil {
				error_message(writer, request, "Cannot read messages")
			}

			d := struct {
				Demand   demands.Demand
				Username string
				Messages []messages.Message
			}{
				Demand:   demand,
				Username: username,
				Messages: sm,
			}
			generateHTML(writer, d, "layout/base", "private/navbar", "private/demand")
		}
	}
}

// POST /demand/post
// Create the demand message
//ToDO Refact name of this func to PostMessage
func PostDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("PostDemand")
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			logger.Danger(err, "Cannot parse form")
		}
		userD, err := sess.User()
		if err != nil {
			logger.Danger(err, "Cannot get user from session")
		}
		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		demand, err := demands.DemandByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read demand")
		}
		if _, err := messages.CreateMessage(&userD, demand, body); err != nil {
			logger.Danger(err, "Cannot create message")
		}
		url := fmt.Sprint("/demand/read?id=", uuid)
		//fmt.Println(url)
		http.Redirect(writer, request, url, 302)
	}
}

//DemandList function provides a list of all demands together with pagination
func DemandList(writer http.ResponseWriter, request *http.Request) {

	var limit int = 8

	//parse request and obtain offset
	u, err := strconv.Atoi(strings.TrimPrefix(request.URL.Path, "/demand/list/"))
	if err != nil {
		error_message(writer, request, "Cannot parse request")
		return
	}
	//offset calculation
	offset := (u - 1) * limit

	// retrieve demands; limit specifies size of a list
	dmnds, err := demands.Demands(limit, offset)
	if err != nil {
		error_message(writer, request, "Cannot get demands")
		return
	}

	// number of demands in database
	x, err := demands.DemandsNum()
	if err != nil {
		error_message(writer, request, "Cannot retrieve a list of demands")
		return
	}

	// m contains number of offers
	m := make(map[string]int)
	for _, v := range dmnds {
		m[v.Uuid] = v.NumReplies()
	}

	//paginator : total number of pages calculated using PageNum
	p := pagination.NewPaginator(limit, u, pagination.PageNum(x, limit))
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

	sess, err := session(writer, request)
	//get username of active user
	username, err := user.UsernamebySession(sess.UserId)
	if err != nil {
		username = "My Account"
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

	if err != nil {
		generateHTMLwithFunc(writer, d, funcMap, "layout/base", "public/navbar", "demandList")
	} else {
		generateHTMLwithFunc(writer, d, funcMap, "layout/base", "private/navbar", "demandList")
	}

}
