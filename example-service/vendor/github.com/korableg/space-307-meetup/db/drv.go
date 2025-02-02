package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

func init() {
	fmt.Println("YAY! db driver was registered 😻")
	sql.Register("fooDB", &drv{})
	inject()
}

type drv struct {
}

func (*drv) Open(name string) (driver.Conn, error) {
	return nil, nil
}
