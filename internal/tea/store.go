package tea

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/antunesleo/tea/internal/map_utils"
)

type StoredRequest struct {
	RequestBody  json.RawMessage
	Headers      map[string]string
	ResponseBody json.RawMessage
	Method       string
	URL          string
}

type RequestsStore struct {
	storedRequests []*StoredRequest
}

func (rs *RequestsStore) Register(storedRequest *StoredRequest) {
	rs.storedRequests = append(rs.storedRequests, storedRequest)
}

func (rs *RequestsStore) MatchRequest(underTestReq *UnderTestRequest) (bool, StoredRequest) {
	for _, storedReq := range rs.storedRequests {
		if storedReq.Method != underTestReq.Method {
			continue
		}

		var storedBody, underTestBody interface{}
		if err := json.Unmarshal(storedReq.RequestBody, &storedBody); err != nil {
			fmt.Println("Error unmarshaling rawMessage1:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if err := json.Unmarshal(underTestReq.RequestBody, &underTestBody); err != nil {
			fmt.Println("Error unmarshaling rawMessage2:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if !reflect.DeepEqual(storedBody, underTestBody) {
			continue
		}

		if !map_utils.MapsEqual(storedReq.Headers, underTestReq.Headers) {
			continue
		}

		if storedReq.Method != underTestReq.Method {
			continue
		}

		if storedReq.URL != underTestReq.URL {
			continue
		}

		return true, *storedReq
	}
	return false, StoredRequest{}
}
