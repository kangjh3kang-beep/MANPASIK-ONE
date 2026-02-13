package clients

import "context"

// NotificationClient sends notifications to users
type NotificationClient interface {
	SendNotification(ctx context.Context, userID, notifType, title, body, priority, channel string) error
}

// PrescriptionClient creates prescriptions
type PrescriptionClient interface {
	CreatePrescriptionFromReservation(ctx context.Context, userID, doctorName, diagnosis string, medications []MedicationItem) (string, error)
}

// MedicationItem for prescription creation
type MedicationItem struct {
	Name     string
	Dosage   string
	Duration string
}

// HealthScoreClient gets health scores
type HealthScoreClient interface {
	GetHealthScore(ctx context.Context, userID string) (float64, string, error)
}

// MeasurementClient gets measurement data
type MeasurementClient interface {
	GetLatestMeasurements(ctx context.Context, userID string, limit int) ([]MeasurementSummary, error)
}

// MeasurementSummary represents a measurement for cross-service use
type MeasurementSummary struct {
	SessionID   string
	BiomarkerID string
	Value       float64
	Unit        string
	MeasuredAt  string
}

// SubscriptionClient checks subscription status
type SubscriptionClient interface {
	CheckAccess(ctx context.Context, userID, feature string) (bool, string, error)
}
