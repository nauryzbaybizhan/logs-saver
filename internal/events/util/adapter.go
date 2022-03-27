package util

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"emperror.dev/errors"
	eventEntities "github.com/rome314/idkb-events/internal/events/entities"
)

func RawToEvent(input eventEntities.RawEvent) (event *eventEntities.Event, err error) {
	if input.UserId == "" {
		err = errors.New("user id not provided")
		return
	}

	if input.ApiKey == "" {
		err = errors.New("api key not provided")
		return
	}

	if input.UserAgent == "" {
		err = errors.New("user agent not provided")
		return
	}

	decodedUrl, err := url.QueryUnescape(input.Url)
	if err != nil {
		err = errors.WithMessage(err, "decoding input url")
		return
	}

	u, err := url.ParseRequestURI(decodedUrl)
	if err != nil {
		err = errors.WithMessage(err, "validating input url")
		return
	}

	splitedIp := strings.Split(input.Ip, "/")
	rawIp := splitedIp[0]
	if len(splitedIp) == 2 {
		rawIp = splitedIp[1]
	}

	ip := net.ParseIP(rawIp)
	if ip == nil {
		err = errors.Errorf("invalid ip provided: %s", rawIp)
		return
	}

	// timestamp, err := strconv.ParseFloat(input.RequestTime, 64)
	// if err != nil {
	// 	err = errors.WithMessage(err, "invalid request time provided")
	// 	return
	// }

	requestTime, err := parseTime(int64(input.RequestTime))
	if err != nil {
		err = errors.WithMessage(err, "invalid request time provided")
		return
	}

	event = &eventEntities.Event{
		Url:         u.String(),
		UserId:      input.UserId,
		Ip:          ip,
		ApiKey:      input.ApiKey,
		UserAgent:   input.UserAgent,
		RequestTime: requestTime,
	}
	return event, nil
}

func parseTime(timestamp int64) (time.Time, error) {
	if timestamp <= 0 {
		return time.Time{}, errors.New("less or equal 0")
	}

	digits := getDigitsCount(timestamp)
	trinaries := digits / 3

	switch trinaries {
	case 3:
		return time.Unix(timestamp, 0), nil
	case 4:
		return time.UnixMilli(timestamp), nil
	case 5:
		return time.UnixMicro(timestamp), nil
	case 6:
		return time.Unix(0, timestamp), nil
	default:
		return time.Time{}, errors.New("not enough precision")
	}
}
func getDigitsCount(n int64) int {
	return len(strconv.Itoa(int(n)))
}
