package pid

import "github.com/google/uuid"

type PID struct {
	ID uuid.UUID
}

func NewPID() (PID, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return PID{}, err
	}
	return PID{ID: u}, nil
}

func (pid *PID) Equal(other *PID) bool {
	return pid != nil && other != nil && pid.ID == other.ID
}
