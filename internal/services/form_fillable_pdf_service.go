package services

import (
	"bytes"
	"fmt"
	"strconv"

	"tormenta20-builder/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// FormFillablePDFService cria PDFs editáveis como os do D&D Beyond
type FormFillablePDFService struct{}

func NewFormFillablePDFService() *FormFillablePDFService {
	return &FormFillablePDFService{}
}

func (s *FormFillablePDFService) GenerateEditableCharacterSheet(personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)

	if options.Layout == "double" {
		return s.generateDoublePageFormSheet(pdf, personagem, options)
	}

	return s.generateSinglePageFormSheet(pdf, personagem, options)
}

func (s *FormFillablePDFService) generateSinglePageFormSheet(pdf *gofpdf.Fpdf, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	pdf.AddPage()

	// Configurar fonte Arial padrão
	pdf.SetFont("Arial", "", 12)

	// Cabeçalho
	s.addFormHeader(pdf, personagem)

	// Informações básicas com campos editáveis
	s.addFormBasicInfo(pdf, personagem)

	// Atributos com campos editáveis
	s.addFormAttributes(pdf, personagem, options.ShowCalculations)

	// Seções adicionais baseadas nas opções
	currentY := pdf.GetY()
	if currentY < 200 && len(options.ExtraSections) > 0 {
		for _, section := range options.ExtraSections {
			switch section {
			case "skills":
				s.addFormSkills(pdf)
			case "inventory":
				s.addFormInventory(pdf)
			case "notes":
				s.addFormNotes(pdf)
			}
		}
	}

	// Verificar se há erros
	if pdf.Error() != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", pdf.Error())
	}

	// Obter bytes do PDF usando um buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("erro ao escrever PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *FormFillablePDFService) generateDoublePageFormSheet(pdf *gofpdf.Fpdf, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	// Primeira página
	pdf.AddPage()

	// Configurar fonte Arial padrão
	pdf.SetFont("Arial", "", 12)

	s.addFormHeader(pdf, personagem)
	s.addFormBasicInfo(pdf, personagem)
	s.addFormAttributes(pdf, personagem, options.ShowCalculations)
	s.addFormSkills(pdf)

	// Segunda página
	pdf.AddPage()
	s.addFormInventory(pdf)
	s.addFormNotes(pdf)
	s.addFormHistory(pdf)

	if pdf.Error() != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", pdf.Error())
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("erro ao escrever PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *FormFillablePDFService) addFormHeader(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	// Título centralizado
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "FICHA DE PERSONAGEM - TORMENTA20", "", 1, "C", false, 0, "")
	pdf.Ln(5)
}

func (s *FormFillablePDFService) addFormBasicInfo(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	// Reset para fonte normal
	pdf.SetFont("Arial", "", 10)

	y := pdf.GetY()

	// Nome do Personagem
	pdf.Text(10, y, "Nome do Personagem:")
	s.addEditableField(pdf, "character_name", personagem.Nome, 45, y-3, 80, 6)

	// Nível
	pdf.Text(135, y, "Nível:")
	s.addEditableField(pdf, "level", strconv.Itoa(personagem.Nivel), 150, y-3, 20, 6)

	// Experiência
	pdf.Text(175, y, "Experiência:")
	s.addEditableField(pdf, "experience", "", 190, y-3, 25, 6)

	pdf.SetY(y + 10)
	y = pdf.GetY()

	// Segunda linha - Raça, Classe, Origem
	racaNome := ""
	if personagem.RacaID > 0 {
		racaNome = personagem.Raca.Nome
	}

	classeNome := ""
	if personagem.ClasseID > 0 {
		classeNome = personagem.Classe.Nome
	}

	origemNome := ""
	if personagem.OrigemID > 0 {
		origemNome = personagem.Origem.Nome
	}

	pdf.Text(10, y, "Raça:")
	s.addEditableField(pdf, "race", racaNome, 25, y-3, 45, 6)

	pdf.Text(75, y, "Classe:")
	s.addEditableField(pdf, "class", classeNome, 90, y-3, 45, 6)

	pdf.Text(140, y, "Origem:")
	s.addEditableField(pdf, "origin", origemNome, 155, y-3, 45, 6)

	pdf.SetY(y + 15)
}

