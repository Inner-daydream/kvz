package templating

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/inner-daydream/kvz/internal/kv"
	"gopkg.in/yaml.v2"
)

type TemplatingService interface {
	Render(templateContent string) (Template, error)
}

func walkParseTree(node parse.Node, visit func(parse.Node)) {
	visit(node)
	switch node := node.(type) {
	case *parse.ListNode:
		for _, n := range node.Nodes {
			walkParseTree(n, visit)
		}
	case *parse.IfNode:
		walkParseTree(node.List, visit)
		if node.ElseList != nil {
			walkParseTree(node.ElseList, visit)
		}
	case *parse.RangeNode:
		walkParseTree(node.List, visit)
		if node.ElseList != nil {
			walkParseTree(node.ElseList, visit)
		}
	case *parse.WithNode:
		walkParseTree(node.List, visit)
		if node.ElseList != nil {
			walkParseTree(node.ElseList, visit)
		}
	}
}

func getTemplateVars(templateContent string) ([]string, error) {
	tpl, err := template.New("template").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var templateVars []string
	visitNode := func(node parse.Node) {
		if node == nil {
			return
		}
		switch n := node.(type) {
		case *parse.ActionNode:
			for _, cmd := range n.Pipe.Cmds {
				for _, arg := range cmd.Args {
					switch arg := arg.(type) {
					case *parse.FieldNode:
						templateVars = append(templateVars, arg.Ident...)

					case *parse.IdentifierNode:
						templateVars = append(templateVars, arg.Ident)
					}
				}
			}
		}
	}

	walkParseTree(tpl.Tree.Root, visitNode)
	return templateVars, nil
}

type TemplateMetadata struct {
	RenderLocation string `yaml:"render_location"`
}

type Template struct {
	Content  string
	Metadata TemplateMetadata
}

type templatingService struct {
	s kv.KvService
}

func (s *templatingService) Render(templateContent string) (Template, error) {
	parts := strings.SplitN(templateContent, "\n---\n", 2)
	if len(parts) < 2 {
		return Template{}, fmt.Errorf("template does not contain metadata")
	}
	var metadata TemplateMetadata
	if err := yaml.Unmarshal([]byte(parts[0]), &metadata); err != nil {
		return Template{}, fmt.Errorf("failed to parse template metadata: %w", err)
	}

	templateVars, err := getTemplateVars(parts[1])
	if err != nil {
		return Template{}, err
	}

	data := make(map[string]interface{})
	for _, varName := range templateVars {
		value, err := s.s.Get(varName)
		if err != nil {
			return Template{}, fmt.Errorf("failed to get value for variable %s: %w", varName, err)
		}
		data[varName] = value
	}

	tpl, err := template.New("template").Parse(parts[1])
	if err != nil {
		return Template{}, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return Template{}, fmt.Errorf("failed to execute template: %w", err)
	}

	return Template{
		Content:  buf.String(),
		Metadata: metadata,
	}, nil
}

func NewService(s kv.KvService) TemplatingService {
	return &templatingService{s}
}
