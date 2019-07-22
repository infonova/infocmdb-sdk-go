package cmdb

import (
	"regexp"
	"strconv"
)

func (i *InfoCMDB) AttributeBasedRelation(ciId int, attributeName string, ciRelationTypeName string, triggerType string, swapCiColumns bool) (relationCisAdded []int, relationCisRemoved [] int, err error) {

	currentCiRelations, err := i.GetListOfCiIdsByCiRelation(ciId, ciRelationTypeName, CI_RELATION_DIRECTION_ALL)
	if err != nil {
		return
	}

	var currentCiValues []string
	if triggerType != "ci_attribute_delete" {
		var value string
		value, _, err = i.GetCiAttributeValueCi(ciId, attributeName)
		if err != nil {
			return
		}

		currentCiValues = regexp.MustCompile(",\\s?").Split(value, -1)
	}

	// add relations
	for _, valueCiIdString := range currentCiValues {
		valueCiId, _ := strconv.Atoi(valueCiIdString)
		add := true;
		for _, relationCiId := range currentCiRelations {
			if valueCiId == relationCiId {
				add = false
				break
			}
		}

		if add == true {
			if swapCiColumns == true {
				err = i.CreateCiRelation(valueCiId, ciId, ciRelationTypeName, CI_RELATION_DIRECTION_OMNIDIRECTIONAL);
			} else {
				err = i.CreateCiRelation(ciId, valueCiId, ciRelationTypeName, CI_RELATION_DIRECTION_OMNIDIRECTIONAL);
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
		for _, valueCiIdString := range currentCiValues {
			valueCiId, _ := strconv.Atoi(valueCiIdString)
			if valueCiId == relationCiId {
				remove = false
				break
			}
		}

		if remove == true {
			err = i.DeleteCiRelation(ciId, relationCiId, ciRelationTypeName)
			if err != nil {
				return
			}

			relationCisRemoved = append(relationCisRemoved, relationCiId)
		}
	}

	return
}
