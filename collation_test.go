package art_test

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/Clement-Jean/go-art"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func TestCollateInsertVeryLong(t *testing.T) {
	tr := art.NewCollationSortedTree[[]byte, int]()
	key1 := []byte{16, 0, 0, 0, 7, 10, 0, 0, 0, 2, 17, 10, 0, 0, 0, 120, 10, 0, 0, 0, 120, 10, 0,
		0, 0, 216, 10, 0, 0, 0, 202, 10, 0, 0, 0, 194, 10, 0, 0, 0, 224, 10, 0, 0, 0,
		230, 10, 0, 0, 0, 210, 10, 0, 0, 0, 206, 10, 0, 0, 0, 208, 10, 0, 0, 0, 232,
		10, 0, 0, 0, 124, 10, 0, 0, 0, 124, 2, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 182, 92, 236, 147, 171, 101, 150, 195, 112, 185, 218, 108, 246,
		139, 164, 234, 195, 58, 177, 0, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 180, 93, 46, 151, 9, 212, 190, 95, 102, 178, 217, 44, 178, 235,
		29, 190, 218, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213, 251, 173, 202,
		211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0, 8, 18, 180, 93,
		46, 151, 9, 212, 190, 95, 102, 183, 219, 229, 214, 59, 125, 182, 71,
		108, 180, 220, 238, 150, 91, 117, 150, 201, 84, 183, 128, 8, 16, 0, 0,
		0, 2, 12, 185, 89, 44, 213, 251, 173, 202, 211, 95, 185, 89, 110, 118,
		251, 173, 202, 199, 101, 0, 8, 18, 180, 93, 46, 151, 9, 212, 190, 95,
		108, 176, 217, 47, 50, 219, 61, 134, 207, 97, 151, 88, 237, 246, 208,
		8, 18, 255, 255, 255, 219, 191, 198, 134, 5, 223, 212, 72, 44, 208,
		250, 180, 14, 1, 0, 0, 8}
	key2 := []byte{16, 0, 0, 0, 7, 10, 0, 0, 0, 2, 17, 10, 0, 0, 0, 120, 10, 0, 0, 0, 120, 10, 0,
		0, 0, 216, 10, 0, 0, 0, 202, 10, 0, 0, 0, 194, 10, 0, 0, 0, 224, 10, 0, 0, 0,
		230, 10, 0, 0, 0, 210, 10, 0, 0, 0, 206, 10, 0, 0, 0, 208, 10, 0, 0, 0, 232,
		10, 0, 0, 0, 124, 10, 0, 0, 0, 124, 2, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 182, 92, 236, 147, 171, 101, 150, 195, 112, 185, 218, 108, 246,
		139, 164, 234, 195, 58, 177, 0, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213,
		251, 173, 202, 211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0,
		8, 18, 180, 93, 46, 151, 9, 212, 190, 95, 102, 178, 217, 44, 178, 235,
		29, 190, 218, 8, 16, 0, 0, 0, 2, 12, 185, 89, 44, 213, 251, 173, 202,
		211, 95, 185, 89, 110, 118, 251, 173, 202, 199, 101, 0, 8, 18, 180, 93,
		46, 151, 9, 212, 190, 95, 102, 183, 219, 229, 214, 59, 125, 182, 71,
		108, 180, 220, 238, 150, 91, 117, 150, 201, 84, 183, 128, 8, 16, 0, 0,
		0, 3, 12, 185, 89, 44, 213, 251, 133, 178, 195, 105, 183, 87, 237, 150,
		155, 165, 150, 229, 97, 182, 0, 8, 18, 161, 91, 239, 50, 10, 61, 150,
		223, 114, 179, 217, 64, 8, 12, 186, 219, 172, 150, 91, 53, 166, 221,
		101, 178, 0, 8, 18, 255, 255, 255, 219, 191, 198, 134, 5, 208, 212, 72,
		44, 208, 250, 180, 14, 1, 0, 0, 8}

	tr.Insert(key1, 299)
	tr.Insert(key2, 302)
}

