package lampa_categories

import "testing"

func TestSplitCommaSeparatedParameter(t *testing.T) {
	rawQuery := "apikey=craftokey&query=Made+in+Abyss%3A+Journey%27s+Dawn+2019&categories=2000%2C5070&type=search"

	got, report := splitCommaSeparatedParameter(rawQuery, "categories")
	want := "apikey=craftokey&query=Made+in+Abyss%3A+Journey%27s+Dawn+2019&categories=5070&type=search"

	if got != want {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", want, got)
	}

	if report.status != "rewritten" {
		t.Fatalf("unexpected status: %s", report.status)
	}
}

func TestSplitCommaSeparatedParameterWithUnicodeQuery(t *testing.T) {
	rawQuery := "apikey=craftokey&query=%E5%8A%87%E5%A0%B4%E7%89%88%E7%B7%8F%E9%9B%86%E7%B7%A8%E3%80%90%E5%89%8D%E7%B7%A8%E3%80%91%E3%83%A1%E3%82%A4%E3%83%89%E3%82%A4%E3%83%B3%E3%82%A2%E3%83%93%E3%82%B9+%E6%97%85%E7%AB%8B%E3%81%A1%E3%81%AE%E5%A4%9C%E6%98%8E%E3%81%91&categories=2000%2C5070&type=search"

	got, report := splitCommaSeparatedParameter(rawQuery, "categories")
	want := "apikey=craftokey&query=%E5%8A%87%E5%A0%B4%E7%89%88%E7%B7%8F%E9%9B%86%E7%B7%A8%E3%80%90%E5%89%8D%E7%B7%A8%E3%80%91%E3%83%A1%E3%82%A4%E3%83%89%E3%82%A4%E3%83%B3%E3%82%A2%E3%83%93%E3%82%B9+%E6%97%85%E7%AB%8B%E3%81%A1%E3%81%AE%E5%A4%9C%E6%98%8E%E3%81%91&categories=5070&type=search"

	if got != want {
		t.Fatalf("unexpected query\nwant: %s\n got: %s", want, got)
	}

	if report.status != "rewritten" {
		t.Fatalf("unexpected status: %s", report.status)
	}
}
