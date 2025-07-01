package reports

import (
	db_ "awesomeProject1/internal/db"
	"awesomeProject1/internal/models"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"log"
	"net/http"
	"strconv"
)

func GraphHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "graph.png"
		if err := GraphForRecords(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		http.ServeFile(w, r, path)
	}
}

func GraphForRecords() error {
	p := plot.New()
	p.Title.Text = "Timeout and Status Over Time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Status"
	recordsTable := db_.SelectRecord()
	pointsTime := pointsForTimeout(recordsTable)
	pointsStatus := pointsForStatus(recordsTable)
	if err := plotutil.AddLinePoints(p,
		"Time", pointsTime,
		"Status", pointsStatus,
	); err != nil {
		return err
	}
	return p.Save(4*vg.Inch, 4*vg.Inch, "graph.png")
}
func pointsForTimeout(recordsTable []models.Record) plotter.XYs {
	pts := make(plotter.XYs, 0, len(recordsTable))
	for _, record := range recordsTable {
		pts = append(pts, plotter.XY{
			X: float64(record.Created_at.Unix()),
			Y: float64(record.Timeout),
		})
	}
	return pts
}

func pointsForStatus(recordsTable []models.Record) plotter.XYs {
	pts := make(plotter.XYs, 0, len(recordsTable))
	for _, record := range recordsTable {
		status, err := strconv.Atoi(record.Status)
		if err != nil {
			log.Println("status parse error:", err)
			continue
		}
		pts = append(pts, plotter.XY{
			X: float64(record.Created_at.Unix()),
			Y: float64(status),
		})
	}
	return pts
}
