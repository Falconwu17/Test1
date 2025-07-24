package reports

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"net/http"
)

func GraphHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		if err := Graph(buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(buf.Bytes())
	}
}

func Graph(buf *bytes.Buffer) error {
	p := plot.New()
	p.Title.Text = "График: Температура и Нагрузка"
	p.X.Label.Text = "Время"
	p.Y.Label.Text = "Значение"
	p.X.Tick.Marker = plot.TimeTicks{Format: "15:04:05"}
	p.Add(plotter.NewGrid())
	p.Legend.Top = true
	p.Legend.Left = true

	recordsTable, err := db_.SelectRecord(1000, 0)
	if err != nil {
		log.Println("Ошибка получения Record:", err)
		return err
	}
	entriesTable, err := db_.SelectEntry(1000, 0)
	if err != nil {
		log.Println("Ошибка получения Entry:", err)
		return err
	}

	pointsRecord := interpolateXY(pointsForRecords(recordsTable), 10)

	if len(pointsRecord) > 0 {
		line, err := plotter.NewLine(pointsRecord)
		if err == nil {
			line.Color = color.RGBA{R: 200, B: 200, A: 255}
			p.Add(line)
			p.Legend.Add("Record.Timeout", line)
		}
	}

	entrySeries := pointsByMetric(entriesTable)

	colors := map[string]color.Color{
		"temperature": color.RGBA{R: 255, A: 255},
		"loading":     color.RGBA{G: 255, A: 255},
	}

	for metric, points := range entrySeries {
		sorted := interpolateXY(points, 10)

		lpLine, pointDots, err := plotter.NewLinePoints(sorted)
		if err != nil {
			log.Println("Line error:", err)
			continue
		}

		lpLine.LineStyle.Width = vg.Points(2)
		lpLine.Color = colors[metric]

		pointDots.Shape = draw.CircleGlyph{}
		pointDots.Color = colors[metric]

		p.Add(lpLine, pointDots)
		p.Legend.Add("Entry."+metric, lpLine)
	}

	if len(pointsRecord) == 0 && len(entrySeries) == 0 {
		return errors.New("нет данных для построения графика")
	}

	writerTo, err := p.WriterTo(12*vg.Inch, 6*vg.Inch, "png")
	if err != nil {
		return err
	}
	_, err = writerTo.WriteTo(buf)
	return err
}

func pointsForRecords(recordsTable []models.Record) plotter.XYs {
	pts := make(plotter.XYs, 0, len(recordsTable))
	for _, record := range recordsTable {
		pts = append(pts, plotter.XY{
			X: float64(record.CreatedAt.Unix()),
			Y: float64(record.Timeout),
		})
	}
	return pts
}

func pointsByMetric(entries []models.Entry) map[string]plotter.XYs {
	series := make(map[string]plotter.XYs)
	for _, entry := range entries {
		if !json.Valid(entry.Data) {
			log.Printf("❌ Invalid JSON в entry ID=%d", entry.Id)
			continue
		}
		var payload map[string]float64
		if err := json.Unmarshal(entry.Data, &payload); err != nil {
			log.Printf("❌ Ошибка JSON entry ID=%d: %v", entry.Id, err)
			continue
		}
		for key, val := range payload {
			series[key] = append(series[key], plotter.XY{
				X: float64(entry.CreatedAt.Unix()),
				Y: val,
			})
		}
	}
	return series
}

func interpolateXY(data plotter.XYs, steps int) plotter.XYs {
	if len(data) < 2 {
		return data
	}
	var result plotter.XYs
	for i := 0; i < len(data)-1; i++ {
		x0, y0 := data[i].X, data[i].Y
		x1, y1 := data[i+1].X, data[i+1].Y

		for j := 0; j < steps; j++ {
			t := float64(j) / float64(steps)
			x := x0 + t*(x1-x0)
			y := y0 + t*(y1-y0) + 0.1*(1-t)*t*(y1-y0) // лёгкая волна
			result = append(result, plotter.XY{X: x, Y: y})
		}
	}
	result = append(result, data[len(data)-1])
	return result
}
