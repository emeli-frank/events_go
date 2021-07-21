package events

import (
	"events/pkg/validation"
	"fmt"
	"time"
)

var StartTimeRule = func(required bool) validation.Rule {
	return validation.RuleFunc(
		func(v interface{}) error {
			if required && v == nil {
				return validation.NewError("start time is required")
			}

			start, ok := v.(*time.Time)
			if !ok {
				return fmt.Errorf("value is not *time.Time, %T gotten instead", v)
			}

			now := time.Now()
			if start != nil {
				if start.Before(now) {
					return validation.NewError("start time has passed")
				}
			}
			return nil
		},
	)
}

var EndTimeRule = func(start *time.Time, required bool) validation.Rule {
	return validation.RuleFunc(
		func(v interface{}) error {
			if required && v == nil {
				return validation.NewError("end time is required")
			}

			end, ok := v.(*time.Time)
			if !ok {
				return fmt.Errorf("value is not *time.Time, %T gotten instead", v)
			}

			if start != nil {
				if end != nil && end.Before(*start) {
					return validation.NewError("end time cannot be set before start time")
				}
			} else {
				if end != nil {
					return validation.NewError("end time cannot be set without a start time")
				}
			}
			return nil
		},
	)
}
