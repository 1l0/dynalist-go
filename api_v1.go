// Dynalist API for Go
//
//  Set an env before using:
//   export DYNALIST_TOKEN=your_secret_token
package dynalist

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

const (
	baseurl = `https://dynalist.io/api/v1/`
)

type API struct {
	Token      string       `json:"token"`
	RateLimit  time.Time    `json:"-"`
	BurstLimit int          `json:"-"`
	client     *http.Client `json:"-"`
}

func New() (*API, error) {
	token := os.Getenv("DYNALIST_TOKEN")
	if token == "" {
		return nil, errors.New("failed to get $DYNALIST_TOKEN")
	}
	return &API{
		Token:  token,
		client: &http.Client{},
	}, nil
}

func (api *API) FileList() (*Response, error) {
	var res Response
	err := api.post(baseurl+"file/list", api, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (api *API) FileEdit(changes []*Change) (*Response, error) {
	var res Response
	param := struct {
		Token   string    `json:"token"`
		Changes []*Change `json:"changes"`
	}{
		Token:   api.Token,
		Changes: changes,
	}
	err := api.post(baseurl+"file/edit", &param, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (api *API) DocRead(fileID string) (*Response, error) {
	var res Response
	param := struct {
		Token  string `json:"token"`
		FileID string `Json:"file_id"`
	}{
		Token:  api.Token,
		FileID: fileID,
	}
	err := api.post(baseurl+"doc/read", &param, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (api *API) DocEdit(fileID string, changes []*Change) (*Response, error) {
	var res Response
	param := struct {
		Token   string    `json:"token"`
		FileID  string    `Json:"file_id"`
		Changes []*Change `json:"changes"`
	}{
		Token:   api.Token,
		FileID:  fileID,
		Changes: changes,
	}
	err := api.post(baseurl+"doc/edit", &param, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (api *API) InboxAdd(change *Change) (*Response, error) {
	var res Response
	param := struct {
		Change
		Token string `json:"token"`
	}{
		Token:  api.Token,
		Change: *change,
	}
	err := api.post(baseurl+"inbox/add", &param, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (api *API) post(url string, in, out interface{}) error {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if res.Body == nil {
		return errors.New("no body in the response")
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

type Limit struct {
	Rate  time.Duration
	Burst int
}

func (api *API) LimitFileList() *Limit {
	return &Limit{
		Rate:  time.Minute / 6,
		Burst: 10,
	}
}

func (api *API) LimitFileEdit() *Limit {
	return &Limit{
		Rate:  time.Minute / 60,
		Burst: 50,
	}
}

func (api *API) LimitDocRead() *Limit {
	return &Limit{
		Rate:  time.Minute / 60,
		Burst: 50,
	}
}

func (api *API) LimitDocEdit() *Limit {
	return &Limit{
		Rate:  time.Minute / 60,
		Burst: 20,
	}
}

func (api *API) LimitChange() *Limit {
	return &Limit{
		Rate:  time.Minute / 240,
		Burst: 500,
	}
}

func (api *API) LimitInboxAdd() *Limit {
	return &Limit{
		Rate:  time.Minute / 6,
		Burst: 10,
	}
}

type Change struct {
	Action   Action `json:"action"`
	Index    int    `json:"index,omitempty"`
	NodeID   string `json:"node_id,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
	Content  string `json:"content,omitempty"`
	Type     Type   `json:"type,omitempty"`
	FileID   string `json:"file_id,omitempty"`
	Title    string `json:"title,omitempty"`
	Note     string `json:"note,omitempty"`
	Checked  bool   `json:"checked,omitempty"`
}

func NewChange(action Action) *Change {
	return &Change{Action: action}
}

type Action string

const (
	ActionInsert Action = "action"
	ActionEdit   Action = "edit"
	ActionMove   Action = "move"
	ActionDelete Action = "delete"
)

type Response struct {
	Code       Code   `json:"_code"`
	Msg        string `json:"_msg"`
	RootFileID string `json:"root_file_id,omitempty"`
	Files      []File `json:"files,omitempty"`
	Results    []bool `json:"results,omitempty`
	Title      string `json:"title,omitempty`
	Nodes      []Node `json:"nodes,omitempty`
}

type Code string

const (
	CodeOK Code = "Ok"
	//Your request is not valid JSON.
	CodeInvalid Code = "Invalid"
	//You've hit the limit on how many requests you can send.
	CodeTooManyRequests Code = "TooManyRequests"
	//Your secret token is invalid.
	CodeInvalidToken Code = "InvalidToken"
	//Server unable to handle the request.
	CodeLockFail Code = "LockFail"
	//You don't have permission to access this document.
	CodeUnauthorized Code = "Unauthorized"
	//The document you're requesting is not found.
	CodeNotFound Code = "NotFound"
	//The node (item) you're requesting is not found.
	CodeNodeNotFound Code = "NodeNotFound"
	//Inbox location is not configured, or invalid.
	CodeNoInbox Code = "NoInbox"
)

type File struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	Type       Type       `json:"type"`
	Permission Permission `json:"permission"`
	Collapsed  bool       `json:"collapsed,omitempty"`
	Children   []string   `json:"children,omitempty"`
}

type Type string

const (
	TypeDocument Type = "document"
	TypeFolder   Type = "folder"
)

type Permission int

const (
	PermissionNoAccess Permission = iota
	PermissionReadOnly
	PermissionEditRights
	PermissionManage
	PermissionOwner
)

type Node struct {
	ID        string   `json:"id"`
	Content   string   `json:"content"`
	Note      string   `json:"note,omitempty"`
	Checked   bool     `json:"checked,omitempty"`
	Collapsed bool     `json:"collapsed,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	Children  []string `json:"children,omitempty"`
}
