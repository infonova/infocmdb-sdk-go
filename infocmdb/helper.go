package infocmdb

import (
	"regexp"
	"strconv"

	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/cmdb"
)

func (c *Client) AttributeBasedRelation(sourceCiId int, attributeName string, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {

	var currentCiValues []string
	if triggerType != "ci_attribute_delete" {
		var value string
		value, _, err = c.GetCiAttributeValueCi(sourceCiId, attributeName)
		if err != nil {
			return
		}

		currentCiValues = regexp.MustCompile(",\\s?").Split(value, -1)
	}

	destinationCiIds := make([]int, len(currentCiValues))
	for i, s := range currentCiValues {
		destinationCiIds[i], _ = strconv.Atoi(s)
	}

	return c.CiBasedRelation(sourceCiId, destinationCiIds, ciRelationTypeName, triggerType, swapCiColumns)
}

func (c *Client) CiBasedRelation(srcCiId int, destCiId []int, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {

	currentCiRelations, err := c.GetListOfCiIdsByCiRelation(srcCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_ALL)
	if err != nil {
		return
	}
	// add relations
	for _, valueCiId := range destCiId {
		add := true
		for _, relationCiId := range currentCiRelations {
			if valueCiId == relationCiId {
				add = false
				break
			}
		}

		if add == true {
			if swapCiColumns == true {
				err = c.CreateCiRelation(valueCiId, srcCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
			} else {
				err = c.CreateCiRelation(srcCiId, valueCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
			}
			if err != nil {
				return
			}

			relationCisAdded = append(relationCisAdded, valueCiId)
		}
	}

	// remove relations
	for _, relationCiId := range currentCiRelations {
		remove := true
		for _, valueCiId := range destCiId {
			if valueCiId == relationCiId {
				remove = false
				break
			}
		}

		if remove == true {
			err = c.DeleteCiRelation(srcCiId, relationCiId, ciRelationTypeName)
			if err != nil {
				return
			}

			relationCisRemoved = append(relationCisRemoved, relationCiId)
		}
	}

	return
}
