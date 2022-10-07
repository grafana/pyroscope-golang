package pyroscope

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/pyroscope-io/client/upstream/remote"
)

type Config struct {
	ApplicationName string // e.g backend.purchases
	Tags            map[string]string
	ServerAddress   string // e.g http://pyroscope.services.internal:4040
	AuthToken       string // specify this token when using pyroscope cloud
	SampleRate      uint32 // todo this one is not used
	Logger          Logger
	ProfileTypes    []ProfileType
	DisableGCRuns   bool // this will disable automatic runtime.GC runs between getting the heap profiles
	ManualFlush     bool // disable periodic profiler flush every 10 seconds. Library user is expected to call Flush periodically, manually.
}

type Profiler struct {
	session  *Session
	uploader *remote.Remote
}

// Start starts continuously profiling go code
func Start(cfg Config) (*Profiler, error) {
	if len(cfg.ProfileTypes) == 0 {
		cfg.ProfileTypes = DefaultProfileTypes
	}
	if cfg.SampleRate == 0 {
		cfg.SampleRate = DefaultSampleRate
	}
	if cfg.Logger == nil {
		cfg.Logger = noopLogger
	}

	// Override the address to use when the environment variable is defined.
	// This is useful to support adhoc push ingestion.
	if address, ok := os.LookupEnv("PYROSCOPE_ADHOC_SERVER_ADDRESS"); ok {
		cfg.ServerAddress = address
	}

	rc := remote.Config{
		AuthToken: cfg.AuthToken,
		Address:   cfg.ServerAddress,
		Threads:   4,
		Timeout:   30 * time.Second,
		Logger:    cfg.Logger,
	}
	uploader, err := remote.NewRemote(rc)
	if err != nil {
		return nil, err
	}

	sc := SessionConfig{
		Upstream:       uploader,
		Logger:         cfg.Logger,
		AppName:        cfg.ApplicationName,
		Tags:           cfg.Tags,
		ProfilingTypes: cfg.ProfileTypes,
		DisableGCRuns:  cfg.DisableGCRuns,
		SampleRate:     cfg.SampleRate,
		UploadRate:     10 * time.Second,
	}

	cfg.Logger.Infof("starting profiling session:")
	cfg.Logger.Infof("  AppName:        %+v", sc.AppName)
	cfg.Logger.Infof("  Tags:           %+v", sc.Tags)
	cfg.Logger.Infof("  ProfilingTypes: %+v", sc.ProfilingTypes)
	cfg.Logger.Infof("  DisableGCRuns:  %+v", sc.DisableGCRuns)
	cfg.Logger.Infof("  SampleRate:     %+v", sc.SampleRate)
	cfg.Logger.Infof("  UploadRate:     %+v", sc.UploadRate)
	s, err := NewSession(sc)
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}
	uploader.Start()
	if err = s.start(); err != nil {
		return nil, fmt.Errorf("start session: %w", err)
	}

	return &Profiler{session: s, uploader: uploader}, nil
}

// Stop stops continious profiling session and uploads the remaining profiling data
func (p *Profiler) Stop() {
	p.session.Stop()
	p.uploader.Stop()
}

func (p *Profiler) Flush(syncUpload bool) {
	p.session.Flush(syncUpload)
}

type LabelSet = pprof.LabelSet

var Labels = pprof.Labels

func TagWrapper(ctx context.Context, labels LabelSet, cb func(context.Context)) {
	pprof.Do(ctx, labels, func(c context.Context) { cb(c) })
}
