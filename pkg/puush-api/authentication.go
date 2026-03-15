package puush

import "time"

type Account struct {
	Type            AccountType
	DiskUsage       int64
	SubscriptionEnd *time.Time
}

type Credentials struct {
	Username *string
	Password *string
	Key      *string
}

func (c *Credentials) HasApiKey() bool {
	return c.Key != nil
}

func (c *Credentials) HasLoginCredentials() bool {
	return c.Username != nil && c.Password != nil
}

func (c *Credentials) IsValid() bool {
	return c.HasApiKey() || c.HasLoginCredentials()
}
