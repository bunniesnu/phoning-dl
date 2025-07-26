package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
)

func getPNXML(apiKey, accessToken string, id int) (*PNXMLInfo, error) {
	endpoint := "/fan/v1.0/lives/" + strconv.Itoa(id) + "/play-info-v3"
	params := map[string]string{
		"countryCode": "KR",
	}
	res, err := Phoning("GET", apiKey, accessToken, endpoint, params)
	if err != nil {
		return nil, err
	}

    raw, ok := res["data"].(map[string]any)["lipPlayback"]
    if !ok {
        return nil, fmt.Errorf("missing lipPlayback field")
    }
    lipJSON, ok := raw.(string)
    if !ok {
        return nil, fmt.Errorf("lipPlayback is not a JSON string (got %T)", raw)
    }

    var lipMap map[string]any
    if err := json.Unmarshal([]byte(lipJSON), &lipMap); err != nil {
        return nil, fmt.Errorf("failed to parse lipPlayback JSON: %w", err)
    }

	pnxmlData := new(PNXMLInfo)

	lipBytes, err := json.Marshal(lipMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lipMap: %w", err)
	}
	var pnxmlJSON PNXMLJSON
	if err := json.Unmarshal(lipBytes, &pnxmlJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into PNXMLJSON: %w", err)
	}
	
	pnxmlData.MaxHeight = int(pnxmlJSON.Period[0].AdaptationSet[0].MaxHeight)
	for _, adaptation := range pnxmlJSON.Period[0].AdaptationSet {
		for _, representation := range adaptation.Representation {
			metaData := MetaData{
				Bitrate: int(representation.BandWidth),
				FPS:     representation.FrameRate,
				Codec:   representation.Codec,
				Width:   int(representation.Width),
				Height:  int(representation.Height),
				URL:     representation.BaseURL[0].Value,
			}
			pnxmlData.MetaDatas = append(pnxmlData.MetaDatas, metaData)
		}
	}
	for _, sup := range pnxmlJSON.Period[0].SupplementalProperty[0].Any {
		if len(sup.Cover) > 0 {
			pnxmlData.ImageURL = sup.Cover[0].Value
		}
	}
	if pnxmlData.ImageURL == "" {
		slog.Warn("No cover image found in PNXML data", "liveId", id)
	}

	return pnxmlData, nil
}