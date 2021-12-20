# csvconverter
Converter of xml data to csv written in go

`csvconverter.ConvertToCsv` takes xml file name as an argument and returns map of string arrays.  
There are 3 keys in the map:  
- 5gdu  
- 5g  
- 4g

Example of running this  

```
package main

import (
	"fmt"

	"github.com/DmitryZayats/csvconverter"
)

func main() {
	var SiteName = make(map[string][]string)
	SiteName = csvconverter.ConvertToCsv("/home/dmitry/Dev/Go/sax/CM_Export-12-13-2021.xml")
	fmt.Printf("Found %d 4G sites\n", len(SiteName["4g"]))
	fmt.Printf("Found %d 5G sites\n", len(SiteName["5g"]))
	fmt.Printf("Found %d 5GDUs\n", len(SiteName["5gdu"]))
}
```
Once we have a map of string arrays - we can iterate over them and write data to text files or whatever is the requirement.
