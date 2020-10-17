package convert

import "testing"

func TestConvertCSV2YML_HappyCase(t *testing.T) {
	csvFile := "data.csv"
	ymlFileOut := "data.yml"
	if err := ConvertCSV2YML(csvFile, ymlFileOut); err != nil {
		t.Errorf("ConvertCSV2YML() should not error: %v", err)
	}
}

func TestConvertCSV2YML_HappyCase_FileOutYmlEmpty(t *testing.T) {
	csvFile := "data.csv"
	ymlFileOut := ""
	if err := ConvertCSV2YML(csvFile, ymlFileOut); err != nil {
		t.Errorf("ConvertCSV2YML() should not error: %v", err)
	}
}

func TestConvertCSV2YML_UnhappyCase_FileNotFound(t *testing.T) {
	csvFile := "yen.csv"
	ymlFileOut := "data.yml"
	if err := ConvertCSV2YML(csvFile, ymlFileOut); err == nil {
		t.Errorf("ConvertCSV2YML() should error")
	}
}

func TestConvertCSV2YML_UnhappyCase_FileIncorrectFormat(t *testing.T) {
	csvFile := "data.jpg"
	ymlFileOut := "data.yml"
	if err := ConvertCSV2YML(csvFile, ymlFileOut); err == nil {
		t.Errorf("ConvertCSV2YML() should error")
	}
}
