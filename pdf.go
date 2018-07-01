package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"rsc.io/pdf"
)

type pdfStatement struct {
	credit     float64
	debit      float64
	oldBalance float64
	operations []pdfOperation
}

type pdfParser struct{}

func newPDFParser() *pdfParser {
	return &pdfParser{}
}

func (p *pdfParser) parse(path string) (s pdfStatement, err error) {
	// Open file
	var f *os.File
	astilog.Debugf("main: opening %s", path)
	if f, err = os.Open(path); err != nil {
		err = errors.Wrapf(err, "main: opening %s failed", path)
		return
	}
	defer f.Close()

	// Stat file
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		err = errors.Wrapf(err, "main: stating %s failed", path)
		return
	}

	// Create reader
	var r *pdf.Reader
	if r, err = pdf.NewReader(f, fi.Size()); err != nil {
		err = errors.Wrap(err, "main: creating pdf reader failed")
		return
	}

	// Loop through pages
	for idxPage := 1; idxPage <= r.NumPage(); idxPage++ {
		// Loop through lines
		astilog.Debugf("main: parsing page %d/%d", idxPage, r.NumPage())
		var h *pdfHeader
		var o pdfOperation
		for _, l := range p.lines(r, idxPage) {
			// Find header
			if h == nil {
				h = l.header()
				continue
			}

			// Parse values
			t := l.parse(h)

			// New value
			if (!t.date.IsZero() && !o.date.IsZero()) ||
				(t.debit > 0 && o.debit > 0) ||
				(t.credit > 0 && o.credit > 0) ||
				(t.debit > 0 && o.credit > 0) ||
				(t.credit > 0 && o.debit > 0) {
				s.operations = append(s.operations, o)
				o = pdfOperation{}
			}

			// Update value
			if !t.date.IsZero() {
				o.date = t.date
			}
			if len(t.description) > 0 {
				o.description = append(o.description, t.description...)
			}
			if t.debit > 0 {
				o.debit = t.debit
			}
			if t.credit > 0 {
				o.credit = t.credit
			}
		}

		// Add last value
		if !o.date.IsZero() {
			s.operations = append(s.operations, o)
		}
	}

	// Retrieve old balance and global credit/debit
	for idxOp := 0; idxOp < len(s.operations); idxOp++ {
		o := s.operations[idxOp]
		if s.oldBalance == 0 && len(o.description) == 1 && strings.HasPrefix(o.description[0], "Ancien solde au ") {
			s.oldBalance = o.credit - o.debit
			s.operations = append(s.operations[:idxOp], s.operations[idxOp+1:]...)
			idxOp--
		} else if (s.credit == 0 || s.debit == 0) && len(o.description) == 1 && o.description[0] == "Total des opérations" {
			s.credit = o.credit
			s.debit = o.debit
			s.operations = append(s.operations[:idxOp], s.operations[idxOp+1:]...)
			idxOp--
		}
	}
	return
}

type pdfLine []pdf.Text

type pdfHeader struct {
	xCredit, xDate, xDebit, xOperations float64
}

func (l pdfLine) header() (h *pdfHeader) {
	// Not the header
	if l.concat() != "DateOpérationsDébit(\xc2\xa4)Crédit(\xc2\xa4)" {
		return
	}

	// Create header
	h = &pdfHeader{
		xCredit:     l[21].X + l[21].W,
		xDate:       l[0].X,
		xOperations: l[4].X,
	}
	h.xDebit = h.xCredit - (l[len(l)-1].X + l[len(l)-1].W - h.xCredit)

	// Add margin of error
	h.xCredit += 1
	h.xDate -= 1
	h.xDebit += 1
	h.xOperations -= 1
	return
}

func (l pdfLine) concat() (s string) {
	for _, t := range l {
		s += t.S
	}
	return
}

type pdfOperation struct {
	credit      float64
	date        time.Time
	debit       float64
	description []string
}

func (l pdfLine) parse(h *pdfHeader) (o pdfOperation) {
	// Get values
	var credit, date, debit, operations string
	var pt pdf.Text
	for _, t := range l {
		if t.X < h.xOperations {
			date += t.S
		} else if t.X < h.xDebit {
			if len(pt.S) > 0 {
				if t.X > pt.X+pt.W+1 {
					operations += " "
				}
			}
			operations += t.S
		} else if t.X < h.xCredit {
			debit += t.S
		} else {
			credit += t.S
		}
		pt = t
	}

	// Parse date
	var err error
	if len(date) > 0 {
		if o.date, err = time.Parse("02/01", date); err != nil {
			o.date = time.Time{}
		}
	}

	// Parse operations
	if len(operations) > 0 {
		o.description = append(o.description, strings.TrimSpace(operations))
	}

	// Parse debit
	if len(debit) > 0 {
		if o.debit, err = strconv.ParseFloat(strings.Replace(debit, ",", ".", -1), 64); err != nil {
			o.debit = 0
		}
	}

	// Parse credit
	if len(credit) > 0 {
		if o.credit, err = strconv.ParseFloat(strings.Replace(credit, ",", ".", -1), 64); err != nil {
			o.credit = 0
		}
	}
	return
}

func (p *pdfParser) lines(r *pdf.Reader, idxPage int) (ls []pdfLine) {
	// Get text
	ts := r.Page(idxPage).Content().Text

	// Index y
	var ys []float64
	var ty = make(map[float64][]pdf.Text)
	for _, t := range ts {
		if _, ok := ty[t.Y]; !ok {
			ys = append(ys, t.Y)
		}
		ty[t.Y] = append(ty[t.Y], t)
	}

	// Sort y
	sort.Float64s(ys)

	// Loop through y in reverse order
	for idxY := len(ys) - 1; idxY >= 0; idxY-- {
		// Index x
		var xs []float64
		var tx = make(map[float64][]pdf.Text)
		for _, t := range ty[ys[idxY]] {
			if _, ok := tx[t.X]; !ok {
				xs = append(xs, t.X)
			}
			tx[t.X] = append(tx[t.X], t)
		}

		// Sort x
		sort.Float64s(xs)

		// Loop through x
		l := pdfLine{}
		for _, x := range xs {
			l = append(l, tx[x]...)
		}

		// Append line
		ls = append(ls, l)
	}
	return
}
