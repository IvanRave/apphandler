package apphandler

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

// appHandler
// func as a paramter
// error as a result
// Result - a map (will be converted to json/othertype to send to user
type AppHandlerType func(map[string]string, int32, int32) (
	interface{}, error)

// ServeHTTP implements Handler interface
// Param ah - data of appHandlerType instead Handler type
func (ah AppHandlerType) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {

	// access only from localhost (nginx)
	// headers moved to nginx
	//addAccessControl(w, r)
	
	// Stop here if its Preflighted OPTIONS request
	// move to nginx
	//if r.Method == "OPTIONS" {
	//	return
	//}

	// 2. PARSE middleware
	parseErr := r.ParseForm()
	if parseErr != nil {
		handleServerError(r, w, parseErr)
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
				handleNonAuth(r, w,
					"authTokenIsExpired", apiKey)
				return
			}
			
			handleServerError(r, w, err)
			return
		}


		if tkn.Valid == false {
			// handle 401 response
			handleNonAuth(r, w,
				"authTokenIsInvalid", apiKey)
			return
		}

		if uidFloat, isUid := tkn.Claims["uid"].(float64);
		isUid == false {	
			handleNonAuth(r, w, "authTokenUidIsEmpty", apiKey)
			return
		} else {
			// only int32 supported for UID
			uid = int32(uidFloat)
		}

		if permsFloat, isPerms := tkn.Claims["perms"].(float64);
		isPerms == false {
			handleNonAuth(r, w,
				"authTokenPermsIsEmpty", apiKey)
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
		if clerrApp, ok := errApp.(*clerr); ok {
			if clerrApp.ErrKey == "permissionError" {
				// 401
				handleNonAuth(r, w, "notEnoughPermissions", "")
			} else {			
				// 422
				handleClientError(r, w, clerrApp)
			}
		}else {
			// 500
			handleServerError(r, w, errApp)
		}
	} else {
		// 200, 204
		handleSuccess(r, w, rdata)
	}
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
