package core

import (
	"regexp"
	"strconv"

	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb"
)

func (i *InfoCMDB) AttributeBasedRelation(sourceCiId int, attributeName string, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {

	var currentCiValues []string
	if triggerType != "ci_attribute_delete" {
		var value string
		value, _, err = i.GetCiAttributeValueCi(sourceCiId, attributeName)
		if err != nil {
			return
		}

		currentCiValues = regexp.MustCompile(",\\s?").Split(value, -1)
	}

	destinationCiIds := make([]int, len(currentCiValues))
	for i, s := range currentCiValues {
		destinationCiIds[i], _ = strconv.Atoi(s)
	}

	return i.CiBasedRelation(sourceCiId, destinationCiIds, ciRelationTypeName, triggerType, swapCiColumns)
}

func (i *InfoCMDB) CiBasedRelation(srcCiId int, destCiId []int, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved []int, err error) {

	currentCiRelations, err := i.GetListOfCiIdsByCiRelation(srcCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_ALL)
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
				err = i.CreateCiRelation(valueCiId, srcCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
			} else {
				err = i.CreateCiRelation(srcCiId, valueCiId, ciRelationTypeName, cmdb.CI_RELATION_DIRECTION_OMNIDIRECTIONAL)
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
			err = i.DeleteCiRelation(srcCiId, relationCiId, ciRelationTypeName)
			if err != nil {
				return
			}

			relationCisRemoved = append(relationCisRemoved, relationCiId)
		}
	}

	return
}
