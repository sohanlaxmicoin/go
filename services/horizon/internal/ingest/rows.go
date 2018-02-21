package ingest

func (r effectRow) GetParams() []interface{} {
	return []interface{}{
		r.AccountID,
		r.OperationID,
		r.Order,
		r.Type,
		r.Details,
	}
}

func (r *effectRow) UpdateAccountIDs(accounts map[string]int64) {
	r.AccountID = accounts[r.Address]
}

func (r effectRow) GetAddresses() []string {
	return []string{r.Address}
}

func (r effectRow) GetTableName() TableName {
	return EffectsTableName
}

func (r operationRow) GetParams() []interface{} {
	return []interface{}{
		r.ID,
		r.TxID,
		r.Order,
		r.Source,
		r.Type,
		r.Details,
	}
}

func (r operationRow) UpdateAccountIDs(accounts map[string]int64) {
	return
}

func (r operationRow) GetAddresses() []string {
	return nil
}

func (r operationRow) GetTableName() TableName {
	return OperationsTableName
}

func (r operationParticipantRow) GetParams() []interface{} {
	return []interface{}{
		r.OperationID,
		r.AccountID,
	}
}

func (r *operationParticipantRow) UpdateAccountIDs(accounts map[string]int64) {
	r.AccountID = accounts[r.Address]
}

func (r operationParticipantRow) GetAddresses() []string {
	return []string{r.Address}
}

func (r operationParticipantRow) GetTableName() TableName {
	return OperationParticipantsTableName
}

func (r tradeRow) GetParams() []interface{} {
	return []interface{}{
		r.OperationID,
		r.Order,
		r.LedgerCloseAt,
		r.OfferID,
		r.BaseAccountID,
		r.BaseAssetID,
		r.BaseAmount,
		r.CounterAccountID,
		r.CounterAssetID,
		r.CounterAmount,
		r.BaseIsSeller,
	}
}

func (r *tradeRow) UpdateAccountIDs(accounts map[string]int64) {
	r.BaseAccountID = accounts[r.BaseAddress]
	r.CounterAccountID = accounts[r.CounterAddress]
}

func (r tradeRow) GetAddresses() []string {
	return []string{r.BaseAddress, r.CounterAddress}
}

func (r tradeRow) GetTableName() TableName {
	return TradesTableName
}

func (r transactionParticipantRow) GetParams() []interface{} {
	return []interface{}{
		r.TransactionID,
		r.AccountID,
	}
}

func (r *transactionParticipantRow) UpdateAccountIDs(accounts map[string]int64) {
	r.AccountID = accounts[r.Address]
}

func (r transactionParticipantRow) GetAddresses() []string {
	return []string{r.Address}
}

func (r transactionParticipantRow) GetTableName() TableName {
	return TransactionParticipantsTableName
}
