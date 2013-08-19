package gocli

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type Router struct {
	Actions map[string]*Action
}

func (cli *Router) Register(path string, action *Action) {
	if cli.Actions == nil {
		cli.Actions = make(map[string]*Action)
	}
	cli.Actions[path] = action
}

func (cli *Router) Usage() string {
	keys := []string{}
	maxParts := 0
	for key := range cli.Actions {
		partsCount := len(strings.Split(key, "/"))
		if partsCount > maxParts {
			maxParts = partsCount
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := []string{"USAGE"}
	table := NewTable()
	for _, key := range keys {
		parts := strings.Split(key, "/")
		action := cli.Actions[key]

		// fill up parts
		for i := (maxParts - len(parts)); i > 0; i-- {
			parts = append(parts, "")
		}

		parts = append(parts, action.Usage, action.Description)
		table.Add(parts...)
		if action.Args != nil {
			usage := action.Args.Usage()
			if usage != "" {
				lines := strings.Split(usage, "\n")
				for _, line := range lines {
					usageParts := []string{}
					for j := 0; j < 3; j++ {
						usageParts = append(usageParts, "")
					}
					current := append(usageParts, line)
					table.Add(current...)
				}
			}
		}
	}
	out = append(out, table.String())
	return strings.Join(out, "\n")
}

func AddActionUsage(parts []string, table *Table, action *Action) {
	parts = append(parts, action.Usage, action.Description)
	table.Add(parts...)
	if action.Args != nil {
		usage := action.Args.Usage()
		if usage != "" {
			lines := strings.Split(usage, "\n")
			for _, line := range lines {
				usageParts := []string{}
				for j := 0; j < 3; j++ {
					usageParts = append(usageParts, "")
				}
				current := append(usageParts, line)
				table.Add(current...)
			}
		}
	}
}

func (cli *Router) Handle(raw []string) error {
	for i := len(raw); i > 0; i-- {
		parts := raw[1:i]
		path := strings.Join(parts, "/")
		if action, ok := cli.Actions[path]; ok {
			args := action.Args
			if args == nil {
				args = &Args{}
			}
			e := args.Parse(raw[i:])
			if e == nil {
				e = action.Handler(args)
			}
			if e != nil {
				table := NewTable()
				AddActionUsage(parts, table, action)
				fmt.Println(table.String())
				os.Exit(0)
			}
			return nil
		}
	}
	fmt.Println(cli.Usage())
	return nil
}