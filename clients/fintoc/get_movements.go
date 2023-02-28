package fintoc

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mtavano/devkit/errors"
)

const (
	getAccountMovementsURL = "/accounts/%s/movements?%s"
)

type GetAccountMovementsRequest struct {
	MaxItems int
	Since    *time.Time
	Until    *time.Time
}

// GetAccountMovements will fetch the
func (cl *Client) GetAccountMovements(req *GetAccountMovementsRequest) ([]*Movement, *Pages, error) {
	if req == nil {
		return nil, nil, errors.New("fintoc: Client.GetAccountMovements invalid request error")
	}

	urlValues := url.Values{}
	urlValues.Add("link_token", cl.linkToken)
	urlValues.Add("per_page", fmt.Sprintf("%d", req.MaxItems))
	path := fmt.Sprintf(
		getAccountMovementsURL,
		cl.accountID,
		urlValues.Encode(),
	)

	res, err := cl.makeRequest(http.MethodGet, path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements cl.makeRequest error")
	}

	movements := make([]*Movement, 0)
	err = cl.scanBody(res, &movements)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements cl.scanBody error")
	}

	nextPageLink, lastPageLink, err := getLinkFromString(res.Header.Get("Link"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements getLinkFromString error")
	}

	return movements, &Pages{
		Next: nextPageLink,
		Last: lastPageLink,
	}, nil
}

func (cl *Client) GetAccountMovementsByPage(path string) ([]*Movement, *Pages, error) {
	urlValues := url.Values{}
	urlValues.Add("link_token", cl.linkToken)
	urlValues.Add("per_page", fmt.Sprintf("%d", 300))

	res, err := cl.makeRequest(http.MethodGet, path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements cl.makeRequest error")
	}

	movements := make([]*Movement, 0)
	err = cl.scanBody(res, &movements)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements cl.scanBody error")
	}

	nextPageLink, lastPageLink, err := getLinkFromString(res.Header.Get("Link"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "fintoc: Client.GetAccountMovements getLinkFromString error")
	}

	return movements, &Pages{
		Next: nextPageLink,
		Last: lastPageLink,
	}, nil
}

func getLinkFromString(headerString string) (string, string, error) {
	re := regexp.MustCompile("<(.*?)>")
	links := strings.Split(headerString, ",")
	currentPageLink := ""
	lastPageLink := ""

	if len(links) == 2 {
		return "", "", nil
	} else if len(links) == 3 {
		currentPage := links[1]
		linkMatch := re.FindStringSubmatch(currentPage)
		currentPageLink = linkMatch[1]

		lastPage := links[2]
		lastLinkMatch := re.FindStringSubmatch(lastPage)
		lastPageLink = lastLinkMatch[1]
	} else if len(links) == 4 {
		currentPage := links[2]
		linkMatch := re.FindStringSubmatch(currentPage)
		currentPageLink = linkMatch[1]

		lastPage := links[3]
		lastLinkMatch := re.FindStringSubmatch(lastPage)
		lastPageLink = lastLinkMatch[1]
	} else {
		return "", "", errors.New("Invalid or unsupported string parse")
	}

	return strings.Replace(currentPageLink, "https://api.fintoc.com/v1", "", -1), strings.Replace(lastPageLink, "https://api.fintoc.com/v1", "", -1), nil
}
