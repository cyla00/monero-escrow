// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/cyla00/monero-escrow/components"

func Signin() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" class=\"no-js\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>FideXMR Login</title><link rel=\"stylesheet\" href=\"/assets/css/styles.css\"><link rel=\"stylesheet\" href=\"/assets/css/print.css\" media=\"print\"><meta name=\"description\" content=\"Page description\"><meta property=\"og:title\" content=\"Unique page title - My Site\"><meta property=\"og:description\" content=\"Page description\"><meta property=\"og:image\" content=\"https://www.fidexmr.com/image.jpg\"><meta property=\"og:image:alt\" content=\"Image description\"><meta property=\"og:locale\" content=\"en_GB\"><meta property=\"og:type\" content=\"website\"><meta name=\"twitter:card\" content=\"summary_large_image\"><meta property=\"og:url\" content=\"https://www.fidexmr.com/page\"><link rel=\"canonical\" href=\"https://www.fidexmr.com/page\"><link rel=\"icon\" href=\"/favicon.ico\"><link rel=\"icon\" href=\"/favicon.svg\" type=\"image/svg+xml\"><link rel=\"apple-touch-icon\" href=\"/apple-touch-icon.png\"><link rel=\"manifest\" href=\"/my.webmanifest\"><meta name=\"theme-color\" content=\"#727c83\"></head><body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.Header().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("sign in</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}