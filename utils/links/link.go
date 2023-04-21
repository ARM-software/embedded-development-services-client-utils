package links

import "encoding/json"

// SerialiseLink serialises links into JSON strings
func SerialiseLink(links any) (str string) {
	if links == nil {
		return
	}
	linkStr, err := json.Marshal(links)
	if err != nil {
		return
	}
	str = string(linkStr)
	return
}
