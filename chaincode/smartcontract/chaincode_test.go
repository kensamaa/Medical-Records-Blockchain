package smartcontract

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
)

// ---
// Mocks for the Fabric Transaction Context and Client Identity
// ---

// MockClientIdentity implements the minimal interface required for testing.
type MockClientIdentity struct {
	id string
}

// GetID returns a dummy client ID.
func (mci *MockClientIdentity) GetID() (string, error) {
	return mci.id, nil
}

// Ensure MockClientIdentity satisfies the expected interface.
// (In production, the interface is provided by Fabric’s cid package)

// MockTransactionContext implements the TransactionContextInterface.
type MockTransactionContext struct {
	stub           *shimtest.MockStub
	clientIdentity *MockClientIdentity
}

// GetStub returns the underlying chaincode stub.
func (m *MockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return m.stub
}

// GetClientIdentity returns the mock client identity.
func (m *MockTransactionContext) GetClientIdentity() *MockClientIdentity {
	return m.clientIdentity
}

// ---
// Helper function to set up a test instance of our chaincode and mock context.
// ---

func setupTest(t *testing.T) (*SmartContract, *MockTransactionContext) {
	// Create a new instance of our smart contract
	chaincode := new(SmartContract)

	// Create a mock stub (this simulates the ledger state)
	stub := shimtest.NewMockStub("mockStub", chaincode)

	// Create a mock transaction context with a dummy client identity.
	ctx := &MockTransactionContext{
		stub:           stub,
		clientIdentity: &MockClientIdentity{id: "mockClient"},
	}

	return chaincode, ctx
}

// ---
// Unit Test: CreateMedicalRecord
// ---

func TestCreateMedicalRecord(t *testing.T) {
	cc, ctx := setupTest(t)

	// Prepare a dummy medical record.
	record := MedicalRecord{
		ID:          "record1",
		PatientID:   "patient1",
		DoctorID:    "doctor1",
		HospitalID:  "hospital1",
		Diagnosis:   "Flu",
		Treatment:   "Rest",
		Medications: []string{"med1", "med2"},
		// CreatedAt and UpdatedAt will be set in the function.
	}
	recordJSON, err := json.Marshal(record)
	assert.NoError(t, err)

	// Invoke the CreateMedicalRecord function.
	err = cc.CreateMedicalRecord(ctx, string(recordJSON))
	assert.NoError(t, err)

	// Retrieve the record from the mock ledger.
	storedBytes := ctx.stub.State["record1"]
	assert.NotNil(t, storedBytes)

	var storedRecord MedicalRecord
	err = json.Unmarshal(storedBytes, &storedRecord)
	assert.NoError(t, err)

	// Validate fields (timestamps are set by the chaincode).
	assert.Equal(t, record.ID, storedRecord.ID)
	assert.Equal(t, record.PatientID, storedRecord.PatientID)
	assert.Equal(t, record.DoctorID, storedRecord.DoctorID)
	assert.Equal(t, record.HospitalID, storedRecord.HospitalID)
	assert.Equal(t, record.Diagnosis, storedRecord.Diagnosis)
	assert.Equal(t, record.Treatment, storedRecord.Treatment)
	assert.ElementsMatch(t, record.Medications, storedRecord.Medications)
	assert.NotEmpty(t, storedRecord.CreatedAt)
	assert.NotEmpty(t, storedRecord.UpdatedAt)
}

// ---
// Unit Test: ReadMedicalRecord
// ---

func TestReadMedicalRecord(t *testing.T) {
	cc, ctx := setupTest(t)

	// Prepare and insert a dummy record directly into the ledger.
	record := MedicalRecord{
		ID:          "record2",
		PatientID:   "patient2",
		DoctorID:    "doctor2",
		HospitalID:  "hospital2",
		Diagnosis:   "Cold",
		Treatment:   "Medication",
		Medications: []string{"medA", "medB"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	recordBytes, err := json.Marshal(record)
	assert.NoError(t, err)
	ctx.stub.State["record2"] = recordBytes

	// Invoke the ReadMedicalRecord function.
	returnedRecord, err := cc.ReadMedicalRecord(ctx, "record2")
	assert.NoError(t, err)

	// Validate that the returned record matches the stored record.
	assert.Equal(t, record.ID, returnedRecord.ID)
	assert.Equal(t, record.PatientID, returnedRecord.PatientID)
	assert.Equal(t, record.DoctorID, returnedRecord.DoctorID)
	assert.Equal(t, record.HospitalID, returnedRecord.HospitalID)
	assert.Equal(t, record.Diagnosis, returnedRecord.Diagnosis)
	assert.Equal(t, record.Treatment, returnedRecord.Treatment)
	assert.ElementsMatch(t, record.Medications, returnedRecord.Medications)
}

// ---
// Unit Test: UpdateMedicalRecord
// ---

func TestUpdateMedicalRecord(t *testing.T) {
	cc, ctx := setupTest(t)

	// Start with an existing record.
	record := MedicalRecord{
		ID:          "record3",
		PatientID:   "patient3",
		DoctorID:    "doctor3",
		HospitalID:  "hospital3",
		Diagnosis:   "Headache",
		Treatment:   "Painkillers",
		Medications: []string{"medX"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	origBytes, err := json.Marshal(record)
	assert.NoError(t, err)
	ctx.stub.State["record3"] = origBytes

	// Prepare an update: change the diagnosis and treatment.
	updatedRecord := record
	updatedRecord.Diagnosis = "Migraine"
	updatedRecord.Treatment = "Stronger Painkillers"
	updatedJSON, err := json.Marshal(updatedRecord)
	assert.NoError(t, err)

	// Invoke the UpdateMedicalRecord function.
	err = cc.UpdateMedicalRecord(ctx, string(updatedJSON))
	assert.NoError(t, err)

	// Retrieve the updated record from the ledger.
	storedBytes := ctx.stub.State["record3"]
	assert.NotNil(t, storedBytes)
	var storedRecord MedicalRecord
	err = json.Unmarshal(storedBytes, &storedRecord)
	assert.NoError(t, err)

	// Verify that the updates were applied.
	assert.Equal(t, "Migraine", storedRecord.Diagnosis)
	assert.Equal(t, "Stronger Painkillers", storedRecord.Treatment)
	// The UpdatedAt timestamp should have changed.
	assert.NotEqual(t, record.UpdatedAt, storedRecord.UpdatedAt)
}

// ---
// (Optional) Print output from verifyAccess for debugging.
// ---

func Example_verifyAccess() {
	// Create a dummy context.
	cc, ctx := setupTest(nil)
	err := verifyAccess(ctx, "test", "patientX")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Output will include a log from verifyAccess. (This is for demonstration.)
}
