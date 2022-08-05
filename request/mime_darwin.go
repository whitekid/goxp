//go:build darwin

package request

func init() {
	setMimeTypes(map[string]string{
		".vcf": "text/vcard; charset=utf-8",
	})
}
