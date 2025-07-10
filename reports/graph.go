package reports

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"bytes"
	"encoding/json"
	"errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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
	p.Title.Text = "Данные Record и Entry по времени"
	p.X.Label.Text = "Время"
	p.Y.Label.Text = "Значение Value JSON"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02 15:04:05"}

	recordsTable, err := db_.SelectRecord(100, 0)
	if err != nil {
		log.Println("Ошибка построение графики с лимитом", err)
		return err
	}
	entriesTable, err := db_.SelectEntry(100, 0)
	if err != nil {
		log.Println("Ошибка построение графики с лимитом", err)
		return err
	}
	pointsRecord := pointsForRecords(recordsTable)
	pointsEntry := pointsForEntry(entriesTable)
	p.Add(plotter.NewGrid())
	if err := plotutil.AddLinePoints(p, "Record", pointsRecord,
		"Entry", pointsEntry); err != nil {
		return err
	}
	writerTo, err := p.WriterTo(8*vg.Inch, 5*vg.Inch, "png")
	if err != nil {
		return err
	}
	if len(pointsRecord) == 0 && len(pointsEntry) == 0 {
		return errors.New("нет данных для построения графика")
	}

	_, err = writerTo.WriteTo(buf)
	return err

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
		if !json.Valid(entry.Data) {
			log.Printf("⚠ Invalid JSON in entry ID=%d: %s", entry.Id, string(entry.Data))
			continue
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
