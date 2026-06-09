package lampa_categories

import "testing"

func TestSplitCommaSeparatedParameter(t *testing.T) {
	rawQuery := "apikey=craftokey&query=Made+in+Abyss%3A+Journey%27s+Dawn+2019&categories=2000%2C5070&type=search"

	got := splitCommaSeparatedParameter(rawQuery, "categories")
	want := "apikey=craftokey&query=Made+in+Abyss%3A+Journey%27s+Dawn+2019&categories=2000&categories=5070&type=search"

	if got != want {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", want, got)
	}
}
