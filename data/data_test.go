// test data.go
package data

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

//test NewUpdateMI
func TestNewUpdateMI(t *testing.T) {
	expected := UpdateMI{
		UpdateTimeStamp: "2022-05-10 12:12:12 500",
		InstID:          "cu",
		Value:           3459.2,
	}
	actual := NewUpdateMI("2022-05-10 12:12:12 500", "cu", 3459.2)

	assert.Equal(t, expected, actual, fmt.Sprintf("Expected: %v, Actual: %v", expected, actual))
}

// benchmark NewUpdateMI @  0.3838 ns/op	   0 B/op	       0 allocs/op
// mac m1 is better
// BenchmarkNewUpdateMI-8 @ 0.3192 ns/op	       0 B/op	       0 allocs/op
func BenchmarkNewUpdateMI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUpdateMI("2022-05-10 12:12:12 500", "cu", 3459.2)
	}
}

// test goroutine read UpdateMI
func TestGoroutineReadUpdateMIwithLock(t *testing.T) {
	tmpdata := UpdateMI{
		UpdateTimeStamp: "2022-05-10 12:12:12 500",
		InstID:          "cu",
		Value:           3459.2,
	}
	mutex := sync.Mutex{}

	process := func(i int, wg *sync.WaitGroup, m *sync.Mutex) {
		m.Lock()
		fmt.Println("started Goroutine ", i)
		defer wg.Done()
		fmt.Printf("Goroutine %d ended, data is %f\n", i, tmpdata.Value)
		m.Unlock()
	}

	no := 10
	var wg sync.WaitGroup
	for i := 0; i < no; i++ {
		wg.Add(1)
		go process(i, &wg, &mutex)
	}
	wg.Wait()
	fmt.Println("All go routines finished executing")

}

// 33406	     32050 ns/op	     633 B/op	      23 allocs/op
func BenchmarkGoroutineReadUpdateMIwithLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tmpdata := UpdateMI{
			UpdateTimeStamp: "2022-05-10 12:12:12 500",
			InstID:          "cu",
			Value:           3459.2,
		}
		mutex := sync.Mutex{}
		process := func(i int, wg *sync.WaitGroup, m *sync.Mutex) {
			m.Lock()
			fmt.Println("started Goroutine ", i)
			defer wg.Done()
			fmt.Printf("Goroutine %d ended, data is %f\n", i, tmpdata.Value)
			m.Unlock()
		}
		no := 10
		var wg sync.WaitGroup
		for i := 0; i < no; i++ {
			wg.Add(1)
			go process(i, &wg, &mutex)
		}
		wg.Wait()
		fmt.Println("All go routines finished executing")
	}
}

// del lock nearly 10% speed up
// 42810	     28421 ns/op	     466 B/op	      22 allocs/op
func BenchmarkGoroutineReadUpdateMI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tmpdata := UpdateMI{
			UpdateTimeStamp: "2022-05-10 12:12:12 500",
			InstID:          "cu",
			Value:           3459.2,
		}
		process := func(i int, wg *sync.WaitGroup) {
			fmt.Println("started Goroutine ", i)
			defer wg.Done()
			fmt.Printf("Goroutine %d ended, data is %f\n", i, tmpdata.Value)
		}
		no := 10
		var wg sync.WaitGroup
		for i := 0; i < no; i++ {
			wg.Add(1)
			go process(i, &wg)
		}
		wg.Wait()
		fmt.Println("All go routines finished executing")
	}
}
