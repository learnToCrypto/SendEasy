package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/learnToCrypto/lakoposlati/internal/logger"
	"github.com/learnToCrypto/lakoposlati/internal/user"
)

// GET /demand/new
// Show the new demand form page
func NewDemand(writer http.ResponseWriter, request *http.Request) {
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
		if _, err := userD.CreateDemand(object, collection, delivery, timeframe, status); err != nil {
			logger.Danger(err, "Cannot create demand")
		}
		http.Redirect(writer, request, "/demand/list", 302)
	}
}

// GET /demand/read
// Show the details of the thread, including the posts and the form to write a post
func ReadDemand(writer http.ResponseWriter, request *http.Request) {
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
// Create the demand
func PostDemand(writer http.ResponseWriter, request *http.Request) {
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
		fmt.Println(url)
		http.Redirect(writer, request, url, 302)
	}
}

func DemandList(writer http.ResponseWriter, request *http.Request) {
	type Data struct {
		Demand        user.Demand // Must be exported!
		NumOfMessages int         // Must be exported!
		NumofDemands  int
	}

	demands, err := user.Demands()
	if err != nil {
		error_message(writer, request, "Cannot get demands")
	} else {
		d := make([]Data, len(demands))
		for i, v := range demands {
			d[i] = Data{Demand: v, NumOfMessages: v.NumReplies(), NumofDemands: len(demands)}
		}
		x := len(demands)
		fmt.Println(request.URL.Path)
		i, err := strconv.Atoi(strings.TrimPrefix(request.URL.Path, "/demand/list/"))
		fmt.Println(i)
		start := ((i - 1) * 10)
		end := i * 10
		if end > x {
			end = (i-1)*10 + (x % 10)
		}
		if err != nil {
			error_message(writer, request, "Cannot get demands")
		} else {
			_, err := session(writer, request)
			if err != nil {
				generateHTML(writer, d[start:end], "layout", "public.navbar", "demandList")
			} else {
				generateHTML(writer, d[start:end], "layout", "private.navbar", "demandList")
			}
		}
	}
}

//    {{range $1}}
//	 <input type="radio" name={{.Name}} value={{.Value}} {{if .IsDisabled}} disabled=true {{end}} {{if .IsChecked}}checked{{end}}> {{.Text}}
//		{{end}}
