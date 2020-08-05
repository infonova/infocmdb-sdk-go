package infocmdb

import (
	"errors"
	"strconv"
)

type CiRelationDirection string

const (
	CI_RELATION_DIRECTION_ALL             CiRelationDirection = "all"
	CI_RELATION_DIRECTION_DIRECTED_FROM                       = "directed_from"
	CI_RELATION_DIRECTION_DIRECTED_TO                         = "directed_to"
	CI_RELATION_DIRECTION_BIDIRECTIONAL                       = "bidirectional"
	CI_RELATION_DIRECTION_OMNIDIRECTIONAL                     = "omnidirectional"
)

func NewCiRelationDirection(directionId int) (direction CiRelationDirection, err error) {
	switch directionId {
	case 1:
		direction = CI_RELATION_DIRECTION_DIRECTED_FROM
	case 2:
		direction = CI_RELATION_DIRECTION_DIRECTED_TO
	case 3:
		direction = CI_RELATION_DIRECTION_BIDIRECTIONAL
	case 4:
		direction = CI_RELATION_DIRECTION_OMNIDIRECTIONAL
	default:
		err = errors.New("invalid direction id: " + strconv.Itoa(directionId))
	}
	return
}

func (direction CiRelationDirection) GetId() (directionId int, err error) {
	switch direction {
	case CI_RELATION_DIRECTION_DIRECTED_FROM:
		directionId = 1
	case CI_RELATION_DIRECTION_DIRECTED_TO:
		directionId = 2
	case CI_RELATION_DIRECTION_BIDIRECTIONAL:
		directionId = 3
	case CI_RELATION_DIRECTION_OMNIDIRECTIONAL:
		directionId = 4
	default:
		err = errors.New("invalid direction: " + string(direction))
	}
	return
}
