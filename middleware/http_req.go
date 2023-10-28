package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

func HttpReq(url, data string, connId uint64) string {
	reader := bytes.NewReader(([]byte)(data))
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)

	if err != nil {
		Log.Errorf("[%v] get err %s", connId, err.Error())
		return ""
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	return buf.String()
}

// HTTPReqWebVersion http request with web-version header
func HTTPReqWebVersion(url, data, webVersion string, connId uint64) string {
	reader := bytes.NewReader(([]byte)(data))
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")
	if webVersion == "" {
		webVersion = "0"
	}
	req.Header.Add("web-version", webVersion)

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)

	if err != nil {
		Log.Errorf("[%v] get err %s", connId, err.Error())
		return ""
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	return buf.String()
}

// HTTPReqWebVersion2 http request with web-version header
func HTTPReqWebVersion2(url, data, webVersion string, connId uint64) (string, error) {
	reader := bytes.NewReader(([]byte)(data))
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", "application/json")
	if webVersion == "" {
		webVersion = "0"
	}
	req.Header.Add("web-version", webVersion)

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)

	if err != nil {
		Log.Errorf("ERROR:[%v] get err %s", connId, err.Error())
		return "", err
	}
	defer resp.Body.Close()

	Log.Debugf("ContentLength:%d", resp.ContentLength)
	if resp.ContentLength < 0 {
		return "", errors.New("ContentLength is -1")
	}
	buf := make([]byte, resp.ContentLength)
	nReadIndex := int64(0)
	for nReadIndex < resp.ContentLength {
		n, err := resp.Body.Read(buf[nReadIndex:])
		if err != nil && err != io.EOF {
			Log.Errorf("ERROR:%s", err.Error())
			return "", err
		}
		nReadIndex += int64(n)
		if n == 0 {
			Log.Debugf("LEN:%d", nReadIndex)
			break
		}
	}

	return string(buf), nil
}
