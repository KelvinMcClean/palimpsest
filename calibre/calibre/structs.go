package calibre

type Book struct {
	CalibreID int
	Title     string
	Path      string
}

type Format struct {
	BookID int
	Format string
	Name   string
	Size   int64
}

type Author struct {
	ID   int
	Name string
	Sort string
}

type Identifier struct {
	BookID int
	Type   string
	Value  string
}
