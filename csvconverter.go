package csvconverter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unsafe"

	"github.com/eliben/gosax"
)

func ConvertToCsv(filename string) map[string][]string {
	counter := 0
	currentxpath := ""
	btsnamefound := false
	timezonefound := false
	pValue := ""
	version := ""
	var mrbts = map[string]map[string]string{}
	mrbtsre := regexp.MustCompile(`MRBTS-\d+`)
	btsnamere := regexp.MustCompile(`^\d+_(5GDU_)*`)
	timezonere := regexp.MustCompile(`^\(.+\) `)
	btsdure := regexp.MustCompile(`BTSDU`)
	btscure := regexp.MustCompile(`BTSCU`)
	bts5gre := regexp.MustCompile(`BTS5G`)
	SiteName := make(map[string][]string)

	scb := gosax.SaxCallbacks{
		StartElement: func(name string, attrs []string) {
			nvp := make(map[string]string)
			if name == "managedObject" {
				counter++
				for i := 0; i < len(attrs)-1; i = i + 2 {
					nvp[attrs[i]] = attrs[i+1]
				}
				currentxpath = nvp["distName"]
				if _, ok := nvp["version"]; ok {
					version = nvp["version"]
				}
			} else if name == "p" {
				for i := 0; i < len(attrs)-1; i = i + 2 {
					nvp[attrs[i]] = attrs[i+1]
				}
				if nvp["name"] == "btsName" {
					btsnamefound = true
				} else if nvp["name"] == "timeZone" {
					timezonefound = true
				}
			}
		},

		EndElement: func(name string) {
			if btsnamefound {
				btsnamefound = false
				mrbtsid := strings.Split(currentxpath, "-")[1]
				if _, ok := mrbts[mrbtsid]["timeZone"]; !ok {
					mrbts[mrbtsid] = map[string]string{}
				}
				mrbts[mrbtsid]["btsName"] = btsnamere.ReplaceAllString(pValue, "")
				mrbts[mrbtsid]["version"] = version
			}
			if timezonefound {
				timezonefound = false
				mrbtsid := strings.Split(currentxpath, "/")[0]
				if len(mrbtsre.FindString(mrbtsid)) != 0 {
					_, ok := mrbts[strings.Split(mrbtsid, "-")[1]]["btsName"]
					if !ok {
						mrbts[strings.Split(mrbtsid, "-")[1]] = map[string]string{}
					}
					mrbts[strings.Split(mrbtsid, "-")[1]]["timeZone"] = timezonere.ReplaceAllString(pValue, "")
				}
			}
		},

		Characters: func(contents string) {
			// pValue = contents
		},

		CharactersRaw: func(ch unsafe.Pointer, chlen int) {
			pValue = gosax.UnpackString(ch, chlen)
		},

		EndDocument: func() {
			for key, _ := range mrbts {
				// fmt.Printf("%s,%s,%s,%s\n", key, mrbts[key]["btsName"], mrbts[key]["timeZone"], mrbts[key]["version"])
				// 5G DU data
				if btsdure.FindStringIndex(mrbts[key]["version"]) != nil {
					duid, _ := strconv.Atoi(key[len(key)-3:])
					SiteName["5gdu"] = append(SiteName["5gdu"], fmt.Sprintf("%s-%d,%d,%s,N", key[:7], duid, duid, mrbts[key]["btsName"]))
					// 5G Data
				} else if btscure.FindStringIndex(mrbts[key]["version"]) != nil || bts5gre.FindStringIndex(mrbts[key]["version"]) != nil {
					SiteName["5g"] = append(SiteName["5g"], fmt.Sprintf("N,%s,%s,%s", mrbts[key]["btsName"], key, mrbts[key]["timeZone"]))
					// 4G
				} else {
					SiteName["4g"] = append(SiteName["4g"], fmt.Sprintf("N,%s,%s,%s", mrbts[key]["btsName"], key, mrbts[key]["timeZone"]))
				}
			}
		},
	}

	err := gosax.ParseFile(filename, scb)
	if err != nil {
		panic(err)
	}
	return SiteName
}
