## 1 The Purpose of Chaincode

Chaincode in Hyperledger Fabric is essentially the smart contract. It contains the rules for how data can be created, read, updated, or deleted (CRUD) on the blockchain ledger. In our case, the chaincode manages medical records securely.

## 2 Structure of the Code

Our chaincode is split into several files:

models.go – Defines the data structure (i.e., the MedicalRecord model).
chaincode.go – Contains the main smart contract logic with functions like Create, Read, and Update.
access_control.go – Contains helper functions for role-based access control.
utils.go – (Optional) Contains helper functions, such as encryption/decryption utilities.

## 4. Understanding chaincode.go

Explanation:

Embedding contractapi.Contract: By embedding this, our SmartContract inherits necessary functionalities that let Fabric know how to call our functions.
Purpose: This struct acts as a container for our chaincode methods (like create, read, update).

### b. CreateMedicalRecord Function

Explanation:

Input Parameter: The function receives a JSON string (recordJSON) representing the medical record.
Unmarshaling: The JSON string is converted into a Go struct (MedicalRecord), so we can work with it in our code.
Access Control: The verifyAccess function (which we'll discuss later) checks if the caller is allowed to perform this action.
Timestamps: The current time is set for both CreatedAt and UpdatedAt.
Storing the Record: The record is marshaled back into JSON and stored on the ledger using ctx.GetStub().PutState, which saves it under a unique key (record.ID).

### c. ReadMedicalRecord Function

Explanation:

Retrieving Data: GetState(id) retrieves the record from the ledger using the record’s unique ID.
Error Handling: The function checks if the record exists and handles errors in unmarshaling.
Access Control: The access control check ensures only authorized users can read the record.

### d. UpdateMedicalRecord Function

Explanation:

Input and Validation: The function first unmarshals the updated record data from JSON and then checks if the record exists on the ledger.
Access Control: It verifies that the caller is allowed to update this record.
Timestamp Update: It refreshes the UpdatedAt timestamp to indicate when the record was modified.
Saving the Update: The updated record is then marshaled into JSON and written back to the ledger.

## 5. Understanding access_control.go

Explanation:

Client Identity: The function uses ctx.GetClientIdentity() to get the identity of the caller. In a production system, you’d inspect the certificate attributes (like role or organization) to enforce rules.
Logging and Placeholder: Currently, it only logs the action. Later, you could expand this to restrict access based on specific roles.

## 6. Utilities for Encryption in utils.go

Sometimes you want to secure sensitive data further by encrypting it before writing it to the ledger. The utilities in utils.go provide simple AES-256 encryption/decryption functions
Explanation:

AES-256: This symmetric encryption method uses a 256-bit key to secure data.
Encryption Process: A random Initialization Vector (IV) is generated and prepended to the ciphertext. The CFB mode (Cipher Feedback) is used to encrypt the data.
Decryption Process: The IV is extracted from the ciphertext, and then the same process is reversed to retrieve the original data.
