// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Index() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Natzone</title><style>\n        body {\n            background: #2d302f;\n            color: #ffffff;\n            font-family: CM-sans-serif;\n            font-size: larger;\n            height: 100vh;\n            display: flex;\n            justify-content: center;\n            align-items: center;\n            margin: 0px;\n        }\n\n        bingus {\n            color: #fe7879;\n            text-decoration: none;\n        }\n\n        a {\n            color: inherit;\n            text-decoration: none;\n            text-align: center;\n        }\n\n        .center {\n            text-align: center;\n        }\n\n        a:focus,\n        a:hover {\n            color: #fe7879;\n        }\n\n        nav {\n            display: grid;\n            grid-row-gap: 1em;\n            grid-template-columns: 1fr 1fr;\n            margin: 0 1em;\n            min-width: 24em;\n            padding: 1em 0;\n        }\n\n        ul {\n            list-style-type: none;\n            margin: 0;\n            padding: 5px;\n            white-space: nowrap;\n        }\n\n        li:first-child {\n            text-align: center;\n        }\n\n        body>* {\n            margin: 15px;\n        }\n    </style></head><body><img src=\"/images/LainLaugh.gif\" width=\"200\" height=\"200\"><div id=\"sysinfo\" style=\"float: right\"><nav><ul><li>nat@navi</li><hr><li><bingus>OS:</bingus> Linux x86_64</li><li><bingus>Resolution:</bingus> 1920x1080</li><li><bingus>CPU:</bingus> 8 core processor</li><li><bingus>RAM:</bingus> 8GB</li><li><bingus>GPU:</bingus> Intel Mesa</li></ul><ul><li>Links</li><hr><li class=\"center\"><a href=\"https://boards.4channel.org/wsg/catalog\">/wsg/</a></li><li class=\"center\"><a href=\"https://inv.tux.pizza/feed/popular\">invidious</a></li><li class=\"center\"><a href=\":3\">blackboard</a></li><li class=\"center\"><a href=\"https://chat.openai.com/auth/login?next=%2Fc%2Fdc8e28c0-aea0-4b00-9fb3-7d85a9b0c9c8\">chatgpt</a></li><li class=\"center\"><a href=\"https://app.element.io/\">Element</a></li></ul></nav></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
