package apphandler

import (
    "fmt"
	"net/http"
	"encoding/json"
)

// setOutMime for response
func setOutMime(w http.ResponseWriter, outMime string){
	if outMime == "json" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	} 
	// no other types
}

// Send 400 or 422 response (or something, different than right response)
func handleClientError(w http.ResponseWriter, err IClerr, outMime string){
	
	str, parseErr := err.ToJson()
	
	if parseErr != nil {
		handleServerError(w, parseErr, outMime)	
		return
	}
	
	w.WriteHeader(422)
	w.Write([]byte(str))
}

// handleServerError writes error to client and send a notif to admin
func handleServerError(w http.ResponseWriter, err error, outMime string){
	
	clerr := Clerr {
		ErrKey: "unexpectedError",
		//Details: nil
	}
	
	bstr, _ := clerr.ToJson()
	   	
	w.WriteHeader(500)
	w.Write(bstr)
	
	fmt.Printf("systemError %v", err)
	fmt.Println()
	// TODO: #33! Send an error to admin
}

// appHandler
// func as a paramter
// error as a result
// Result - a map (will be converted to json/othertype to send to user
type AppHandlerType func(map[string]string) (IReq, IClerr, error)

// ServeHTTP implements Handler interface
// Param ah - data of appHandlerType instead Handler type
func (ah AppHandlerType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check cors requests
	if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin) // "http://google.com"
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS") //, PUT, DELETE
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept") 
            //X-CSRF-Token,  Content-Length, Accept-Encoding, Authorization
            
        // Non-simple requests
        // http://stackoverflow.com/questions/10636611/how-does-access-control-allow-origin-header-work
    }
    
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }
	
	// Returned format stores in url parameter (not in Accept header)
	//acceptTypes := r.Header["Accept"]
	var outMime string = "json" //by default
	
	parseErr := r.ParseForm()
	if parseErr != nil {
		handleServerError(w, parseErr, outMime)
		return
	}
	
	// body or url params - body in priority
	inParams := BeautyMap(r.Form)
	
	setOutMime(w, outMime)
	
	// Log a request
	fmt.Printf("inParams: %v", inParams)
	fmt.Println()
	
	// TODO: #33! Convert r.Form to normal map (without arrays)
	// all arrays will be sended, using _ divider in one param
	// like favarr=123_234_2345&otherparam=123
	
	// Execute required function
	// If errors occured - return error to client
	// w, r - convert to data for controllers
	rdata, errClient, errServer := ah(inParams);
	if errClient != nil {
		// Handle error client
		// Send 400 or 422 response (or something, different than right response)
		//
		handleClientError(w, errClient, outMime)
		return
	}
	
    if  errServer != nil {
    	// Controller methods can return client or server errors
    	//
    	// If a server error - need additional description and place
    	// and some other parameters to collect info about error
    	// to log the error (and send to admin)
    	//
    	// If a client error (like your text is already posted)
    	// - id (from 10 to 255)
    	// - title (attached from enum)
    	// - initial function (do not send to user or omit)
    	// - description (attached from enum)
    	// - a map of helper parameters (like price=23 is not allowed)
    	//
    	// For server errors (like unstable db connections)
    	// - id (1 - common, possible from 2-9, like db or net errors, 
    	//   or temporary unavailable)
    	// - title (attached from enum)
    	// - initial function
    	// - description (standard error.message)
    	// - parameters, associated with this error
    	//
    	// TODO: #33! define error
    	//if errCode = 1 {
    	handleServerError(w, errServer, outMime)
    	//}
        return
    } 
    
    //rdata = nil
    
    if rdata == nil {
    	//w.Write([]byte("sdf"))
	    //w.WriteHeader(http.StatusNoContent)
	    //
	    w.WriteHeader(204) // HTTP 500
    	//w.Write([]byte("Hello123!"))
	    return
    }
    
    var responseData string
	
	responseJson, errJson := json.Marshal(rdata)
	
	if errJson != nil {
		handleServerError(w, errJson, outMime)
		return
	}
	
	responseData = string(responseJson)
    
    // Write writes the data to the connection as part of an HTTP reply.
    // If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
    // before writing the data.  If the Header does not contain a
    // Content-Type line, Write adds a Content-Type set to the result of passing
    // the initial 512 bytes of written data to DetectContentType.
    w.Write([]byte(responseData))
}