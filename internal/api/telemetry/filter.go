package telemetry

import (
	"bytes"

	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/telemetry"
)

// compileFilter recursively builds the whole filter into functions which avoid extra allocations and parsing.
func compileFilter(rawFilter *pb.LogFilter) storage.Filter {
	if rawFilter == nil {
		// Nil filter is properly recognized by the storage
		return nil
	}

	switch f := rawFilter.Filter.(type) {
	case *pb.LogFilter_Text_:
		return compileTextFilter(f.Text)
	case *pb.LogFilter_Scoped_:
		return compileScopedFilter(f.Scoped)
	case *pb.LogFilter_And_:
		return compileAndFilter(f.And)
	case *pb.LogFilter_Or_:
		return compileOrFilter(f.Or)
	case *pb.LogFilter_Not_:
		return compileNotFilter(f.Not)
	}
	return nil
}

func compileTextFilter(textFilter *pb.LogFilter_Text) storage.Filter {
	b := []byte(textFilter.GetValue())
	return func(m storage.Message, _ storage.FlatMapping) bool {
		return bytes.Contains(m, b)
	}
}

func compileScopedFilter(scopedFilter *pb.LogFilter_Scoped) storage.Filter {
	return func(_ storage.Message, fm storage.FlatMapping) bool {
		v, ok := fm(scopedFilter.GetField())
		return ok && v == scopedFilter.GetValue()
	}
}

func compileAndFilter(andFilter *pb.LogFilter_And) storage.Filter {
	a := compileFilter(andFilter.GetA())
	b := compileFilter(andFilter.GetB())
	return func(m storage.Message, fm storage.FlatMapping) bool {
		return a(m, fm) && b(m, fm)
	}
}

func compileOrFilter(orFilter *pb.LogFilter_Or) storage.Filter {
	a := compileFilter(orFilter.GetA())
	b := compileFilter(orFilter.GetB())
	return func(m storage.Message, fm storage.FlatMapping) bool {
		return a(m, fm) || b(m, fm)
	}
}

func compileNotFilter(notFilter *pb.LogFilter_Not) storage.Filter {
	a := compileFilter(notFilter.GetA())
	return func(m storage.Message, fm storage.FlatMapping) bool {
		return !a(m, fm)
	}
}
