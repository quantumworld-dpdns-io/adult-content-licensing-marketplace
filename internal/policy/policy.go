package policy

type CompliancePolicy struct {
	Frameworks        []string `json:"frameworks"`
	KYCAML            string   `json:"kyc_aml"`
	AgeVerification   string   `json:"age_verification"`
	DMCA              string   `json:"dmca"`
	GDPRCCPA          string   `json:"gdpr_ccpa"`
	AITrainingControl string   `json:"ai_training_control"`
}

type ThreatModelPolicy struct {
	Methods        []string `json:"methods"`
	PriorityRisks  []string `json:"priority_risks"`
	OwaspReference string   `json:"owasp_reference"`
}

func Compliance() CompliancePolicy {
	return CompliancePolicy{
		Frameworks:        []string{"KYC/AML", "Age Verification", "GDPR/CCPA", "2257", "DMCA"},
		KYCAML:            "kyc-lite",
		AgeVerification:   "required",
		DMCA:              "notice-and-takedown workflow required",
		GDPRCCPA:          "data access/deletion workflow required",
		AITrainingControl: "contractual + technical flags",
	}
}

func ThreatModel() ThreatModelPolicy {
	return ThreatModelPolicy{
		Methods:        []string{"STRIDE", "PASTA"},
		PriorityRisks:  []string{"api_abuse", "data_breach", "fraud"},
		OwaspReference: "OWASP Top 10 2021",
	}
}
