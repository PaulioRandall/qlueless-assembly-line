package ventures

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	h "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/uhttp"
	u "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/utils"
	v "github.com/PaulioRandall/go-qlueless-assembly-api/internal/pkg/ventures"
)

// InjectDummyVentures injects dummy Ventures so the API testing can performed.
// This function is expected to be removed once a database and formal test data
// has been crafted.
func InjectDummyVentures() {
	ventures.Add(v.NewVenture{
		Description: "White wizard",
		State:       "Not started",
		Extra:       "colour: white; power: 9000",
	})
	ventures.Add(v.NewVenture{
		Description: "Green lizard",
		State:       "In progress",
		OrderIDs:    "4,5,6,7,8",
	})
	ventures.Add(v.NewVenture{
		Description: "Pink gizzard",
		State:       "Finished",
		OrderIDs:    "1,2,3",
	})
	ventures.Add(v.NewVenture{
		Description: "Eddie Izzard",
		State:       "In Progress",
		OrderIDs:    "4,5,6",
	})
	ventures.Add(v.NewVenture{
		Description: "The Count of Tuscany",
		State:       "In Progress",
		OrderIDs:    "4,5,6",
	})
	ventures.Update(&v.ModVenture{
		IDs:   "3",
		Props: "is_alive",
		Values: v.Venture{
			IsAlive: false,
		},
	})
}

// writeSuccessReply writes a success response.
func writeSuccessReply(res *http.ResponseWriter, req *http.Request, code int, data interface{}, msg string) {
	h.AppendJSONHeader(res, "")
	(*res).WriteHeader(code)
	reply := h.PrepResponseData(req, data, msg)
	json.NewEncoder(*res).Encode(reply)
}

// findVenture finds the Venture with the specified ID.
func findVenture(id string, res *http.ResponseWriter, req *http.Request) (v.Venture, bool) {
	ven, ok := ventures.Get(id)
	if !ok {
		h.WriteBadRequest(res, req, fmt.Sprintf("Thing '%s' not found", id))
		return v.Venture{}, false
	}
	return ven, true
}

// decodeNewVenture decodes a NewVenture from a Request.Body.
func decodeNewVenture(res *http.ResponseWriter, req *http.Request) (v.NewVenture, bool) {
	ven, err := v.DecodeNewVenture(req.Body)
	if err != nil {
		h.WriteBadRequest(res, req, "Unable to decode request body into a Venture")
		return v.NewVenture{}, false
	}
	return ven, true
}

// validateNewVenture validates a NewVenture that has yet to be assigned an ID.
func validateNewVenture(ven v.NewVenture, res *http.ResponseWriter, req *http.Request) bool {
	errMsgs := ven.Validate()
	if len(errMsgs) != 0 {
		h.WriteBadRequest(res, req, strings.Join(errMsgs, " "))
		return false
	}
	return true
}

// decodeModVentures decodes modifications to Ventures from a Request.Body.
func decodeModVentures(res *http.ResponseWriter, req *http.Request) (*v.ModVenture, bool) {
	mv, err := v.DecodeModVenture(req.Body)
	if err != nil {
		h.WriteBadRequest(res, req,
			"Unable to decode request body into a Venture update")
		return nil, false
	}
	return &mv, true
}

// validateModVentures validates a Venture update.
func validateModVentures(mv *v.ModVenture, res *http.ResponseWriter, req *http.Request) bool {
	errMsgs := mv.Validate()
	if len(errMsgs) != 0 {
		h.WriteBadRequest(res, req, strings.Join(errMsgs, " "))
		return false
	}
	return true
}

// deleteVenture deletes a Venture from the data store.
func deleteVenture(id string, res *http.ResponseWriter, req *http.Request) (v.Venture, bool) {
	ven, ok := ventures.Delete(id)
	if !ok {
		h.WriteBadRequest(res, req, fmt.Sprintf("Thing '%s' not found", id))
		return v.Venture{}, false
	}
	return ven, true
}

// ventureIdCsvToSlice validates then parses a CSV string of IDs into a slice
func ventureIdCsvToSlice(idCsv string, res *http.ResponseWriter, req *http.Request) ([]string, bool) {
	idCsv = u.StripWhitespace(idCsv)

	if idCsv == "" {
		h.WriteBadRequest(res, req, "Query parameter 'ids' is missing or empty")
		return nil, false
	}

	if !u.IsPositiveIntCSV(idCsv) {
		h.WriteBadRequest(res, req, fmt.Sprintf("Could not parse query parameter"+
			" 'ids=%s' into a list of Venture IDs", idCsv))
		return nil, false
	}

	ids := strings.Split(idCsv, ",")
	return ids, true
}
