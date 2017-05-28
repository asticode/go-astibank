package main

import (
	"encoding/json"

	"sort"

	"github.com/asticode/go-astichartjs"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/pkg/errors"
)

// handleMessageChartsList handles the "charts.all" message
func handleMessageChartsAll(w *astilectron.Window, m bootstrap.MessageIn) {
	// Process errors
	var err error
	defer processMessageError(w, &err)

	// Unmarshal
	var id string
	if err = json.Unmarshal(m.Payload, &id); err != nil {
		err = errors.Wrapf(err, "unmarshaling %s failed", m.Payload)
		return
	}

	// Fetch account
	var a *Account
	if a, err = data.Accounts.One(id); err != nil {
		err = errors.Wrapf(err, "fetching account %s failed", id)
		return
	}

	// Send
	if err = w.Send(bootstrap.MessageOut{Name: "charts.all", Payload: buildCharts(a)}); err != nil {
		err = errors.Wrap(err, "sending message failed")
		return
	}
}

// buildCharts builds charts
func buildCharts(a *Account) (cs []astichartjs.Chart) {
	// Loop through operations
	var categories, dates []string
	var datesMap = make(map[string]bool)
	var d = make(map[string]map[string]map[string]float64)
	for _, operation := range a.Operations.All() {
		// New category
		if _, ok := d[operation.Category]; !ok {
			categories = append(categories, operation.Category)
			d[operation.Category] = make(map[string]map[string]float64)
		}

		// New date for category
		var date = operation.Date.Format("01/2006")
		if _, ok := d[operation.Category][date]; !ok {
			d[operation.Category][date] = make(map[string]float64)
		}

		// New date
		if _, ok := datesMap[date]; !ok {
			dates = append(dates, date)
			datesMap[date] = true
		}

		// Update sum
		d[operation.Category][date][operation.Subject] += operation.Amount
	}
	sort.Strings(categories)
	sort.Strings(dates)

	// Build average chart
	cs = append(cs, buildChartAverage(dates, d))

	// Build monthly balance
	cs = append(cs, buildChartMonthlyBalance(dates, d))

	// Build monthly sum charts
	for _, category := range categories {
		cs = append(cs, buildChartMonthlySum(category, dates, d[category]))
	}
	return
}

// buildChartAverage builds the average chart
// d  is indexed by category then by date then by subject
func buildChartAverage(dates []string, d map[string]map[string]map[string]float64) (c astichartjs.Chart) {
	// Init
	c = astichartjs.Chart{
		Data: astichartjs.Data{
			Datasets: []astichartjs.Dataset{{
				BorderWidth: 1,
			}},
			Labels: []string{
				categoryAmenities,
				categoryBread,
				categoryFood,
				categoryLoan,
				categoryPleasure,
				categoryTaxes,
				categoryWork,
			},
		},
		Options: astichartjs.Options{
			Responsive: true,
			Scales: astichartjs.Scales{
				XAxes: []astichartjs.Axis{},
				YAxes: []astichartjs.Axis{{
					ScaleLabel: astichartjs.ScaleLabel{
						Display:     true,
						LabelString: "Average(€)",
					},
				}},
			},
			Title: astichartjs.Title{
				Display:  true,
				FontSize: 16,
				Text:     "Average by category",
			},
		},
		Type: astichartjs.ChartTypeBar,
	}

	// Build chart color picker
	var colorPicker = astichartjs.NewChartColorPicker()

	// Loop through categories
	var backgroundColors, borderColors []string
	for _, category := range c.Data.Labels {
		// Init
		var sum float64

		// Loop through dates
		for _, subjects := range d[category] {
			// Loop through subjects
			for _, amount := range subjects {
				sum += amount
			}
		}

		// Add to dataset
		var color = colorPicker.Next()
		backgroundColors = append(backgroundColors, astichartjs.BackgroundColor(color))
		borderColors = append(borderColors, astichartjs.BorderColor(color))
		c.Data.Datasets[0].Data = append(c.Data.Datasets[0].Data, sum/float64(len(dates)))
	}
	c.Data.Datasets[0].BackgroundColor = backgroundColors
	c.Data.Datasets[0].BorderColor = borderColors
	return
}