func (s *FormFillablePDFService) addFormAttributes(pdf *gofpdf.Fpdf, personagem *models.Personagem, showCalculations bool) {
	// Título da seção
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 12)
	} else {
		pdf.SetFont("Arial", "B", 12)
	}

	y := pdf.GetY()
	pdf.Text(90, y, "ATRIBUTOS")
	pdf.SetY(y + 10)

	// Cabeçalho da tabela
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 9)
	} else {
		pdf.SetFont("Arial", "B", 9)
	}

	y = pdf.GetY()
	pdf.Text(15, y, "Atributo")
	pdf.Text(55, y, "Valor")
	pdf.Text(85, y, "Modificador")
	pdf.Text(125, y, "Observações")

	pdf.SetY(y + 8)

	// Reset para fonte normal
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "", 9)
	} else {
		pdf.SetFont("Arial", "", 9)
	}

	// Lista de atributos
	attributes := []struct {
		name  string
		value int
		field string
	}{
		{"Força", personagem.For, "str"},
		{"Destreza", personagem.Des, "dex"},
		{"Constituição", personagem.Con, "con"},
		{"Inteligência", personagem.Int, "int"},
		{"Sabedoria", personagem.Sab, "wis"},
		{"Carisma", personagem.Car, "cha"},
	}

	for _, attr := range attributes {
		y = pdf.GetY()

		// Nome do atributo
		pdf.Text(15, y+4, attr.name)

		// Campo editável para valor
		s.addEditableField(pdf, attr.field+"_value", strconv.Itoa(attr.value), 55, y, 20, 6)

		// Campo editável para modificador (calculado automaticamente)
		modifier := attr.value
		s.addEditableField(pdf, attr.field+"_mod", fmt.Sprintf("%+d", modifier), 85, y, 25, 6)

		// Campo para observações
		s.addEditableField(pdf, attr.field+"_notes", "", 125, y, 70, 6)

		pdf.SetY(y + 8)
	}

	pdf.Ln(5)
}

func (s *FormFillablePDFService) addFormSkills(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Título da seção
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 12)
	} else {
		pdf.SetFont("Arial", "B", 12)
	}

	y := pdf.GetY()
	pdf.Text(85, y, "PERÍCIAS")
	pdf.SetY(y + 10)

	// Cabeçalho
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 8)
	} else {
		pdf.SetFont("Arial", "B", 8)
	}

	y = pdf.GetY()
	pdf.Text(15, y, "Perícia")
	pdf.Text(65, y, "Treinada")
	pdf.Text(90, y, "Atributo")
	pdf.Text(120, y, "Bônus")
	pdf.Text(145, y, "Total")

	pdf.SetY(y + 6)

	// Reset para fonte normal
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "", 7)
	} else {
		pdf.SetFont("Arial", "", 7)
	}

	// Lista de perícias principais (reduzida para caber)
	skills := []string{
		"Acrobacia", "Atletismo", "Atuação", "Conhecimento", "Cura",
		"Diplomacia", "Enganação", "Fortitude", "Furtividade", "Guerra",
		"Iniciativa", "Intimidação", "Intuição", "Investigação", "Luta",
		"Misticismo", "Percepção", "Pontaria", "Reflexos", "Religião",
		"Sobrevivência", "Vontade",
	}

	for i, skill := range skills {
		if pdf.GetY() > 280 {
			break // Evitar sair da página
		}

		y = pdf.GetY()

		// Nome da perícia
		pdf.Text(15, y+3, skill)

		// Checkbox para treinada
		s.addCheckbox(pdf, fmt.Sprintf("skill_%d_trained", i), 68, y, 4, 4)

		// Campo para atributo
		s.addEditableField(pdf, fmt.Sprintf("skill_%d_attr", i), "", 90, y, 20, 4)

		// Campo para bônus
		s.addEditableField(pdf, fmt.Sprintf("skill_%d_bonus", i), "", 120, y, 20, 4)

		// Campo para total
		s.addEditableField(pdf, fmt.Sprintf("skill_%d_total", i), "", 145, y, 20, 4)

		pdf.SetY(y + 5)
	}

	pdf.Ln(5)
}

