package graph

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprintf(w, `"%s"`, u.String())
	})
}

func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	str, ok := v.(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("UUID must be a string")
	}
	return uuid.Parse(str)
}

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		fmt.Fprintf(w, `"%s"`, t.Format(time.RFC3339))
	})
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	str, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("Time must be a string")
	}
	return time.Parse(time.RFC3339, str)
}
