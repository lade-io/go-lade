package lade

import (
	"fmt"
	"time"
)

var _ LogService = new(LogClient)

type LogClient struct {
	client *Client
}

type LogStreamOpts struct {
	Follow bool      `qstring:"follow,omitempty"`
	Since  time.Time `qstring:"since,omitempty"`
	Tail   int       `qstring:"tail,omitempty"`
}

type LogService interface {
	AppStream(appID string, opts *LogStreamOpts, handler LogHandler) error
	AddonStream(addonID string, opts *LogStreamOpts, handler LogHandler) error
	ReleaseStream(release *Release, opts *LogStreamOpts, handler LogHandler) error
}

func (l *LogClient) AppStream(appID string, opts *LogStreamOpts, handler LogHandler) error {
	return l.client.doStream("apps/"+appID+"/logs", opts, handler)
}

func (l *LogClient) AddonStream(addonID string, opts *LogStreamOpts, handler LogHandler) error {
	return l.client.doStream("addons/"+addonID+"/logs", opts, handler)
}

func (l *LogClient) ReleaseStream(release *Release, opts *LogStreamOpts, handler LogHandler) error {
	path := fmt.Sprintf("apps/%d/releases/%d/logs", release.AppID, release.Version)
	return l.client.doStream(path, opts, handler)
}
