package biz

import (
	"fmt"
)

func catchError(detail string, err error) {
	if err != nil {
		fmt.Printf("%s: %v\n", detail, err)
	}
}
