package webapp

import (
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

func (r *WebAppRouter) bankImportersHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "bank_importers")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	bankimporters, err := r.db.GetBankImporters(userID)
	if err != nil {
		r.logger.Error("Failed to get bank importers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.logger.Info("Bank importers", "bankimporters", bankimporters)

	if req.URL.Query().Get("fetchAll") == "true" {
		for i, bankImporter := range bankimporters {
			if bankImporter.Id == req.URL.Query().Get("id") {
				r.logger.Info("Set 'FetchAll' to true", "id", bankImporter.Id)
				bankImporter.FetchAll = true
				if _, err = r.db.UpdateBankImporter(userID, bankImporter.Id, &bankImporter); err != nil {
					r.logger.Error("Failed to update bank importer", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				bankimporters[i] = bankImporter

				// schedule forced import
				background.GetForcedImportChannel(req.Context()) <- background.ForcedImport{
					UserID:         userID,
					BankImporterID: bankImporter.Id,
				}
			}
		}
	}
	data["BankImporters"] = &bankimporters

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

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

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