// buildChartMonthlyBalance builds the monthly balance chart
// d  is indexed by category then by date then by subject
func buildChartMonthlyBalance(dates []string, d map[string]map[string]map[string]float64) (c astichartjs.Chart) {
	// Init
	c = astichartjs.Chart{
		Data: astichartjs.Data{
			Datasets: []astichartjs.Dataset{{
				BackgroundColor: astichartjs.BackgroundColor(astichartjs.ChartColorRed),
				BorderColor:     astichartjs.BorderColor(astichartjs.ChartColorRed),
				BorderWidth:     1,
			}},
			Labels: dates,
		},
		Options: astichartjs.Options{
			Responsive: true,
			Scales: astichartjs.Scales{
				XAxes: []astichartjs.Axis{},
				YAxes: []astichartjs.Axis{{
					ScaleLabel: astichartjs.ScaleLabel{
						Display:     true,
						LabelString: "Balance(€)",
					},
				}},
			},
			Title: astichartjs.Title{
				Display:  true,
				FontSize: 16,
				Text:     "Monthly balance",
			},
		},
		Type: astichartjs.ChartTypeBar,
	}

	// Loop through categories
	var balances = make(map[string]float64)
	for _, ds := range d {
		// Loop through dates
		for date, subjects := range ds {
			// Loop through subjects
			for _, amount := range subjects {
				balances[date] += amount
			}
		}
	}

	// Loop through dates
	for _, date := range dates {
		c.Data.Datasets[0].Data = append(c.Data.Datasets[0].Data, balances[date])
	}
	return
}

// buildChartMonthlySum builds the monthly sum chart
// d  is indexed by date then by subject
func buildChartMonthlySum(category string, dates []string, d map[string]map[string]float64) (c astichartjs.Chart) {
	// Init
	c = astichartjs.Chart{
		Data: astichartjs.Data{
			Labels: dates,
		},
		Options: astichartjs.Options{
			Legend: astichartjs.Legend{
				Display: true,
			},
			Responsive: true,
			Scales: astichartjs.Scales{
				XAxes: []astichartjs.Axis{{
					Stacked: true,
				}},
				YAxes: []astichartjs.Axis{{
					ScaleLabel: astichartjs.ScaleLabel{
						Display:     true,
						LabelString: "Sum(€)",
					},
					Stacked: true,
				}},
			},
			Title: astichartjs.Title{
				Display:  true,
				FontSize: 16,
				Text:     "Monthly sum for " + category,
			},
		},
		Type: astichartjs.ChartTypeBar,
	}

	// Build chart color picker
	var colorPicker = astichartjs.NewChartColorPicker()

	// Get datasets
	var datasets = make(map[string]*astichartjs.Dataset)
	for _, subjects := range d {
		// Loop through subjects
		for subject := range subjects {
			// New dataset
			if _, ok := datasets[subject]; !ok {
				var color = colorPicker.Next()
				datasets[subject] = &astichartjs.Dataset{
					BackgroundColor: astichartjs.BackgroundColor(color),
					BorderColor:     astichartjs.BorderColor(color),
					BorderWidth:     1,
					Label:           subject,
				}
			}
		}
	}

	// Loop through dates
	for _, date := range dates {
		// Loop through subjects
		for subject, dataset := range datasets {
			// No data
			if _, ok := d[date][subject]; !ok {
				dataset.Data = append(dataset.Data, 0)
			} else {
				dataset.Data = append(dataset.Data, d[date][subject])
			}
		}
	}

	// Loop through datasets
	for _, dataset := range datasets {
		c.Data.Datasets = append(c.Data.Datasets, *dataset)
	}
	return
}
