package apphandler

import (
    "fmt"
	"os"
//	"log"
//	"time"
//	"reflect"
	"io/ioutil"
	"strings"
	"net/http"
	"encoding/json"

	//"encoding/base64"
	jwt "github.com/dgrijalva/jwt-go"

	lgr	"github.com/Sirupsen/logrus"
)

const JSON_MIME string = "application/json; charset=utf-8"

// singleton property
var publicKey []byte

func init() {
	lgr.SetFormatter(&lgr.JSONFormatter{})
	lgr.SetOutput(os.Stderr)
	// openssl genrsa -out demo.rsa 1024 # the 1024 is the size of the key we are generating
    // openssl rsa -in demo.rsa -pubout > demo.rsa.pub
	pblKey, errReadKey := ioutil.ReadFile("/opt/jauth/jkey.rsa.pub")
	
	if errReadKey != nil {
		lgr.Fatal(errReadKey)
		// calls os.Exit(1) after logging
		return
	}	

	publicKey = pblKey
}

func cbkJwtParse(token *jwt.Token) (interface{}, error) {

	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

		return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])

	} else {
		return publicKey, nil
	}
	//myLookupKey(token.Header["kid"])
}

// Send 400 or 422 response (or something, different
//   than right response)
func handleClientError(w http.ResponseWriter, err IClerr){
	
	str, parseErr := err.ToJson()
	
	if parseErr != nil {
		handleServerError(w, parseErr)	
		return
	}

	w.Header().Set("Content-Type", JSON_MIME)
	w.WriteHeader(422)
	w.Write([]byte(str))
}

func handleNonAuth(w http.ResponseWriter,
	str string,
	apiKey string){
	
	w.WriteHeader(401)
	w.Write([]byte(str))
	lgr.WithFields(lgr.Fields {
		"err_key": str,
		"api_key": apiKey,
	}).Warn("Error401")
}

// handleServerError writes error to client
// and sends a notif to admin
func handleServerError(w http.ResponseWriter,
	err error){
	
	clerr := Clerr {
		ErrKey: "unexpectedError",
		//Details: nil
	}
	
	bstr, _ := clerr.ToJson()

	w.Header().Set("Content-Type", JSON_MIME)
	w.WriteHeader(500)
	w.Write(bstr)

	lgr.WithFields(lgr.Fields{
		"err": err,
	}).Error("ServerError")
	
	// TODO: #33! Send an error to admin
}

func handleSuccess(w http.ResponseWriter, rdata interface{}){
	if rdata == nil {
    	//w.Write([]byte("sdf"))
	    //w.WriteHeader(http.StatusNoContent)
	    //
		// just header: no content
	    w.WriteHeader(204) // no content type
	    return
    }   
	
	responseJson, errJson := json.Marshal(rdata)
	
	if errJson != nil {
		handleServerError(w, errJson)
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
	w.Header().Set("Content-Type", JSON_MIME)
    w.Write([]byte(string(responseJson)))
}

func calcApiKey(hdr http.Header) string {
	// get Authorization and authorization and other forms
	authHeader := hdr.Get("Authorization")
	if authHeader != "" {
		// extract Bearer
		arrStr := strings.Split(authHeader, " ")

		if len(arrStr) == 2 {
			if strings.ToLower(arrStr[0]) == "bearer" {
				if arrStr[1] != "" {
					return arrStr[1];
				}
			}
		}
	}
	
	return ""
}

func addAccessControl(w http.ResponseWriter, r *http.Request) {
	
	// check cors requests
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		//, PUT, DELETE
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept") 
		//X-CSRF-Token,  Content-Length, Accept-Encoding, Authorization			
		// Non-simple requests
		// http://stackoverflow.com/questions/10636611/how-does-access-control-allow-origin-header-work
	}
}

// Decode JWT specific base64url encoding with padding stripped
// func DecodeSegment(seg string) ([]byte, error) {
// 	if l := len(seg) % 4; l > 0 {
// 		seg += strings.Repeat("=", 4-l)
// 	}
	
