package chrono

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"
)

var zones Zones
var zonesErr error
var zonesOnce sync.Once

type Zones struct {
	zones map[string]Zone
}

type Zone struct {
	loc *time.Location
}

func UTC() Zone {
	return Zone{loc: time.UTC}
}

func Local() Zone {
	localOnce.Do(initLocal)
	return Zone{loc: &localLoc}
}

func LoadZones() (Zones, error) {
	zonesOnce.Do(func() {
		zones, zonesErr = loadZones()
	})
	return zones, zonesErr
}

func LoadZone(name string) (Zone, error) {
	loc, err := time.LoadLocation(name)
	return Zone{loc: loc}, err
}

func loadZones() (Zones, error) {
	sources := zoneSources
	if env := os.Getenv("ZONEINFO"); env != "" {
		sources = append([]string{env}, sources...)
	}

	var firstErr error
	for _, source := range sources {
		zones, err := readTzDataFromDisk(source)
		if err != nil && firstErr == nil {
			firstErr = err
		} else if zones != nil {
			return *zones, nil
		}
	}

	if readEmbeddedTzData != nil {
		// TODO
		/*r, err := zip.NewReader(bytes.NewReader([]byte(embeddedTzData)), int64(len(embeddedTzData)))
		if err != nil && firstErr == nil {
			firstErr = err
		}*/
		panic("embedded timezone data is not supported by github.com/go-chrono/chrono")
	}

	return Zones{}, firstErr
}

func readTzDataFromDisk(source string) (*Zones, error) {
	if len(source) >= 6 && source[len(source)-6:] == "tzdata" {
		// time.loadTzinfoFromTzdata(file, name string) ([]byte, error)
		panic("Android is not supported github.com/go-chrono/chrono")
	}

	if filepath.Ext(source) == ".zip" {
		// TODO load from zip
		panic("zipped timezone files are not supported by github.com/go-chrono/chrono")
	}

	if _, err := os.Open(path.Join(source, "UTC")); err != nil {
		return nil, nil
	}

	var out Zones
	return &out, filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || err != nil {
			return nil
		}

		name, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		switch {
		case
			len(name) == 0,
			name[0] == '/',
			name[0] == '\\',
			strings.Contains(name, "."),
			unicode.IsLower(rune(name[0])):
			return nil
		}

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		zone, err := readTzFileData(name, buf)
		if err != nil {
			return err
		}

		out.zones[name] = zone
		return nil
	})
}

func readTzFileData(name string, data []byte) (Zone, error) {
	loc, err := time.LoadLocationFromTZData(name, data)
	return Zone{loc: loc}, err
}
