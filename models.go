package apphandler

import "encoding/json"

//type IReq interface {
	//ToJson() ([]byte, error)
	
    // Not supported: no possible to store value as keys (like nosql data)	
    // 	ToXml() ([]byte, error)
    // Cannot convert from
    // - map[string]interface{}
    // - map[string]Master
    // - map[string]string
    // these types is not supported in xml.Marshal
//}

// IClerr - interface, like error but only for client error 
// (opposite of server errors)
// type IClerr interface {
// 	ToJson() ([]byte, error)
// }

type clerr struct {
	ErrKey string   		`json:"errkey"`
	Details interface {}	`json:"details"`
}

// IClerr interface
func (appClerr clerr) ToJson() ([]byte, error) {
    return json.Marshal(appClerr)
}

// error interface
func (appClerr clerr) Error() string {
	return "client error: " + string(appClerr.ErrKey)
}

func ErrValidation(details interface{}) (*clerr) {
	return &clerr {
		ErrKey: "validationError",
		Details: details,
	}
}

// if no perms
func ErrPerms(userPerms int32, requiredPerms int32) (*clerr){
	return &clerr {
		ErrKey: "permissionError",
		Details: map[string]int32 {
			"perms": userPerms,
			"required_perms": requiredPerms,
		},
	}
}

// uid and perms can be viewed in JWT token payload
// required uid - in query
// required perms - client info (to request required perms
//     after unsuccessful result)
func ErrUid(uid int32, requiredUid int32) (*clerr){
	return &clerr {
		ErrKey: "permissionError",
		Details: map[string]int32 {
			"uid": uid,
			"required_uid": requiredUid,
		},
	}
}

func ErrDuplicateKey(propName string) (*clerr){
	return &clerr {
		ErrKey: "duplicateKeyError",
		Details: map[string]string {
			"property": propName,
		},
	}
}

func ErrForeignKey(propName string) (*clerr){
	return &clerr {
		ErrKey: "foreignKeyError",
		Details: map[string]string {
			"property": propName,
		},
	}
}

// noSuchServRubric
// nstgIdIsEmpty
// masterProfileIsNotFound
// noAccess
