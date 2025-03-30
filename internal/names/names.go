package names

import (
	"fmt"
	"github.com/sqc157400661/jobx/config"
)

func PreLockKey(uid string) string {
	return fmt.Sprintf("%s%s", config.PreLockPrefix, uid)
}
