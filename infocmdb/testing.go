package infocmdb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

// Type of the precondition that will be tested for existence.
type PreconditionType string

const (
	TYPE_CI_TYPE   = "ci type"
	TYPE_ATTRIBUTE = "attribute"
	TYPE_RELATION  = "relation"
)

// Single workflow requirement.
type Precondition struct {
	Type PreconditionType
	Name string
}

// List of workflow requirements.
type Preconditions []Precondition

// Executes an existence check test for every supplied precondition.
func (w Workflow) TestPreconditions(t *testing.T, preconditions Preconditions) {
	cmdb := NewClient()
	cmdbClientErr := cmdb.LoadConfig(w.config)
	if cmdbClientErr != nil {
		log.Fatalf("Failed to Login: %v", cmdbClientErr)
	}

	for _, precondition := range preconditions {
		testName := fmt.Sprintf("%v \"%v\" exists", precondition.Type, precondition.Name)

		t.Run(testName, func(t *testing.T) {
			testPrecondition(t, cmdb, precondition)
		})
	}
}

func testPrecondition(t *testing.T, cmdb *Client, precondition Precondition) {
	switch precondition.Type {
	case TYPE_CI_TYPE:
		if _, err := cmdb.GetCiTypeIdByCiTypeName(precondition.Name); err != nil {
			t.Error(err)
		}

	case TYPE_ATTRIBUTE:
		if _, err := cmdb.GetAttributeIdByAttributeName(precondition.Name); err != nil {
			t.Error(err)
		}

	case TYPE_RELATION:
		if _, err := cmdb.GetCiRelationTypeIdByRelationTypeName(precondition.Name); err != nil {
			t.Error(err)
		}
	}
}
