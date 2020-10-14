package convert

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type City struct {
	ID        string
	Name      string
	Districts []*District
}

type District struct {
	ID    string
	Name  string
	Wards []*Ward
}

type Ward struct {
	ID   string
	Name string
}

// ConvertCSV2YML: return nil if succeed and not nil if failed
func ConvertCSV2YML(csvFile string, ymlFileOut string) error {

	// Open the file
	csvfile, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("Couldn't open the csv file: %v", err)
	}
	defer csvfile.Close()

	// Parse the file
	r := csv.NewReader(csvfile)

	var i int
	var cities []*City

	// map city with its districts
	mapCities := make(map[string]*[]*District)
	// map district with its wards
	mapDistricts := make(map[string]*[]*Ward)

	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Read record at line %d error: %v", i, err)
		}
		i++
		if i == 1 {
			// skip header
			continue
		}
		ward := &Ward{
			ID:   record[5],
			Name: record[4],
		}
		// check district existed
		if district, ok := mapDistricts[record[3]]; !ok {
			district := &District{
				ID:    record[3],
				Name:  record[2],
				Wards: []*Ward{ward},
			}
			mapDistricts[record[3]] = &district.Wards

			// check city existed
			if city, ok := mapCities[record[1]]; !ok {
				city := City{
					ID:        record[1],
					Name:      record[0],
					Districts: []*District{district},
				}
				cities = append(cities, &city)
				mapCities[record[1]] = &city.Districts
			} else {
				*city = append(*city, district)
			}
		} else {
			*district = append(*district, ward)
		}
	}

	d, err := yaml.Marshal(&cities)
	if err != nil {
		return fmt.Errorf("Marshal yml error: %v", err)
	}

	ymlFileOut = strings.Trim(ymlFileOut, " ")
	if ymlFileOut == "" {
		ymlFileOut = "data.yml"
	}
	if err := ioutil.WriteFile(ymlFileOut, d, 0644); err != nil {
		return fmt.Errorf("Write yml file failed: %v", err)
	}

	return nil
}