// 	return base64.URLEncoding.DecodeString(seg)
// }

// appHandler
// func as a paramter
// error as a result
// Result - a map (will be converted to json/othertype to send to user
type AppHandlerType func(map[string]string, int32, int32) (
	interface{},
	error)

// ServeHTTP implements Handler interface
// Param ah - data of appHandlerType instead Handler type
func (ah AppHandlerType) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {

	addAccessControl(w, r)
	
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	// 2. PARSE middleware
	parseErr := r.ParseForm()
	if parseErr != nil {
		handleServerError(w, parseErr)
		return
	}
	
	// body or url params - body in priority
	inParams := BeautyMap(r.Form)
	
	// Returned format stores in url parameter
	// (not in Accept header)
	//acceptTypes := r.Header["Accept"]

	
	// Log a request

	// 3. AUTH middleware
	// - check auth-token
	// - transform it to permScope (if exists)
	// - send permScope to next middleware
	// This MDW doesn't sends AuthToken to next MDW (only perms)
	
	// TODO: #33! Convert r.Form to normal map (without arrays)
	// all arrays will be sended, using _ divider in one param
	// like favarr=123_234_2345&otherparam=123	
	
	apiKey := calcApiKey(r.Header);

	// user id: calc from apiKey
	var uid int32 = 0
	var perms int32 = 0
	// Check apiKey for all requests (even for non-authed)
	// a client sends a token only for authed requests
	if apiKey != "" {

		tkn, err := jwt.Parse(apiKey, cbkJwtParse)

		// Check expired time: automatically inside jwt library
		// https://github.com/dgrijalva/jwt-go/blob/master/jwt.go#L140
		if err != nil {
			if err.Error() == "token is expired" {
				handleNonAuth(w, "authTokenIsExpired", apiKey)
				return
			}
			
			handleServerError(w, err)
			return
		}


		if tkn.Valid == false {
			// handle 401 response
			handleNonAuth(w, "authTokenIsInvalid", apiKey)
			return
		}

		if uidFloat, isUid := tkn.Claims["uid"].(float64);
		isUid == false {	
			handleNonAuth(w, "authTokenUidIsEmpty", apiKey)
			return
		} else {
			// only int32 supported for UID
			uid = int32(uidFloat)
		}

		if permsFloat, isPerms := tkn.Claims["perms"].(float64);
		isPerms == false {	
			handleNonAuth(w, "authTokenPermsIsEmpty", apiKey)
			return
		} else {
			// only int32 supported for PERMS
			perms = int32(permsFloat)
		}
		
		//parts := strings.Split(tkn.Raw, ".")

		//dcd, _ := DecodeSegment(parts[1])


		// exp: int64 unixtimestamp
	}
	
	// translate apiKey to userId + perms (roles)
	// from JWT or TableSession
	// define perms from DB or JWT?
	// Are perms can be different per sessions?
	// To get perms from DB - DB perms required (cycle)
	
	
	// Execute required function
	// If errors occured - return error to client
	// w, r - convert to data for controllers
	// 4. MAIN middleware
	// - check permScope, if required
	// - check inParams
	// - execute main methods
	// ? nil or int32 = 0


	// 5. RESULT middleware
	if 	rdata, errApp := ah(inParams, uid, perms);
	errApp != nil {
		if clerrApp, ok := errApp.(*Clerr); ok {
			// 4xx (422)
			handleClientError(w, clerrApp)
		}else {
			// 5xx (500)
			handleServerError(w, errApp)
		}
	} else {
		// 2xx (200, 204)
		handleSuccess(w, rdata)
	}

	// Logging (after execution)
	lgr.WithFields(lgr.Fields{
		"url": r.URL,
		"uid": uid,
		"params": inParams,
		"perms": perms,
	}).Info("request data")
}

// if  errServer != nil {
//  if tmpClerr, ok := errServer.(*Clerr); ok {
// 		handleClientError(w, tmpClerr, outMime)
// 		return
// 	}

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

// 	//}
//     return
// } 

//rdata = nil
