package wired

import (
	"fmt"
	"strings"
)

// topoSort orders providers so every provider appears after all providers it
// depends on. It also resolves each provider's DepVars slice to the local
// variable name supplying that parameter (either another provider's VarName
// or an injector parameter name).
//
// Stable property: when several providers are equally ready, the one that
// appeared first in the original wire(...) call wins. This keeps emit output
// deterministic and matches the author's intent ("this is the order I want
// you to consider them").
func topoSort(providers []ProviderSpec, injectorParams []ParamSpec) ([]ProviderSpec, error) {
	if len(providers) == 0 {
		return nil, nil
	}

	typeIdx := make(map[string]int, len(providers))
	for i, p := range providers {
		if existing, dup := typeIdx[p.ReturnType]; dup {
			return nil, fmt.Errorf("wired: duplicate provider for type %q (positions %d and %d)",
				p.ReturnType, existing, i)
		}
		typeIdx[p.ReturnType] = i
	}

	nameUsed := map[string]int{}
	for _, ip := range injectorParams {
		if ip.Name != "" {
			nameUsed[ip.Name] = 1
		}
	}
	for i := range providers {
		base := providers[i].VarName
		if base == "" {
			base = "v"
		}
		name := base
		for nameUsed[name] > 0 {
			nameUsed[base]++
			name = fmt.Sprintf("%s%d", base, nameUsed[base])
		}
		nameUsed[name] = 1
		providers[i].VarName = name
	}

	typeToVar := make(map[string]string, len(providers)+len(injectorParams))
	for _, ip := range injectorParams {
		if ip.Name == "" || ip.Type == "" {
			continue
		}
		typeToVar[ip.Type] = ip.Name
	}
	for _, p := range providers {
		typeToVar[p.ReturnType] = p.VarName
	}

	for i := range providers {
		deps := make([]string, len(providers[i].Params))
		for j, paramType := range providers[i].Params {
			varName, ok := typeToVar[paramType]
			if !ok {
				return nil, fmt.Errorf("wired: factory %s needs %q but no provider or injector param supplies it",
					providers[i].FuncRef, paramType)
			}
			deps[j] = varName
		}
		providers[i].DepVars = deps
	}

	edges := make([][]int, len(providers))
	indeg := make([]int, len(providers))
	for i := range providers {
		for _, paramType := range providers[i].Params {
			j, ok := typeIdx[paramType]
			if !ok {
				continue
			}
			if j == i {
				return nil, fmt.Errorf("wired: factory %s depends on its own output type %q",
					providers[i].FuncRef, paramType)
			}
			edges[j] = append(edges[j], i)
			indeg[i]++
		}
	}

	ready := make([]int, 0, len(providers))
	for i := range providers {
		if indeg[i] == 0 {
			ready = append(ready, i)
		}
	}
	out := make([]ProviderSpec, 0, len(providers))
	for len(ready) > 0 {
		idx := ready[0]
		ready = ready[1:]
		out = append(out, providers[idx])
		for _, next := range edges[idx] {
			indeg[next]--
			if indeg[next] == 0 {
				ready = append(ready, next)
			}
		}
	}
	if len(out) != len(providers) {
		var unsat []string
		for i, deg := range indeg {
			if deg > 0 {
				unsat = append(unsat, providers[i].ReturnType)
			}
		}
		return nil, fmt.Errorf("wired: dependency cycle detected involving %s",
			strings.Join(unsat, ", "))
	}
	return out, nil
}
