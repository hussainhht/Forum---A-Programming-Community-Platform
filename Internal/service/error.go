package service

import (
	// "forum/Internal/model"
	// "html/template"
	"net/http"
)

func (service *Service) HandleError(w http.ResponseWriter, err int) {
	// if err == 0 {
	// 	return
	// }
	// // http.Error(w, http.StatusText(err), err)
	// tmpl, tmplErr := template.ParseFiles("./web/templates/error.html")
	// if tmplErr != nil {
	// 	http.Error(w, "Failed to load error template", http.StatusInternalServerError)
	// 	return
	// }
	// pageData := model.PageData{
	// 	ErrorMsg:  http.StatusText(err),
	// 	ErrorCode: err,
	// }
	// w.WriteHeader(err)
	// tmpl.Execute(w, pageData)

}
