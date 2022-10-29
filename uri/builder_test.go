package uri

import "testing"

func TestBuilder(t *testing.T) {
	b := NewURLBuilder()

	b.SetURL("http://www.google.com")

	b.AddParam("q", "golang")

	if b.Build() != "http://www.google.com?q=golang" {
		t.Error("URL is not correct")
	}
} 
