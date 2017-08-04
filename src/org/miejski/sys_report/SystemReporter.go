package sys_report

import (
	"time"
	"fmt"
	"runtime"
)

func StartReporting() {
	go func() {
		for {
			doEvery(5*time.Second, func(t time.Time) {
				value := runtime.NumGoroutine()
				fmt.Println(fmt.Sprintf("Current goroutines count = %d", value))
			})
		}
	}()
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
