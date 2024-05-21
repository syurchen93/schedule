package manager

import (
	"fmt"
	"github.com/fogleman/gg"
	"os"
)

var imgDir string

const (
	CharWidth = 23
	S         = 20 // Font size
	Padding   = 20 // Padding
	RowHeight = 30 // Row height
)

func InitImageGenerator(imgDirArg string) {
	imgDir = imgDirArg
}

func GetStandingsImage(compId int, standingsData []StandingsData) (string, error) {
	filePathBase := fmt.Sprintf("%sstandings_%d", imgDir, compId)
	filepath := fmt.Sprintf("%s.png", filePathBase)
	if _, err := os.Stat(filepath); err == nil {
		return filepath, nil
	}

	err := createCompetitionStandingsImage(standingsData, filePathBase)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func createCompetitionStandingsImage(standings []StandingsData, imgPath string) error {
	maxLengths := make([]int, 9)
	totalStandings := 0
	for _, group := range standings {
		totalStandings += len(group.Standings)
		for _, standing := range group.Standings {
			cells := []string{
				fmt.Sprintf("%d", standing.Position),
				standing.GetTeamNameWithCode(),
				fmt.Sprintf("%d", standing.Points),
				fmt.Sprintf("%d", standing.Played),
				fmt.Sprintf("%d", standing.Won),
				fmt.Sprintf("%d", standing.Drawn),
				fmt.Sprintf("%d", standing.Lost),
				fmt.Sprintf("%d", standing.GoalsDiff),
				standing.Form,
			}
			for i, cell := range cells {
				if len(cell) > maxLengths[i] {
					maxLengths[i] = len(cell)
				}
			}
		}
	}

	totalWidth := 0
	for _, length := range maxLengths {
		totalWidth += length
	}

	width := totalWidth * CharWidth
	height := totalStandings * (RowHeight + Padding/2)

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace("tgbot/fonts/DejaVuSans-Bold.ttf", S)
	if err != nil {
		return err
	}

	y := Padding
	for _, group := range standings {
		dc.DrawString(group.GroupName, Padding, float64(y))
		y += RowHeight

		headers := []string{"R", "Team", "Pts", "P", "W", "D", "L", "GD", "Form"}
		x := Padding
		for i, header := range headers {
			dc.DrawString(header, float64(x), float64(y))
			x += maxLengths[i] * S
		}
		y += RowHeight

		for _, standing := range group.Standings {
			cells := []string{
				fmt.Sprintf("%d", standing.Position),
				standing.TeamName,
				fmt.Sprintf("%d", standing.Points),
				fmt.Sprintf("%d", standing.Played),
				fmt.Sprintf("%d", standing.Won),
				fmt.Sprintf("%d", standing.Drawn),
				fmt.Sprintf("%d", standing.Lost),
				fmt.Sprintf("%d", standing.GoalsDiff),
				standing.Form,
			}
			x := Padding
			for i, cell := range cells {
				dc.DrawString(cell, float64(x), float64(y))
				x += maxLengths[i] * S
			}
			y += RowHeight
		}
	}

	dc.Stroke()
	return dc.SavePNG(imgPath + ".png")
}
