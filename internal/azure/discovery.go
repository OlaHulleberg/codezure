package azure

import (
    "encoding/json"
    "fmt"
    "os/exec"
)

type Subscription struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type OpenAIResource struct {
    Name          string `json:"name"`
    ResourceGroup string `json:"resourceGroup"`
    Location      string `json:"location"`
}

type Deployment struct {
    Name      string `json:"name"`
    ModelName string `json:"properties.modelName"`
}

func ListSubscriptions() ([]Subscription, error) {
    if err := requireAz(); err != nil { return nil, err }
    out, err := exec.Command("az", "account", "list", "-o", "json").Output(); if err != nil { return nil, err }
    var subs []Subscription
    if err := json.Unmarshal(out, &subs); err != nil { return nil, err }
    return subs, nil
}

func ListOpenAIResources(subscription string) ([]OpenAIResource, error) {
    if err := requireAz(); err != nil { return nil, err }
    if err := runCmd("az", "account", "set", "--subscription", subscription); err != nil { return nil, err }
    out, err := exec.Command("az", "cognitiveservices", "account", "list", "--subscription", subscription, "-o", "json").Output(); if err != nil { return nil, err }
    var raw []map[string]any
    if err := json.Unmarshal(out, &raw); err != nil { return nil, err }
    var res []OpenAIResource
    for _, r := range raw {
        kind, _ := r["kind"].(string)
        // Filter for OpenAI and AIServices (unified service that includes OpenAI)
        if kind != "OpenAI" && kind != "AIServices" {
            continue
        }
        name, _ := r["name"].(string)
        group, _ := r["resourceGroup"].(string)
        location, _ := r["location"].(string)
        res = append(res, OpenAIResource{Name: name, ResourceGroup: group, Location: location})
    }
    return res, nil
}

func ListDeployments(subscription, resource, group string) ([]Deployment, error) {
    if err := requireAz(); err != nil { return nil, err }
    if err := runCmd("az", "account", "set", "--subscription", subscription); err != nil { return nil, err }
    out, err := exec.Command("az", "cognitiveservices", "account", "deployment", "list", "--name", resource, "--resource-group", group, "--subscription", subscription, "-o", "json").Output(); if err != nil { return nil, err }
    var raw []map[string]any
    if err := json.Unmarshal(out, &raw); err != nil { return nil, err }
    var deps []Deployment
    for _, d := range raw {
        name, _ := d["name"].(string)
        props, _ := d["properties"].(map[string]any)
        model := ""
        if props != nil {
            if modelObj, ok := props["model"].(map[string]any); ok {
                model, _ = modelObj["name"].(string)
            }
        }
        deps = append(deps, Deployment{Name: name, ModelName: model})
    }
    return deps, nil
}

// No extension management — 'az cognitiveservices' is part of core CLI on modern versions.

func ThinkingLevels() []string { return []string{"low", "medium", "high"} }

func FormatThinkingModel(base string, level string) string {
    if level == "" { return base }
    return fmt.Sprintf("%s:%s", base, level) // simple convention
}
