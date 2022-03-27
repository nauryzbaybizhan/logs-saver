package eventsSupport

import (
	"encoding/json"
	"fmt"
	"github.com/rome314/idkb-events/internal/events/repository"
	"io/ioutil"
	"net/http"

	"emperror.dev/errors"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

type apiProvider struct {
}

func NewApiInfoProvider() repository.IpInfoProvider {
	return &apiProvider{}
}

func (a *apiProvider) GetIpInfo(ip string) (info *eventEntities.IpInfo, err error) {
	info = &eventEntities.IpInfo{}

	url := fmt.Sprintf("http://idkb.com/%s", ip)

	client := &http.Client{}

	res, err := client.Get(url)
	if err != nil {
		err = errors.WithMessage(err, "making request")
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithMessage(err, "reading body")
		return
	}

	if err = json.Unmarshal(body, info); err != nil {
		err = errors.WithMessage(err, "unmarshalling")
		return
	}
	return info, nil

}
