package exporter

import (
	"echo-react-serve/config"
	"echo-react-serve/server/models/dto"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExportToExcel exports the ScholarshipFormResponse to an Excel file
func FormsToExcel(forms []dto.ScholarshipFormResponse, out io.Writer) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("[EXCEL] Failed to close excel file:", err)
		}
	}()
	sheetName := "Sheet1"

	// Write headers
	headers := []string{
		"Nama", "NIK", "NISN", "Asal Sekolah", "Tempat Lahir", "Tanggal Lahir", "Alamat",
		"Nama Ayah", "NIK Ayah", "Pekerjaan Ayah", "Penghasilan Ayah",
		"Nama Ibu", "NIK Ibu", "Pekerjaan Ibu", "Penghasilan Ibu",
		"Nomor HP", "Tabungan 3 Bulan Terakhir", "Daya Listrik",
		"Tujuan Kampus", "Tujuan Negara", "Tujuan Prodi", "Besaran Biaya Studi", "Status LoA", "Deadline Konfirmasi", "Deadline Pembayaran",
		"Nomor Peserta Program Sosial", "", "", "", "",
		"jumlah Keluarga Inti",
		"Target", "", "", "", "",
		"Catatan", "Catatan LoA",
		"File-file Pendukung", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
	}
	nestedHeaders := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"PIP", "PKH", "KKS", "DTKS", "PPKE", "", "",
		"Universitas", "Negara", "Besaran Biaya Studi", "Tanggal Pengumuman", "Status", "",
		"Loa", "Biaya Studi", "Akta Lahir", "PIP", "PKH", "KKS", "DTKS", "PPKE", "Slip Gaji Ayah", "Slip Gaji Ibu", "SPT Ayah", "SPT Ibu", "Rekening Koran", "Meteran Listrik", "SPTJM",
	}
	f.SetSheetRow(sheetName, "A1", &headers)
	f.SetSheetRow(sheetName, "A2", &nestedHeaders)

	// Merge parent headers
	f.MergeCell(sheetName, "Z1", "AD1")  // Merge "Nomor Peserta Program Sosial" header
	f.MergeCell(sheetName, "AF1", "AJ1") // Merge "Target" header
	f.MergeCell(sheetName, "AM1", "BA1") // Merge "File-file Pendukung" header

	// Add bottom border to the nested header row (row 2)
	for col := 1; col <= len(headers); col++ {
		cell, _ := excelize.CoordinatesToCellName(col, 2)
		addBottomBorder(f, sheetName, cell)
	}

	// Write data
	rowIndex := 3 // Start from row 3 (after headers)
	for _, form := range forms {
		writeDataToSheet(f, sheetName, &rowIndex, form)
		fmt.Println("current row: ", rowIndex)
		for col := 1; col <= len(headers); col++ {
			cell, _ := excelize.CoordinatesToCellName(col, rowIndex-1)
			addBottomBorder(f, sheetName, cell)
		}
	}

	// write the file to output
	return f.Write(out)
}

func addBottomBorder(f *excelize.File, sheet, cell string) error {
	// Retrieve the existing style ID for the cell
	prevStyleId, err := f.GetCellStyle(sheet, cell)
	if err != nil {
		return fmt.Errorf("error getting cell style for cell %s: %v", cell, err)
	}

	// Retrieve the existing style definition
	style, err := f.GetStyle(prevStyleId)
	if err != nil {
		return fmt.Errorf("error getting style for cell %s: %v", cell, err)
	}

	// If no existing style, create a new one
	if style == nil {
		style = &excelize.Style{}
	}

	// Add the bottom border to the style
	if style.Border == nil {
		style.Border = []excelize.Border{}
	}
	style.Border = append(style.Border, excelize.Border{Type: "bottom", Color: "000000", Style: 1})

	// Create a new style with the updated properties
	newStyleId, err := f.NewStyle(style)
	if err != nil {
		return fmt.Errorf("error creating new style for cell %s: %v", cell, err)
	}

	// Apply the new style to the cell
	if err := f.SetCellStyle(sheet, cell, cell, newStyleId); err != nil {
		return fmt.Errorf("error setting style for cell %s: %v", cell, err)
	}

	return nil
}

