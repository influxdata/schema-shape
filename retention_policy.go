package main

// NewRetentionPolicy creates RetentionPolicies
func NewRetentionPolicy(args []interface{}) *RetentionPolicy {
	return &RetentionPolicy{
		Name:               iToS(args[0]),
		Duration:           iToS(args[1]),
		ShardGroupDuration: iToS(args[2]),
		Replication:        iToS(args[3]),
		Default:            args[4].(bool),
	}
}

// RetentionPolicy is a RetentionPolicy
type RetentionPolicy struct {
	Name               string
	Duration           string
	ShardGroupDuration string
	Replication        string
	Default            bool
}