func (s *FormFillablePDFService) addFormInventory(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Título
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 12)
	} else {
		pdf.SetFont("Arial", "B", 12)
	}

	y := pdf.GetY()
	pdf.Text(85, y, "INVENTÁRIO")
	pdf.SetY(y + 10)

	// Cabeçalho
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 9)
	} else {
		pdf.SetFont("Arial", "B", 9)
	}

	y = pdf.GetY()
	pdf.Text(15, y, "Item")
	pdf.Text(100, y, "Qtd")
	pdf.Text(130, y, "Peso")
	pdf.Text(160, y, "Valor")

	pdf.SetY(y + 8)

	// Linhas do inventário
	for i := 0; i < 15; i++ {
		if pdf.GetY() > 280 {
			break
		}

		y = pdf.GetY()

		s.addEditableField(pdf, fmt.Sprintf("item_%d_name", i), "", 15, y, 80, 6)
		s.addEditableField(pdf, fmt.Sprintf("item_%d_qty", i), "", 100, y, 20, 6)
		s.addEditableField(pdf, fmt.Sprintf("item_%d_weight", i), "", 130, y, 25, 6)
		s.addEditableField(pdf, fmt.Sprintf("item_%d_value", i), "", 160, y, 30, 6)

		pdf.SetY(y + 8)
	}
}

func (s *FormFillablePDFService) addFormNotes(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 230 {
		pdf.AddPage()
	}

	// Título
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 12)
	} else {
		pdf.SetFont("Arial", "B", 12)
	}

	y := pdf.GetY()
	pdf.Text(85, y, "ANOTAÇÕES")
	pdf.SetY(y + 10)

	// Área de texto grande para anotações
	for i := 0; i < 8; i++ {
		if pdf.GetY() > 280 {
			break
		}

		y = pdf.GetY()
		s.addEditableField(pdf, fmt.Sprintf("notes_%d", i), "", 15, y, 180, 6)
		pdf.SetY(y + 8)
	}
}

func (s *FormFillablePDFService) addFormHistory(pdf *gofpdf.Fpdf) {
	// Título
	if pdf.Error() == nil {
		pdf.SetFont("DejaVu", "B", 12)
	} else {
		pdf.SetFont("Arial", "B", 12)
	}

	y := pdf.GetY()
	pdf.Text(75, y, "HISTÓRICO DO PERSONAGEM")
	pdf.SetY(y + 10)

	// Área de texto para histórico
	for i := 0; i < 10; i++ {
		if pdf.GetY() > 280 {
			break
		}

		y = pdf.GetY()
		s.addEditableField(pdf, fmt.Sprintf("history_%d", i), "", 15, y, 180, 6)
		pdf.SetY(y + 8)
	}
}

// addEditableField adiciona um campo de texto editável
func (s *FormFillablePDFService) addEditableField(pdf *gofpdf.Fpdf, name, value string, x, y, w, h float64) {
	// Desenhar bordas do campo
	pdf.Rect(x, y, w, h, "D")

	// Adicionar texto se houver valor
	if value != "" {
		pdf.Text(x+1, y+h-1, value)
	}

	// Nota: gofpdf não suporta nativamente campos de formulário PDF
	// Para campos verdadeiramente editáveis, seria necessário usar uma biblioteca
	// mais avançada como pdfcpu ou uma integração com ferramentas externas
}

// addCheckbox adiciona uma checkbox editável
func (s *FormFillablePDFService) addCheckbox(pdf *gofpdf.Fpdf, name string, x, y, w, h float64) {
	// Desenhar quadrado para checkbox
	pdf.Rect(x, y, w, h, "D")
}
