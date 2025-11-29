package cache

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	base := filepath.Dir(file)
	return filepath.Join(base, "..", "..", "..", "..")
}

func GetCodeCSVPath() string {
	return filepath.Join(repoRoot(), "python", "adata", "adata", "stock", "cache", "code.csv")
}

func GetCalendarCSVPath(year int) string {
	return filepath.Join(repoRoot(), "python", "adata", "adata", "stock", "cache", "calendar", strconv.Itoa(year)+".csv")
}

func GetAllConceptCodeEastCSVPath() string {
	return filepath.Join(repoRoot(), "python", "adata", "adata", "stock", "info", "cache", "all_concept_code_east.csv")
}

func GetAllCodeCSVPath() string {
	return filepath.Join(repoRoot(), "python", "adata", "adata", "stock", "info", "cache", "all_code.csv")
}

func CalendarYears() []int {
	return []int{2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025}
}

func LoadIndexCodeRelTHS() (map[string]string, error) {
	p := filepath.Join(repoRoot(), "python", "adata", "adata", "stock", "cache", "index_code_rel_ths.py")
	b, err := os.ReadFile(p)
	if err != nil {
		return map[string]string{}, err
	}
	text := string(b)
	re := regexp.MustCompile(`"([^"]+)"\s*:\s*"([^"]+)"`)
	out := map[string]string{}
	for _, m := range re.FindAllStringSubmatch(text, -1) {
		if len(m) == 3 {
			out[m[1]] = m[2]
		}
	}
	return out, nil
}
