package gribtest

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/jpmono416/grib/griblib"
)

func Test_read_section0(t *testing.T) {

	testData := griblib.Section0{
		Discipline:    2,
		Edition:       2,
		MessageLength: 4,
		Indicator:     0x47524942,
		Reserved:      6,
	}

	section0, readError := griblib.ReadSection0(toIoReader(testData))

	if readError != nil {
		t.Fatal(readError)
	}

	if testData != section0 {
		t.Error("Deserialized struct is not equal to original struct")
	}
}
func Test_read_section1(t *testing.T) {

	testData := griblib.Section1{
		LocalTablesVersion:   1,
		MasterTablesVersion:  2,
		OriginatingCenter:    3,
		OriginatingSubCenter: 4,
		ProductionStatus:     5,
		ReferenceTime: griblib.Time{
			Day:    1,
			Hour:   2,
			Minute: 3,
			Month:  4,
			Second: 5,
			Year:   6,
		},
		ReferenceTimeSignificance: 7,
		Type:                      8,
	}

	section1, readError := griblib.ReadSection1(toIoReader(testData), binary.Size(testData))

	if readError != nil {
		t.Fatal(readError)
	}

	if testData != section1 {
		t.Error("Deserialized section1 struct is not equal to original struct")
	}
}

// create a reader from a struct for testing purposes
func toIoReader(data interface{}) (reader io.Reader) {
	var binBuf bytes.Buffer

	binary.Write(&binBuf, binary.BigEndian, data)

	reader = bytes.NewReader(binBuf.Bytes())

	return

}
