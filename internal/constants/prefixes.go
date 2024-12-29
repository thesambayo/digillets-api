package constants

// https://gist.github.com/fnky/76f533366f75cf75802c8052b577e2a5
// Prefix constants used for identifying various entities.
// These prefixes are used to generate unique identifiers for different types of objects
// and transactions within the system. Each prefix is associated with a specific entity type.

const (
	// PrefixUserID is used for users IDs.
	PrefixUserID = "usr_"

	// PrefixWalletID is used for users IDs.
	PrefixWalletID = "wllt_"

	// PrefixTransactionID is used for transaction IDs.
	PrefixTransactionID = "txn_"
)
