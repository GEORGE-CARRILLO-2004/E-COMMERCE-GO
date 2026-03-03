package vo

import "golan/domain/payment/errors"

type PaymentStatus struct {
	value string
}

func PendingStatus() PaymentStatus    { return PaymentStatus{value: "PENDING"} }
func ProcessingStatus() PaymentStatus { return PaymentStatus{value: "PROCESSING"} }
func CompletedStatus() PaymentStatus  { return PaymentStatus{value: "COMPLETED"} }
func FailedStatus() PaymentStatus     { return PaymentStatus{value: "FAILED"} }
func RefundedStatus() PaymentStatus   { return PaymentStatus{value: "REFUNDED"} }

func ReconstitutePaymentStatus(val string) PaymentStatus {
	return PaymentStatus{value: val}
}

func (s PaymentStatus) String() string {
	return s.value
}

func (s PaymentStatus) TransitionTo(newStatus string) (PaymentStatus, error) {
	switch s.value {
	case "PENDING":
		if newStatus == "PROCESSING" || newStatus == "FAILED" {
			return PaymentStatus{value: newStatus}, nil
		}
	case "PROCESSING":
		if newStatus == "COMPLETED" || newStatus == "FAILED" {
			return PaymentStatus{value: newStatus}, nil
		}
	case "COMPLETED":
		if newStatus == "REFUNDED" {
			return PaymentStatus{value: newStatus}, nil
		}
	}
	return s, errors.ErrInvalidPaymentState
}
