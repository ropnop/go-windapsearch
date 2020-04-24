package windapsearch

import (
	"fmt"
	"strings"
)

func wrap(err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "Invalid Credentials") {
		return fmt.Errorf("Invalid Credentials")
	}
	if strings.Contains(err.Error(), "to perform this operation a successful bind must be completed") {
		return fmt.Errorf("A successful bind is required for this operation. Please provide valid credentials")
	}
	return err
}