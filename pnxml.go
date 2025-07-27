package main

import (
	"context"
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
		if adaptation.MimeType != "video/mp4" {
			continue
		}
		getMetaData := func(rep Representation, ctx context.Context) (MetaData, error) {
			size, err := getFileSizeWithContext(rep.BaseURL[0].Value, ctx)
			if err != nil {
				slog.Error("Failed to get file size", "error", err, "url", rep.BaseURL[0].Value)
				return MetaData{}, err
			}
			return MetaData{
				Bitrate:   rep.BandWidth,
				FPS:       rep.FrameRate,
				Codec:     rep.Codec,
				Width:     rep.Width,
				Height:    rep.Height,
				URL:       rep.BaseURL[0].Value,
				Size:      size,
			}, nil
		}
		concurrentRes, err := concurrentExecuteAny(getMetaData, adaptation.Representation, len(adaptation.Representation))
		if err != nil {
			slog.Error("Failed to decode representations", "error", err)
			return nil, err
		}
		pnxmlData.MetaDatas = make([]MetaData, 0, len(concurrentRes))
		for _, res := range concurrentRes {
			pnxmlData.MetaDatas = append(pnxmlData.MetaDatas, res)
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