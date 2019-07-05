package matching

import (
	"fmt"
)

type DoubleSubmitError struct{}
type InvalidOrderError struct{}
type NotEnoughVolume struct{}
type ExceedSumLimit struct{}

func (d *DoubleSubmitError) Error() string {
	return fmt.Sprintf("err")
}
func (i *InvalidOrderError) Error() string {
	return fmt.Sprintf("err")
}
func (n *NotEnoughVolume) Error() string {
	return fmt.Sprintf("err")
}
func (doubleSeubmitError *ExceedSumLimit) Error() string {
	return fmt.Sprintf("err")
}
