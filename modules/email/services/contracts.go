package services

import (
	"path/filepath"

	templateServices "github.com/AgileExecutives/serverbase/modules/templates/services"
)

// RegisterEmailContracts registers template contracts for the email module
func RegisterEmailContracts(contractRegistrar *templateServices.ContractRegistrar, tenantID uint) error {
	contractsDir := "modules/email/contracts"

	contracts := []string{
		// add contract file names if present, example:
		// "email_verification_contract.json",
	}

	for _, contractFile := range contracts {
		contractPath := filepath.Join(contractsDir, contractFile)
		if err := contractRegistrar.RegisterContractFromFile(tenantID, "email", contractPath); err != nil {
			return err
		}
	}
	return nil
}
