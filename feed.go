package main

type Feed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Link Link `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func (f *Feed) GetLinks() []string {
	links := make([]string, len(f.Entries))
	for i, entry := range f.Entries {
		links[i] = entry.Link.Href
	}
	return links
}
