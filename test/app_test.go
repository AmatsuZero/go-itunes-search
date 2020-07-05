package test

import (
	"github.com/go-pg/pg"
	. "go-itunes-search"
	"go-itunes-search/app"
	"testing"
)

// App Specific API
func TestNewApp(t *testing.T) {
	entry, _ := Lookup().ID(989673964).Result()
	app.NewDetailedApp(entry, "CN").Print()
}

// Test US Store
func TestEntry_DetailUS(t *testing.T) {
	entry, _ := Lookup().Country(US).ID(1061097588).Result()
	app.NewDetailedApp(entry, "US").Print()
}

func TestApp_Save(t *testing.T) {
	Pg := pg.Connect(&pg.Options{
		Addr:     ":5432",
		Database: "haha",
		User:     "haha",
	})

	idList := []int64{
		989673964,
		1110193350,
		1110195252,
		1110194837,
	}

	for _, id := range idList {
		appByID, err := app.NewAppByID(id, "CN")
		if err != nil {
			t.Error(err)
		}
		err = appByID.Save(Pg)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestApp_Save2(t *testing.T) {

}
