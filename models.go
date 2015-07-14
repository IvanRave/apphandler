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
func ErrAccess(details string) (*clerr){
	return &clerr {
		ErrKey: "accessError",
		Details: details,
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
