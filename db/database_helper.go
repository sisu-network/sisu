package db

// getQueryQuestionMark returns a string in a form (?, ?, ?), (?, ?, ?), (?, ?, ?) to allow
// multiple row insertion.
func getQueryQuestionMark(rowCount, fieldCount int) string {
	s := ""

	for i := 0; i < rowCount; i++ {
		q := ""
		for j := 0; j < fieldCount; j++ {
			q = q + "?"
			if j < fieldCount-1 {
				q = q + ", "
			}
		}

		s = s + "(" + q + ")"
		if i < rowCount-1 {
			s = s + ", "
		}
	}

	return s
}
