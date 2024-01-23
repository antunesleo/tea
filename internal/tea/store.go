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

		var data1, data2 interface{}
		if err := json.Unmarshal(storedReq.RequestBody, &data1); err != nil {
			fmt.Println("Error unmarshaling rawMessage1:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if err := json.Unmarshal(underTestReq.RequestBody, &data2); err != nil {
			fmt.Println("Error unmarshaling rawMessage2:", err)
			return false, StoredRequest{} // TODO: refactor to return error
		}

		if !reflect.DeepEqual(data1, data2) {
			continue
		}

		if !map_utils.MapsEqual(storedReq.Headers, underTestReq.Headers) {
			continue
		}

		return true, *storedReq
	}
	return false, StoredRequest{}
}