func toFilesValue(files []dto.FIleResponse) string {
	var fileName string
	if len(files) == 0 {
		return "-"
	}
	if len(files) == 1 {
		fileName = files[0].Name
	} else {
		// get the file type from one of the files
		fileName = strings.Split(files[0].Path, "/")[2] + ".zip" // [bucket]/[identifier]/[file_type]/[file_name]
	}
	var paths []string
	for _, file := range files {
		// fileUrl := config.Envs.App.FileServerUrl + file.Path
		paths = append(paths, file.Path)
	}
	return fmt.Sprintf("=HYPERLINK(\"%s\", \"%s\")", config.Envs.App.FileServerUrl+strings.Join(paths, "[|]"), fileName)
}

func setFilesCell(f *excelize.File, sheet, cell, value string) error {
	// Set the formula in the cell
	if value == "-" {
		return f.SetCellValue(sheet, cell, value)
	}
	if err := f.SetCellFormula(sheet, cell, value); err != nil {
		return fmt.Errorf("error setting formula for cell %s: %v", cell, err)
	}

	// Retrieve the existing style ID for the cell
	prevStyleId, err := f.GetCellStyle(sheet, cell)
	if err != nil {
		return fmt.Errorf("error getting cell style for cell %s: %v", cell, err)
	}

	// Retrieve the existing style definition
	style, err := f.GetStyle(prevStyleId)
	if err != nil {
		return fmt.Errorf("error getting style for cell %s: %v", cell, err)
	}

	// If no existing style, create a new one
	if style == nil {
		style = &excelize.Style{}
	}

	// Update the font properties
	linkStyle := excelize.Font{
		Color:     "#0000FF", // Ensure the color has a # prefix
		Underline: "single",
	}
	style.Font = &linkStyle

	// Create a new style with the updated properties
	newStyleId, err := f.NewStyle(style)
	if err != nil {
		return fmt.Errorf("error creating new style for cell %s: %v", cell, err)
	}

	// Apply the new style to the cell
	if err := f.SetCellStyle(sheet, cell, cell, newStyleId); err != nil {
		return fmt.Errorf("error setting style for cell %s: %v", cell, err)
	}

	return nil
}

