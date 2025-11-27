package netgear

import (
	"bytes"
	"io"
	"net/url"
	"testing"
)

func TestNetgear_DeviceInfo(t *testing.T) {
	n, err := NewNetgear("http://example.com", "", "")
	if err != nil {
		t.Fatal(err)
	}

	n.requestHandler = func(_ string, _ *url.URL) (io.Reader, error) {
		var deviceInfoResponse = `
{
 "deviceInfo": {
  "name": "",
  "mac": "28:94:01:75:0D:62",
  "servicePortIP": "192.168.0.239",
  "avuiVer": "2.2.10.20",
  "poe": true,
  "STP": 32768,
  "fanMode": 2,
  "units": 1,
  "details": [{
    "unit": 1,
    "management": true,
    "active": true,
    "model": "M4250-8G2XF-PoE+",
    "fwVer": "13.0.5.14",
    "bootVer": "1.0.0.11",
    "sn": "XXXXXX",
    "upTime": "1 days, 8 hrs, 40 mins, 35 secs"
   }],
  "fan": [{
    "unit": 1,
    "details": [{
      "id": 1,
      "desc": "FAN-1",
      "speed": 0,
      "dutyLevel": 0,
      "state": 3
     }]
   }],
  "sensor": [{
    "unit": 1,
    "details": [{
      "id": 1,
      "desc": "sensor-System1",
      "temp": 43,
      "maxTemp": 81,
      "state": 2
     }, {
      "id": 2,
      "desc": "sensor-MAC",
      "temp": 44,
      "maxTemp": 81,
      "state": 2
     }, {
      "id": 3,
      "desc": "sensor-System2",
      "temp": 43,
      "maxTemp": 81,
      "state": 2
     }]
   }],
  "cpu": [{
    "unit": 1,
    "usage": "5.20%"
   }],
  "memory": [{
    "unit": 1,
    "usage": "32.65%"
   }]
 },
 "resp": {
  "respCode": 0,
  "respMsg": "Success",
  "status": "success"
 }
}`
		buff := bytes.NewBuffer([]byte(deviceInfoResponse))
		return buff, nil
	}

	dInfo, err := n.DeviceInfo()
	if err != nil {
		t.Error(err)
	}

	details := dInfo.DeviceInfo.Details
	if len(details) != 1 {
		t.Errorf("unexpected lenght %d", len(details))
	}
	uptime := details[0].Uptime
	if uptime != "1 days, 8 hrs, 40 mins, 35 secs" {
		t.Errorf("unexpected uptime %q", uptime)
	}
}

func TestNetgear_ModeBasic(t *testing.T) {
	// n, err := NewNetgear("http://example.com", "", "")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// n.ModeBasic()
	// TODO
}
