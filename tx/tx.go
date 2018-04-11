package tx

// Node is one of AffectedNodes in Meta.
// The key is one of "CreatedNode", "ModifiedNode" or "DeletedNode".
// The value is one of RippleState, AccountRoot, Offer or DirectoryNode which are LedgerEntryType's.
type Node map[string]map[string]interface{}

// NodeCommon contains the common fields for all Nodes.
type NodeCommon struct {
	LedgerEntryType   string
	LedgerIndex       string
	PreviousTxnID     string
	PreviousTxnLgrSeq uint32
}

// Balance is balance in RippleState.
type Balance struct {
	Currency string
	Issuer   string
	Value    string
}

// RippleStateFinalFields is the FinalFields in RippleState.
type RippleStateFinalFields struct {
	Balance   Balance
	Flags     uint32
	HighLimit Balance
	HighNode  string
	LowLimit  Balance
	LowNode   string
}

// RippleStatePreviousFields is the PreviousFields in RippleState.
type RippleStatePreviousFields struct {
	Balance Balance
}

// RippleState is one of the LedgerEntryType.
type RippleState struct {
	NodeCommon     `mapstructure:",squash"`
	FinalFields    RippleStateFinalFields
	PreviousFields RippleStatePreviousFields
}

// AccountRootFinalFields is the FinalFields in AccountRoot.
type AccountRootFinalFields struct {
	Account    string
	Balance    string
	Flags      uint32
	OwnerCount uint32
	Sequence   uint32
}

// AccountRootPreviousFields is the PreviousFields in AccountRoot.
type AccountRootPreviousFields struct {
	Balance    string
	OwnerCount uint32
	Sequence   uint32
}

// AccountRoot is one of the LedgerEntryType.
type AccountRoot struct {
	NodeCommon     `mapstructure:",squash"`
	FinalFields    AccountRootFinalFields
	PreviousFields AccountRootPreviousFields
}

// Meta includes all AffectedNodes in an XRP Ledger transaction.
type Meta struct {
	AffectedNodes     []Node
	TransactionIndex  uint32
	TransactionResult string
}

// Transaction includes informations about an XRP Ledger transaction.
type Transaction struct {
	Account            string
	Fee                string
	Flags              uint32
	LastLedgerSequence uint32
	OfferSequence      uint32
	Sequence           uint32
	SigningPubKey      string
	TransactionType    string
	TxnSignature       string
	Date               uint32
	Hash               string
	Ledger_index       uint32
	Meta               Meta
	Validated          bool
}

// Response is the response of rippled tx method. https://ripple.com/build/rippled-apis/#tx
type Response struct {
	Result Transaction
	Status string
	Type   string
}

// Request is the request of rippled tx method.
type Request struct {
	Command     string `json:"command"`
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary"`
}

// NewRequest returns a tx.Request.
func NewRequest(transaction string) *Request {
	return &Request{
		Command:     "tx",
		Transaction: transaction,
		Binary:      false,
	}
}
