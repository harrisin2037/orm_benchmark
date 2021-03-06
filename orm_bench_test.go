package orm_benchmark

import (
	"fmt"
	"io"
	"testing"

	"github.com/joho/godotenv"
)

type MatchStringOnly func(pat, str string) (bool, error)

var MatchFunc MatchStringOnly = func(pat, str string) (bool, error) { return true, nil }

func (f MatchStringOnly) MatchString(pat, str string) (bool, error)   { return f(pat, str) }
func (f MatchStringOnly) StartCPUProfile(w io.Writer) error           { return nil }
func (f MatchStringOnly) StopCPUProfile()                             {}
func (f MatchStringOnly) WriteProfileTo(string, io.Writer, int) error { return nil }
func (f MatchStringOnly) ImportPath() string                          { return "" }
func (f MatchStringOnly) StartTestLog(io.Writer)                      {}
func (f MatchStringOnly) StopTestLog() error                          { return nil }
func (f MatchStringOnly) SetPanicOnExit0(v bool)                      {}

type BenchORMInstance interface {
	ConnTest(*testing.B)
}

type ORMTestingBenchFunc func(b *testing.B, ormInstance BenchORMInstance)

var allORMs = map[string]func() BenchORMInstance{
	"pgx": PgxTestBenchORMInstance,
}

type Test struct {
	name string
	run  func(b *testing.B, ins BenchORMInstance)
}

var tests = []Test{
	{"connection test", testConn},
}

/*
	As testing main will never return as it calls os.Exit(), all benchmark tests cannot be run,
	so need to assembling the dynamic benchmark cases (values of struct testing.InternalBenchmark,
	ref: https://pkg.go.dev/testing#InternalBenchmark), call testing.Main() (ref: https://pkg.go.dev/testing#Main)
	which properly parses command line flags, creates and sets up testing.M (ref: https://pkg.go.dev/testing#M),
	and prepares and calls testing.RunBenchmarks() (ref: https://pkg.go.dev/testing#RunBenchmarks).
	This way your dynamic benchmarks will still be runnable by go test.

	Notes: testing.Main() will never return as it calls os.Exit().
	If you want to perform further logging and calculations on the benchmark results,
	you may also call testing.MainStart().Run() (which is what testing.Main() does),
	and you may pass the exit code which is returned by M.Run() to os.Exit().
*/
func TestMain(m *testing.M) {

	godotenv.Load(".env")

	benchmarks := []testing.InternalBenchmark{}

	for _, test := range tests {
		for ormName, ormInstance := range allORMs {
			bm := testing.InternalBenchmark{
				Name: fmt.Sprintf("%T[test case = %s, orm = %s]", test, test.name, ormName),
				F: func(b *testing.B) {
					b.ResetTimer()
					test.run(b, ormInstance())
				},
			}
			benchmarks = append(benchmarks, bm)
		}
	}

	testing.MainStart(MatchFunc, []testing.InternalTest{}, benchmarks, nil).Run()
}

func testConn(b *testing.B, ins BenchORMInstance) {
	ins.ConnTest(b)
}
