package webapp

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

func (r *WebAppRouter) bankImportersHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "bank_importers")

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if ok {
		data["UserID"] = userID

		// accounts, err := r.db.GetAccounts(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get accounts", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// currencies, err := r.db.GetCurrencies(userID)
		// if err != nil {
		// 	r.logger.Error("Failed to get currencies", "error", err)
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		bankimporters, err := r.db.GetBankImporters(userID)
		if err != nil {
			r.logger.Error("Failed to get bank importers", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.logger.Info("Bank importers", "bankimporters", bankimporters)

		data["BankImporters"] = &bankimporters
	}

	if err := tmpl.ExecuteTemplate(w, "bank_importers.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//nolint:funlen
func (r *WebAppRouter) bankImporterUploadHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "bank_importers")

	session, _ := r.cookies.Get(req, "session-name")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err = req.ParseMultipartForm(10 << 20); err != nil {
		r.logger.Error("Failed to parse form", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, file, err := req.FormFile("file")
	if err != nil {
		r.logger.Error("Failed to get file", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fileData, fileName, err := readFileHeader(file)
	if err != nil {
		r.logger.Error("Failed to read file", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	format := fileName[strings.LastIndex(fileName, ".")+1:]

	bankImporterID := req.FormValue("id")
	if bankImporterID == "" {
		r.logger.Error("No bank importer ID", "error", err)
		http.Error(w, "No bank importer ID", http.StatusBadRequest)
		return
	}

	parser := api.NewBankImportersAPIServiceImpl(r.logger, r.db)
	lastImport, err := parser.Upload(userID, bankImporterID, format, fileData)
	if err != nil {
		r.logger.Error("Failed to upload bank importer", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["LastImport"] = lastImport
	r.logger.Info("Imported successfully", "lastImport", lastImport)

	if err := tmpl.ExecuteTemplate(w, "bank_importers_upload.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readFileHeader(fileHeader *multipart.FileHeader) ([]byte, string, error) {
	formFile, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	defer formFile.Close()

	fileData, err := io.ReadAll(formFile)
	if err != nil {
		return nil, "", err
	}
	return fileData, fileHeader.Filename, nil
}
