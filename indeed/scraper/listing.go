package scraper

type Listing struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	CompanyURL  string `json:"companyUrl"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Lat         string `json:"latitude"`
	Lon         string `json:"longitude"`
}

type IListing interface {
	GetLink() string
}

func (l *Listing) GetLink() string {
	return l.Link
}

func (l *Listing) IsEmpty() bool {
	return l.Link == "" &&
		l.Title == "" &&
		l.Company == "" &&
		l.CompanyURL == "" &&
		l.Location == "" &&
		l.Description == ""
}
