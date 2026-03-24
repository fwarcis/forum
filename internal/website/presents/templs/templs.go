package templs

import (
	"bufio"
	"html/template"
	"io"
	"unicode"

	"forum/internal/common/system/projfs"
)

var templatesDir = projfs.TemplatesDir()

var layout = template.Must(template.ParseFiles(
	string(templatesDir.Join("layout.html"))))

type layoutData struct {
	any

	Title string
}

func newLayoutData(title string, data any) layoutData {
	titleChars := []rune(title)
	if len(titleChars) <= 1 {
		return layoutData{Title: string(title), any: data}
	}
	firstCharUpperTitle := string(unicode.ToUpper(titleChars[0])) + string(titleChars[1:])
	return layoutData{Title: firstCharUpperTitle, any: data}
}

type ErrorData struct {
	Message string
	Details string
}

func ExecuteError(w io.Writer, data ErrorData) error {
	return ExecuteWithLayout(w, "error", data)
}

func ExecuteWithLayout(w io.Writer, blockBaseName string, data any) error {
	pageBuf := bufio.NewWriter(w)
	layoutClone, err := layout.Clone()
	if err != nil {
		return err
	}

	resTmpl, err := layoutClone.ParseFiles(
		string(templatesDir.Join(blockBaseName + ".html")))
	if err != nil {
		return err
	}

	err = resTmpl.ExecuteTemplate(
		pageBuf, "layout.html", newLayoutData(blockBaseName, data))
	if err != nil {
		return err
	}

	err = pageBuf.Flush()
	if err != nil {
		return err
	}
	return nil
}
