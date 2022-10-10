package logs

import (
	"github.com/cirruslabs/echelon"
	"github.com/cirruslabs/echelon/renderers"
	"github.com/cirruslabs/echelon/renderers/config"
	"io"
	"os"
)

func GetLogger(verbose bool, plain bool, logWriter io.Writer) (*echelon.Logger, func()) {
	var defaultSimpleRenderer = renderers.NewSimpleRenderer(logWriter, nil)
	var renderer echelon.LogRendered = defaultSimpleRenderer

	cancelFunc := func() {}

	if !plain {
		c := config.NewDefaultEmojiRenderingConfig()
		c.SuccessStatus = "✓"
		c.FailureStatus = "✗"
		c.ProgressIndicatorFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		interactiveRenderer := renderers.NewInteractiveRenderer(os.Stdout, c)
		go interactiveRenderer.StartDrawing()
		cancelFunc = func() {
			interactiveRenderer.StopDrawing()
		}
		renderer = interactiveRenderer
	}

	logger := echelon.NewLogger(echelon.InfoLevel, renderer)

	if verbose {
		logger = echelon.NewLogger(echelon.DebugLevel, renderer)
	}

	return logger, cancelFunc
}
