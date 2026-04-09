package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"tormenta20-builder/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// FormFillablePDFService cria PDFs com a ficha padrão do Tormenta 20
type FormFillablePDFService struct{}

func NewFormFillablePDFService() *FormFillablePDFService {
	return &FormFillablePDFService{}
}

func (s *FormFillablePDFService) GenerateEditableCharacterSheet(personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(8, 8, 8)
	pdf.SetAutoPageBreak(true, 10)

	if options.Layout == "double" {
		return s.generateDoublePageFormSheet(pdf, personagem, options)
	}

	return s.generateSinglePageFormSheet(pdf, personagem, options)
}

func (s *FormFillablePDFService) generateSinglePageFormSheet(pdf *gofpdf.Fpdf, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	pdf.AddPage()

	s.addFormHeader(pdf, personagem)
	s.addFormBasicInfo(pdf, personagem)
	s.addCombatStats(pdf, personagem)
	s.addFormAttributes(pdf, personagem, options.ShowCalculations)

	for _, section := range options.ExtraSections {
		switch section {
		case "skills":
			s.addFormSkillsFilled(pdf, personagem)
		case "inventory":
			s.addFormInventory(pdf)
		case "notes":
			s.addFormNotes(pdf)
		case "history":
			s.addFormHistory(pdf)
		}
	}

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

func (s *FormFillablePDFService) generateDoublePageFormSheet(pdf *gofpdf.Fpdf, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	// Pagina 1: Info basica, combate, atributos, pericias
	pdf.AddPage()
	s.addFormHeader(pdf, personagem)
	s.addFormBasicInfo(pdf, personagem)
	s.addCombatStats(pdf, personagem)
	s.addFormAttributes(pdf, personagem, options.ShowCalculations)
	s.addFormSkillsFilled(pdf, personagem)

	// Pagina 2: Poderes, inventario, notas, historico
	pdf.AddPage()
	s.addPowersSection(pdf, personagem)
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

// ========== HEADER ==========

func (s *FormFillablePDFService) addFormHeader(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	// Fundo do titulo
	pdf.SetFillColor(139, 0, 0) // Vermelho escuro (tematico T20)
	pdf.Rect(8, 8, 194, 14, "F")

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(8, 9)
	pdf.CellFormat(194, 6, "FICHA DE PERSONAGEM", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(8, 15)
	pdf.CellFormat(194, 5, "TORMENTA 20", "", 1, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(4)
}

// ========== INFO BASICA ==========

func (s *FormFillablePDFService) addFormBasicInfo(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	y := pdf.GetY()

	// Linha 1: Nome, Nivel, Experiencia
	s.drawLabelField(pdf, "NOME", personagem.Nome, 8, y, 120, 10)
	s.drawLabelField(pdf, "NIVEL", strconv.Itoa(personagem.Nivel), 130, y, 30, 10)
	s.drawLabelField(pdf, "EXPERIENCIA", "", 162, y, 40, 10)

	y += 12

	// Linha 2: Raca, Classe, Origem
	racaNome := "-"
	if personagem.Raca.Nome != "" {
		racaNome = personagem.Raca.Nome
	}
	classeNome := "-"
	if personagem.Classe.Nome != "" {
		classeNome = personagem.Classe.Nome
	}
	origemNome := "-"
	if personagem.Origem.Nome != "" {
		origemNome = personagem.Origem.Nome
	}

	s.drawLabelField(pdf, "RACA", racaNome, 8, y, 64, 10)
	s.drawLabelField(pdf, "CLASSE", classeNome, 73, y, 64, 10)
	s.drawLabelField(pdf, "ORIGEM", origemNome, 138, y, 64, 10)

	y += 12

	// Linha 3: Divindade, Tamanho, Deslocamento
	divNome := "-"
	if personagem.Divindade != nil && personagem.Divindade.Nome != "" && personagem.Divindade.Nome != "-" {
		divNome = personagem.Divindade.Nome
	}
	tamanho := "-"
	if personagem.Raca.Tamanho != "" {
		tamanho = personagem.Raca.Tamanho
	}
	deslocamento := "-"
	if personagem.Raca.Deslocamento > 0 {
		deslocamento = fmt.Sprintf("%dm", personagem.Raca.Deslocamento)
	}

	s.drawLabelField(pdf, "DIVINDADE", divNome, 8, y, 94, 10)
	s.drawLabelField(pdf, "TAMANHO", tamanho, 103, y, 49, 10)
	s.drawLabelField(pdf, "DESLOCAMENTO", deslocamento, 153, y, 49, 10)

	pdf.SetY(y + 14)
}

// ========== STATS DE COMBATE ==========

func (s *FormFillablePDFService) addCombatStats(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	y := pdf.GetY()

	// Titulo da secao
	s.drawSectionTitle(pdf, "COMBATE", y)
	y += 7

	boxW := 45.0
	boxH := 18.0
	gap := 3.5

	// PV
	pdf.SetFillColor(220, 38, 38) // Vermelho
	pdf.Rect(8, y, boxW, boxH, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(10, y+5, "PONTOS DE VIDA (PV)")
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(22, y+14, strconv.Itoa(personagem.PVTotal))
	pdf.SetTextColor(0, 0, 0)

	// PM
	x := 8 + boxW + gap
	pdf.SetFillColor(37, 99, 235) // Azul
	pdf.Rect(x, y, boxW, boxH, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(x+2, y+5, "PONTOS DE MANA (PM)")
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(x+14, y+14, strconv.Itoa(personagem.PMTotal))
	pdf.SetTextColor(0, 0, 0)

	// Defesa
	x = 8 + 2*(boxW+gap)
	pdf.SetFillColor(22, 163, 74) // Verde
	pdf.Rect(x, y, boxW, boxH, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(x+2, y+5, "DEFESA")
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(x+14, y+14, strconv.Itoa(personagem.Defesa))
	pdf.SetTextColor(0, 0, 0)

	// Iniciativa (mod DES)
	x = 8 + 3*(boxW+gap)
	pdf.SetFillColor(168, 85, 247) // Roxo
	pdf.Rect(x, y, boxW, boxH, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(x+2, y+5, "INICIATIVA")
	pdf.SetFont("Arial", "B", 16)
	modDes := personagem.Des
	pdf.Text(x+14, y+14, fmt.Sprintf("%+d", modDes))
	pdf.SetTextColor(0, 0, 0)

	pdf.SetY(y + boxH + 4)
}

// ========== ATRIBUTOS ==========

func (s *FormFillablePDFService) addFormAttributes(pdf *gofpdf.Fpdf, personagem *models.Personagem, showCalculations bool) {
	y := pdf.GetY()
	s.drawSectionTitle(pdf, "ATRIBUTOS", y)
	y += 7

	type attrInfo struct {
		nome  string
		sigla string
		valor int
	}

	attrs := []attrInfo{
		{"Forca", "FOR", personagem.For},
		{"Destreza", "DES", personagem.Des},
		{"Constituicao", "CON", personagem.Con},
		{"Inteligencia", "INT", personagem.Int},
		{"Sabedoria", "SAB", personagem.Sab},
		{"Carisma", "CAR", personagem.Car},
	}

	colW := 31.3
	boxH := 22.0

	for i, attr := range attrs {
		x := 8 + float64(i)*colW + float64(i)*1.0

		// Fundo
		pdf.SetFillColor(243, 244, 246) // Cinza claro
		pdf.Rect(x, y, colW, boxH, "FD")

		// Sigla
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(100, 100, 100)
		pdf.Text(x+2, y+5, attr.sigla)

		// Valor (grande)
		pdf.SetFont("Arial", "B", 18)
		pdf.SetTextColor(0, 0, 0)
		valStr := fmt.Sprintf("%+d", attr.valor)
		pdf.Text(x+colW/2-5, y+16, valStr)

		// Nome pequeno
		pdf.SetFont("Arial", "", 5)
		pdf.SetTextColor(120, 120, 120)
		pdf.Text(x+2, y+20, attr.nome)
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(y + boxH + 4)
}

// ========== PERICIAS PREENCHIDAS ==========

func (s *FormFillablePDFService) addFormSkillsFilled(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	if pdf.GetY() > 230 {
		pdf.AddPage()
	}

	y := pdf.GetY()
	s.drawSectionTitle(pdf, "PERICIAS", y)
	y += 7

	// Mapa de pericias treinadas do personagem
	treinadas := make(map[string]bool)
	for _, p := range personagem.Pericias {
		treinadas[strings.ToLower(p.Nome)] = true
	}

	// Lista completa de pericias T20 com atributos
	type periciaInfo struct {
		nome     string
		atributo string
	}

	pericias := []periciaInfo{
		{"Acrobacia", "DES"}, {"Adestramento", "CAR"}, {"Atletismo", "FOR"},
		{"Atuacao", "CAR"}, {"Cavalgar", "DES"}, {"Conhecimento", "INT"},
		{"Cura", "SAB"}, {"Diplomacia", "CAR"}, {"Enganacao", "CAR"},
		{"Fortitude", "CON"}, {"Furtividade", "DES"}, {"Guerra", "INT"},
		{"Iniciativa", "DES"}, {"Intimidacao", "CAR"}, {"Intuicao", "SAB"},
		{"Investigacao", "INT"}, {"Luta", "FOR"}, {"Misticismo", "INT"},
		{"Navegacao", "INT"}, {"Nobreza", "INT"}, {"Oficio", "INT"},
		{"Percepcao", "SAB"}, {"Pilotagem", "DES"}, {"Pontaria", "DES"},
		{"Reflexos", "DES"}, {"Religiao", "SAB"}, {"Sobrevivencia", "SAB"},
		{"Vontade", "SAB"},
	}

	// Cabecalho
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(50, 50, 50)
	pdf.SetTextColor(255, 255, 255)
	pdf.Rect(8, y, 47, 5, "F")
	pdf.Text(10, y+3.5, "Pericia")
	pdf.Rect(55, y, 12, 5, "F")
	pdf.Text(57, y+3.5, "Treino")
	pdf.Rect(67, y, 12, 5, "F")
	pdf.Text(68, y+3.5, "Atrib.")
	pdf.Rect(79, y, 12, 5, "F")
	pdf.Text(80, y+3.5, "Bonus")
	pdf.Rect(91, y, 12, 5, "F")
	pdf.Text(93, y+3.5, "Total")
	pdf.SetTextColor(0, 0, 0)

	y += 6

	// Mapa de atributos para valores
	attrMap := map[string]int{
		"FOR": personagem.For, "DES": personagem.Des, "CON": personagem.Con,
		"INT": personagem.Int, "SAB": personagem.Sab, "CAR": personagem.Car,
	}

	// Duas colunas
	colStartX := []float64{8, 108}
	halfLen := len(pericias) / 2
	if len(pericias)%2 != 0 {
		halfLen++
	}

	for col := 0; col < 2; col++ {
		startIdx := col * halfLen
		endIdx := startIdx + halfLen
		if endIdx > len(pericias) {
			endIdx = len(pericias)
		}

		baseY := pdf.GetY()
		if col == 1 {
			baseY = y // reset to same starting Y
		}

		for i := startIdx; i < endIdx; i++ {
			p := pericias[i]
			rowY := baseY + float64(i-startIdx)*5.0
			x := colStartX[col]

			if rowY > 280 {
				break
			}

			isTreinada := treinadas[strings.ToLower(p.nome)]
			modAttr := attrMap[p.atributo]
			bonusTreino := 0
			if isTreinada {
				bonusTreino = personagem.Nivel / 2
				if bonusTreino < 2 {
					bonusTreino = 2
				}
			}
			total := modAttr + bonusTreino

			// Fundo alternado
			if i%2 == 0 {
				pdf.SetFillColor(248, 248, 248)
				pdf.Rect(x, rowY, 95, 5, "F")
			}

			pdf.SetFont("Arial", "", 7)
			pdf.Text(x+1, rowY+3.5, p.nome)

			// Checkbox treinada
			if isTreinada {
				pdf.SetFillColor(22, 163, 74)
				pdf.Rect(x+48, rowY+0.5, 3.5, 3.5, "F")
				pdf.SetFont("Arial", "B", 6)
				pdf.SetTextColor(255, 255, 255)
				pdf.Text(x+48.7, rowY+3.2, "T")
				pdf.SetTextColor(0, 0, 0)
			} else {
				pdf.Rect(x+48, rowY+0.5, 3.5, 3.5, "D")
			}

			pdf.SetFont("Arial", "", 6)
			pdf.Text(x+54, rowY+3.5, p.atributo)

			pdf.SetFont("Arial", "", 7)
			pdf.Text(x+66, rowY+3.5, fmt.Sprintf("%+d", modAttr))

			if isTreinada {
				pdf.SetFont("Arial", "B", 7)
			} else {
				pdf.SetFont("Arial", "", 7)
			}
			pdf.Text(x+78, rowY+3.5, fmt.Sprintf("%+d", total))
		}
	}

	// Pular para apos a coluna mais longa
	maxRows := halfLen
	pdf.SetY(y + float64(maxRows)*5.0 + 4)
}

// ========== PODERES ==========

func (s *FormFillablePDFService) addPowersSection(pdf *gofpdf.Fpdf, personagem *models.Personagem) {
	y := pdf.GetY()
	s.drawSectionTitle(pdf, "HABILIDADES E PODERES", y)
	y += 7

	// Habilidades de Raca
	if len(personagem.Raca.Habilidades) > 0 {
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(34, 139, 34)
		pdf.Text(10, y, fmt.Sprintf("Habilidades de Raca (%s):", personagem.Raca.Nome))
		y += 5
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 7)
		for _, hab := range personagem.Raca.Habilidades {
			if y > 275 {
				pdf.AddPage()
				y = 15
			}
			pdf.Text(12, y, fmt.Sprintf("- %s", hab.Nome))
			y += 4
		}
		y += 3
	}

	// Habilidades de Classe
	if len(personagem.Classe.Habilidades) > 0 {
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(37, 99, 235)
		pdf.Text(10, y, fmt.Sprintf("Habilidades de Classe (%s):", personagem.Classe.Nome))
		y += 5
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Arial", "", 7)
		for _, hab := range personagem.Classe.Habilidades {
			if hab.Nivel <= personagem.Nivel {
				if y > 275 {
					pdf.AddPage()
					y = 15
				}
				nivel := ""
				if hab.Nivel > 1 {
					nivel = fmt.Sprintf(" (Nv.%d)", hab.Nivel)
				}
				pdf.Text(12, y, fmt.Sprintf("- %s%s", hab.Nome, nivel))
				y += 4
			}
		}
		y += 3
	}

	// Linhas vazias para poderes adicionais
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.Text(10, y, "Poderes Adicionais:")
	y += 5
	for i := 0; i < 8; i++ {
		if y > 275 {
			break
		}
		pdf.Line(10, y, 200, y)
		y += 6
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetY(y + 2)
}

// ========== INVENTARIO ==========

func (s *FormFillablePDFService) addFormInventory(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 230 {
		pdf.AddPage()
	}

	y := pdf.GetY()
	s.drawSectionTitle(pdf, "INVENTARIO", y)
	y += 7

	// Cabecalho
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(50, 50, 50)
	pdf.SetTextColor(255, 255, 255)
	pdf.Rect(8, y, 100, 5, "F")
	pdf.Text(10, y+3.5, "Item")
	pdf.Rect(108, y, 25, 5, "F")
	pdf.Text(110, y+3.5, "Qtd")
	pdf.Rect(133, y, 35, 5, "F")
	pdf.Text(135, y+3.5, "Peso")
	pdf.Rect(168, y, 34, 5, "F")
	pdf.Text(170, y+3.5, "Valor (T$)")
	pdf.SetTextColor(0, 0, 0)

	y += 6
	pdf.SetFont("Arial", "", 8)

	for i := 0; i < 12; i++ {
		if y > 275 {
			break
		}
		if i%2 == 0 {
			pdf.SetFillColor(248, 248, 248)
			pdf.Rect(8, y, 194, 5.5, "F")
		}
		pdf.Rect(8, y, 100, 5.5, "D")
		pdf.Rect(108, y, 25, 5.5, "D")
		pdf.Rect(133, y, 35, 5.5, "D")
		pdf.Rect(168, y, 34, 5.5, "D")
		y += 5.5
	}

	// Carga maxima
	pdf.SetFont("Arial", "B", 7)
	pdf.Text(135, y+3, fmt.Sprintf("Carga Maxima: ______ kg"))

	pdf.SetY(y + 6)
}

// ========== ANOTACOES ==========

func (s *FormFillablePDFService) addFormNotes(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	y := pdf.GetY()
	s.drawSectionTitle(pdf, "ANOTACOES", y)
	y += 7

	for i := 0; i < 6; i++ {
		if y > 280 {
			break
		}
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(10, y, 200, y)
		pdf.SetDrawColor(0, 0, 0)
		y += 6
	}

	pdf.SetY(y + 2)
}

// ========== HISTORICO ==========

func (s *FormFillablePDFService) addFormHistory(pdf *gofpdf.Fpdf) {
	if pdf.GetY() > 240 {
		pdf.AddPage()
	}

	y := pdf.GetY()
	s.drawSectionTitle(pdf, "HISTORICO DO PERSONAGEM", y)
	y += 7

	for i := 0; i < 8; i++ {
		if y > 280 {
			break
		}
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(10, y, 200, y)
		pdf.SetDrawColor(0, 0, 0)
		y += 6
	}

	pdf.SetY(y + 2)
}

// ========== HELPERS ==========

// drawSectionTitle desenha um titulo de secao com linha horizontal
func (s *FormFillablePDFService) drawSectionTitle(pdf *gofpdf.Fpdf, title string, y float64) {
	pdf.SetFillColor(139, 0, 0)
	pdf.Rect(8, y, 194, 6, "F")
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	pdf.Text(10, y+4.5, title)
	pdf.SetTextColor(0, 0, 0)
}

// drawLabelField desenha um campo com label e valor preenchido
func (s *FormFillablePDFService) drawLabelField(pdf *gofpdf.Fpdf, label, value string, x, y, w, h float64) {
	// Borda do campo
	pdf.SetDrawColor(180, 180, 180)
	pdf.Rect(x, y, w, h, "D")
	pdf.SetDrawColor(0, 0, 0)

	// Label (pequeno, acima)
	pdf.SetFont("Arial", "", 5)
	pdf.SetTextColor(120, 120, 120)
	pdf.Text(x+1, y+3, label)

	// Valor
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(0, 0, 0)
	pdf.Text(x+2, y+h-1.5, value)
}

// getAtributosLivres retorna os atributos livres escolhidos
func (s *FormFillablePDFService) getAtributosLivres(personagem *models.Personagem) []string {
	var atributos []string
	if personagem.AtributosLivres != "" && personagem.AtributosLivres != "[]" {
		json.Unmarshal([]byte(personagem.AtributosLivres), &atributos)
	}
	return atributos
}

// addCheckbox adiciona uma checkbox
func (s *FormFillablePDFService) addCheckbox(pdf *gofpdf.Fpdf, name string, x, y, w, h float64) {
	pdf.Rect(x, y, w, h, "D")
}
