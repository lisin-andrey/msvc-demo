package consts

// ----- COMMON RELATION URLS -----
const (
	// ReadinessPathURL - relation path of readiness
	ReadinessPathURL = "/-/ready"
	// HealthyPathURL - relation path of liveness
	HealthyPathURL = "/-/healthy"
	// MetricsPathURL - relation path of metrics
	MetricsPathURL = "/-/metrics"
)

const (
	// TestReponseSuccess - test response of rest service for liveness
	TestReponseSuccess = "OK"
)
