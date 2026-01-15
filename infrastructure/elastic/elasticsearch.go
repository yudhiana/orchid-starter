package elastic

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"orchid-starter/config"

	"net/http"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/yudhiana/logos"
)

type esUtil struct {
	esClient *elasticsearch.Client
}

var esInstance *esUtil
var esOnce sync.Once

// NewESConnection initialize elasticsearch connection
func NewESConnection(localConfig *config.LocalConfig) *elasticsearch.Client {
	esOnce.Do(func() {
		logos.NewLogger().Info("Initialize elasticsearch connection...")
		cfg := elasticsearch.Config{
			Addresses: strings.Split(localConfig.EsConfig.ESAddresses, ","),
			Transport: &http.Transport{
				IdleConnTimeout:     time.Duration(localConfig.EsConfig.ESIdleTimeOut) * time.Second,
				MaxIdleConns:        localConfig.EsConfig.ESMaxIdleConns,
				MaxIdleConnsPerHost: localConfig.EsConfig.ESMaxIdleConnsPerHost,
				MaxConnsPerHost:     localConfig.EsConfig.ESMaxConnsPerHost,
			},
			MaxRetries: 3,
			RetryOnStatus: []int{
				http.StatusBadGateway,
				http.StatusServiceUnavailable,
				http.StatusGatewayTimeout,
			},
			RetryBackoff: func(attempt int) time.Duration {
				return time.Duration(attempt) * 100 * time.Millisecond
			},
			Logger: getLogger(localConfig.ElasticsearchDebug),
		}

		logos.NewLogger().Info("Connect to elasticsearch", "addresses", cfg.Addresses)
		client, errConnect := elasticsearch.NewClient(cfg)
		if errConnect != nil {
			logos.NewLogger().Error("Failed to connect to elasticsearch", "error", errConnect)
			panic(errConnect)
		}

		// ping elasticsearch to check elasticsearch  established connection
		if response, errPing := client.Ping(); errPing != nil {
			logos.NewLogger().Error("Failed Ping to elasticsearch", "error", errPing)
			panic(errPing)
		} else {
			logos.NewLogger().Info("Ping to elasticsearch", "status", response.Status())
		}

		esInstance = &esUtil{
			esClient: client,
		}
	})
	return esInstance.esClient
}

func getLogger(debug bool) (esDebug elastictransport.Logger) {
	if !debug {
		return
	}

	return &elastictransport.ColorLogger{
		Output:             os.Stdout,
		EnableRequestBody:  true, // show NDJSON you’re sending
		EnableResponseBody: true, // show ES response
	}
}

type mockTransport struct {
	response string
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	body := `{
		"name":"mock-es",
		"cluster_name":"elasticsearch",
		"cluster_uuid":"uuid",
		"version":{"number":"8.12.0"},
		"tagline":"You Know, for Search"
	}`

	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	resp.Header.Set("X-Elastic-Product", "Elasticsearch")
	return resp, nil
}

func GetESMockConnection() *elasticsearch.Client {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Transport: &mockTransport{
			response: `{"hits":{"total":{"value":1}}}`,
		},
	})

	if err != nil {
		panic(err)
	}
	return es
}
