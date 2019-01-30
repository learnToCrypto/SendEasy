package main

import (
	"fmt"
	"github.com/learnToCrypto/lakoposlati/data"
	"net/http"
)

// GET /demand/new
// Show the new demand form page
func newDemand(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "new.demand")
	}
}

// POST /
// Create the demand (insert in database) / demand.go
func createDemand(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		object := request.PostFormValue("object")
		collection := request.PostFormValue("collection")
		delivery := request.PostFormValue("delivery")
		timeframe := request.PostFormValue("timeframe")
		status := 0
		if _, err := user.CreateDemand(object, collection, delivery, timeframe, status); err != nil {
			danger(err, "Cannot create demand")
		}
		http.Redirect(writer, request, "/demand/list", 302)
	}
}

// GET /demand/read
// Show the details of the thread, including the posts and the form to write a post
func readDemand(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	uuid := vals.Get("id")
	demand, err := data.DemandByUUID(uuid)
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
func postDemand(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/login", 302)
	} else {
		err = request.ParseForm()
		if err != nil {
			danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			danger(err, "Cannot get user from session")
		}
		body := request.PostFormValue("body")
		uuid := request.PostFormValue("uuid")
		demand, err := data.DemandByUUID(uuid)
		if err != nil {
			error_message(writer, request, "Cannot read demand")
		}
		if _, err := user.CreateMessage(demand, body); err != nil {
			danger(err, "Cannot create message")
		}
		url := fmt.Sprint("/demand/read?id=", uuid)
		fmt.Println(url)
		http.Redirect(writer, request, url, 302)
	}
}


func demandList(writer http.ResponseWriter, request *http.Request) {
	type Data struct {
		Demand        data.Demand // Must be exported!
		NumOfMessages int         // Must be exported!
	}
	demands, err := data.Demands()
	fmt.Println("Len(demands) = ", len(demands))
	if err != nil {
		error_message(writer, request, "Cannot get demands")
	} else {
		d := make([]Data, len(demands))
		for i, v := range demands {
			d[i] = Data{Demand: v, NumOfMessages: v.NumReplies()}
		}
		_, err := session(writer, request)
   //pagination
	//	postsPerPage := 20
		if err != nil {
			generateHTML(writer, d, "layout", "public.navbar", "demandList")

		} else {
			generateHTML(writer, d, "layout", "private.navbar", "demandList")
		}
	}
}
