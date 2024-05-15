// This file is auto-generated, don't edit it. Thanks.
package client

import (
	"github.com/alibabacloud-go/tea/tea"
)

type QuerySpecificIPBelongRequest struct {
	// {'en':'IP address, use English comma to separate two items. Every IP address needs to following regular expression rule of   ((2[0-4]\\d|25[0-5]|1\\d\\d|0|[1-9]\\d?)\\.){3}(2[0-4]\\d|25[0-5]|1\\d\\d|0|[1-9]\\d?).   No more than 20 IP addresses are allowed', 'zh_CN':'ip地址，以英文逗号分隔，每个ip都需要符合正则((2[0-4]\\d|25[0-5]|1\\d\\d|0|[1-9]\\d?)\\.){3}(2[0-4]\\d|25[0-5]|1\\d\\d|0|[1-9]\\d?)，ip个数不能超过20'}
	Ip []*string `json:"ip,omitempty" xml:"ip,omitempty" require:"true" type:"Repeated"`
}

func (s QuerySpecificIPBelongRequest) String() string {
	return tea.Prettify(s)
}

func (s QuerySpecificIPBelongRequest) GoString() string {
	return s.String()
}

func (s *QuerySpecificIPBelongRequest) SetIp(v []*string) *QuerySpecificIPBelongRequest {
	s.Ip = v
	return s
}

type QuerySpecificIPBelongResponse struct {
	// {'en':'checkList', 'zh_CN':'结果数据'}
	CheckList []*QuerySpecificIPBelongResponseCheckList `json:"checkList,omitempty" xml:"checkList,omitempty" require:"true" type:"Repeated"`
}

func (s QuerySpecificIPBelongResponse) String() string {
	return tea.Prettify(s)
}

func (s QuerySpecificIPBelongResponse) GoString() string {
	return s.String()
}

func (s *QuerySpecificIPBelongResponse) SetCheckList(v []*QuerySpecificIPBelongResponseCheckList) *QuerySpecificIPBelongResponse {
	s.CheckList = v
	return s
}

type QuerySpecificIPBelongResponseCheckList struct {
	// {'en':'yes: the IP belongs to Our system,
	//         no: the IP does not belong to Our system', 'zh_CN':'yes：ip属于我司，no：ip不属于我司'}
	Response *string `json:"response,omitempty" xml:"response,omitempty" require:"true"`
	// {'en':'IP addresses', 'zh_CN':'ip地址'}
	Ip *string `json:"ip,omitempty" xml:"ip,omitempty" require:"true"`
}

func (s QuerySpecificIPBelongResponseCheckList) String() string {
	return tea.Prettify(s)
}

func (s QuerySpecificIPBelongResponseCheckList) GoString() string {
	return s.String()
}

func (s *QuerySpecificIPBelongResponseCheckList) SetResponse(v string) *QuerySpecificIPBelongResponseCheckList {
	s.Response = &v
	return s
}

func (s *QuerySpecificIPBelongResponseCheckList) SetIp(v string) *QuerySpecificIPBelongResponseCheckList {
	s.Ip = &v
	return s
}

type Paths struct {
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

type Parameters struct {
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

type RequestHeader struct {
}

func (s RequestHeader) String() string {
	return tea.Prettify(s)
}

func (s RequestHeader) GoString() string {
	return s.String()
}

type ResponseHeader struct {
}

func (s ResponseHeader) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeader) GoString() string {
	return s.String()
}
