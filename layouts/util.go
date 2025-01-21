package layouts

import "github.com/a-h/templ"

type LayoutComponent = func(header templ.Component, meta templ.Component) templ.Component
