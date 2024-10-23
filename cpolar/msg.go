package cpolar

type Tunnel struct {
	Name           string `json:"name,omitempty"`
	PublicUrl      string `json:"public_url,omitempty"`
	Proto          string `json:"proto,omitempty"`
	Addr           string `json:"addr,omitempty"`
	Type           string `json:"type,omitempty"`
	CreateDatetime string `json:"create_datetime,omitempty"`
}

type Items struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Configuration  any      `json:"configuration"`
	Status         string   `json:"status"`
	PublicUrl      string   `json:"public_url"`
	PublishTunnels []Tunnel `json:"publish_tunnels"`
}

type Data struct {
	Token string  `json:"token,omitempty"`
	Total int     `json:"total,omitempty"`
	Items []Items `json:"items,omitempty"`
}

type Response struct {
	Data    Data   `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
