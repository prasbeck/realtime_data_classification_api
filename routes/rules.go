package routes

import (
    "encoding/json"
    "net/http"

    "rtmc/dsl"
)

func CreateRuleHandler(w http.ResponseWriter, r *http.Request) {
    var rule dsl.RuleStorage
    if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if err := dsl.NewRuleStorage().AddRule(rule); err != nil {
        if err == dsl.ErrMaxRulesExceeded {
            http.Error(w, "Maximum number of rules exceeded", http.StatusBadRequest)
        } else {
            http.Error(w, "Error adding rule", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(rule)
}

func GetRulesHandler(w http.ResponseWriter, r *http.Request) {
    rules := dsl.ruleStorage.GetRules()
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rules)
}

func GetRuleHandler(w http.ResponseWriter, r *http.Request) {
    ruleName := r.URL.Query().Get("name")
    rule, err := dsl.ruleStorage.GetRule(ruleName)
    if err != nil {
        if err == dsl.ErrRuleNotFound {
            http.Error(w, "Rule not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error retrieving rule", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rule)
}

func UpdateRuleHandler(w http.ResponseWriter, r *http.Request) {
    ruleName := r.URL.Query().Get("name")
    var newRule dsl.Rule
    if err := json.NewDecoder(r.Body).Decode(&newRule); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if err := dsl.ruleStorage.UpdateRule(ruleName, newRule); err != nil {
        if err == dsl.ErrRuleNotFound {
            http.Error(w, "Rule not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error updating rule", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(newRule)
}

func DeleteRuleHandler(w http.ResponseWriter, r *http.Request) {
    ruleName := r.URL.Query().Get("name")
    if err := dsl.ruleStorage.DeleteRule(ruleName); err != nil {
        if err == dsl.ErrRuleNotFound {
            http.Error(w, "Rule not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error deleting rule", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
