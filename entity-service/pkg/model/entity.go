package model

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	// InvalidEntityID - mark invalid entity id
	InvalidEntityID int32 = -1

	// KeyNameID - key name id. Create operation will return ID of created entity in the result map
	KeyNameID = "id"
)

//EntityCmd - entity model to save data
type EntityCmd struct {
	Name         string    `json:"name"`
	Descr        string    `json:"descr,omitempty"`
	LastUpdated  time.Time `json:"last_updated,omitempty"`
	LastOperator string    `json:"last-operator"`
}

//ToString - convert object to string
func (v EntityCmd) String() string {
	return fmt.Sprintf("Name: [%s], Descr: [%s], LastUpdated: [%s], LastOperator: [%s]",
		v.Name, v.Descr, v.LastUpdated, v.LastOperator)
}

// Validate - check valid data
func (v EntityCmd) Validate() error {
	var em string
	if len(v.Name) == 0 {
		em += "Name is empty\n"
	}
	// if v.LastUpdated.IsZero() {
	// 	em += "LastUpdated is undefined\n"
	// }
	if len(v.LastOperator) == 0 {
		em += "LastOperator is empty\n"
	}
	if len(em) > 0 {
		return errors.New(em)
	}
	return nil
}

// ValidateForUpdate - check valid data for Update operation (Name field can be empty in this case)
func (v EntityCmd) ValidateForUpdate() error {
	var em string
	if len(v.LastOperator) == 0 {
		em += "LastOperator is empty\n"
	}
	if len(em) > 0 {
		return errors.New(em)
	}
	return nil
}

//EntityQuery - entity model to show data
type EntityQuery struct {
	ID           int32     `json:"id"`
	Name         string    `json:"name"`
	Descr        string    `json:"descr"`
	Created      time.Time `json:"created"`
	LastUpdated  time.Time `json:"last_updated"`
	LastOperator string    `json:"last-operator"`
}

//String - convert object to string
func (v EntityQuery) String() string {
	return fmt.Sprintf("ID: %d, Name: [%s], Descr: [%s], Created: [%s], LastUpdated: [%s], LastOperator: [%s]",
		v.ID, v.Name, v.Descr, v.Created, v.LastUpdated, v.LastOperator)
}
