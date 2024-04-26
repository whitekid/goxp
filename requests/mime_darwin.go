//go:build darwin

package requests

func init() {
	setMimeTypes(map[string]string{
		".vcf": "text/vcard; charset=utf-8",
	})
}
