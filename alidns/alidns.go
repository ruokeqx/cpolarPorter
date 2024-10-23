// Package alidns implements a DNS provider for solving the DNS-01 challenge using Alibaba Cloud DNS.
package alidns

import (
	"errors"
	"fmt"
	"time"

	"github.com/ruokeqx/cpolarPorter/env"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

const (
	defaultRegionID    = "cn-hangzhou"
	defaultTTL         = 600
	defaultHTTPTimeout = 10 * time.Second
)

// Environment variables names.
const (
	envNamespace = "ALICLOUD_"

	EnvRAMRole       = envNamespace + "RAM_ROLE"
	EnvAccessKey     = envNamespace + "ACCESS_KEY"
	EnvSecretKey     = envNamespace + "SECRET_KEY"
	EnvSecurityToken = envNamespace + "SECURITY_TOKEN"
	EnvRegionID      = envNamespace + "REGION_ID"

	EnvTTL         = envNamespace + "TTL"
	EnvHTTPTimeout = envNamespace + "HTTP_TIMEOUT"
)

// Config is used to configure the creation of the DNSProvider.
type Config struct {
	RAMRole       string
	APIKey        string
	SecretKey     string
	SecurityToken string
	RegionID      string
	TTL           int
	HTTPTimeout   time.Duration
}

// NewDefaultConfig returns a default configuration for the DNSProvider.
func NewDefaultConfig() *Config {
	return &Config{
		TTL:         env.GetOrDefaultInt(EnvTTL, defaultTTL),
		HTTPTimeout: env.GetOrDefaultSecond(EnvHTTPTimeout, defaultHTTPTimeout),
	}
}

// DNSProvider implements the challenge.Provider interface.
type DNSProvider struct {
	config *Config
	client *alidns.Client
}

// NewDNSProvider returns a DNSProvider instance configured for Alibaba Cloud DNS.
// - If you're using the instance RAM role, the RAM role environment variable must be passed in: ALICLOUD_RAM_ROLE.
// - Other than that, credentials must be passed in the environment variables:
// ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY, and optionally ALICLOUD_SECURITY_TOKEN.
func NewDNSProvider() (*DNSProvider, error) {
	config := NewDefaultConfig()
	config.RegionID = env.GetOrFile(EnvRegionID)

	values, err := env.Get(EnvRAMRole)
	if err == nil {
		config.RAMRole = values[EnvRAMRole]
		return NewDNSProviderConfig(config)
	}

	values, err = env.Get(EnvAccessKey, EnvSecretKey)
	if err != nil {
		return nil, fmt.Errorf("alicloud: %w", err)
	}

	config.APIKey = values[EnvAccessKey]
	config.SecretKey = values[EnvSecretKey]
	config.SecurityToken = env.GetOrFile(EnvSecurityToken)

	return NewDNSProviderConfig(config)
}

// NewDNSProviderConfig return a DNSProvider instance configured for alidns.
func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("alicloud: the configuration of the DNS provider is nil")
	}

	if config.RegionID == "" {
		config.RegionID = defaultRegionID
	}

	var credential auth.Credential
	switch {
	case config.RAMRole != "":
		credential = credentials.NewEcsRamRoleCredential(config.RAMRole)
	case config.APIKey != "" && config.SecretKey != "" && config.SecurityToken != "":
		credential = credentials.NewStsTokenCredential(config.APIKey, config.SecretKey, config.SecurityToken)
	case config.APIKey != "" && config.SecretKey != "":
		credential = credentials.NewAccessKeyCredential(config.APIKey, config.SecretKey)
	default:
		return nil, errors.New("alicloud: ram role or credentials missing")
	}

	conf := sdk.NewConfig().WithTimeout(config.HTTPTimeout)

	client, err := alidns.NewClientWithOptions(config.RegionID, conf, credential)
	if err != nil {
		return nil, fmt.Errorf("alicloud: credentials failed: %w", err)
	}

	return &DNSProvider{config: config, client: client}, nil
}

// Present creates a TXT record
func (d *DNSProvider) Present(DomainName, RR, value string) error {
	recordAttributes, err := d.newTxtRecord(DomainName, RR, value)
	if err != nil {
		return err
	}

	_, err = d.client.AddDomainRecord(recordAttributes)
	if err != nil {
		return fmt.Errorf("alicloud: API call failed: %w", err)
	}
	return nil
}

// CleanUp removes the TXT record matching the specified parameters.
func (d *DNSProvider) CleanUp(DomainName, RR string) error {
	records, err := d.findTxtRecords(DomainName, RR)
	if err != nil {
		return fmt.Errorf("alicloud: %w", err)
	}

	for _, rec := range records {
		request := alidns.CreateDeleteDomainRecordRequest()
		request.RecordId = rec.RecordId
		_, err = d.client.DeleteDomainRecord(request)
		if err != nil {
			return fmt.Errorf("alicloud: %w", err)
		}
	}
	return nil
}

func (d *DNSProvider) newTxtRecord(DomainName, RR, value string) (*alidns.AddDomainRecordRequest, error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.RR = RR
	request.Type = "TXT"
	request.Value = value
	request.DomainName = DomainName
	request.TTL = requests.NewInteger(d.config.TTL)

	return request, nil
}

func (d *DNSProvider) findTxtRecords(DomainName, RR string) ([]alidns.Record, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = DomainName
	request.PageSize = requests.NewInteger(500)

	var records []alidns.Record

	result, err := d.client.DescribeDomainRecords(request)
	if err != nil {
		return records, fmt.Errorf("API call has failed: %w", err)
	}

	for _, record := range result.DomainRecords.Record {
		if record.RR == RR && record.Type == "TXT" {
			records = append(records, record)
		}
	}
	return records, nil
}