func TestCollateInsert(t *testing.T) {
	tr := art.NewCollationSortedTree[string, int]()

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestCollateInsertSearchWords(t *testing.T) {
	tr := art.NewCollationSortedTree[string, int]()

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("error seeking file: %s", err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := tr.Search(line); !ok {
			t.Fatalf("word %s not found", line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestCollateInsertSearchUUIDs(t *testing.T) {
	tr := art.NewCollationSortedTree[string, int]()

	file, err := os.Open("testdata/uuid.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("error seeking file: %s", err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := tr.Search(line); !ok {
			t.Fatalf("word %s not found", line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestCollateInsertDeleteWords(t *testing.T) {
	tr := art.NewCollationSortedTree[string, int]()

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Insert(line, len(line))
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("error seeking file: %s", err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tr.Delete(line)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := tr.Search(line); ok {
			t.Fatalf("word %s found", line)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading file: %s", err)
	}
}

func TestCollateMinimum(t *testing.T) {
	var words []string
	tr := art.NewCollationSortedTree[string, int]()
	c := collate.New(language.Und)

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	for _, word := range words {
		tr.Insert(word, len(word))
	}

	var res []string
	for key, _ := range tr.All() {
		res = append(res, key)
	}

	c.SortStrings(words)

	if key, _, ok := tr.Minimum(); ok == true {
		if key != words[0] {
			t.Fatalf("expected word %q, got %q", words[0], key)
		}
	} else {
		t.Fatal("minimum not found")
	}
}

func TestCollateMaximum(t *testing.T) {
	var words []string
	tr := art.NewCollationSortedTree[string, int]()
	c := collate.New(language.Und)

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	for _, word := range words {
		tr.Insert(word, len(word))
	}

	var res []string
	for key, _ := range tr.All() {
		res = append(res, key)
	}

	c.SortStrings(words)

	if key, _, ok := tr.Maximum(); ok == true {
		if key != words[len(words)-1] {
			t.Fatalf("expected word %q, got %q", words[len(words)-1], key)
		}
	} else {
		t.Fatal("maximum not found")
	}
}

func TestCollateIterAll(t *testing.T) {
	var words []string
	tr := art.NewCollationSortedTree[string, int]()
	c := collate.New(language.Und)

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	for _, word := range words {
		tr.Insert(word, len(word))
	}

	var res []string
	for key, _ := range tr.All() {
		res = append(res, key)
	}

	c.SortStrings(words)

	if !slices.Equal(words, res) {
		t.Fatal("slices are not the same")
	}
}

func TestCollateIterBackward(t *testing.T) {
	var words []string
	tr := art.NewCollationSortedTree[string, int]()
	c := collate.New(language.Und)

	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}

	for _, word := range words {
		tr.Insert(word, len(word))
	}

	var res []string
	for key, _ := range tr.Backward() {
		res = append(res, key)
	}

	c.SortStrings(words)
	slices.Reverse(words)

	if !slices.Equal(words, res) {
		t.Fatal("slices are not the same")
	}
}

func TestCollateInsertRune(t *testing.T) {
	tr := art.NewCollationSortedTree[[]rune, int]()
	key := []rune("你好")

	tr.Insert(key, 1)
	if _, ok := tr.Search(key); !ok {
		t.Fatal("key was not found")
	}
}

func TestCollateAll(t *testing.T) {
	c := collate.New(language.English, collate.Numeric)
	tr := art.NewCollationSortedTree(art.WithCollator[string, int](c))
	expected := []string{"1", "11", "9"}

	c.SortStrings(expected)
	tr.Insert("1", 1)
	tr.Insert("11", 1)
	tr.Insert("9", 1)

	var got []string
	for k, _ := range tr.All() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCollateBackward(t *testing.T) {
	c := collate.New(language.English, collate.Numeric)
	tr := art.NewCollationSortedTree(art.WithCollator[string, int](c))
	expected := []string{"1", "11", "9"}

	c.SortStrings(expected)
	slices.Reverse(expected)
	tr.Insert("1", 1)
	tr.Insert("11", 1)
	tr.Insert("9", 1)

	var got []string
	for k, _ := range tr.Backward() {
		got = append(got, k)
	}

	if !slices.Equal(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestCollatePrefix(t *testing.T) {
	tests := []struct {
		prefix         string
		keys, expected []string
	}{
		{
			"empty",
			[]string{},
			[]string{},
		},
		{
			"api",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "api.foo", "api"},
		},
		{
			"a",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
		}, {
			"b",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{},
		},
		{
			"api.",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "api.foo"},
		},
		{
			"api.foo.bar",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{"api.foo.bar"},
		},
		{
			"api.end",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{},
		},
		{
			"",
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
			[]string{"api.foo.bar", "api.foo.baz", "api.foe.fum", "abc.123.456", "api.foo", "api"},
		},
		{
			"this:key:has",
			[]string{
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:common:prefix:1",
			},
			[]string{
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:common:prefix:1",
			},
		},
		{
			"ele",
			[]string{"elector", "electibles", "elect", "electible"},
			[]string{"elector", "electibles", "elect", "electible"},
		},
		{
			"long.api.url.v1",
			[]string{"long.api.url.v1.foo", "long.api.url.v1.bar", "long.api.url.v2.foo"},
			[]string{"long.api.url.v1.foo", "long.api.url.v1.bar"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("prefix-%s", tt.prefix), func(t *testing.T) {
			tree := art.NewCollationSortedTree[string, string]()

			for _, k := range tt.keys {
				tree.Insert(k, k)
			}

			actual := []string{}
			for key, _ := range tree.Prefix(tt.prefix) {
				actual = append(actual, key)
			}

			slices.Sort(tt.expected)

			if !slices.Equal(tt.expected, actual) {
				t.Fatalf("slices are not the same!")
			}
		})
	}
}

func TestCollationRange(t *testing.T) {
	tests := []struct {
		name, start, end string
		keys, expected   []string
	}{
		{
			name:  "empty_start",
			start: "", end: "bc",
			keys:     []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
			expected: []string{"aa", "ab", "ac", "ba", "bb", "bc"},
		},
		{
			name:  "empty_end",
			start: "bc", end: "",
			keys:     []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
			expected: []string{"bc", "ca", "cb", "cc"},
		},
		{
			name:  "empty_start_end",
			start: "", end: "",
			keys:     []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
			expected: []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
		},
		{
			name:  "simple",
			start: "ba", end: "bc",
			keys:     []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
			expected: []string{"ba", "bb", "bc"},
		},
		{
			name:  "simple_start_end_no_common_prefix",
			start: "ba", end: "cc",
			keys:     []string{"aa", "ab", "ac", "ba", "bb", "bc", "ca", "cb", "cc"},
			expected: []string{"ba", "bb", "bc", "ca", "cb", "cc"},
		},
		{
			name:  "this:key:has:a:long:common:prefix",
			start: "this:key:has:a:long:common:prefix:3",
			end:   "this:key:has:a:long:common:prefix:5",
			keys: []string{
				"this:key:has:a:long:common:prefix:1",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:common:prefix:3",
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:common:prefix:8",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:common:prefix:10",
			},
			expected: []string{
				"this:key:has:a:long:common:prefix:3",
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
			},
		},
		{
			name:  "this:key:has:a:long",
			start: "this:key:has:a:long:common:prefix:3",
			end:   "this:key:has:a:long:common:prefix:5",
			keys: []string{
				"this:key:has:a:long:common:prefix:1",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:prefix:8",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:prefix:10",
			},
			expected: []string{
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
			},
		},
		{
			name:  "this:key:has:a:long_end",
			start: "this:key:has:a:long:common:prefix:3",
			end:   "this:key:has:a:long:prefix:8",
			keys: []string{
				"this:key:has:a:long:common:prefix:1",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:prefix:8",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:prefix:10",
			},
			expected: []string{
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:prefix:10",
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:prefix:8",
			},
		},
		{
			name:  "this:key:has:a:start>end",
			start: "this:key:has:a:long:prefix:10",
			end:   "this:key:has:a:long:common:prefix:3",
			keys: []string{
				"this:key:has:a:long:common:prefix:1",
				"this:key:has:a:long:common:prefix:2",
				"this:key:has:a:long:prefix:3",
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:prefix:8",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:prefix:10",
			},
			expected: []string{
				"this:key:has:a:long:common:prefix:4",
				"this:key:has:a:long:common:prefix:5",
				"this:key:has:a:long:common:prefix:6",
				"this:key:has:a:long:common:prefix:7",
				"this:key:has:a:long:common:prefix:9",
				"this:key:has:a:long:prefix:10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("range-%s", tt.name), func(t *testing.T) {
			tr := art.NewCollationSortedTree[string, int]()

			for _, key := range tt.keys {
				tr.Insert(key, len(key))
			}

			var res []string
			for key, _ := range tr.Range(tt.start, tt.end) {
				res = append(res, key)
			}

			if !slices.Equal(tt.expected, res) {
				fmt.Printf("%v %v\n", tt.expected, res)
				t.Fatal("slices are not the same")
			}
		})
	}
}
