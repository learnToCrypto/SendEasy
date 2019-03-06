package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/learnToCrypto/lakoposlati/internal/logger"
	"github.com/learnToCrypto/lakoposlati/internal/pagination"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

// GET /demand/new
// Show the new demand form page
func NewDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("NewDemand")
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "new.demand")
	}
}

// POST /
// Create the demand (insert in database) / demand.go
func CreateDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("CreateDemand")
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
		object := request.PostFormValue("object")
		collection := request.PostFormValue("collection")
		delivery := request.PostFormValue("delivery")
		timeframe := request.PostFormValue("timeframe")
		status := 0
		if object == "" || collection == "" || delivery == "" || timeframe == "" {
			error_message(writer, request, "Cannot create demand - not all fields filled")
		} else if _, err := userD.CreateDemand(object, collection, delivery, timeframe, status); err != nil {
			logger.Danger(err, "Cannot create demand")
		}
		http.Redirect(writer, request, "/demand/list/1", 302) // change /1 with proper pagination, last demands first
	}
}

// GET /demand/read
// Show the details of the thread, including the posts and the form to write a post
func ReadDemand(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("ReadDemand")
	vals := request.URL.Query()
	uuid := vals.Get("id")
	demand, err := user.DemandByUUID(uuid)
	if err != nil {
		error_message(writer, request, "Cannot read demand")
	} else {
		_, err := session(writer, request)
		if err != nil {
			generateHTML(writer, &demand, "layout", "public.navbar", "public.demand")
		} else {
			generateHTML(writer, &demand, "layout", "private.navbar", "private.demand")
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
		demand, err := user.DemandByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read demand")
		}
		if _, err := userD.CreateMessage(demand, body); err != nil {
			logger.Danger(err, "Cannot create message")
		}
		url := fmt.Sprint("/demand/read?id=", uuid)
		//fmt.Println(url)
		http.Redirect(writer, request, url, 302)
	}
}

func DemandList(writer http.ResponseWriter, request *http.Request) {

	type Data struct {
		Demands   []user.Demand  // Must be exported!
		MsgNum    map[string]int // Must be exported!
		Paginator pagination.Paginator
	}

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
	demands, err := user.Demands(limit, offset)
	if err != nil {
		error_message(writer, request, "Cannot get demands")
		return
	}

	// number of demands in database
	x, err := user.DemandsNum()
	if err != nil {
		error_message(writer, request, "Cannot retrieve a list of demands")
		return
	}

	// m contains number of offers
	m := make(map[string]int)
	for _, v := range demands {
		m[v.Uuid] = v.NumReplies()
	}

	//paginator : total number of pages calculated using PageNum
	p := pagination.NewPaginator(limit, u, pagination.PageNum(x, limit))
	// generate a slice of int containing page links
	p.GeneratePageLink()

	d := Data{
		Demands:   demands,
		MsgNum:    m,
		Paginator: *p,
	}

	funcMap := template.FuncMap{
		"FormatTime": func(t time.Time) string {
			tn := time.Now().Local()
			th := t.Local()
			diff := tn.Sub(th)
			return humanizeDuration(diff)
		},
	}

	_, err = session(writer, request)
	if err != nil {
		generateHTMLwithFunc(writer, d, funcMap, "layout", "public.navbar", "demandList")
	} else {
		generateHTMLwithFunc(writer, d, funcMap, "layout", "private.navbar", "demandList")
	}

}
