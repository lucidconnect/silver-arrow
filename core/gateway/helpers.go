package gateway

import (
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/pkg/errors"
)

// creats the calldata that scopes a kernel executor to a validator
func setValidatorExecutor(sessionKey, validatorAddress, executorAddress, ownerAddress string, chain int64) ([]byte, error) {
	mode := erc4337.ENABLE_MODE
	validator, err := erc4337.InitialiseValidator(validatorAddress, executorAddress, sessionKey, mode, chain)
	if err != nil {
		return nil, err
	}

	enableData, err := validator.GetEnableData()
	if err != nil {
		return nil, err
	}

	callData, err := validator.SetExecution(enableData, ownerAddress)
	if err != nil {
		err = errors.Wrap(err, "validator.SetExecution():")
		return nil, err
	}
	return callData, nil
}
