package authy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// OneTouchStatus is the type of the OneTouch statuses.
type OneTouchStatus string

var (
	// OneTouchStatusApproved is the approved status of an approval request
	OneTouchStatusApproved OneTouchStatus = "approved"

	// OneTouchStatusPending is the pending status of an approval request
	OneTouchStatusPending OneTouchStatus = "pending"

	// OneTouchStatusDenied is the denied status of an approval request
	OneTouchStatusDenied OneTouchStatus = "denied"

	// OneTouchStatusExpired is the expired status of an approval request
	OneTouchStatusExpired OneTouchStatus = "expired"
)

// ApprovalRequest is the approval request response.
type ApprovalRequest struct {
	HTTPResponse *http.Response

	Status   OneTouchStatus `json:"status"`
	UUID     string         `json:"uuid"`
	Notified bool           `json:"notified"`
}

// NewApprovalRequest returns an instance of ApprovalRequest.
func NewApprovalRequest(response *http.Response) (*ApprovalRequest, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jsonResponse := struct {
		Success         bool             `json:"success"`
		ApprovalRequest *ApprovalRequest `json:"approval_request"`
		Message         string           `json:"message"`
	}{}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	approvalRequest := jsonResponse.ApprovalRequest
	if jsonResponse.Success == false {
		return nil, errors.New(
			fmt.Sprintf(
				"Invalid approval request response: %s\n",
				jsonResponse.Message))
	}
	approvalRequest.HTTPResponse = response

	return approvalRequest, nil
}

// Valid returns true if the approval request was valid.
func (request *ApprovalRequest) Valid() bool {
	if request.HTTPResponse.StatusCode == 200 {
		return true
	}

	return false
}
