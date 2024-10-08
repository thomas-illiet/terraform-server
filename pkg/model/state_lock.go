package model

// StateLock gets sent by Terraform as locking payload.
type StateLock struct {
	ID        string
	Operation string
	Info      string
	Who       string
	Version   string
	Created   string
	Path      string
}
