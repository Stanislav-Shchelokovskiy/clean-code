package cas

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	storage := NewStorage(70)

	var wg sync.WaitGroup

	wg.Go(func() {
		diff := 40
		for {
			value := storage.Get()

			if value.val < diff {
				return
			}

			someCalculations()

			value.val -= diff
			fmt.Println("new value", value.val)

			err := storage.Set(value)
			if err == nil {
				return
			}
		}
	})

	wg.Go(func() {
		diff := 30
		for {
			value := storage.Get()

			if value.val < diff {
				return
			}

			someCalculations()

			value.val -= diff
			fmt.Println("new value", value.val)

			err := storage.Set(value)
			if err == nil {
				return
			}
		}
	})

	wg.Wait()

	value := storage.Get()
	require.Equal(t, 0, value.val)
}

func someCalculations() {
	time.Sleep(time.Duration((rand.Intn(1000) + 1)) * time.Millisecond)
}
