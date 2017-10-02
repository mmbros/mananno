
// Generated by gentmpl; *** DO NOT EDIT ***
// Created: 2017-10-02 23:11:05
// Params: no_cache=true, no_go_format=true, asset_manager="go-bindata", func_map="funcMap"

package templates

import (
"html/template"
	"io"
	"path/filepath"
)


	// type definitions
	type (
		// templateEnum is the type of the Templates
		templateEnum uint8
		// PageEnum is the type of the Pages
		PageEnum uint8
	)
	// number of templates
	const templatesLen = 4
	// PageEnum constants
	const (
		PageAcestreamidChannel PageEnum = iota
		PageAcestreamidChannels
		PageArenavisionChannel
		PageArenavisionSchedule
		PageCorsaroIndex
		PageTestTransmission
		)
	


func file2path(file string) string {
	const templatesFolder = "tmpl"
	var path string
	switch {
	case len(file) == 0, file[0] == '.', file[0] == filepath.Separator:
		path = file
	default:
		path = filepath.Join(templatesFolder, file)
	}
	return path
}



// Files returns the files used by the `t` template
func (t templateEnum) Files() []string {
	var (
		// files paths
		files = [...]string{ "partials/_header.tmpl", "partials/_footer.tmpl", "acestreamid/main.tmpl", "arenavision/schedule.tmpl", "ilcorsaronero/_base.tmpl", "ilcorsaronero/index.tmpl", "test/transmission.tmpl"  }
		// template-index to array of file-index
		ti2afi = [...][]uint8{
		{ 0, 1, 2 }, // acestreamid
		{ 0, 1, 3 }, // av_schedule
		{ 4, 5 }, // corsaro
		{ 0, 1, 6 }, // test_transmission
		}
	)
	// get the template files indexes
	idxs := ti2afi[t]
	// build the array of files
	astr := make([]string, len(idxs))
	for j, idx := range idxs {
		astr[j] = files[idx]
	}
	return astr
}


// Files returns the files used by the template of the page
func (page PageEnum) Files() []string {
	// from page to template indexes
	var p2t = [...]templateEnum{0, 0, 1, 1, 2, 3}
	// get the template of the page
	t := p2t[page]
	return t.Files()
}


	
// Template returns the template.Template of the page
func (page PageEnum) Template() *template.Template {
files := page.Files()


// use go-bindata MustAsset func to load templates
tmpl := template.New(filepath.Base(files[0])).Funcs(funcMap)
for _, file := range files {
	tmpl.Parse(string(MustAsset(file2path(file))))
}
return tmpl

}



// Base returns the template name of the page
func (page PageEnum) Base() string {
	var bases = [...]string{ "channel", "channels", "schedule", "", "main" }
	
		
	var pi2bi = [...]PageEnum{ 0, 1, 0, 2, 3, 4 }
	return bases[pi2bi[page]]
	
}


// Execute applies a parsed page template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output, execution
// stops, but partial results may already have been written to the output writer.
// A template may be executed safely in parallel.
func (page PageEnum) Execute(wr io.Writer, data interface{}) error {
	tmpl := page.Template()
	name := page.Base()
	if name != "" {
		return tmpl.ExecuteTemplate(wr, name, data)
	}
	return tmpl.Execute(wr, data)
}


/*
func main(){
	var page = PageAcestreamidChannel
	wr := os.Stdout

	if err := page.Execute(wr, nil); err != nil {
		fmt.Print(err)
	}
}
*/

