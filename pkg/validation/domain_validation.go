package validation

const (
	noIpAddressesFound = "no ip addresses found for domain: [%s]"
	domainLookUpError  = "error when lookup domain: [%s] for ips"
)

type DomainValidation struct {
	DomainToValidate string
}

func (dv *DomainValidation) Validate() error {
	return dv.validateDomain(dv.DomainToValidate)
}

func (dv *DomainValidation) validateDomain(domain string) error {
	//todo validation

	return nil
}
