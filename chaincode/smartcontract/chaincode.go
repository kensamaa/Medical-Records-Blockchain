package smartcontract

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing medical records.
type SmartContract struct {
	contractapi.Contract
}

// CreateMedicalRecord creates a new medical record.
func (s *SmartContract) CreateMedicalRecord(ctx contractapi.TransactionContextInterface, recordJSON string) error {
	var record MedicalRecord
	if err := json.Unmarshal([]byte(recordJSON), &record); err != nil {
		return fmt.Errorf("failed to unmarshal record: %v", err)
	}

	// Access control: Verify if the client has permissions to create a record.
	if err := verifyAccess(ctx, "create", record.PatientID); err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	record.CreatedAt = now
	record.UpdatedAt = now

	recordBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record: %v", err)
	}

	return ctx.GetStub().PutState(record.ID, recordBytes)
}

// ReadMedicalRecord retrieves a medical record from the ledger using its ID.
func (s *SmartContract) ReadMedicalRecord(ctx contractapi.TransactionContextInterface, id string) (*MedicalRecord, error) {
	recordBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if recordBytes == nil {
		return nil, fmt.Errorf("record %s does not exist", id)
	}

	var record MedicalRecord
	if err := json.Unmarshal(recordBytes, &record); err != nil {
		return nil, fmt.Errorf("failed to unmarshal record: %v", err)
	}

	// Access control: Verify if the client is allowed to view this record.
	if err := verifyAccess(ctx, "read", record.PatientID); err != nil {
		return nil, err
	}

	return &record, nil
}

// UpdateMedicalRecord updates an existing medical record.
func (s *SmartContract) UpdateMedicalRecord(ctx contractapi.TransactionContextInterface, recordJSON string) error {
	var updatedRecord MedicalRecord
	if err := json.Unmarshal([]byte(recordJSON), &updatedRecord); err != nil {
		return fmt.Errorf("failed to unmarshal record: %v", err)
	}

	existingRecordBytes, err := ctx.GetStub().GetState(updatedRecord.ID)
	if err != nil {
		return fmt.Errorf("failed to read existing record: %v", err)
	}
	if existingRecordBytes == nil {
		return fmt.Errorf("record %s does not exist", updatedRecord.ID)
	}

	// Optionally, enforce that only permitted users can update.
	if err := verifyAccess(ctx, "update", updatedRecord.PatientID); err != nil {
		return err
	}

	// Update the timestamp.
	updatedRecord.UpdatedAt = time.Now().Format(time.RFC3339)

	newRecordBytes, err := json.Marshal(updatedRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal updated record: %v", err)
	}

	return ctx.GetStub().PutState(updatedRecord.ID, newRecordBytes)
}
