package orm_benchmark

import (
	"testing"

	"github.com/joho/godotenv"
)

type BenchORMInstance interface {
	ConnTest(*testing.B)
}

type ORMTestingBenchFunc func(b *testing.B, ormInstance BenchORMInstance)

var allORMs = map[string]func() BenchORMInstance{
	"pgx": PgxTestBenchORMInstance,
}

func TestMain(m *testing.M) {
	godotenv.Load(".env")
}

func BenchmarkORMsTestWithCases(b *testing.B) {

	type Test struct {
		run func(b *testing.B, ins BenchORMInstance)
	}

	tests := []Test{
		{testConn},
	}

	for i := 0; i < len(tests); i++ {
		benchTestORMs(b, tests[i].run)
	}
}

func ORMsTestConn(b *testing.B) {
	benchTestORMs(b, testConn)
}

func testConn(b *testing.B, ins BenchORMInstance) {
	ins.ConnTest(b)
}

func benchTestORMs(b *testing.B, f ORMTestingBenchFunc) {
	for ormName, ormInstance := range allORMs {
		b.Run(ormName, func(b *testing.B) {
			b.ResetTimer()
			f(b, ormInstance())
		})
	}
}
