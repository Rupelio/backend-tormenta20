package services

import (
	"fmt"
	"strconv"

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
	Layout           string            `json:"layout"` // "single" ou "double"
	IncludeImage     bool              `json:"include_image"`
	CustomColors     map[string]string `json:"custom_colors"`
	ExtraSections    []string          `json:"extra_sections"` // "history", "notes", "inventory"
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
	// Cabeçalho da ficha
	s.addHeader(mrt, personagem)

	// Informações básicas
	s.addBasicInfo(mrt, personagem)

	// Atributos
	s.addAttributes(mrt, personagem, options.ShowCalculations)

	// Perícias
	if len(options.ExtraSections) > 0 {
		for _, section := range options.ExtraSections {
			switch section {
			case "skills":
				s.addSkills(mrt, personagem)
			case "inventory":
				s.addInventory(mrt)
			case "notes":
				s.addNotes(mrt)
			}
		}
	}

	doc, err := mrt.Generate()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return doc.GetBytes(), nil
}

func (s *PDFService) generateDoublePageSheet(mrt core.Maroto, personagem *models.Personagem, options PDFExportOptions) ([]byte, error) {
	// Primeira página
	s.addHeader(mrt, personagem)
	s.addBasicInfo(mrt, personagem)
	s.addAttributes(mrt, personagem, options.ShowCalculations)
	s.addSkills(mrt, personagem)

	// Segunda página
	// Nota: A biblioteca maroto/v2 adiciona páginas automaticamente quando necessário
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
	mrt.AddRow(15,
		col.New(12).Add(
			text.New("FICHA DE PERSONAGEM - TORMENTA20", props.Text{
				Top:   3,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  16,
			}),
		),
	)

	mrt.AddRow(2) // Espaçamento
}

func (s *PDFService) addBasicInfo(mrt core.Maroto, personagem *models.Personagem) {
	// Nome do personagem
	mrt.AddRow(8,
		col.New(6).Add(
			text.New("Nome do Personagem:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New(personagem.Nome, props.Text{
				Top:  4,
				Size: 12,
			}),
		),
		col.New(3).Add(
			text.New("Nível:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New(strconv.Itoa(personagem.Nivel), props.Text{
				Top:  4,
				Size: 12,
			}),
		),
		col.New(3).Add(
			text.New("Experiência:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New("________", props.Text{
				Top:  4,
				Size: 12,
			}),
		),
	)

	// Raça, Classe, Origem
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

	mrt.AddRow(8,
		col.New(4).Add(
			text.New("Raça:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New(racaNome, props.Text{
				Top:  4,
				Size: 11,
			}),
		),
		col.New(4).Add(
			text.New("Classe:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New(classeNome, props.Text{
				Top:  4,
				Size: 11,
			}),
		),
		col.New(4).Add(
			text.New("Origem:", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Size:  10,
			}),
			text.New(origemNome, props.Text{
				Top:  4,
				Size: 11,
			}),
		),
	)

	mrt.AddRow(3) // Espaçamento
}

func (s *PDFService) addAttributes(mrt core.Maroto, personagem *models.Personagem, showCalculations bool) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("ATRIBUTOS", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  12,
			}),
		),
	)

	// Cabeçalho dos atributos
	mrt.AddRow(5,
		col.New(3).Add(
			text.New("Atributo", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Valor", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Modificador", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(5).Add(
			text.New("Observações", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
	)

	// Lista de atributos
	attributes := []struct {
		name  string
		value int
	}{
		{"Força", personagem.For},
		{"Destreza", personagem.Des},
		{"Constituição", personagem.Con},
		{"Inteligência", personagem.Int},
		{"Sabedoria", personagem.Sab},
		{"Carisma", personagem.Car},
	}

	for _, attr := range attributes {
		modifier := attr.value // No Tormenta20, modificador = valor do atributo

		mrt.AddRow(6,
			col.New(3).Add(
				text.New(attr.name, props.Text{
					Size: 9,
				}),
			),
			col.New(2).Add(
				text.New(strconv.Itoa(attr.value), props.Text{
					Size:  9,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New(fmt.Sprintf("%+d", modifier), props.Text{
					Size:  9,
					Align: align.Center,
				}),
			),
			col.New(5).Add(
				text.New("", props.Text{
					Size: 9,
				}),
			),
		)
	}

	mrt.AddRow(3) // Espaçamento
}

func (s *PDFService) addSkills(mrt core.Maroto, personagem *models.Personagem) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("PERÍCIAS", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  12,
			}),
		),
	)

	// Cabeçalho das perícias
	mrt.AddRow(5,
		col.New(4).Add(
			text.New("Perícia", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
			}),
		),
		col.New(2).Add(
			text.New("Treinada", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Atributo", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Bônus", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Total", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
	)

	// Lista básica de perícias do Tormenta20
	skills := []string{
		"Acrobacia", "Adestramento", "Atletismo", "Atuação", "Cavalgar",
		"Conhecimento", "Cura", "Diplomacia", "Enganação", "Fortitude",
		"Furtividade", "Guerra", "Iniciativa", "Intimidação", "Intuição",
		"Investigação", "Luta", "Misticismo", "Navegação", "Nobreza",
		"Ofício", "Percepção", "Pilotagem", "Pontaria", "Reflexos",
		"Religião", "Sobrevivência", "Vontade",
	}

	for _, skill := range skills {
		mrt.AddRow(4,
			col.New(4).Add(
				text.New(skill, props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("☐", props.Text{
					Size:  8,
					Align: align.Center,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
		)
	}

	mrt.AddRow(3) // Espaçamento
}

func (s *PDFService) addInventory(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("INVENTÁRIO", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  12,
			}),
		),
	)

	// Cabeçalho do inventário
	mrt.AddRow(5,
		col.New(6).Add(
			text.New("Item", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
			}),
		),
		col.New(2).Add(
			text.New("Quantidade", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Peso", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
		col.New(2).Add(
			text.New("Valor", props.Text{
				Style: fontstyle.Bold,
				Size:  9,
				Align: align.Center,
			}),
		),
	)

	// Linhas do inventário
	for i := 0; i < 15; i++ {
		mrt.AddRow(4,
			col.New(6).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
			col.New(2).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
		)
	}
}

func (s *PDFService) addNotes(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("ANOTAÇÕES", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  12,
			}),
		),
	)

	// Área de anotações
	for i := 0; i < 10; i++ {
		mrt.AddRow(4,
			col.New(12).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
		)
	}
}

func (s *PDFService) addHistory(mrt core.Maroto) {
	mrt.AddRow(6,
		col.New(12).Add(
			text.New("HISTÓRICO DO PERSONAGEM", props.Text{
				Top:   1,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size:  12,
			}),
		),
	)

	// Área de histórico
	for i := 0; i < 8; i++ {
		mrt.AddRow(4,
			col.New(12).Add(
				text.New("", props.Text{
					Size: 8,
				}),
			),
		)
	}
}
