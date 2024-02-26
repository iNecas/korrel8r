// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

// package engine implements generic correlation logic to correlate across domains.
package engine

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	sprig "github.com/go-task/slim-sprig"
	"github.com/korrel8r/korrel8r/internal/pkg/logging"
	"github.com/korrel8r/korrel8r/pkg/graph"
	"github.com/korrel8r/korrel8r/pkg/korrel8r"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var log = logging.Log()

// Engine combines a set of domains and a set of rules, so it can perform correlation.
type Engine struct {
	domains       []korrel8r.Domain
	domainMap     map[string]korrel8r.Domain
	stores        map[string][]korrel8r.Store
	storeConfigs  map[string][]korrel8r.StoreConfig
	rules         []korrel8r.Rule
	templateFuncs template.FuncMap
}

func New(domains ...korrel8r.Domain) *Engine {
	e := &Engine{
		domains:      slices.Clone(domains), // Predicatable order for Domains()
		domainMap:    map[string]korrel8r.Domain{},
		stores:       map[string][]korrel8r.Store{},
		storeConfigs: map[string][]korrel8r.StoreConfig{},
	}
	// FIXME document template funcs
	e.templateFuncs = template.FuncMap{
		"get":       e.get,
		"className": korrel8r.ClassName,
		"ruleName":  korrel8r.RuleName,
	}

	maps.Copy(e.templateFuncs, sprig.TxtFuncMap())

	for _, d := range domains {
		e.domainMap[d.Name()] = d
		e.addTemplateFuncs(d)
	}
	return e
}

// Domain returns the named domain or nil if not found.
func (e *Engine) Domain(name string) korrel8r.Domain { return e.domainMap[name] }
func (e *Engine) Domains() []korrel8r.Domain         { return e.domains }
func (e *Engine) DomainErr(name string) (korrel8r.Domain, error) {
	if d := e.Domain(name); d != nil {
		return d, nil
	}
	return nil, korrel8r.DomainNotFoundErr{Domain: name}
}

// StoresFor returns the known stores for a domain.
func (e *Engine) StoresFor(d korrel8r.Domain) []korrel8r.Store { return e.stores[d.Name()] }

// StoreConfigsFor returns store configurations added with AddStoreConfig
func (e *Engine) StoreConfigsFor(d korrel8r.Domain) []korrel8r.StoreConfig {
	return e.storeConfigs[d.Name()]
}

// StoreErr returns the default (first) store for domain, or an error.
func (e *Engine) StoreErr(d korrel8r.Domain) (korrel8r.Store, error) {
	stores := e.StoresFor(d)
	if len(stores) == 0 {
		return nil, korrel8r.StoreNotFoundErr{Domain: d}
	}
	return stores[0], nil
}

// TemplateFuncser can be implemented by Domain or Store implementations to contribute
// domain-specific template functions to template rules generated by the Engine.
// See text/template.Template.Funcs for details.
type TemplateFuncser interface{ TemplateFuncs() map[string]any }

// AddStore adds a store to the engine.
func (e *Engine) AddStore(s korrel8r.Store) error {
	domain := s.Domain().Name()
	e.stores[domain] = append(e.stores[domain], s)
	e.addTemplateFuncs(s)
	return nil
}

// AddStoreConfig saves the store configuration and creates a store.
//
// If there is an error, it is returned, and the error key in the configuration is set.
func (e *Engine) AddStoreConfig(sc korrel8r.StoreConfig) (err error) {
	defer func() {
		if err != nil {
			sc[korrel8r.StoreKeyError] = err.Error()
		}
	}()
	d, err := e.DomainErr(sc[korrel8r.StoreKeyDomain])
	if err != nil {
		return err
	}
	e.storeConfigs[d.Name()] = append(e.storeConfigs[d.Name()], sc)
	if err := e.expandStoreConfig(sc); err != nil {
		return err
	}
	store, err := d.Store(sc)
	if err != nil {
		return err
	}
	if err := e.AddStore(store); err != nil {
		return err
	}
	return nil
}

// expandStoreConfig expands templates in store config values providing the engine's
func (e *Engine) expandStoreConfig(sc korrel8r.StoreConfig) error {
	for k, v := range sc {
		t, err := template.New(k + ": " + v).Funcs(e.TemplateFuncs()).Parse(v)
		if err != nil {
			return err
		}
		w := &strings.Builder{}
		err = t.Execute(w, nil)
		if err != nil {
			return err
		}
		sc[k] = w.String()
	}
	return nil
}

func (e *Engine) addTemplateFuncs(v any) {
	// Stores and Domains may implement TemplateFuncser if they provide template helper functions for rules
	if tf, ok := v.(TemplateFuncser); ok {
		maps.Copy(e.templateFuncs, tf.TemplateFuncs())
	}
}

// Class parses a full class name and returns the
func (e *Engine) Class(fullname string) (korrel8r.Class, error) {
	d, c, ok := korrel8r.SplitClassName(fullname)
	if !ok {
		return nil, fmt.Errorf("invalid class name: %v", fullname)
	} else {
		return e.DomainClass(d, c)
	}
}

func (e *Engine) DomainClass(domain, class string) (korrel8r.Class, error) {
	d, err := e.DomainErr(domain)
	if err != nil {
		return nil, err
	}
	c := d.Class(class)
	if c == nil {
		return nil, korrel8r.ClassNotFoundErr{Class: class, Domain: d}
	}
	return c, nil
}

// Query parses a query string to a query object.
func (e *Engine) Query(query string) (korrel8r.Query, error) {
	d, _, _, ok := korrel8r.SplitClassData(query)
	if !ok {
		return nil, fmt.Errorf("invalid query string: %v", query)
	}
	domain, err := e.DomainErr(d)
	if err != nil {
		return nil, err
	}
	return domain.Query(query)
}

func (e *Engine) Rules() []korrel8r.Rule { return e.rules }

func (e *Engine) AddRules(rules ...korrel8r.Rule) { e.rules = append(e.rules, rules...) }

// Graph creates a new graph of the rules and classes of this engine.
func (e *Engine) Graph() *graph.Graph { return graph.NewData(e.rules...).NewGraph() }

// TemplateFuncs returns template helper functions for stores and domains known to this engine.
// See text/template.Template.Funcs
func (e *Engine) TemplateFuncs() map[string]any { return e.templateFuncs }

// Get finds the store for the query.Class() and gets into result.
func (e *Engine) Get(ctx context.Context, q korrel8r.Query, c *korrel8r.Constraint, result korrel8r.Appender) error {
	for _, store := range e.StoresFor(q.Class().Domain()) {
		if err := store.Get(ctx, q, c, result); err != nil {
			return err
		}
	}
	return nil
}

// Follower creates a follower. Constraint can be nil.
func (e *Engine) Follower(ctx context.Context, c *korrel8r.Constraint) *Follower {
	return &Follower{Engine: e, Context: ctx, Constraint: c}
}

// FIXME Document template funs.
//
//	get DOMAIN:CLASS:QUERY
//	  Executes QUERY and returns a list of result objects. Type of results depends DOMAIN and CLASS.
func (e *Engine) get(query string) ([]korrel8r.Object, error) {
	q, err := e.Query(query)
	if err != nil {
		return nil, err
	}
	results := korrel8r.NewResult(q.Class())
	err = e.Get(context.Background(), q, nil, results)
	return results.List(), err
}
