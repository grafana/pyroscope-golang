package main

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"

	"runtime"
	"sync"

	"github.com/pyroscope-io/client/pyroscope"
)

//go:noinline
func work(n int) {
	// revive:disable:empty-block this is fine because this is a example app, not real production code
	for i := 0; i < n; i++ {
	}
	fmt.Printf("work\n")
	// revive:enable:empty-block
}

var m sync.Mutex

func fastFunction(c context.Context, wg *sync.WaitGroup) {
	m.Lock()
	defer m.Unlock()

	pyroscope.TagWrapper(c, pyroscope.Labels("function", "fast"), func(c context.Context) {
		work(200000000)
	})
	wg.Done()
}

func slowFunction(c context.Context, wg *sync.WaitGroup) {
	m.Lock()
	defer m.Unlock()

	// standard pprof.Do wrappers work as well
	pprof.Do(c, pprof.Labels("function", "slow"), func(c context.Context) {
		work(800000000)
	})
	wg.Done()
}

func main() {
	sa := os.Getenv("SERVER_ADDRESS")
	if sa == "" {
		sa = "https://localhost:4317"
	}
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	pyroscope.Start(pyroscope.Config{
		ApplicationName:   "simple.golang.app-new",
		ServerAddress:     sa,
		Logger:            pyroscope.StandardLogger,
		AuthToken:         os.Getenv("PYROSCOPE_AUTH_TOKEN"),
		TenantID:          os.Getenv("PYROSCOPE_TENANT_ID"),
		BasicAuthUser:     os.Getenv("PYROSCOPE_BASIC_AUTH_USER"),
		BasicAuthPassword: os.Getenv("PYROSCOPE_BASIC_AUTH_PASSWORD"),
		UseOTLP:           true,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		HTTPHeaders: map[string]string{"X-Extra-Header": "extra-header-value"},
	})

	pyroscope.TagWrapper(context.Background(), pyroscope.Labels("foo", "bar"), func(c context.Context) {
		for {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go fastFunction(c, &wg)
			go slowFunction(c, &wg)
			wg.Wait()
		}
	})

}
