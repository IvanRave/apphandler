package apphandler

import (
	"net/http"
	"encoding/json"

	lgr	"github.com/Sirupsen/logrus"
)

// Error: 5xx (500)
// Warn: 4xx (400, 401, 422)
// Info: 2xx (200, 204)

// Return JSON in all cases, instead 204 (no data to return)

const tagRqst string = "byt.rqst"
const jsonMime string = "application/json; charset=utf-8"

// r.URL.String() 
//    String reassembles the URL into a valid URL string

func HandleNonAuth(r *http.Request,
	w http.ResponseWriter,
	errKey string,
	apiKey string){

	statusCode := 401

	errObj := clerr {
		ErrKey: errKey,
		//Details: nil
	}

	// no errors possible: one field only
	bstr, _ := errObj.ToJson()

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(statusCode)
	w.Write(bstr)
	
	lgr.WithFields(lgr.Fields {
		"tag": tagRqst,
		"status_code": statusCode,
		"err_key": errKey,
		"msg": "http401 apiKey: " + apiKey,
		"url": r.URL.String(),
	}).Warn() // client error = warn + rqst
}

// handleServerError writes error to client
// and sends a notif to admin
func HandleServerError(r *http.Request,
	w http.ResponseWriter,
	err error){

	errKey :=  "unexpectedError"
	statusCode := 500
	
	errObj := clerr {
		ErrKey: errKey,
		//Details: nil
	}

	// no errors possible: one field only
	bstr, _ := errObj.ToJson()

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(statusCode)
	w.Write(bstr)

	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"status_code": statusCode,
		"err_key": errKey,
		"msg": err.Error(),
		"url": r.URL.String(),
	}).Error() // server error + rqst
	// TODO: #33! Send an error to admin
}

// Send 400 or 422 response (or something, different
//   than right response)
func HandleClientError(r *http.Request,
	w http.ResponseWriter,
	err *clerr){
	
	str, parseErr := err.ToJson()
	
	if parseErr != nil {
		HandleServerError(r, w, parseErr)	
		return
	}

	statusCode := 422

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(statusCode)
	w.Write(str)
	
	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"status_code": statusCode,
		"err_key": err.ErrKey,
		"msg": string(str),
		"url": r.URL.String(),
	}).Warn()	// client error
}

func HandleNotFound(r *http.Request,
	w http.ResponseWriter) {

	statusCode := 404

	w.WriteHeader(statusCode) // no content type

	// Logging (after execution)
	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"status_code": statusCode,
		"err_key": "NotFoundError",
		"msg": "notFound",
		"url": r.URL.String(),
		//"uid": uid,
		//"params": strParams,
		//"perms": perms,
	}).Warn()  // rqst + info
}

func Handle204(r *http.Request,
	w http.ResponseWriter) {

	statusCode := 204
	
	//w.WriteHeader(http.StatusNoContent)
	//
	// just header: no content
	w.WriteHeader(statusCode) // no content type
		
	// Logging (after execution)
	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"status_code": statusCode,
		"err_key": "",
		"msg": "success",
		"url": r.URL.String(),
		//"uid": uid,
		//"params": strParams,
		//"perms": perms,
	}).Info()  // rqst + info
}

func HandleSuccess(r *http.Request,
	w http.ResponseWriter,
	rdata interface{}){

	if rdata == nil {
		Handle204(r, w)
		return
    }
	
	responseJson, errJson := json.Marshal(rdata)
	
	if errJson != nil {
		HandleServerError(r, w, errJson)
		return
	}
    
    // Write writes the data to the connection
	//    as part of an HTTP reply.
    // If WriteHeader has not yet been called,
	//    Write calls WriteHeader(http.StatusOK)
    // before writing the data.  If the Header does not contain a
    // Content-Type line, Write adds a Content-Type
	//    set to the result of passing
    // the initial 512 bytes of written data to DetectContentType
	w.Header().Set("Content-Type", jsonMime)
    w.Write([]byte(string(responseJson)))

	// Logging (after execution)
	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"status_code": 200,
		"err_key": "",
		"msg": "success",
		"url": r.URL.String(),
		//"uid": uid,
		//"params": strParams,
		//"perms": perms,
	}).Info()  // rqst + info
}
