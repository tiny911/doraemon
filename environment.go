package doraemon

// Env deploy environment
type Env string

const (
	// Local local environment
	Local Env = "local"

	// Dev environment
	Dev Env = "dev"

	// TestAuto auto test environment
	TestAuto Env = "test_auto"

	// Test test environment
	Test Env = "test"

	// Sandbox sandbox environment
	Sandbox Env = "sandbox"

	// AppRelease release environment
	AppRelease Env = "app_release"

	// OnlinePre onlinePre environment
	OnlinePre Env = "online_pre"

	// Online online environment
	Online Env = "online"
)
