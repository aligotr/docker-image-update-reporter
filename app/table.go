package main

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func tableRender(tableData chan tableDataType) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	currentTime := time.Now().Format("02.01.2006, 15:04")
	t.SetTitle(currentTime)

	columnNameContainer := "Контейнер"
	t.AppendHeader(table.Row{columnNameContainer, "Текущая версия", "Обновление", "Ссылка на репозиторий"})

	for res := range tableData {
		updateText := ""
		switch res.hasUpdate {
		case "error":
			updateText = text.FgRed.Sprint(res.hasUpdate)
		case "true":
			updateText = text.FgGreen.Sprint("V")
		}

		t.AppendRow(table.Row{
			res.name,
			res.version,
			updateText,
			res.link,
		})
	}

	t.SortBy([]table.SortBy{
		{Name: columnNameContainer, Mode: table.Asc},
	})

	t.Render()
}
