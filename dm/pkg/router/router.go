// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"regexp"
	"strings"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb-tools/pkg/filter"
	tablerouter "github.com/pingcap/tidb-tools/pkg/table-router"

	"github.com/pingcap/tiflow/dm/pkg/terror"
)

type (
	TableRule       = tablerouter.TableRule
	Table           = filter.Table
	FilterRule      = filter.Rules
	TableExtractor  = tablerouter.TableExtractor
	SchemaExtractor = tablerouter.SchemaExtractor
	SourceExtractor = tablerouter.SourceExtractor
)

type FilterType = int32

const (
	TblFilter FilterType = iota + 1
	SchmFilter
)

type filterWrapper struct {
	filter *filter.Filter
	typ    FilterType
	target Table

	rawRule *TableRule
}

type RouteTable struct {
	filters       []*filterWrapper
	caseSensitive bool
}

func NewRouter(caseSensitive bool, rules []*TableRule) (*RouteTable, error) {
	r := &RouteTable{
		filters:       make([]*filterWrapper, 0),
		caseSensitive: caseSensitive,
	}
	for _, rule := range rules {
		if err := r.AddRule(rule); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (r *RouteTable) AddRule(rule *TableRule) error {
	err := rule.Valid()
	if err != nil {
		return errors.Trace(err)
	}
	if !r.caseSensitive {
		rule.ToLower()
	}
	newFilter := &filterWrapper{
		rawRule: rule,
	}
	newFilter.target = Table{
		Schema: rule.TargetSchema,
		Name:   rule.TargetTable,
	}
	if len(rule.TablePattern) == 0 {
		// raw schema rule
		newFilter.typ = SchmFilter
		rawFilter, err := filter.New(r.caseSensitive, &FilterRule{
			DoDBs: []string{rule.SchemaPattern},
		})
		if err != nil {
			return errors.Annotatef(err, "add rule %+v into table router", rule)
		}
		newFilter.filter = rawFilter
	} else {
		newFilter.typ = TblFilter
		rawFilter, err := filter.New(r.caseSensitive, &FilterRule{
			DoTables: []*Table{
				{
					Schema: rule.SchemaPattern,
					Name:   rule.TablePattern,
				},
			},
			DoDBs: []string{rule.SchemaPattern},
		})
		if err != nil {
			return errors.Annotatef(err, "add rule %+v into table router", rule)
		}
		newFilter.filter = rawFilter
	}
	r.filters = append(r.filters, newFilter)
	return nil
}

func (r *RouteTable) Route(schema, table string) (string, string, error) {
	curTable := &Table{
		Schema: schema,
		Name:   table,
	}
	tblRules := make([]*filterWrapper, 0)
	schmRules := make([]*filterWrapper, 0)
	for _, filterWrapper := range r.filters {
		if filterWrapper.filter.Match(curTable) {
			if filterWrapper.typ == TblFilter {
				tblRules = append(tblRules, filterWrapper)
			} else {
				schmRules = append(schmRules, filterWrapper)
			}
		}
	}
	var (
		targetSchema string
		targetTable  string
	)
	if table == "" || len(tblRules) == 0 {
		// 1. no need to match table or
		// 2. match no table
		if len(schmRules) > 1 {
			return "", "", terror.ErrWorkerRouteTableDupMatch.Generate(schema, table)
		}
		if len(schmRules) == 1 {
			targetSchema, targetTable = schmRules[0].target.Schema, schmRules[0].target.Name
		}
	} else {
		if len(tblRules) > 1 {
			return "", "", terror.ErrWorkerRouteTableDupMatch.Generate(schema, table)
		}
		targetSchema, targetTable = tblRules[0].target.Schema, tblRules[0].target.Name
	}
	if len(targetSchema) == 0 {
		targetSchema = schema
	}
	if len(targetTable) == 0 {
		targetTable = table
	}
	return targetSchema, targetTable, nil
}

func (r *RouteTable) AllRules() ([]TableRule, []TableRule) {
	var (
		schmRouteRules  []TableRule
		tableRouteRules []TableRule
	)
	for _, filter := range r.filters {
		if filter.typ == SchmFilter {
			schmRouteRules = append(schmRouteRules, *filter.rawRule)
		} else {
			tableRouteRules = append(tableRouteRules, *filter.rawRule)
		}
	}
	return schmRouteRules, tableRouteRules
}

func (r *RouteTable) FetchExtendColumn(schema, table, source string) ([]string, []string) {
	var cols []string
	var vals []string
	rules := []*filterWrapper{}
	curTable := &Table{
		Schema: schema,
		Name:   table,
	}
	for _, filter := range r.filters {
		if filter.filter.Match(curTable) {
			rules = append(rules, filter)
		}
	}
	var (
		schemaRules = make([]*TableRule, 0, len(rules))
		tableRules  = make([]*TableRule, 0, len(rules))
	)
	for i := range rules {
		rule := rules[i].rawRule
		if rule.TablePattern == "" {
			schemaRules = append(schemaRules, rule)
		} else {
			tableRules = append(tableRules, rule)
		}
	}
	if len(tableRules) == 0 && len(schemaRules) == 0 {
		return cols, vals
	}
	var rule *TableRule
	if len(tableRules) == 0 {
		rule = schemaRules[0]
	} else {
		rule = tableRules[0]
	}
	if rule.TableExtractor != nil {
		cols = append(cols, rule.TableExtractor.TargetColumn)
		vals = append(vals, extractVal(table, rule.TableExtractor))
	}

	if rule.SchemaExtractor != nil {
		cols = append(cols, rule.SchemaExtractor.TargetColumn)
		vals = append(vals, extractVal(schema, rule.SchemaExtractor))
	}

	if rule.SourceExtractor != nil {
		cols = append(cols, rule.SourceExtractor.TargetColumn)
		vals = append(vals, extractVal(source, rule.SourceExtractor))
	}
	return cols, vals
}

func extractVal(s string, ext interface{}) string {
	var params []string
	switch e := ext.(type) {
	case *tablerouter.TableExtractor:
		if regExpr, err := regexp.Compile(e.TableRegexp); err == nil {
			params = regExpr.FindStringSubmatch(s)
		}
	case *tablerouter.SchemaExtractor:
		if regExpr, err := regexp.Compile(e.SchemaRegexp); err == nil {
			params = regExpr.FindStringSubmatch(s)
		}
	case *tablerouter.SourceExtractor:
		if regExpr, err := regexp.Compile(e.SourceRegexp); err == nil {
			params = regExpr.FindStringSubmatch(s)
		}
	}
	var val strings.Builder
	for idx, param := range params {
		if idx > 0 {
			val.WriteString(param)
		}
	}
	return val.String()
}
