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
	"net/url"
	"strings"

	"github.com/dyninc/qstring"
	"github.com/fatih/structs"
	"github.com/r3labs/sse"
	"nhooyr.io/websocket"
)

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
	doList(path string, params, out interface{}) error
	doUpdate(path string, params, out interface{}) error
	doRequest(method, path, ctype string, params, out interface{}) error
	doStream(path string, params interface{}, handler LogHandler) error
	doWebsocket(path string, handler ConnHandler) error
}

func (c *Client) doByID(path, id string, params, out interface{}) error {
	if id == "" {
		return ErrNotFound
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
	if params != nil {
		query, err := qstring.MarshalString(params)
		if err != nil {
			return err
		}
		path += "?" + query
	}
	return c.doRequest(http.MethodGet, path, "", nil, out)
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

	ctx, cancel := context.WithTimeout(req.Context(), defaultTimeout)
	defer cancel()

	req = req.WithContext(ctx)
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

func (c *Client) doStream(path string, params interface{}, handler LogHandler) error {
	query, err := qstring.MarshalString(params)
	if err != nil {
		return err
	}
	if query != "" {
		path += "?" + query
	}
	client := sse.NewClient(c.apiURL + path)
	client.Connection.Transport = c.httpClient.Transport
	ctx, cancel := context.WithCancel(context.Background())
	return client.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
		if len(msg.Data) > 0 {
			entry := new(LogEntry)
			err := json.Unmarshal(msg.Data, entry)
			if err != nil {
				return
			}
			handler(cancel, entry)
		}
	})
}

func (c *Client) doWebsocket(path string, handler ConnHandler) error {
	u, err := url.Parse(c.apiURL + path)
	if err != nil {
		return err
	}
	u.Scheme = strings.Replace(u.Scheme, "http", "ws", 1)

	ctx := context.Background()
	ws, resp, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPClient: c.httpClient,
	})
	if err != nil {
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			err = ErrServerError
		}
		return err
	}

	conn := websocket.NetConn(ctx, ws, websocket.MessageText)
	defer conn.Close()

	err = handler(conn)
	var ce websocket.CloseError
	if errors.As(err, &ce) {
		if ce.Reason == "" {
			return nil
		}
		return errors.New(ce.Reason)
	}
	return err
}
