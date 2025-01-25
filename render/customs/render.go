package customs

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

// This file takes in the *os.File from whoever is rendering .html files
// and makes it so that any HTML tags that are in the registeredComponents
// will be rendered according to the ComponentFunc and replaced in the file.

type ComponentFunc func(component) (templ.Component, error)

var registeredComponents = map[string]ComponentFunc{
	"Code":     code,
	"BlogPost": blogPost,
}

type component struct {
	Element    string
	Attributes map[string]string
	Children   string // any nested components are already html by the time their parent gets rendered
}

// processComponents walks through the content, looking for <Component> tags.
// Unregistered or improperly closed tags are output as-is. Registered components
// with proper closing tags are processed.
func RenderCustomComponents(content string) (string, error) {
	var output strings.Builder

	for {
		// Find the next '<'
		startIdx := strings.Index(content, "<")
		if startIdx == -1 {
			// No more tags
			output.WriteString(content)
			break
		}

		// Write everything up to this '<' into the output.
		output.WriteString(content[:startIdx])
		content = content[startIdx:] // Move to the start of the tag.

		var isSelfClosing bool
		// Find the closing '>' of this tag
		endIdx := strings.Index(content, ">")
		if endIdx == -1 {
			// No closing '>', output the rest as is and stop.
			output.WriteString(content)
			break
		}

		// Found "/>". Don't put that slash in the rest of the processing!
		if content[endIdx-1] == '/' {
			endIdx = endIdx - 1
			isSelfClosing = true
		}

		// Check if there's another `<` before we reach `>`
		innerLT := strings.Index(content[1:endIdx], "<")
		if innerLT != -1 {
			// Output just '<' as text and continue.
			output.WriteString("<")
			content = content[1:]
			continue
		}

		// The raw tag (without the < >)
		tag := content[1:endIdx]
		elem := isComponent(tag)
		if elem == "" {
			// Not a recognized component; output verbatim
			output.WriteString(content[:endIdx+1])
			content = content[endIdx+1:]
			continue
		}

		// Get component metadata (element, attributes).
		comp, err := parseStartTag(tag, elem)
		if err != nil {
			// If there's any unexpected parse error, output verbatim and continue.
			output.WriteString(content[:endIdx+1])
			content = content[endIdx+1:]
			continue
		}

		var closingTag string
		if isSelfClosing {
			closingTag = " />"
		} else {
			closingTag = fmt.Sprintf("</%s>", elem)
		}

		// Look for the closing tag for this component (e.g. </Code>).
		closingIdx := strings.Index(content, closingTag)
		if closingIdx == -1 {
			// Proper closing tag not found, output verbatim and skip editing.
			output.WriteString(content[:endIdx+1])
			content = content[endIdx+1:]
			continue
		}

		// Extract the children (the substring after the tag to just before the closing tag).
		comp.Children = content[endIdx+1 : closingIdx]
		content = content[closingIdx+len(closingTag):] // Advance past the closing tag.

		// Recursively process children that may contain further components.
		if strings.Contains(comp.Children, "<") {
			parsedChildren, err := RenderCustomComponents(comp.Children)
			if err != nil {
				return "", err
			}
			comp.Children = parsedChildren
		}

		// Look up the component's rendering function.
		compFunc, ok := registeredComponents[comp.Element]
		if !ok {
			// If somehow not registered, output the raw text we had and move on.
			output.WriteString("<" + tag + ">")
			output.WriteString(comp.Children)
			output.WriteString(closingTag)
			continue
		}

		// Render the component.
		result, err := compFunc(comp)
		if err != nil {
			return "", err
		}

		// Output the rendered component.
		result.Render(context.Background(), &output)
	}

	return output.String(), nil
}

func isComponent(tag string) string {
	for elem := range registeredComponents {
		if strings.HasPrefix(tag, elem) {
			return elem
		}
	}
	return ""
}

func parseStartTag(tag string, elem string) (component, error) {
	comp := component{Element: elem}
	parts := strings.Fields(tag[len(elem):])
	attributes := make(map[string]string)
	for _, attr := range parts {
		kv := strings.SplitN(attr, "=", 2)
		if len(kv) == 2 {
			attributes[kv[0]] = strings.Trim(kv[1], "\"")
		}
	}
	comp.Attributes = attributes
	return comp, nil
}
