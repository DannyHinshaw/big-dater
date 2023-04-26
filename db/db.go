package db

import (
	"fmt"

	"github.com/dannyhinshaw/big-dater/dater"
)

const (
	portEST = 6543
	portPST = 7654
	portUTC = 8765

	fmtConnStr = "port=%d user=postgres sslmode=disable"
)

var (
	ConnEST = fmt.Sprintf(fmtConnStr, portEST)
	ConnPST = fmt.Sprintf(fmtConnStr, portPST)
	ConnUTC = fmt.Sprintf(fmtConnStr, portUTC)
)

type TableNamer interface {
	TableName() string
}

func TargetStr(origin dater.Origin, t TableNamer) string {
	return fmt.Sprintf("%s.%s", origin, t.TableName())
}
