package threadsafeset

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

var TheWords = func() []string {
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic("error opening words file: " + err.Error())
	}
	defer f.Close()
	results := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); len(line) > 0 {
			results = append(results, line)
		}
	}
	slices.Sort(results)
	return results
}()

func RandomWord() string { return TheWords[rand.Intn(len(TheWords))] }

func TestThreadsafeSet(t *testing.T) {
	set := New[string]()
	if want, got := 0, set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	}
	for _, word := range TheWords {
		if added := set.Add(word); !added {
			t.Fatalf("error adding %s...\n", word)
		}
	}
	if want, got := len(TheWords), set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	} else if want, got := TheWords, set.Slice(); !slices.Equal(want, got) {
		t.Fatalf("error: wanted %v; got %v\n", want, got)
	}
	for _, word := range TheWords {
		if contains := set.Contains(word); !contains {
			t.Fatalf("error: didn't find %s in set...\n", word)
		} else if dropped := set.Drop(word); !dropped {
			t.Fatalf("error dropping %s...\n", word)
		}
	}
	if want, got := 0, set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	}
	set = NewFromSlice(TheWords)
	if want, got := len(TheWords), set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	}
	if contains := set.ContainsSlice(TheWords); len(contains) != len(TheWords) {
		t.Fatalf("error: expected %d bools; got %d\n", len(TheWords), len(contains))
	} else {
		for i, contains := range contains {
			if !contains {
				t.Fatalf("error: %s (index %d) not found\n", TheWords[i], i)
			}
		}
	}
	if dropped := set.DropSlice(TheWords); len(dropped) != len(TheWords) {
		t.Fatalf("error: wanted %d bools; got %d\n", len(TheWords), len(dropped))
	} else {
		for i, dropped := range dropped {
			if !dropped {
				t.Fatalf("error dropping %s (index %d)\n", TheWords[i], i)
			}
		}
	}
	if want, got := 0, set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	} else if added := set.AddSlice(TheWords); len(added) != len(TheWords) {
		t.Fatalf("error: wanted %d bools; got %d\n", len(TheWords), len(added))
	} else {
		for i, added := range added {
			if !added {
				t.Fatalf("error adding %s (index %d)\n", TheWords[i], i)
			}
		}
	}
	if want, got := len(TheWords), set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	}
	set.Reset()
	if want, got := 0, set.Len(); want != got {
		t.Fatalf("error: wanted %d keys; got %d\n", want, got)
	}
}

func BenchmarkAdd(b *testing.B) {
	set := New[string]()
	for i := 0; i < b.N; i++ {
		set.Add(RandomWord())
	}
}

func BenchmarkDrop(b *testing.B) {
	set := NewFromSlice(TheWords)
	for i := 0; i < b.N; i++ {
		set.Drop(RandomWord())
	}
}

func BenchmarkContains(b *testing.B) {
	set := NewFromSlice(TheWords)
	for i := 0; i < b.N; i++ {
		set.Contains(RandomWord())
	}
}
