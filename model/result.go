//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

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
