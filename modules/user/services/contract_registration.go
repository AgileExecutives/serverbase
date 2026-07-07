package services

import (
	"path/filepath"

	templateServices "github.com/AgileExecutives/serverbase/modules/templates/services"
)

// RegisterUserContracts registers contracts for the user module
func RegisterUserContracts(contractRegistrar *templateServices.ContractRegistrar, tenantID uint) error {
	contractsDir := "modules/user/contracts"
	contracts := []string{}
	for _, contractFile := range contracts {
		contractPath := filepath.Join(contractsDir, contractFile)
		if err := contractRegistrar.RegisterContractFromFile(tenantID, "user", contractPath); err != nil {
			return err
		}
	}
	return nil
}
