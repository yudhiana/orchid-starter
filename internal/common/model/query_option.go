package model

// QueryOption --
type QueryOption struct {
	Condition   string        `json:"condition"`   // ex : SELECT * FROM `table`  WHERE (field='field')
	Page        int           `json:"page"`        // ex : SELECT * FROM `table`  LIMIT 10 OFFSET 10
	Limit       int           `json:"limit"`       // ex : SELECT * FROM `table`  LIMIT 10 OFFSET 0
	Order       string        `json:"order"`       // ex : SELECT * FROM `table`  ORDER BY ID ASC
	Bind        *string       `json:"bind"`        // binding query filter. ex : `condition: "key=:key"` it will replace by ` bind: "{'key':'value'}"`
	Filter      []Filter      `json:"filter"`      // filtering key on process filter
	ExtraParams []ExtraParams `json:"extraParams"` // Ex: /company/{compID}/address/billing/{billingAddressID}
}

// Filter --
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PaginationRest struct {
	Page   int    `json:"page"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Order  string `json:"order"`
}

// Pagination --
type Pagination struct {
	Page        int  `json:"page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	PerPage     int  `json:"limit"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// ValidateFilter --
type ValidateFilter struct {
	Condition   string                 `json:"condition"`   // ex : SELECT * FROM `table`  WHERE (field='field')
	Page        int                    `json:"page"`        // ex : SELECT * FROM `table`  LIMIT 10 OFFSET 10
	Limit       int                    `json:"limit"`       // ex : SELECT * FROM `table`  LIMIT 10 OFFSET 0
	Bind        map[string]interface{} `json:"bind"`        // binding query filter. ex : `condition: "key=:key"` it will replace by ` bind: "{'key':'value'}"`
	Order       string                 `json:"order"`       // ex : SELECT * FROM `table`  order by id asc
	Filter      []Filter               `json:"filter"`      // filtering key on process filter
	ExtraParams interface{}            `json:"extraParams"` // Ex: /company/{compID}/address/billing/{billingAddressID}
}

// ExtraParams can be used for URL path parameters such as company ID, address ID
// Ex: /company/{compID}/address/billing/{billingAddressID}
type ExtraParams struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
