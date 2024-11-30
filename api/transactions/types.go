package transactions

type TransactionsType struct {
	Hash  string `json:"hash" bson:"hash" validate:"required"`
	Nonce uint64 `json:"nonce" bson:"nonce"`
	From  string `json:"from" bson:"from"`
	To    string `json:"to" bson:"to"`
}
