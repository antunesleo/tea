package tea

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/antunesleo/tea/internal/map_utils"
)

type ExpectedRequest struct {
	Body    json.RawMessage
	Headers map[string]string
	Method  string
	URL     string
}

type WantedResponse struct {
	Body       json.RawMessage
	StatusCode int
}

type StoredRequest struct {
	ExpectedRequest ExpectedRequest
	WantedResponse  WantedResponse
}

type RequestsStore struct {
	storedRequests []*StoredRequest
}

func (rs *RequestsStore) Register(storedRequest *StoredRequest) {
	rs.storedRequests = append(rs.storedRequests, storedRequest)
}

func (rs *RequestsStore) MatchRequest(underTestReq *UnderTestRequest) (bool, StoredRequest) {
	for _, storedReq := range rs.storedRequests {
		if storedReq.ExpectedRequest.Method != underTestReq.Method {
			continue
		}

		if len(storedReq.ExpectedRequest.Body) != 0 {
			var storedBody, underTestBody interface{}
			if err := json.Unmarshal(storedReq.ExpectedRequest.Body, &storedBody); err != nil {
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
		}

		if !map_utils.MapsEqual(storedReq.ExpectedRequest.Headers, underTestReq.Headers) {
			continue
		}

		if storedReq.ExpectedRequest.Method != underTestReq.Method {
			continue
		}

		if storedReq.ExpectedRequest.URL != underTestReq.URL {
			continue
		}

		return true, *storedReq
	}
	return false, StoredRequest{}
}

func NewRequestsStore() *RequestsStore {
	return &RequestsStore{}
}
