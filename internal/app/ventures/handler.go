package ventures

import (
	"fmt"
	"net/http"

	. "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg"
	v "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
)

const (
	httpMethods = "GET, POST, PUT, DELETE, HEAD, OPTIONS"
)

var ventures = v.NewVentureStore()

// VenturesHandler handles requests to do with collections of, or individual,
// Ventures.
func VenturesHandler(res http.ResponseWriter, req *http.Request) {
	LogRequest(req)

	id := req.FormValue("id")
	switch {
	case req.Method == "GET" && id == "":
		get_AllVentures(&res, req)
	case req.Method == "GET":
		get_OneVenture(id, &res, req)
	//case req.Method == "POST":
	//post_NewThing(&res, req)
	//case req.Method == "PUT":
	//put_OneThing(&res, req)
	case req.Method == "HEAD":
		fallthrough
	case req.Method == "OPTIONS":
		AppendCORSHeaders(&res, httpMethods)
		WriteEmptyJSONReply(&res, "")
	default:
		MethodNotAllowed(&res, req)
	}
}

// get_AllVentures handles client requests for all living Ventures.
func get_AllVentures(res *http.ResponseWriter, req *http.Request) {
	vens := ventures.GetAllAlive()
	m := fmt.Sprintf("Found %d Ventures", len(vens))
	data := PrepResponseData(req, vens, m)

	AppendCORSHeaders(res, httpMethods)
	AppendJSONHeader(res, "")
	WriteJSONReply(res, req, data, "")
}

// get_OneVenture handles client requests for a specific Venture.
func get_OneVenture(id string, res *http.ResponseWriter, req *http.Request) {
	ven, ok := findVenture(id, res, req)
	if !ok {
		return
	}

	m := fmt.Sprintf("Found Venture '%s'", id)
	data := PrepResponseData(req, ven, m)

	AppendCORSHeaders(res, httpMethods)
	AppendJSONHeader(res, "")
	WriteJSONReply(res, req, data, "")
}

// findVenture finds the Venture with the specified ID
func findVenture(id string, res *http.ResponseWriter, req *http.Request) (v.Venture, bool) {
	ven, ok := ventures.Get(id)
	if !ok || !ven.IsAlive {
		r := WrappedReply{
			Message: fmt.Sprintf("Thing '%s' not found", id),
		}
		AppendCORSHeaders(res, httpMethods)
		Write4XXReply(res, req, 404, r)
		return v.Venture{}, false
	}
	return ven, true
}

// InjectDummyVentures injects dummy Ventures so the API testing can performed.
// This function is expected to be removed once a database and formal test data
// has been crafted
func InjectDummyVentures() {
	ventures.Add(v.Venture{
		Description: "White wizard",
		State:       "Not started",
		Extra:       "colour: white; power: 9000",
		IsAlive:     true,
	})
	ventures.Add(v.Venture{
		Description: "Green lizard",
		State:       "In progress",
		OrderIDs:    "4,5,6,7,8",
		IsAlive:     true,
	})
	ventures.Add(v.Venture{
		Description: "Pink gizzard",
		State:       "Finished",
		OrderIDs:    "1,2,3",
		IsAlive:     false,
	})
}
