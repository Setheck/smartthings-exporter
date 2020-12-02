package smartthings

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	API = "https://api.smartthings.com/v1"
)

type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Error     *Error `json:"error"`
}

type Error struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Target  string   `json:"target"`
	Details []*Error `json:"details"`
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Code, e.Message, e.Target)
}

type ListResponse struct {
	Items       json.RawMessage `json:"items"`
	PagingLinks `json:"_links"`
}

type PagingLinks struct {
	Next     map[string]string `json:"next"`
	Previous map[string]string `json:"previous"`
}

type Location struct {
	LocationID       string
	Name             string
	CountryCode      string
	Latitude         float32
	Longitude        float32
	RegionRadius     int
	TemperatureScale string
	TimeZoneID       string
	Locale           string
	Parent           map[string]string
}

type Profile struct {
	ID string
}

type App struct {
	InstalledAppID string
	ExternalID     string
	*Profile
}

type Device struct {
	DeviceID               string       `json:"deviceId"`
	Name                   string       `json:"name,omitempty"`
	Label                  string       `json:"label,omitempty"`
	ManufacturerName       string       `json:"manufacturerName,omitempty"`
	PresentationID         string       `json:"presentationId,omitempty"`
	DeviceManufacturerCode string       `json:"deviceManufacturerCode,omitempty"`
	LocationID             string       `json:"locationId,omitempty"`
	OwnerID                string       `json:"ownerId,omitempty"`
	RoomID                 string       `json:"roomId,omitempty"`
	DeviceTypeID           string       `json:"deviceTypeId,omitempty"`
	DeviceTypeName         string       `json:"deviceTypeName,omitempty"`
	Components             []*Component `json:"components,omitempty"`
	ChildDevices           []*Device    `json:"childDevices,omitempty"`
	Profile                *Profile     `json:"profile,omitempty"`
	App                    *App         `json:"app,omitempty"`
	Ble                    *Ble         `json:"ble,omitempty"`
	BleD2D                 *BleD2D      `json:"bleD2D,omitempty"`
	Dth                    *Dth         `json:"dth,omitempty"`
	Ir                     *Ir          `json:"ir,omitempty"`
	IrOcf                  *IrOcf       `json:"irOcf,omitempty"`
	Viper                  *Viper       `json:"viper,omitempty"`
	Type                   string       `json:"type,omitempty"`
	RestrictionTier        int          `json:"restrictionTier,omitempty"`
}

func ToString(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Ble struct {
}
type BleD2D struct {
}
type Dth struct {
}
type Ir struct {
}
type IrOcf struct {
}
type Viper struct {
}

type Capability struct {
	ID      string
	Version int
}

type Component struct {
	ID           string
	Label        string
	Capabilities []*Capability
	Categories   []map[string]string
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (client *Client) apiGet(endpoint string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, API+endpoint, nil)
	if err != nil {
		return nil, err
	}

	queryParams.Encode()
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed request: %s", resp.Status)
	}

	return resp, nil
}

func (client *Client) ListAllCapabilities() ([]*Capability, error) {
	resp, err := client.apiGet("/capabilities", nil)
	if err != nil {
		return nil, err
	}

	response, err := parseListResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	var capabilities []*Capability
	err = json.Unmarshal(response.Items, &capabilities)
	return capabilities, err
}

func (client *Client) ListDevices() ([]*Device, error) {
	resp, err := client.apiGet("/devices", nil)
	if err != nil {
		return nil, err
	}

	response, err := parseListResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	err = json.Unmarshal(response.Items, &devices)
	return devices, err
}

func parseListResponse(input io.ReadCloser) (*ListResponse, error) {
	raw, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}
	var errResponse *ErrorResponse
	if err := json.Unmarshal(raw, &errResponse); err == nil {
		if errResponse.Error != nil {
			return nil, errResponse.Error
		}
	}

	var listResponse *ListResponse
	if err := json.Unmarshal(raw, &listResponse); err != nil {
		return nil, err
	}
	return listResponse, nil
}
