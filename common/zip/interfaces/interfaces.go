package interfaces

import (
	"github.com/vbauerster/mpb/v8"
)

type ProgressBar interface {
	IncrBy(n int)
}

type MultiProgressBar interface {
	AddBar(total int64, options ...mpb.BarOption) *mpb.Bar
}

type Removable interface {
	Remove() error
}

type Retrivable interface {
	GetName() string
	GetPath() string
}

type Downloadable interface {
	SetDownloadBar(MultiProgressBar)
	Download(chan<- int) error
	IncrementDownloadBar(int)
}

type Extractable interface {
	SetExtractionBar(MultiProgressBar, int64)
	ExtractTo(string, chan<- int) error
	IncrementExtractionBar(int)
}

type ZipI interface {
	Retrivable
	Downloadable
	Extractable
	Removable
}