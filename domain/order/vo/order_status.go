package vo

import "golan/domain/order/errors"

type OrderStatus struct {
	value string
}

func PendingStatus() OrderStatus   { return OrderStatus{value: "PENDING"} }
func ConfirmedStatus() OrderStatus { return OrderStatus{value: "CONFIRMED"} }
func PaidStatus() OrderStatus      { return OrderStatus{value: "PAID"} }
func ShippedStatus() OrderStatus   { return OrderStatus{value: "SHIPPED"} }
func DeliveredStatus() OrderStatus { return OrderStatus{value: "DELIVERED"} }
func CancelledStatus() OrderStatus { return OrderStatus{value: "CANCELLED"} }

func ReconstituteOrderStatus(val string) OrderStatus {
	return OrderStatus{value: val}
}

func (s OrderStatus) String() string {
	return s.value
}

func (s OrderStatus) CanBePaid() bool {
	return s.value == "PENDING"
}

func (s OrderStatus) CanModifyItems() bool {
	return s.value == "PENDING" || s.value == "CONFIRMED"
}

func (s OrderStatus) TransitionTo(newStatus string) (OrderStatus, error) {
	switch s.value {
	case "PENDING":
		if newStatus == "CONFIRMED" || newStatus == "PAID" || newStatus == "CANCELLED" {
			return OrderStatus{value: newStatus}, nil
		}
	case "CONFIRMED":
		if newStatus == "PAID" || newStatus == "CANCELLED" {
			return OrderStatus{value: newStatus}, nil
		}
	case "PAID":
		if newStatus == "SHIPPED" || newStatus == "CANCELLED" {
			return OrderStatus{value: newStatus}, nil
		}
	case "SHIPPED":
		if newStatus == "DELIVERED" {
			return OrderStatus{value: newStatus}, nil
		}
	}
	return s, errors.ErrInvalidOrderState
}
