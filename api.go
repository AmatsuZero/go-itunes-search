package itunes_search

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/**************************************************************
* Structure Define
**************************************************************/

// iTunesResult represent iTunes response outer most structure
type iTunesResult struct {
	ResultCount int     `json:"resultCount"`
	Results     []Entry `json:"results"`
}

// Params holds iTunes API params.
// See following url for more details:
// https://affiliate.itunes.apple.com/resources/documentation/itunes-store-web-service-search-api/
type Params struct {
	url.Values
	endpoint string
}

/**************************************************************
* Enter Point
**************************************************************/

// Lookup begins the API chain
func Lookup() Params {
	return Params{make(url.Values), LookupURL}.Country(CN)
}

// Search begins the API chain with a series of search terms
func Search(terms []string) Params {
	return Params{make(url.Values), SearchURL}.Country(CN).Terms(terms)
}

// SearchOne begins the API chain with one term
func SearchOne(term string) Params {
	return Params{make(url.Values), SearchURL}.Country(CN).Term(term)
}

/**************************************************************
* Chain Method
**************************************************************/

func (p Params) SetParam(k, v string) Params {
	p.Values.Set(k, v)
	return p
}

func (p Params) Term(term string) Params {
	p.Values.Set("term", term)
	return p
}

func (p Params) Terms(terms []string) Params {
	p.Values.Set("term", strings.Join(terms, "+"))
	return p
}

func (p Params) Country(country string) Params {
	p.Values.Set("country", country)
	return p
}

func (p Params) Entity(entity string) Params {
	p.Values.Set("entity", entity)
	return p
}

func (p Params) Entities(entities []string) Params {
	p.Values["entity"] = entities
	return p
}

func (p Params) AddEntity(entity string) Params {
	p.Values.Add("entity", entity)
	return p
}

func (p Params) Media(media string) Params {
	p.Values.Set("media", media)
	return p
}

func (p Params) Medias(medias []string) Params {
	p.Values["media"] = medias
	return p
}

func (p Params) AddMedia(media string) Params {
	p.Values.Add("media", media)
	return p
}

func (p Params) ID(id int64) Params {
	p.Values.Set(ITunesID, strconv.FormatInt(id, 10))
	return p
}

func (p Params) BundleID(bundleID string) Params {
	p.Values.Set(BundleID, bundleID)
	return p
}

// App: restrict to application
func (p Params) App() Params {
	p.Values.Set("media", Software)
	return p
}

func (p Params) Limit(n int) Params {
	if n > 200 {
		n = 200
	}

	if n < 1 {
		n = 1
	}
	p.Values.Set("limit", strconv.Itoa(n))
	return p
}

/**************************************************************
* End Point
**************************************************************/

// Results will finally do the request
func (p Params) Results() ([]Entry, error) {
	res, err := http.Get(p.endpoint + p.Encode())
	if err != nil {
		return nil, err
	}

	lr := new(iTunesResult)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.NewDecoder(res.Body).Decode(lr); err != nil {
		return nil, err
	}

	if lr.ResultCount == 0 || lr.Results == nil || len(lr.Results) == 0 {
		return nil, ErrNotFound
	}

	return lr.Results, nil
}

// Result assert there's one result
func (p Params) Result() (*Entry, error) {
	if entries, err := p.Results(); err != nil {
		return nil, err
	} else {
		return &(entries[0]), nil
	}
}
