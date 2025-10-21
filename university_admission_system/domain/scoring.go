package domain

// ScoreCalculator decides how to compute an application score.
type ScoreCalculator interface {
	Compute(app Application, applicant Applicant) float64
}

// DefaultScoreCalculator is a basic scoring strategy.
type DefaultScoreCalculator struct{}

// Compute aggregates GPA and entrance score into a single value.
func (DefaultScoreCalculator) Compute(app Application, applicant Applicant) float64 {
	const gpaWeight = 0.4
	const entranceWeight = 0.6
	return applicant.HighSchoolGPA*25*gpaWeight + applicant.EntranceScore*entranceWeight
}
