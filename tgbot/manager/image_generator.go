package manager

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"time"

	"schedule/util"
)

var imgDir string

const (
	CharWidth = 22
	S         = 20 // Font size
	Padding   = 20 // Padding
	RowHeight = 30 // Row height

	ImageLifetime = 24 * time.Hour

	TeamLogoUrl = "https://media.api-sports.io/football/teams/%d.png"

	TeamLogoSubdir = "team"

	TeamLogoIconPrefix = "icon_"
	IconWidth          = 28
	IconHeight         = 28
	IconPadding        = 5
)

func InitImageGenerator(imgDirArg string) {
	imgDir = imgDirArg
}

func GetStandingsImage(compId int, standingsData []StandingsData, locale string) (string, error) {
	filePathBase := fmt.Sprintf("%sstandings_%d_%s", imgDir, compId, locale)
	filepath := fmt.Sprintf("%s.png", filePathBase)
	if checkIfUpToDateImageExists(filepath) {
		return filepath, nil
	}

	err := createCompetitionStandingsImage(standingsData, filePathBase, locale)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func checkIfUpToDateImageExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		oneDayAgo := time.Now().Add(-ImageLifetime)
		if fileInfo.ModTime().Before(oneDayAgo) {
			os.Remove(filePath)
			return false
		}
		return true
	}
	return false
}

func createCompetitionStandingsImage(standings []StandingsData, imgPath string, locale string) error {
	headers := []string{
		util.Translate(locale, "StandingRank"),
		util.Translate(locale, "StandingTeam"),
		util.Translate(locale, "StandingPoints"),
		util.Translate(locale, "StandingPlayed"),
		util.Translate(locale, "StandingWon"),
		util.Translate(locale, "StandingDrawn"),
		util.Translate(locale, "StandingLost"),
		util.Translate(locale, "StandingGoalsDiff"),
		util.Translate(locale, "StandingForm"),
	}

	maxLengths, totalStandings := generateStandingsDimensions(standings, headers)

	totalWidth := 0
	for _, length := range maxLengths {
		totalWidth += length
	}

	width := totalWidth * CharWidth
	height := totalStandings * (RowHeight + Padding)

	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace("tgbot/fonts/DejaVuSans-Bold.ttf", S)
	if err != nil {
		return err
	}

	y := Padding
	for id, group := range standings {
		dc.DrawString(group.GroupName, Padding, float64(y))
		y += RowHeight

		if id == 0 {
			x := Padding
			for i, header := range headers {
				if i == 1 {
					x += IconWidth + IconPadding
				}
				dc.DrawString(header, float64(x), float64(y))
				x += maxLengths[i] * S
			}
			y += RowHeight
		}

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
				if i == 1 {
					teamLogoPath, err := GetTeamLogoIconImage(int(standing.TeamId))
					if err != nil {
						panic(err)
					}
					img, err := gg.LoadImage(teamLogoPath)
					if err != nil {
						panic(err)
					}
					dc.DrawImage(img, x, y-IconHeight+IconPadding)
					x += IconWidth + IconPadding
				}

				dc.DrawString(cell, float64(x), float64(y))
				x += maxLengths[i] * S
			}
			y += RowHeight
		}
	}

	dc.Stroke()
	return dc.SavePNG(imgPath + ".png")
}

func GetTeamLogoImage(teamId int) (string, error) {
	filePath := fmt.Sprintf("%s%s/%d.png", imgDir, TeamLogoSubdir, teamId)
	_, err := os.Stat(filePath)
	if err == nil {
		return filePath, nil
	}

	teamLogoUrl := fmt.Sprintf(TeamLogoUrl, teamId)
	err = downloadImage(teamLogoUrl, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GetTeamLogoIconImage(teamId int) (string, error) {
	iconFilePath := fmt.Sprintf("%s%s/%s%d.png", imgDir, TeamLogoSubdir, TeamLogoIconPrefix, teamId)
	_, err := os.Stat(iconFilePath)
	if err == nil {
		return iconFilePath, nil
	}

	teamLogoPath, err := GetTeamLogoImage(teamId)
	if err != nil {
		return "", err
	}

	file, err := os.Open(teamLogoPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	resizedImg := resize.Resize(IconWidth, IconHeight, img, resize.Lanczos3)

	newFile, err := os.Create(iconFilePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	err = png.Encode(newFile, resizedImg)

	return iconFilePath, err
}

func downloadImage(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateStandingsDimensions(standings []StandingsData, headers []string) ([]int, int) {
	maxLengths := make([]int, 9)
	totalStandings := 0

	addToMaxLengths(maxLengths, headers)
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
			addToMaxLengths(maxLengths, cells)
		}
	}

	return maxLengths, totalStandings
}

func addToMaxLengths(maxLengths []int, cells []string) {
	for i, cell := range cells {
		totalLen := len(cell)
		if totalLen > maxLengths[i] {
			maxLengths[i] = totalLen
		}
	}
}
