package main

import "fmt"

type InvalidFieldsError struct {
	affectedField string
	reason        string
}

func (m *InvalidFieldsError) Error() string {
	return fmt.Sprintf("Cannot process query param: <%s>. Reason: <%s>", m.affectedField, m.reason)
}
