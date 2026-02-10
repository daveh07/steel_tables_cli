// Package columns defines steel property column configurations.
package columns

import (
	"fmt"
	"strconv"

	"steel_tables/internal/models"
)

// ColumnInfo defines a column with its name and value formatter.
type ColumnInfo struct {
	Name      string
	Formatter func(models.SteelProperty) string
}

// unitMap maps column names to their display units.
var unitMap = map[string]string{
	"Weight": "kg/m", "d": "mm", "bf": "mm", "tf": "mm", "tw": "mm", "r1": "mm", "d1": "mm",
	"tw__1": "mm", "tf__1": "mm", "Ag": "mm²", "Ix": "10⁶mm⁴", "Zx": "10³mm³", "Sx": "10³mm³",
	"rx": "mm", "Iy": "10⁶mm⁴", "Zy": "mm³", "Sy": "mm³", "ry": "mm", "J": "10³mm⁴", "Iw": "10⁹mm⁶",
	"flange": "mm", "web": "mm", "Zex": "mm³", "Zey": "mm³", "Zy5": "mm³", "Fu": "MPa", "r2": "mm",
	"ZeyD": "mm³", "In": "10³mm⁴", "Ip": "10³mm⁴", "ZexC": "mm³", "x5": "mm", "y5": "mm", "nL": "mm",
	"pB": "mm", "pT": "mm",
}

// GetHeaderWithUnit returns the column name with its unit.
func GetHeaderWithUnit(columnName string) string {
	if unit, exists := unitMap[columnName]; exists {
		return fmt.Sprintf("%s (%s)", columnName, unit)
	}
	return columnName
}

// FormatInterface formats interface{} values for display.
func FormatInterface(value interface{}) string {
	if value == nil {
		return "-"
	}
	switch v := value.(type) {
	case string:
		if v == "" || v == "-" {
			return "-"
		}
		return v
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%.1f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		str := fmt.Sprintf("%v", v)
		if str == "" {
			return "-"
		}
		return str
	}
}

// GetAll returns all available column definitions.
func GetAll() []ColumnInfo {
	return []ColumnInfo{
		{"Grade", func(p models.SteelProperty) string { return fmt.Sprintf("%d", p.Grade) }},
		{"Weight", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Weight) }},
		{"d", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.D) }},
		{"bf", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Bf) }},
		{"tf", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Tf) }},
		{"tw", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Tw) }},
		{"r1", func(p models.SteelProperty) string { return FormatInterface(p.R1) }},
		{"d1", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.D1) }},
		{"tw__1", func(p models.SteelProperty) string { return FormatInterface(p.Tw1) }},
		{"tf__1", func(p models.SteelProperty) string { return FormatInterface(p.Tf1) }},
		{"Ag", func(p models.SteelProperty) string { return fmt.Sprintf("%.0f", p.Ag) }},
		{"Ix", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Ix) }},
		{"Zx", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Zx) }},
		{"Sx", func(p models.SteelProperty) string { return fmt.Sprintf("%.0f", p.Sx) }},
		{"rx", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Rx) }},
		{"Iy", func(p models.SteelProperty) string { return fmt.Sprintf("%.2f", p.Iy) }},
		{"Zy", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Zy) }},
		{"Sy", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Sy) }},
		{"ry", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Ry) }},
		{"J", func(p models.SteelProperty) string { return fmt.Sprintf("%.0f", p.J) }},
		{"Iw", func(p models.SteelProperty) string { return FormatInterface(p.Iw) }},
		{"flange", func(p models.SteelProperty) string { return FormatInterface(p.Flange) }},
		{"web", func(p models.SteelProperty) string { return FormatInterface(p.Web) }},
		{"kf", func(p models.SteelProperty) string { return FormatInterface(p.Kf) }},
		{"C,N,S", func(p models.SteelProperty) string { return FormatInterface(p.CNS) }},
		{"Zex", func(p models.SteelProperty) string { return fmt.Sprintf("%.0f", p.Zex) }},
		{"C,N,S__1", func(p models.SteelProperty) string { return FormatInterface(p.CNS2) }},
		{"Zey", func(p models.SteelProperty) string { return fmt.Sprintf("%.1f", p.Zey) }},
		{"2tf", func(p models.SteelProperty) string { return FormatInterface(p.TwoTf) }},
		{"Zy5", func(p models.SteelProperty) string {
			if p.Zy5 == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.Zy5)
		}},
		{"TanAlpha", func(p models.SteelProperty) string {
			if p.TanAlpha == 0 {
				return "-"
			}
			return fmt.Sprintf("%.3f", p.TanAlpha)
		}},
		{"αb", func(p models.SteelProperty) string { return FormatInterface(p.AlphaB) }},
		{"Fu", func(p models.SteelProperty) string { return FormatInterface(p.Fu) }},
		{"r2", func(p models.SteelProperty) string { return FormatInterface(p.R2) }},
		{"ZeyD", func(p models.SteelProperty) string {
			if p.ZeyD == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.ZeyD)
		}},
		{"In", func(p models.SteelProperty) string {
			if p.In == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2f", p.In)
		}},
		{"Ip", func(p models.SteelProperty) string {
			if p.Ip == 0 {
				return "-"
			}
			return fmt.Sprintf("%.2f", p.Ip)
		}},
		{"ZexC", func(p models.SteelProperty) string {
			if p.ZexC == 0 {
				return "-"
			}
			return fmt.Sprintf("%.0f", p.ZexC)
		}},
		{"x5", func(p models.SteelProperty) string { return FormatInterface(p.X5) }},
		{"y5", func(p models.SteelProperty) string {
			if p.Y5 == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.Y5)
		}},
		{"nL", func(p models.SteelProperty) string {
			if p.NL == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.NL)
		}},
		{"pB", func(p models.SteelProperty) string {
			if p.PB == 0 {
				return "-"
			}
			return fmt.Sprintf("%.1f", p.PB)
		}},
		{"pT", func(p models.SteelProperty) string { return FormatInterface(p.PT) }},
		{"Residual", func(p models.SteelProperty) string {
			if p.Residual == "" {
				return "-"
			}
			return p.Residual
		}},
		{"Type", func(p models.SteelProperty) string { return FormatInterface(p.Type) }},
	}
}

// FilterAvailable returns only columns that have meaningful data.
func FilterAvailable(allColumns []ColumnInfo, properties []models.SteelProperty) []ColumnInfo {
	var availableColumns []ColumnInfo
	for _, col := range allColumns {
		hasNonDashData := false
		hasRealData := false

		for _, p := range properties {
			val := col.Formatter(p)
			if val != "-" {
				hasNonDashData = true
				if val != "" && val != "0" && val != "0.0" && val != "0.00" && val != "0.000" {
					if f, err := strconv.ParseFloat(val, 64); err != nil || f != 0.0 {
						hasRealData = true
					}
				}
			}
		}

		if hasNonDashData && hasRealData {
			availableColumns = append(availableColumns, col)
		}
	}
	return availableColumns
}
