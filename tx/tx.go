package tx

type Node map[string]map[string]interface{}

type Meta struct {
	AffectedNodes     []Node
	TransactionIndex  uint32
	TransactionResult string
}

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

type TransactionResponse struct {
	Result Transaction
	Status string
	Type   string
}

type NodeCommon struct {
	LedgerEntryType   string
	LedgerIndex       string
	PreviousTxnID     string
	PreviousTxnLgrSeq uint32
}

type Balance struct {
	Currency string
	Issuer   string
	Value    string
}

type RippleStateFinalFields struct {
	Balance   Balance
	Flags     uint32
	HighLimit Balance
	HighNode  string
	LowLimit  Balance
	LowNode   string
}

type RippleStatePreviousFields struct {
	Balance Balance
}

type RippleState struct {
	NodeCommon     `mapstructure:",squash"`
	FinalFields    RippleStateFinalFields
	PreviousFields RippleStatePreviousFields
}

type AccountRootFinalFields struct {
	Account    string
	Balance    string
	Flags      uint32
	OwnerCount uint32
	Sequence   uint32
}

type AccountRootPreviousFields struct {
	Balance    string
	OwnerCount uint32
	Sequence   uint32
}

type AccountRoot struct {
	NodeCommon     `mapstructure:",squash"`
	FinalFields    AccountRootFinalFields
	PreviousFields AccountRootPreviousFields
}
