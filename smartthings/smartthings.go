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
	Version = "dev"
	Debug   = false

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
	ID               string
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
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Components []*Component `json:"components"`
	Metadata   interface{}  `json:"metadata"`
	Status     string       `json:"status"`
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
	DeviceNetworkType      string       `json:"deviceNetworkType,omitempty"`
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

type ComponentStatus map[string]ComponentAttributes

type ComponentAttributes map[string]ComponentProperties

type ComponentProperties map[string]interface{}

func ToString(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Room struct {
	ID         string `json:"roomId"`
	LocationID string `json:"locationId"`
	Name       string `json:"name"`
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
	ID      string `json:"id"`
	Version int    `json:"version"`
	Status  string `json:"status"`
}

type Component struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Capabilities []*Capability
	Categories   []map[string]string
}

type HealthState struct {
	State           string `json:"state"`
	LastUpdatedDate string `json:"lastUpdatedDate"`
}

type InstalledApp struct {
	InstalledAppID     string            `json:"installedAppId"`
	InstalledAppType   string            `json:"installedAppType"`
	InstalledAppStatus string            `json:"installedAppStatus"`
	DisplayName        string            `json:"displayName"`
	AppId              string            `json:"appId"`
	ReferenceID        string            `json:"referenceId"`
	LocationID         string            `json:"locationId"`
	Owner              map[string]string `json:"owner"`
	Notices            []string          `json:"notices"`
	CreatedDate        string            `json:"createdDate"`
	LastUpdatedDate    string            `json:"lastUpdatedDate"`
}

type Subscription struct {
	ID             string  `json:"id,omitempty"`
	InstalledAppId string  `json:"installedAppId,omitempty"`
	SourceType     string  `json:"sourceType,omitempty"`
	Device         *Device `json:"device,omitempty"`
}

type Cron struct {
	Expression string `json:"expression"`
	Timezone   string `json:"timezone"`
}

type Rules struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Actions    interface{}
	Sequence   map[string]string
	Status     string `json:"status"`
	TimezoneID string `json:"timezoneId"`
}

