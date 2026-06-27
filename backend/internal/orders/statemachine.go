package orders

// Order status values. These mirror the CHECK constraint on orders.status.
const (
	StatusPending         = "pending"
	StatusAwaitingPayment = "awaiting_payment"
	StatusPaid            = "paid"
	StatusFulfilled       = "fulfilled"
	StatusShipped         = "shipped"
	StatusDelivered       = "delivered"
	StatusCancelled       = "cancelled"
	StatusRefunded        = "refunded"
)

// allowedTransitions maps a current status to the set of statuses it may move
// to. cancelled and refunded are terminal (no outgoing transitions).
var allowedTransitions = map[string][]string{
	StatusPending:         {StatusAwaitingPayment, StatusCancelled},
	StatusAwaitingPayment: {StatusPaid, StatusCancelled},
	StatusPaid:            {StatusFulfilled, StatusRefunded, StatusCancelled},
	StatusFulfilled:       {StatusShipped, StatusRefunded},
	StatusShipped:         {StatusDelivered, StatusRefunded},
	StatusDelivered:       {StatusRefunded},
	StatusCancelled:       {},
	StatusRefunded:        {},
}

// canTransition reports whether an order may move from the "from" status to the
// "to" status under the lifecycle state machine.
func canTransition(from, to string) bool {
	for _, next := range allowedTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}
