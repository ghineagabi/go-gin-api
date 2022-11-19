package main

import "fmt"

type InvalidFieldsError struct {
	affectedField string
	reason        string
	location      string
}

func (m *InvalidFieldsError) Error() string {
	return fmt.Sprintf("Cannot process <%s> field: <%s>. Reason: <%s>", m.location, m.affectedField, m.reason)
}
