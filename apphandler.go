// apphandler contains handlers for API results
package apphandler

import (
	"net/http"
	"encoding/json"
	//lgr	"github.com/Sirupsen/logrus"
)

// prev
// Error: 5xx (500)
// Warn: 4xx (400, 401, 422)
// Info: 2xx (200, 204)

// Return JSON in all cases, instead 204 (no data to return)

const jsonMime string = "application/json; charset=utf-8"

// r.URL.String() 
//    String reassembles the URL into a valid URL string

func HandleNonAuth(r *http.Request,
	w http.ResponseWriter,
	errKey string,
	apiKey string) (int16,
		string, // url
		string, // err_key
		string){

	var statusCode int16 = 401

	errObj := clerr {
		ErrKey: errKey,
		//Details: nil
	}

	// no errors possible: one field only
	bstr, _ := errObj.ToJson()

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(int(statusCode))
	w.Write(bstr)

	//	return map[string]interface{}{

	return statusCode,
	r.URL.String(),
	errKey,	
	"http401 apiKey: " + apiKey
	// lgr.WithFields(lgr.Fields {
	// 	"tag": tagRqst,
	// 	"status_code": statusCode,
	// 	"err_key": errKey,
	// 	"msg": "http401 apiKey: " + apiKey,
	// 	"url": r.URL.String(),
	// }).Warn() // client error = warn + rqst
}

// handleServerError writes error to client
// and sends a notif to admin
func HandleServerError(r *http.Request,
	w http.ResponseWriter,
	err error)(int16, string, string, string){

	errKey :=  "unexpectedError"
	var statusCode int16 = 500
	
	errObj := clerr {
		ErrKey: errKey,
		//Details: nil
	}

	// no errors possible: one field only
	bstr, _ := errObj.ToJson()

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(int(statusCode))
	w.Write(bstr)

	return statusCode,
	r.URL.String(),
	errKey,
	err.Error()
	
	// lgr.WithFields(lgr.Fields{
	// 	"tag": tagRqst,
	// 	"status_code": statusCode,
	// 	"err_key": errKey,
	// 	"msg": err.Error(),
	// 	"url": r.URL.String(),
	// }).Error() // server error + rqst
	// TODO: #33! Send an error to admin
}

// Send 400 or 422 response (or something, different
//   than right response)
func HandleClientError(r *http.Request,
	w http.ResponseWriter,
	err *clerr)(int16, string, string, string){
	
	str, parseErr := err.ToJson()
	
	if parseErr != nil {
		return HandleServerError(r, w, parseErr)
	}

	var statusCode int16 = 422

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(int(statusCode))
	w.Write(str)
	
	// lgr.WithFields(lgr.Fields{
	// 	"tag": tagRqst,
	// 	"status_code": statusCode,
	// 	"err_key": err.ErrKey,
	// 	"msg": string(str),
	// 	"url": r.URL.String(),
	// }).Warn()	// client error

	return statusCode,
	r.URL.String(),
	err.ErrKey,	
	string(str)
}

func HandleNotFound(r *http.Request,
	w http.ResponseWriter)(int16, string, string, string) {

	var statusCode int16 = 404

	w.WriteHeader(int(statusCode)) // no content type

	return statusCode,
	r.URL.String(),
	"NotFoundError",	
	"not found"
	
	// Logging (after execution)
	// lgr.WithFields(lgr.Fields{
	// 	"tag": tagRqst,
	// 	"status_code": statusCode,
	// 	"err_key": "NotFoundError",
	// 	"msg": "notFound",
	// 	"url": r.URL.String(),
	// 	//"uid": uid,
	// 	//"params": strParams,
	// 	//"perms": perms,
	// }).Warn()  // rqst + info
}

func Handle204(r *http.Request,
	w http.ResponseWriter)(int16, string, string, string) {

	var statusCode int16 = 204
	
	//w.WriteHeader(http.StatusNoContent)
	//
	// just header: no content
	w.WriteHeader(int(statusCode)) // no content type

	return statusCode,
	r.URL.String(),
	"",	
	"success"
	
	// Logging (after execution)
	// lgr.WithFields(lgr.Fields{
	// 	"tag": tagRqst,
	// 	"status_code": statusCode,
	// 	"err_key": "",
	// 	"msg": "success",
	// 	"url": r.URL.String(),
	// 	//"uid": uid,
	// 	//"params": strParams,
	// 	//"perms": perms,
	// }).Info()  // rqst + info
}

func HandleSuccess(r *http.Request,
	w http.ResponseWriter,
	rdata interface{})(int16, string, string, string){

	if rdata == nil {
		return Handle204(r, w)
    }
	
	responseJson, errJson := json.Marshal(rdata)
	
	if errJson != nil {
		return HandleServerError(r, w, errJson)
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

	
	return 200,
	r.URL.String(),
	"",	
	"success"
	// Logging (after execution)
	// lgr.WithFields(lgr.Fields{
	// 	"tag": tagRqst,
	// 	"status_code": 200,
	// 	"err_key": "",
	// 	"msg": "success",
	// 	"url": r.URL.String(),
	// 	//"uid": uid,
	// 	//"params": strParams,
	// 	//"perms": perms,
	// }).Info()  // rqst + info
}

