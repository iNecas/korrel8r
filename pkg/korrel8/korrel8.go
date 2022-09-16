// package korrel8 correlates observable signals from different domains
//
// Generic types and interfaces to define correlation rules, and correlate objects between different domains.
//
// The main interfaces are:
//
// - Object: A set of attributes representing a signal (e.g. log record, metric time-series, ...)
// The concrete type depends on the domain, for correlation purposes it is equivalent to a JSON object.
//
// - Domain: a set of objects with a common vocabulary (e.g. k8s resources, OpenTracing spans, ...)
//
// - Class: a subset of objects in the same domain with a common schema (e.g. k8s Pod, prometheus Alert)
//
// - Rule: takes a starting object and returns a query for related goal objects.
//
// - Store: a store of objects belonging to the same domain (e.g. a Loki log store, k8s API server)
//
// Signals and resources from different domains have different representations and naming conventions.
// Domain-specific packages implement the interfaces in this package so we can do cross-domain correlation.
//
package korrel8

import (
	"context"
	"time"
)

// Object represents a signal instance.
type Object interface {
	Identifier() Identifier // Identifies this object instance.
	Class() Class           // Class of the object.
	Native() any            // Native representation of the object.
}

// Domain names a set of objects based on the same technology.
type Domain string

// Identifier is a comparable value that identifies an "instance" of a signal.
//
// For example a namespace+name for a k8s resource, or a uri+labels for a metric time series.
type Identifier any

// Class identifies a subset of objects from the same domain with the same schema.
// For example Pod is a class in the k8s domain.
//
// Class implementations must be comparable.
type Class interface {
	Domain() Domain // Domain of this class.
}

// Queries is a collection of query strings.
// Query string format depends on the domain to be queried, for example a k8s GET URoI or a PromQL query string.
type Queries []string

// Get the collection of objects returned by executing all queries against store.
// Results are de-duplicated based on Object.Identifier.
func (r Queries) Get(ctx context.Context, s Store) ([]Object, error) {
	dedup := uniqueObjects{}
	for _, q := range r {
		objs, err := s.Query(ctx, q)
		if err != nil {
			return nil, err
		}
		dedup.add(objs)
	}
	return dedup.list(), nil
}

// Store is a source of signals belonging to a single domain.
type Store interface {
	// Query a query, return the resulting objects.
	Query(ctx context.Context, query string) ([]Object, error)
}

// Rule encapsulates logic to find correlated goal objects from a start object.
//
type Rule interface {
	Start() Class // Class of start object
	Goal() Class  // Class of desired result object(s)

	// Apply the rule to start Object.
	// Return a list of queries for correlated objects in the Goal() domain.
	// The queries include the contraint (which can be nil)
	Apply(Object, *Constraint) (Queries, error)
}

// Constraint to apply to the result of following a rule.
type Constraint struct {
	After  *time.Time // Include only results timestamped after this time.
	Before *time.Time // Include only results timestamped before this time.
}

// Path is a list of rules where the Goal() of each rule is the Start() of the next.
type Path []Rule