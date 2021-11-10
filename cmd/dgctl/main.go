package main

import (
	"github.com/hashed-io/document-graph-cli/cmd"

	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog = zap.NewNop()

//lint:ignore U1000 leveraged at runtime
var tracer = logging.ApplicationLogger("dgctl", "github.com/hashed-io/document-graph-cli/cmd/dgctl", &zlog)

func main() {
	cmd.Execute()
}
