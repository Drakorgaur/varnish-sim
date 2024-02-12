package model

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Result interface {
}

type JsonResult interface {
	ToJson() string
}

type TableResult interface {
	TableData() (name string, rows [][]string)
}

var TableStyleHeader = lipgloss.NewStyle().
	Width(20).Align(lipgloss.Center)

var TableStyleRow = lipgloss.NewStyle().
	Width(20).Align(lipgloss.Center)

func MakeTable(s TableResult) string {
	name, rows := s.TableData()

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return TableStyleHeader
			default:
				return TableStyleRow
			}
		}).
		Rows(rows...)

	t.Headers(name, "")

	return t.Render()
}

func PrintTable(s TableResult) {
	fmt.Println(MakeTable(s))
}
