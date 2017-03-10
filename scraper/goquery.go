package scraper

/*
// NewDocumentFromString returns a goquery.Document from a string.
func NewDocumentFromString(html string) (*goquery.Document, error) {
	// create a reader
	reader := strings.NewReader(html)
	return goquery.NewDocumentFromReader(reader)
}

// NewDocumentFromFile returns a goquery.Document from a file.
func NewDocumentFromFile(path string) (*goquery.Document, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// close file on exit
	defer file.Close()
	// create a buffered reader
	reader := bufio.NewReader(file)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ExtractText returns the data of the Text nodes
// of the selection
func ExtractText(s *goquery.Selection) string {
	var buf bytes.Buffer

	n := s.Nodes[0].FirstChild
	for {
		if n == nil {
			break
		}
		if n.Type == html.TextNode {
			s := strings.TrimSpace(n.Data)
			buf.WriteString(s)
		}
		n = n.NextSibling

	}
	return buf.String()
}

// GetFirstAcestreamLink returns the first Acestream link of the page
func GetFirstAcestreamLink(url string) (string, error) {
	// Load the URL
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	//TODO: check other http.Status (Cache hit ...)
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	// Parse the HTML into nodes
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	link, ok := doc.Find("[href^='acestream://']").Attr("href")
	if !ok {
		return "", errors.New("Acestream link not found")
	}
	return link, nil
}
*/
