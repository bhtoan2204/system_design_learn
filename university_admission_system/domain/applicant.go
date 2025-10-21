package domain

// Applicant represents a candidate who can submit applications.
type Applicant struct {
	ID              string
	FullName        string
	Email           string
	HighSchoolGPA   float64
	EntranceScore   float64
	SubmittedAppIDs []string
}

// CanSubmit verifies whether the applicant is eligible to submit an application.
func (a *Applicant) CanSubmit() bool {
	const minGPA = 2.0
	const minEntranceScore = 40.0
	return a.HighSchoolGPA >= minGPA && a.EntranceScore >= minEntranceScore
}

// TrackSubmission attaches a newly created application to the applicant profile.
func (a *Applicant) TrackSubmission(applicationID string) {
	a.SubmittedAppIDs = append(a.SubmittedAppIDs, applicationID)
}
