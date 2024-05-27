package dsl

import (
    "errors"
    "sync"
)

var (
    ErrRuleNotFound       = errors.New("rule not found")
    ErrMaxRulesExceeded   = errors.New("maximum number of rules exceeded")
    maxRules              = 10 // default maximum number of rules
    ruleStorage           = NewRuleStorage()
)

type RuleStorage struct {
    sync.RWMutex
    rules map[string]RuleStorage
}

func NewRuleStorage() *RuleStorage {
    return &RuleStorage{
        rules: make(map[string]RuleStorage),
    }
}

func SetMaxRules(max int) {
    maxRules = max
}

func (rs *RuleStorage) AddRule(rule RuleStorage) error {
    rs.Lock()
    defer rs.Unlock()
    if len(rs.rules) >= maxRules {
        return ErrMaxRulesExceeded
    }
    rs.rules[rule.Name] = rule
    return nil
}

func (rs *RuleStorage) GetRules() []Rule {
    rs.RLock()
    defer rs.RUnlock()
    rules := make([]Rule, 0, len(rs.rules))
    for _, rule := range rs.rules {
        rules = append(rules, rule)
    }
    return rules
}

func (rs *RuleStorage) GetRule(name string) (Rule, error) {
    rs.RLock()
    defer rs.RUnlock()
    rule, exists := rs.rules[name]
    if !exists {
        return Rule{}, ErrRuleNotFound
    }
    return rule, nil
}

func (rs *RuleStorage) UpdateRule(name string, newRule Rule) error {
    rs.Lock()
    defer rs.Unlock()
    if _, exists := rs.rules[name]; !exists {
        return ErrRuleNotFound
    }
    rs.rules[name] = newRule
    return nil
}

func (rs *RuleStorage) DeleteRule(name string) error {
    rs.Lock()
    defer rs.Unlock()
    if _, exists := rs.rules[name]; !exists {
        return ErrRuleNotFound
    }
    delete(rs.rules, name)
    return nil
}
