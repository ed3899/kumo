package utils

import (
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func AppendBarToDependencies(wg *sync.WaitGroup, deps []*Dependency) {
	progress := mpb.New(mpb.WithWaitGroup(wg), mpb.WithWidth(100), mpb.WithRefreshRate(180*time.Millisecond))

	for _, d := range deps {
		// Config the bar
		downloadBar := progress.AddBar(0,
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),
			),
		)

		zipBar := progress.AddBar(int64(d.ContentLength),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(d.Name),
				decor.Counters(decor.SizeB1024(0), " % .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Percentage(decor.WCSyncSpace),
					"done",
				),
			),
		)

		// Assign the bar to the dependency
		d.DownloadBar = downloadBar
		d.ZipBar = zipBar

	}
}