type Schedule struct {
	ScheduledExecutions []int  `json:"scheduledExecutions"`
	Name                string `json:"name"`
	Cron                *Cron  `json:"cron"`
	InstalledAppID      string `json:"installedAppId"`
	LocationId          string `json:"locationId"`
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (client *Client) ListInstalledApps(params url.Values) ([]*InstalledApp, error) {
	resp, err := client.apiGet("/installedapps", params)
	if err != nil {
		return nil, err
	}

	var installedApps []*InstalledApp
	if _, err := parseListResponse(resp.Body, &installedApps); err != nil {
		return nil, err
	}

	return installedApps, err
}

func (client *Client) ListApps(params url.Values) ([]*App, error) {
	resp, err := client.apiGet("/apps", params)
	if err != nil {
		return nil, err
	}

	var apps []*App
	if _, err := parseListResponse(resp.Body, &apps); err != nil {
		return nil, err
	}

	return apps, err
}

func (client *Client) ListAllDeviceProfiles(params url.Values) ([]*Profile, error) {
	resp, err := client.apiGet("/deviceprofiles", params)
	if err != nil {
		return nil, err
	}

	var profile []*Profile
	if _, err := parseListResponse(resp.Body, &profile); err != nil {
		return nil, err
	}

	return profile, err
}

func (client *Client) ListRooms(locationId string) ([]*Room, error) {
	resp, err := client.apiGet(fmt.Sprintf("/locations/%s/rooms", locationId), nil)
	if err != nil {
		return nil, err
	}

	var rooms []*Room
	if _, err := parseListResponse(resp.Body, &rooms); err != nil {
		return nil, err
	}

	return rooms, err
}

func (client *Client) ListLocations(params url.Values) ([]*Location, error) {
	resp, err := client.apiGet("/locations", params)
	if err != nil {
		return nil, err
	}

	var locations []*Location
	if _, err := parseListResponse(resp.Body, &locations); err != nil {
		return nil, err
	}

	return locations, err
}

func (client *Client) ListAllCapabilities(params url.Values) ([]*Capability, error) {
	resp, err := client.apiGet("/capabilities", params)
	if err != nil {
		return nil, err
	}

	var capabilities []*Capability
	if _, err := parseListResponse(resp.Body, &capabilities); err != nil {
		return nil, err
	}

	return capabilities, err
}

func (client *Client) GetCapabilitiesByIDAndVersion(capabilityId string, capabilityVersion int) ([]*Capability, error) {
	resp, err := client.apiGet(fmt.Sprintf("/capabilities/%s/%d", capabilityId, capabilityVersion), nil)
	if err != nil {
		return nil, err
	}

	var capabilities []*Capability
	if _, err := parseListResponse(resp.Body, &capabilities); err != nil {
		return nil, err
	}

	return capabilities, err
}

func (client *Client) ListDevices() ([]*Device, error) {
	resp, err := client.apiGet("/devices", nil)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	if _, err = parseListResponse(resp.Body, &devices); err != nil {
		return nil, err
	}

	return devices, err
}

func (client *Client) GetFullDeviceStatus(deviceId string) ([]*Component, error) {
	resp, err := client.apiGet(fmt.Sprintf("/devices/%s/status", deviceId), nil)
	if err != nil {
		return nil, err
	}

	var devices []*Component
	if _, err = parseListResponse(resp.Body, &devices); err != nil {
		return nil, err
	}

	return devices, err
}

func (client *Client) GetDeviceComponentStatus(deviceId, componentId string) (ComponentStatus, error) {
	resp, err := client.apiGet(fmt.Sprintf("/devices/%s/components/%s/status", deviceId, componentId), nil)
	if err != nil {
		return nil, err
	}

	var deviceComponentStatus ComponentStatus
	err = parseResponse(resp.Body, &deviceComponentStatus)

	return deviceComponentStatus, err
}

func (client *Client) GetCapabilityStatus(deviceId, componentId, capabilityId string) (*ComponentAttributes, error) {
	resp, err := client.apiGet(fmt.Sprintf(" /devices/%s/components/%s/capabilities/%s/status", deviceId, componentId, capabilityId), nil)
	if err != nil {
		return nil, err
	}

	var attributes *ComponentAttributes
	err = parseResponse(resp.Body, &attributes)

	return attributes, err
}

func (client *Client) ListSubscriptions(installedAppId string) ([]*Subscription, error) {
	resp, err := client.apiGet(fmt.Sprintf("/installedapps/%s/subscriptions", installedAppId), nil)
	if err != nil {
		return nil, err
	}

	var subscriptions []*Subscription
	if _, err := parseListResponse(resp.Body, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, err
}

func (client *Client) ListSchedules(installedAppId string) ([]*Schedule, error) {
	resp, err := client.apiGet(fmt.Sprintf("/installedapps/%s/schedules", installedAppId), nil)
	if err != nil {
		return nil, err
	}

	var schedules []*Schedule
	if _, err := parseListResponse(resp.Body, &schedules); err != nil {
		return nil, err
	}

	return schedules, err
}

func (client *Client) ListRules(params url.Values) ([]*Rules, error) {
	resp, err := client.apiGet("/rules", params)
	if err != nil {
		return nil, err
	}

	var rules []*Rules
	if _, err := parseListResponse(resp.Body, &rules); err != nil {
		return nil, err
	}

	return rules, err
}

func (client *Client) apiGet(endpoint string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, API+endpoint, nil)
	if err != nil {
		return nil, err
	}

	queryParams.Encode()
	req.URL.RawQuery = queryParams.Encode()
	req.Header.Add("User-Agent", fmt.Sprintf("go-smartthings-%s", Version))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed request: %s - %s", req.URL.String(), resp.Status)
	}

	return resp, nil
}

func parseListResponse(input io.ReadCloser, itemsOut interface{}) (*ListResponse, error) {
	raw, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if Debug {
		fmt.Println("raw response:", string(raw))
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

	err = json.Unmarshal(listResponse.Items, &itemsOut)
	return listResponse, err
}

func parseResponse(input io.ReadCloser, Out interface{}) error {
	raw, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}

	if Debug {
		fmt.Println("raw response:", string(raw))
	}

	var errResponse *ErrorResponse
	if err := json.Unmarshal(raw, &errResponse); err == nil {
		if errResponse.Error != nil {
			return errResponse.Error
		}
	}

	if err := json.Unmarshal(raw, &Out); err != nil {
		return err
	}

	return nil
}
