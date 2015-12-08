package rpc

// Initialize - Initizes the RPC library
func Initialize() error {
	return InitializeIO()
}

// Shutdown - shuts down the RPC library
func Shutdown() {
	ShutdownIO()
}
