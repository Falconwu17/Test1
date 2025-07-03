package reports

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"encoding/json"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"net/http"
	"os"
)

func GraphHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "graph.png"
		if err := Graph(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, path)
	}
}

func Graph() error {
	const path = "graph.png"
	_ = os.Remove(path)

	p := plot.New()
	p.Title.Text = "Данные Record и Entry по времени"
	p.X.Label.Text = "Время"
	p.Y.Label.Text = "Значение"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02 15:04:05"}

	recordsTable, err := db_.SelectRecord(100, 0)
	if err != nil {
		return err
	}
	entriesTable := db_.SelectEntry()
	pointsRecord := pointsForRecords(recordsTable)
	pointsEntry := pointsForEntry(entriesTable)

	if err := plotutil.AddLinePoints(p,
		"Record", pointsRecord,
		"Entry", pointsEntry,
	); err != nil {
		return err
	}
	return p.Save(8*vg.Inch, 5*vg.Inch, path)
}

func pointsForRecords(recordsTable []models.Record) plotter.XYs {
	pts := make(plotter.XYs, 0, len(recordsTable))
	for _, record := range recordsTable {
		pts = append(pts, plotter.XY{
			X: float64(record.Created_at.Unix()),
			Y: float64(record.Timeout),
		})
	}
	return pts
}

func pointsForEntry(entriesTable []models.Entry) plotter.XYs {
	pts := make(plotter.XYs, 0, len(entriesTable))
	for _, entry := range entriesTable {
		var payload struct {
			Value float64 `json:"value"`
		}
		if err := json.Unmarshal(entry.Data, &payload); err != nil {
			log.Println("entry.Data parse error:", err)
			continue
		}
		pts = append(pts, plotter.XY{
			X: float64(entry.Created_at.Unix()),
			Y: payload.Value,
		})
	}
	return pts
}
