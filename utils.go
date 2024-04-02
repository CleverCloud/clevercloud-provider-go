package provider

import (
	"net/url"

	"github.com/sirupsen/logrus"
)

func basePath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		logrus.WithError(err).Warnf("invalid URL: '%s'", path)
		return path
	}
	return u.Path
}
