package smartcontract

// MedicalRecord defines the structure for a patient record.
type MedicalRecord struct {
	ID          string   `json:"id"`
	PatientID   string   `json:"patientId"`
	DoctorID    string   `json:"doctorId"`
	HospitalID  string   `json:"hospitalId"`
	Diagnosis   string   `json:"diagnosis"`
	Treatment   string   `json:"treatment"`
	Medications []string `json:"medications"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