// writeDataToSheet writes data to the Excel sheet
func writeDataToSheet(f *excelize.File, sheetName string, rowIndex *int, form dto.ScholarshipFormResponse) {
	// Write Name and Age
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", *rowIndex), form.Nama)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", *rowIndex), form.Nik)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", *rowIndex), form.Nisn)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", *rowIndex), form.AsalSekolah)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", *rowIndex), form.TempatLahir)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", *rowIndex), form.TanggalLahir)
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", *rowIndex), form.Alamat)

	f.SetCellValue(sheetName, fmt.Sprintf("H%d", *rowIndex), form.NamaAyah)
	f.SetCellValue(sheetName, fmt.Sprintf("I%d", *rowIndex), form.NikAyah)
	f.SetCellValue(sheetName, fmt.Sprintf("J%d", *rowIndex), form.PekerjaanAyah)
	f.SetCellValue(sheetName, fmt.Sprintf("K%d", *rowIndex), form.PenghasilanAyah)

	f.SetCellValue(sheetName, fmt.Sprintf("L%d", *rowIndex), form.NamaIbu)
	f.SetCellValue(sheetName, fmt.Sprintf("M%d", *rowIndex), form.NikIbu)
	f.SetCellValue(sheetName, fmt.Sprintf("N%d", *rowIndex), form.PekerjaanIbu)
	f.SetCellValue(sheetName, fmt.Sprintf("O%d", *rowIndex), form.PenghasilanIbu)

	f.SetCellValue(sheetName, fmt.Sprintf("P%d", *rowIndex), form.Nohp)
	f.SetCellValue(sheetName, fmt.Sprintf("Q%d", *rowIndex), form.Tabungan)
	f.SetCellValue(sheetName, fmt.Sprintf("R%d", *rowIndex), form.DayaListrik)

	f.SetCellValue(sheetName, fmt.Sprintf("S%d", *rowIndex), form.TujuanKampus)
	f.SetCellValue(sheetName, fmt.Sprintf("T%d", *rowIndex), form.TujuanNegara)
	f.SetCellValue(sheetName, fmt.Sprintf("U%d", *rowIndex), form.TujuanProdi)

	f.SetCellValue(sheetName, fmt.Sprintf("V%d", *rowIndex), form.BesaranBiayaStudi)
	f.SetCellValue(sheetName, fmt.Sprintf("W%d", *rowIndex), form.StatusLoa)
	f.SetCellValue(sheetName, fmt.Sprintf("X%d", *rowIndex), form.DeadlineKonfirmasi)
	f.SetCellValue(sheetName, fmt.Sprintf("Y%d", *rowIndex), form.DeadlinePembayaran)

	f.SetCellValue(sheetName, fmt.Sprintf("Z%d", *rowIndex), form.NomorPeserta.Pip)
	f.SetCellValue(sheetName, fmt.Sprintf("AA%d", *rowIndex), form.NomorPeserta.Pkh)
	f.SetCellValue(sheetName, fmt.Sprintf("AB%d", *rowIndex), form.NomorPeserta.Kks)
	f.SetCellValue(sheetName, fmt.Sprintf("AC%d", *rowIndex), form.NomorPeserta.Dtks)
	f.SetCellValue(sheetName, fmt.Sprintf("AD%d", *rowIndex), form.NomorPeserta.Ppke)

	f.SetCellValue(sheetName, fmt.Sprintf("AE%d", *rowIndex), form.JumlahKeluarga)

	for i, target := range form.Target {
		f.SetCellValue(sheetName, fmt.Sprintf("AF%d", *rowIndex+i), target.Universitas)
		f.SetCellValue(sheetName, fmt.Sprintf("AG%d", *rowIndex+i), target.Negara)
		f.SetCellValue(sheetName, fmt.Sprintf("AH%d", *rowIndex+i), target.BesaranBiayaStudi)
		f.SetCellValue(sheetName, fmt.Sprintf("AI%d", *rowIndex+i), target.TglPengumuman)
		f.SetCellValue(sheetName, fmt.Sprintf("AJ%d", *rowIndex+i), target.Status)
	}

	f.SetCellValue(sheetName, fmt.Sprintf("AK%d", *rowIndex), form.Catatan)
	f.SetCellValue(sheetName, fmt.Sprintf("AL%d", *rowIndex), form.CatatanLoa)

	// Write files data
	setFilesCell(f, sheetName, fmt.Sprintf("AM%d", *rowIndex), toFilesValue(form.Files.Loa))
	setFilesCell(f, sheetName, fmt.Sprintf("AN%d", *rowIndex), toFilesValue(form.Files.BiayaStudi))
	setFilesCell(f, sheetName, fmt.Sprintf("AO%d", *rowIndex), toFilesValue(form.Files.AktaLahir))
	setFilesCell(f, sheetName, fmt.Sprintf("AP%d", *rowIndex), toFilesValue(form.Files.Program.Pip))
	setFilesCell(f, sheetName, fmt.Sprintf("AQ%d", *rowIndex), toFilesValue(form.Files.Program.Pkh))
	setFilesCell(f, sheetName, fmt.Sprintf("AR%d", *rowIndex), toFilesValue(form.Files.Program.Kks))
	setFilesCell(f, sheetName, fmt.Sprintf("AS%d", *rowIndex), toFilesValue(form.Files.Program.Dtks))
	setFilesCell(f, sheetName, fmt.Sprintf("AT%d", *rowIndex), toFilesValue(form.Files.Program.Ppke))
	setFilesCell(f, sheetName, fmt.Sprintf("AU%d", *rowIndex), toFilesValue(form.Files.SlipAyah))
	setFilesCell(f, sheetName, fmt.Sprintf("AV%d", *rowIndex), toFilesValue(form.Files.SlipIbu))
	setFilesCell(f, sheetName, fmt.Sprintf("AW%d", *rowIndex), toFilesValue(form.Files.SptAyah))
	setFilesCell(f, sheetName, fmt.Sprintf("AX%d", *rowIndex), toFilesValue(form.Files.SptIbu))
	setFilesCell(f, sheetName, fmt.Sprintf("AY%d", *rowIndex), toFilesValue(form.Files.RekeningKoran))
	setFilesCell(f, sheetName, fmt.Sprintf("AZ%d", *rowIndex), toFilesValue(form.Files.MeteranListrik))
	setFilesCell(f, sheetName, fmt.Sprintf("BA%d", *rowIndex), toFilesValue(form.Files.SPTJM))

	// move to the next row
	*rowIndex += len(form.Target)
}
