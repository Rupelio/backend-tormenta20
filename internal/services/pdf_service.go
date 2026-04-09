package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"tormenta20-builder/internal/models"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PDFExportOptions struct {
	Layout           string            `json:"layout"`
	IncludeImage     bool              `json:"include_image"`
	CustomColors     map[string]string `json:"custom_colors"`
	ExtraSections    []string          `json:"extra_sections"`
	ShowCalculations bool              `json:"show_calculations"`
}

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GenerateCharacterSheet(personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		Build()

	mrt := maroto.New(cfg)

	if options.Layout == "double" {
		return s.generateDoublePageSheet(mrt, personagem, options)
	}

	return s.generateSinglePageSheet(mrt, personagem, options)
}

func (s *PDFService) generateSinglePageSheet(mrt core.Maroto, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	s.addHeader(mrt, personagem)
	s.addBasicInfo(mrt, personagem)
	s.addCombatRow(mrt, personagem)
	s.addAttributes(mrt, personagem, options.ShowCalculations)

	for _, section := range options.ExtraSections {
		switch section {
		case "skills":
			s.addSkillsFilled(mrt, personagem)
		case "inventory":
			s.addInventory(mrt)
		case "notes":
			s.addNotes(mrt)
		case "history":
			s.addHistory(mrt)
		}
	}

	doc, err := mrt.Generate()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return doc.GetBytes(), nil
}

func (s *PDFService) generateDoublePageSheet(mrt core.Maroto, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	s.addHeader(mrt, personagem)
	s.addBasicInfo(mrt, personagem)
	s.addCombatRow(mrt, personagem)
	s.addAttributes(mrt, personagem, options.ShowCalculations)
	s.addSkillsFilled(mrt, personagem)
	s.addInventory(mrt)
	s.addNotes(mrt)
	s.addHistory(mrt)

	doc, err := mrt.Generate()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return doc.GetBytes(), nil
}

func (s *PDFService) addHeader(mrt core.Maroto, personagem *models.Personagem) {
	mrt.AddRow(12,
		col.New(12).Add(
			text.New("FICHA DE PERSONAGEM - TORMENTA 20", props.Text{
				Top:   2,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  16,
			}),
		),
	)
	mrt.AddRow(2)
}

func (s *PDFService) addBasicInfo(mrt core.Maroto, personagem *models.Personagem) {
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
	divNome := "-"
	if personagem.Divindade != nil && personagem.Divindade.Nome != "" && personagem.Divindade.Nome != "-" {
		divNome = personagem.Divindade.Nome
	}

	// Nome e nivel
	mrt.AddRow(8,
		col.New(8).Add(
			text.New("Nome: "+personagem.Nome, props.Text{Top: 1, Style: fontstyle.Bold, Size: 11}),
		),
		col.New(4).Add(
			text.New("Nivel: "+strconv.Itoa(personagem.Nivel), props.Text{Top: 1, Style: fontstyle.Bold, Size: 11}),
		),
	)

	// Raca, Classe, Origem
	mrt.AddRow(7,
		col.New(4).Add(
			text.New("Raca: "+racaNome, props.Text{Top: 1, Size: 10}),
		),
		col.New(4).Add(
			text.New("Classe: "+classeNome, props.Text{Top: 1, Size: 10}),
		),
		col.New(4).Add(
			text.New("Origem: "+origemNome, props.Text{Top: 1, Size: 10}),
		),
	)

	// Divindade
	mrt.AddRow(7,
		col.New(6).Add(
			text.New("Divindade: "+divNome, props.Text{Top: 1, Size: 10}),
		),
		col.New(3).Add(
			text.New("Tamanho: "+personagem.Raca.Tamanho, props.Text{Top: 1, Size: 10}),
		),
		col.New(3).Add(
			text.New(fmt.Sprintf("Deslocamento: %dm", personagem.Raca.Deslocamento), props.Text{Top: 1, Size: 10}),
		),
	)

	mrt.AddRow(3)
}

func (s *PDFService) addCombatRow(mrt core.Maroto, personagem *models.Personagem) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("COMBATE", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	mrt.AddRow(10,
		col.New(3).Add(
			text.New(fmt.Sprintf("PV: %d", personagem.PVTotal), props.Text{
				Top: 2, Style: fontstyle.Bold, Align: align.Center, Size: 14,
			}),
		),
		col.New(3).Add(
			text.New(fmt.Sprintf("PM: %d", personagem.PMTotal), props.Text{
				Top: 2, Style: fontstyle.Bold, Align: align.Center, Size: 14,
			}),
		),
		col.New(3).Add(
			text.New(fmt.Sprintf("Defesa: %d", personagem.Defesa), props.Text{
				Top: 2, Style: fontstyle.Bold, Align: align.Center, Size: 14,
			}),
		),
		col.New(3).Add(
			text.New(fmt.Sprintf("Iniciativa: %+d", personagem.Des), props.Text{
				Top: 2, Style: fontstyle.Bold, Align: align.Center, Size: 14,
			}),
		),
	)

	mrt.AddRow(3)
}

