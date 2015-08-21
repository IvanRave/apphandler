package apphandler

import (
	"net/http"
	"encoding/json"

	lgr	"github.com/Sirupsen/logrus"
)

const tagRqst string = "byt.rqst"
const jsonMime string = "application/json; charset=utf-8"

// r.URL.String() 
//    String reassembles the URL into a valid URL string

func handleNonAuth(r *http.Request,
	w http.ResponseWriter,
	errKey string,
	apiKey string){
	
	w.WriteHeader(401)
	w.Write([]byte(errKey))
	
	lgr.WithFields(lgr.Fields {
		"tag": tagRqst,
		"msg": errKey,
		"url": r.URL.String(),
		// "err_key": str,
		// "api_key": apiKey,
	}).Warn() // client error = warn + rqst
}

// handleServerError writes error to client
// and sends a notif to admin
func handleServerError(r *http.Request,
	w http.ResponseWriter,
	err error){
	
	clerr := clerr {
		ErrKey: "unexpectedError",
		//Details: nil
	}
	
	bstr, _ := clerr.ToJson()

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(500)
	w.Write(bstr)

	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"msg": err.Error(),
	}).Error() // server error + rqst
	// TODO: #33! Send an error to admin
}

// Send 400 or 422 response (or something, different
//   than right response)
func handleClientError(r *http.Request,
	w http.ResponseWriter,
	err *clerr){
	
	str, parseErr := err.ToJson()
	
	if parseErr != nil {
		handleServerError(r, w, parseErr)	
		return
	}

	w.Header().Set("Content-Type", jsonMime)
	w.WriteHeader(422)
	w.Write(str)
	
	lgr.WithFields(lgr.Fields{
		"tag": tagRqst,
		"msg": string(str),
	}).Warn()	// client error
}

func handleSuccess(r *http.Request,
	w http.ResponseWriter,
	rdata interface{}){

	if rdata == nil {
    	//w.Write([]byte("sdf"))
	    //w.WriteHeader(http.StatusNoContent)
	    //
		// just header: no content
	    w.WriteHeader(204) // no content type
		
		// Logging (after execution)
		lgr.WithFields(lgr.Fields{
			"tag": tagRqst,
			"msg": "success",
			//"uid": uid,
			//"params": strParams,
			//"perms": perms,
		}).Info()  // rqst + info
		
	    return
    }   
	
	responseJson, errJson := json.Marshal(rdata)
	
	if errJson != nil {
		handleServerError(r, w, errJson)
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
}
