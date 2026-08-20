package test

import (
	"testing"
	"github.com/KelvinMcClean/palimpsest/normalize"
)

func Normalize(t *testing.T, input string, expected string) {
	t.Helper()
	normalized := normalize.Normalize(input)
	if normalized != expected {
		t.Errorf("Normalize(%q) = %q; want %q", input, normalized, expected)
	}
}

func NormalizeAuthorName(t *testing.T, input string, expected string) {
	t.Helper()
	normalized := normalize.NormalizeAuthorName(input)
	if normalized != expected {
		t.Errorf("NormalizeAuthorName(%q) = %q; want %q", input, normalized, expected)
	}
}

func NormalizeBase(t *testing.T, input string, expected string) {
	t.Helper()
	normalized := normalize.NormalizeBase(input)
	if normalized != expected {
		t.Errorf("NormalizeBase(%q) = %q; want %q", input, normalized, expected)
	}
}
func NormalizeBaseAuthorName(t *testing.T, input string, expected string) {
	t.Helper()
	normalized := normalize.NormalizeBaseAuthorName(input)
	if normalized != expected {
		t.Errorf("NormalizeBaseAuthorName(%q) = %q; want %q", input, normalized, expected)
	}
}

func TestBookNormalization(t *testing.T) {
	Normalize(t, "  The Great Gatsby: A Novel  ", "the great gatsby a novel")
	Normalize(t, "To Kill a Mockingbird", "to kill a mockingbird")
	Normalize(t, "Ka —  Dar Oakley in the Ruin of Ymr", "ka dar oakley in the ruin of ymr")
	Normalize(t, "Ka— Dar Oakley in the Ruin of Ymr", "ka dar oakley in the ruin of ymr")
	Normalize(t, "Ka            —Dar Oakley in the Ruin of Ymr", "ka dar oakley in the ruin of ymr")
	Normalize(t, "The Catcher in the Rye(A Coming-of-Age Novel)", "the catcher in the rye a coming of age novel")
	Normalize(t, "War Master's Gate", "war masters gate")
	Normalize(t, "Did Ye Hear Mammy Died? A Memoir", "did ye hear mammy died a memoir")
}

func TestAuthorNormalization(t *testing.T) {
	NormalizeAuthorName(t, "F. Scott Fitzgerald", "f scott fitzgerald")
	NormalizeAuthorName(t, "J. R. R. Tolkien", "jrr tolkien")
	NormalizeAuthorName(t, "Tolkien, J. R. R.", "tolkien jrr")
}
func TestNormalizeBase(t *testing.T) {
	NormalizeBase(t, "  The Great Gatsby: A Novel  ", "the great gatsby")
	NormalizeBase(t, "The Hobbit - An Unexpected Journey", "the hobbit")
	NormalizeBase(t, "To Kill a Mockingbird: A Classic Novel", "to kill a mockingbird")
	NormalizeBase(t, "1984: A Dystopian Novel", "1984")
	NormalizeBase(t, "Pride and Prejudice: A Romantic Novel", "pride and prejudice")
	NormalizeBase(t, "Ka —  Dar Oakley in the Ruin of Ymr", "ka")
	NormalizeBase(t, "The Catcher in the Rye(A Coming-of-Age Novel)", "the catcher in the rye")
	NormalizeBase(t, "The Catcher in the Rye\\A Coming-of-Age Novel)", "the catcher in the rye")
	NormalizeBase(t, "War Master's Gate", "war masters gate")
	NormalizeBase(t, "Did Ye Hear Mammy Died? A Memoir", "did ye hear mammy died")

}

func TestNormalizeBaseAuthorName(t *testing.T) {
	NormalizeBaseAuthorName(t, "F. Scott Fitzgerald", "f scott fitzgerald")
	NormalizeBaseAuthorName(t, "J. R. R. Tolkien", "jrr tolkien")
	NormalizeBaseAuthorName(t, "Tolkien, J. R. R.", "tolkien")
	NormalizeBaseAuthorName(t, "Unknown Author, P.O.", "unknown author")
	NormalizeBaseAuthorName(t, "Unknown Author - 1948", "unknown author")
	NormalizeBaseAuthorName(t, "A.B.C. Reily (1901-1936)", "abc reily")
}