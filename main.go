package main

import (
	"context"
	"github.com/lucasmbaia/baluba/core"
)

func main() {
	core.NewServerTcp(context.Background())
}
