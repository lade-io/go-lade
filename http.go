package lade

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dyninc/qstring"
	"github.com/fatih/structs"
	"github.com/r3labs/sse/v2"
	"golang.org/x/oauth2"
	"gopkg.in/cenkalti/backoff.v1"
	"nhooyr.io/websocket"
)

func init() {
	structs.DefaultTagName = "json"
}

const (
	jsonType = "application/json"
)

type (
	ConnHandler func(net.Conn) error
	LogHandler  func(context.CancelFunc, *LogEntry)
)

var _ httpService = new(Client)

type httpService interface {
	doByID(path, id string, params, out interface{}) error
	doCreate(path string, params, out interface{}) error
	doDelete(path string, params interface{}) error
	doForm(path string, params, out interface{}) error
	doGet(path string, params, out interface{}) error
	doHead(path string, params interface{}) error
	doList(path string, params, out interface{}) error
	doUpdate(path string, params, out interface{}) error
	doRequest(method, path, ctype string, params, out interface{}) error
	doStream(path string, params interface{}, handler LogHandler) error
	doWebsocket(path string, handler ConnHandler) error
}

func addParams(path string, params interface{}) (string, error) {
	if params != nil {
		query, err := qstring.MarshalString(params)
		if err != nil {
			return "", err
		}
		if query != "" {
			path += "?" + query
		}
	}
	return path, nil
}

func (c *Client) doByID(path, id string, params, out interface{}) error {
	if id == "" {
		return ErrNotFound
	}
	if out == nil {
		return c.doHead(path+"/"+id, params)
	}
	return c.doGet(path+"/"+id, params, out)
}

func (c *Client) doCreate(path string, params, out interface{}) error {
	return c.doRequest(http.MethodPost, path, jsonType, params, out)
}

func (c *Client) doDelete(path string, params interface{}) error {
	return c.doRequest(http.MethodDelete, path, jsonType, params, nil)
}

func (c *Client) doForm(path string, params, out interface{}) error {
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	go func() {
		defer pipeWriter.Close()
		defer writer.Close()
		for key, val := range structs.Map(params) {
			switch v := val.(type) {
			case *File:
				part, err := writer.CreateFormFile(key, v.Name)
				if err != nil {
					return
				}
				if _, err = io.Copy(part, v.Body); err != nil {
					return
				}
			}
		}
	}()
	ctype := writer.FormDataContentType()
	return c.doRequest(http.MethodPost, path, ctype, pipeReader, out)
}

func (c *Client) doGet(path string, params, out interface{}) error {
	path, err := addParams(path, params)
	if err != nil {
		return err
	}
	return c.doRequest(http.MethodGet, path, jsonType, nil, out)
}

func (c *Client) doHead(path string, params interface{}) error {
	path, err := addParams(path, params)
	if err != nil {
		return err
	}
	return c.doRequest(http.MethodHead, path, jsonType, nil, nil)
}

func (c *Client) doList(path string, params, out interface{}) error {
	return c.doGet(path, params, out)
}

func (c *Client) doUpdate(path string, params, out interface{}) error {
	return c.doRequest(http.MethodPatch, path, jsonType, params, out)
}

func (c *Client) doRequest(method, path, ctype string, params, out interface{}) error {
	payload, ok := params.(io.Reader)
	if params != nil && !ok {
		data, err := json.Marshal(params)
		if err != nil {
			return err
		}
		payload = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, c.apiURL+path, payload)
	if err != nil {
		return err
	}

	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	req.Header.Set("Accept", jsonType)
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		if len(body) == 0 {
			return nil
		}
		return json.Unmarshal(body, &out)
	}

	err = &APIError{Status: resp.StatusCode}
	if len(body) > 0 {
		json.Unmarshal(body, &err)
	}
	return err
}

func (c *Client) newSSEClient(path string) (context.Context, context.CancelFunc, *sse.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	back := backoff.NewExponentialBackOff()
	back.MaxElapsedTime = 0
	client := sse.NewClient(c.apiURL + path)
	client.Connection = c.httpClient
	client.ReconnectStrategy = back
	client.ReconnectNotify = backoff.Notify(func(err error, delay time.Duration) {
		var e *oauth2.RetrieveError
		if (errors.As(err, &e) && e.Response.StatusCode < 404) ||
			strings.HasSuffix(err.Error(), http.StatusText(404)) {
			cancel()
		}
	})
	return ctx, cancel, client
}

func (c *Client) doStream(path string, params interface{}, handler LogHandler) error {
	path, err := addParams(path, params)
	if err != nil {
		return err
	}
	ctx, cancel, client := c.newSSEClient(path)
	for {
		err = client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
			switch string(msg.Event) {
			case "data":
				entry := new(LogEntry)
				err = json.Unmarshal(msg.Data, entry)
				if err != nil || entry.Source == "ping" {
					return
				}
				handler(cancel, entry)
				client.ReconnectStrategy.Reset()
			case "EOF":
				cancel()
			}
		})
		if ctx.Err() != nil {
			var e *oauth2.RetrieveError
			if errors.As(err, &e) {
				lines := strings.SplitN(e.Error(), "\n", 2)
				return errors.New(lines[0])
			}
			return err
		}
		time.Sleep(client.ReconnectStrategy.NextBackOff())
	}
}

func (c *Client) doWebsocket(path string, handler ConnHandler) error {
	operation := func() error {
		ctx := context.Background()
		ws, resp, err := websocket.Dial(ctx, c.apiURL+path, &websocket.DialOptions{
			HTTPClient: c.httpClient,
		})
		if err != nil {
			if resp.StatusCode == http.StatusInternalServerError {
				return backoff.Permanent(ErrServerError)
			}
			return err
		}
		conn := websocket.NetConn(ctx, ws, websocket.MessageText)
		defer conn.Close()
		return handler(conn)
	}
	return backoff.Retry(operation, backoff.NewExponentialBackOff())
}
