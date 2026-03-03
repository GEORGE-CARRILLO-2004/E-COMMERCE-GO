package vo

import "github.com/google/uuid"

type ReviewID struct {
	value string
}

func NewReviewID() ReviewID {
	return ReviewID{value: uuid.New().String()}
}

func ReconstituteReviewID(val string) ReviewID {
	return ReviewID{value: val}
}

func (id ReviewID) String() string {
	return id.value
}
