package validation

import (
	"errors"
	"fmt"
	"net"
	"passvault/pkg/log"
)

const (
	noIpAddressesFound = "no ip addresses found for domain: [%s]"
	domainLookUpError  = "error when lookup domain: [%s] for ips"
)

type DomainValidation struct {
	DomainToValidate string
	LogManager       log.LogManagerInterface
}

func (dv *DomainValidation) Validate() error {
	return dv.validateDomain(dv.DomainToValidate)
}

func (dv *DomainValidation) validateDomain(domain string) error {
	ips, err := net.LookupIP(dv.DomainToValidate)
	if err != nil {
		return errors.New(fmt.Sprintf(domainLookUpError, domain))
	}

	if len(ips) == 0 {
		return errors.New(fmt.Sprintf(noIpAddressesFound, domain))
	}

	return nil
}
