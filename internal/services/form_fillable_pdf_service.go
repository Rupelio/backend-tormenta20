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
	pdf.SetFont("Arial", "", 11)

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
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "FICHA DE PERSONAGEM", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, "TORMENTA20 - FICHA EDITÁVEL", "", 1, "C", false, 0, "")
	pdf.Ln(4)
}

func (s *FormFillablePDFService) addFormBasicInfo(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	// Reset para fonte normal
	pdf.SetFont("Arial", "", 10)

	y := pdf.GetY()

	// Nome do Personagem (campo largo)
	pdf.Text(10, y, "Nome:")
	s.addEditableField(pdf, "character_name", personagem.Nome, 22, y-3, 110, 8)

	// Nível e Experiência à direita
	pdf.Text(140, y, "Nível:")
	s.addEditableField(pdf, "level", strconv.Itoa(personagem.Nivel), 152, y-3, 18, 8)
	pdf.Text(172, y, "Exp:")
	s.addEditableField(pdf, "experience", "", 184, y-3, 25, 8)

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
	s.addEditableField(pdf, "race", racaNome, 22, y-3, 55, 8)

	pdf.Text(85, y, "Classe:")
	s.addEditableField(pdf, "class", classeNome, 97, y-3, 55, 8)

	pdf.Text(150, y, "Origem:")
	s.addEditableField(pdf, "origin", origemNome, 162, y-3, 47, 8)

	pdf.SetY(y + 15)
}

func (s *FormFillablePDFService) addFormAttributes(pdf *gofpdf.Fpdf, personagem *models.Personagem, showCalculations bool) {
	// Título da seção
	pdf.SetFont("Arial", "B", 12)

	y := pdf.GetY()
	pdf.Text(90, y, "ATRIBUTOS")
	pdf.SetY(y + 8)

	// Cabeçalho da tabela
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFont("Arial", "", 9)

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

	// Grid: 2 colunas x 3 linhas
	baseY := pdf.GetY()
	rowH := 10.0
	leftX := 14.0
	rightX := 104.0

	for i, attr := range attributes {
		var col int
		var row int
		if i < 3 {
			col = 0
			row = i
		} else {
			col = 1
			row = i - 3
		}

		yPos := baseY + float64(row)*rowH
		var xPos float64
		if col == 0 {
			xPos = leftX
		} else {
			xPos = rightX
		}

		// Rótulo
		pdf.SetXY(xPos, yPos)
		pdf.Text(xPos, yPos+4, attr.name)

		// Campo valor
		s.addEditableField(pdf, attr.field+"_value", strconv.Itoa(attr.value), xPos+30, yPos+1, 20, 8)

		// Campo modificador
		s.addEditableField(pdf, attr.field+"_mod", fmt.Sprintf("%+d", attr.value), xPos+54, yPos+1, 18, 8)
	}

	// Avançar Y após a grade de atributos
	pdf.SetY(baseY + 3*rowH + 4)
}

func (s *FormFillablePDFService) addFormSkills(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Título da seção
	pdf.SetFont("Arial", "B", 12)

	y := pdf.GetY()
	pdf.Text(85, y, "PERÍCIAS")
	pdf.SetY(y + 10)

	// Cabeçalho
	pdf.SetFont("Arial", "B", 8)

	y = pdf.GetY()
	pdf.Text(15, y, "Perícia")
	pdf.Text(65, y, "Treinada")
	pdf.Text(90, y, "Atributo")
	pdf.Text(120, y, "Bônus")
	pdf.Text(145, y, "Total")

	pdf.SetY(y + 6)

	// Reset para fonte normal
	pdf.SetFont("Arial", "", 7)

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
	pdf.SetFont("Arial", "B", 12)

	y := pdf.GetY()
	pdf.Text(85, y, "INVENTÁRIO")
	pdf.SetY(y + 10)

	// Cabeçalho
	pdf.SetFont("Arial", "B", 9)

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
	pdf.SetFont("Arial", "B", 12)

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
	pdf.SetFont("Arial", "B", 12)

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