func (s *PDFService) addAttributes(mrt core.Maroto, personagem *models.Personagem, showCalculations bool) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("ATRIBUTOS", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	attrs := []struct {
		name  string
		sigla string
		value int
	}{
		{"Forca", "FOR", personagem.For},
		{"Destreza", "DES", personagem.Des},
		{"Constituicao", "CON", personagem.Con},
		{"Inteligencia", "INT", personagem.Int},
		{"Sabedoria", "SAB", personagem.Sab},
		{"Carisma", "CAR", personagem.Car},
	}

	mrt.AddRow(8,
		col.New(2).Add(text.New(fmt.Sprintf("FOR: %+d", attrs[0].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
		col.New(2).Add(text.New(fmt.Sprintf("DES: %+d", attrs[1].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
		col.New(2).Add(text.New(fmt.Sprintf("CON: %+d", attrs[2].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
		col.New(2).Add(text.New(fmt.Sprintf("INT: %+d", attrs[3].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
		col.New(2).Add(text.New(fmt.Sprintf("SAB: %+d", attrs[4].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
		col.New(2).Add(text.New(fmt.Sprintf("CAR: %+d", attrs[5].value), props.Text{Style: fontstyle.Bold, Size: 11, Align: align.Center})),
	)

	mrt.AddRow(3)
}

func (s *PDFService) addSkillsFilled(mrt core.Maroto, personagem *models.Personagem) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("PERICIAS", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	treinadas := make(map[string]bool)
	for _, p := range personagem.Pericias {
		treinadas[strings.ToLower(p.Nome)] = true
	}

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

	attrMap := map[string]int{
		"FOR": personagem.For, "DES": personagem.Des, "CON": personagem.Con,
		"INT": personagem.Int, "SAB": personagem.Sab, "CAR": personagem.Car,
	}

	// Cabecalho
	mrt.AddRow(5,
		col.New(4).Add(text.New("Pericia", props.Text{Style: fontstyle.Bold, Size: 8})),
		col.New(2).Add(text.New("Treinada", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center})),
		col.New(2).Add(text.New("Atributo", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center})),
		col.New(2).Add(text.New("Mod", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center})),
		col.New(2).Add(text.New("Total", props.Text{Style: fontstyle.Bold, Size: 8, Align: align.Center})),
	)

	for _, p := range pericias {
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

		treinoStr := " "
		if isTreinada {
			treinoStr = "T"
		}

		mrt.AddRow(4,
			col.New(4).Add(text.New(p.nome, props.Text{Size: 7})),
			col.New(2).Add(text.New(treinoStr, props.Text{Size: 7, Align: align.Center})),
			col.New(2).Add(text.New(p.atributo, props.Text{Size: 7, Align: align.Center})),
			col.New(2).Add(text.New(fmt.Sprintf("%+d", modAttr), props.Text{Size: 7, Align: align.Center})),
			col.New(2).Add(text.New(fmt.Sprintf("%+d", total), props.Text{Size: 7, Align: align.Center})),
		)
	}

	mrt.AddRow(3)
}

func (s *PDFService) addInventory(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("INVENTARIO", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	mrt.AddRow(5,
		col.New(6).Add(text.New("Item", props.Text{Style: fontstyle.Bold, Size: 9})),
		col.New(2).Add(text.New("Qtd", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Center})),
		col.New(2).Add(text.New("Peso", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Center})),
		col.New(2).Add(text.New("Valor", props.Text{Style: fontstyle.Bold, Size: 9, Align: align.Center})),
	)

	for i := 0; i < 12; i++ {
		mrt.AddRow(4,
			col.New(6).Add(text.New("", props.Text{Size: 8})),
			col.New(2).Add(text.New("", props.Text{Size: 8})),
			col.New(2).Add(text.New("", props.Text{Size: 8})),
			col.New(2).Add(text.New("", props.Text{Size: 8})),
		)
	}
}

func (s *PDFService) addNotes(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("ANOTACOES", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	for i := 0; i < 8; i++ {
		mrt.AddRow(4, col.New(12).Add(text.New("", props.Text{Size: 8})))
	}
}

func (s *PDFService) addHistory(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("HISTORICO DO PERSONAGEM", props.Text{Top: 1, Style: fontstyle.Bold, Align: align.Center, Size: 11}),
		),
	)

	for i := 0; i < 6; i++ {
		mrt.AddRow(4, col.New(12).Add(text.New("", props.Text{Size: 8})))
	}
}

// getAtributosLivres retorna os atributos livres escolhidos
func (s *PDFService) getAtributosLivres(personagem *models.Personagem) []string {
	var atributos []string
	if personagem.AtributosLivres != "" && personagem.AtributosLivres != "[]" {
		json.Unmarshal([]byte(personagem.AtributosLivres), &atributos)
	}
	return atributos
}
