package main

import (
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	fontSize     = 30
	buttonHeight = 60
)

var ColorBackground = color.RGBA{R: 0, G: 32, B: 36, A: 255}
var ColorLine = color.RGBA{R: 0, G: 255, B: 192, A: 255}
var ColorBorder = color.RGBA{R: 64, G: 64, B: 64, A: 255}
var ColorText = color.RGBA{R: 0, G: 128, B: 96, A: 255}

func main() {
	path, _ := os.UserHomeDir()
	path += "\\AppData\\Local\\stats-balls\\"

	stats := Stats{
		LastDayRecorded: 0,
		Stats: []Stat{
			{Name: "0", Values: []int32{0, 0}, Total: 0, LongestStreak: 0, Max: 0},
		},
		Folder: path,
	}
	stats.Load()

	selected := 0
	if len(stats.Stats) > 1 {
		selected = 1
	}
	input_text := ""
	navigator_scroll := int32(0)
	graph_scroll := int32(0)

	rl.SetTraceLogLevel(rl.LogNone)
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagMsaa4xHint | rl.FlagVsyncHint)
	rl.InitWindow(1280, 720, "Stats")
	rl.SetExitKey(-1)
	rl.MaximizeWindow()

	icon := rl.GenImageGradientLinear(128, 128, 0, rl.Black, ColorBackground)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 16, Y: 112}, rl.Vector2{X: 32, Y: 32}, 8, ColorLine)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 32, Y: 32}, rl.Vector2{X: 48, Y: 80}, 8, ColorLine)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 48, Y: 80}, rl.Vector2{X: 64, Y: 48}, 8, ColorLine)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 64, Y: 48}, rl.Vector2{X: 80, Y: 96}, 8, ColorLine)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 80, Y: 96}, rl.Vector2{X: 96, Y: 16}, 8, ColorLine)
	rl.ImageDrawLineEx(icon, rl.Vector2{X: 96, Y: 16}, rl.Vector2{X: 112, Y: 64}, 8, ColorLine)
	rl.SetWindowIcon(*icon)

	for !rl.WindowShouldClose() {
		if (rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)) && rl.IsKeyPressed(rl.KeyO) {
			exec.Command("explorer", stats.Folder).Start()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		width := int32(rl.GetScreenWidth())
		height := int32(rl.GetScreenHeight())
		mouse_x := rl.GetMouseX()
		mouse_y := rl.GetMouseY()
		graph_width := width - 604
		graph_height := height - 267
		navigator_height := min((height-330)/buttonHeight*buttonHeight, int32((len(stats.Stats)-1)*buttonHeight))
		show_text := input_text
		if len(show_text) == 0 {
			show_text = "stat name"
		}
		max_graph_points := int32(graph_width-100) / 50

		if navigator_scroll+1 > int32(len(stats.Stats)-int(navigator_height/buttonHeight)) {
			navigator_scroll = int32(len(stats.Stats)) - navigator_height/buttonHeight - 1
		}
		if max_graph_points < int32(len(stats.Stats[selected].Values)) && graph_scroll > int32(len(stats.Stats[selected].Values))-max_graph_points {
			graph_scroll = int32(len(stats.Stats[selected].Values)) - max_graph_points
		} else if max_graph_points >= int32(len(stats.Stats[selected].Values)) {
			graph_scroll = 0
		}

		// Draw background and sections
		rl.DrawRectangleGradientV(0, 0, width, height, rl.Black, ColorBackground)
		rl.DrawRectangle(0, 0, width, height, color.RGBA{R: 0, G: 0, B: 0, A: 128})
		rl.DrawRectangle(500, 100, width-600, height-200, ColorBorder)
		rl.DrawRectangleGradientV(502, 102, graph_width, height-204, rl.Black, ColorBackground)
		rl.DrawRectangle(100, 100, 300, height-200, ColorBorder)
		rl.DrawRectangleGradientV(102, 102, 296, height-204, rl.Black, ColorBackground)
		rl.DrawRectangle(100, height-165, 300, 2, ColorBorder)
		rl.DrawRectangle(163, height-165, 2, 64, ColorBorder)
		rl.DrawRectangle(335, height-165, 2, 64, ColorBorder)
		rl.DrawRectangle(100, 163, 300, 2, ColorBorder)
		rl.DrawRectangle(163, 100, 2, 64, ColorBorder)
		rl.DrawRectangle(335, 100, 2, 64, ColorBorder)
		rl.DrawRectangle(500, height-165, width-600, 2, ColorBorder)

		// Draw top buttons
		rl.DrawRectangle(115, 130, 35, 5, ColorText)
		rl.DrawRectangle(350, 130, 35, 5, ColorText)
		rl.DrawRectangle(365, 115, 5, 35, ColorText)
		rl.DrawText(show_text, 250-rl.MeasureText(show_text, 20)/2, 122, 20, ColorText)

		// Draw bottom buttons
		rl.DrawRectangle(115, height-135, 35, 5, ColorText)
		rl.DrawRectangle(350, height-135, 35, 5, ColorText)
		rl.DrawRectangle(365, height-150, 5, 35, ColorText)
		rl.DrawText(strconv.Itoa(int(stats.Get(selected))), 250-rl.MeasureText(strconv.Itoa(int(stats.Get(selected))), 40)/2, height-152, 40, ColorText)

		// Draw navigator
		for i := int(navigator_scroll); i < int(navigator_height/buttonHeight+navigator_scroll); i++ {
			if i+1 == selected {
				rl.DrawText(stats.Stats[i+1].Name, 250-rl.MeasureText(stats.Stats[i+1].Name, fontSize)/2, 165+buttonHeight/2-fontSize/2+(int32(i)-navigator_scroll)*buttonHeight, fontSize, ColorLine)
			} else {
				rl.DrawText(stats.Stats[i+1].Name, 250-rl.MeasureText(stats.Stats[i+1].Name, fontSize)/2, 165+buttonHeight/2-fontSize/2+(int32(i)-navigator_scroll)*buttonHeight, fontSize, ColorText)
			}
		}

		// Draw stat info
		current_streak_val := int32(0)
		for i := len(stats.Stats[selected].Values) - 1; i >= 0; i-- {
			if stats.Stats[selected].Values[i] > 0 {
				current_streak_val++
			} else {
				break
			}
		}
		total := fmt.Sprintf("Total: %d", stats.Stats[selected].Total+stats.Get(selected))
		longest_streak := fmt.Sprintf("Longest streak: %d", max(stats.Stats[selected].LongestStreak, current_streak_val))
		current_streak := fmt.Sprintf("Current streak: %d", current_streak_val)
		max_text := fmt.Sprintf("Max: %d", max(stats.Stats[selected].Max, stats.Get(selected)))

		info_offset := (graph_width - rl.MeasureText(total, 20) - rl.MeasureText(longest_streak, 20) - rl.MeasureText(current_streak, 20) - rl.MeasureText(max_text, 20)) / 5
		rl.DrawText(total, 502+info_offset, height-142, 20, rl.White)
		rl.DrawText(longest_streak, 502+info_offset*2+rl.MeasureText(total, 20), height-142, 20, rl.White)
		rl.DrawText(current_streak, 502+info_offset*3+rl.MeasureText(total, 20)+rl.MeasureText(longest_streak, 20), height-142, 20, rl.White)
		rl.DrawText(max_text, 502+info_offset*4+rl.MeasureText(total, 20)+rl.MeasureText(longest_streak, 20)+rl.MeasureText(current_streak, 20), height-142, 20, rl.White)

		// Draw graph
		offset := (graph_width - 100) / (min(max_graph_points, int32(len(stats.Stats[selected].Values))) - 1)
		lowest_value := int32(2147483647)
		for i := range stats.Stats[selected].Values {
			if stats.Stats[selected].Values[i] < lowest_value {
				lowest_value = stats.Stats[selected].Values[i]
			}
		}
		max_value := max(stats.Stats[selected].Max, stats.Get(selected))
		ratio := (float32(graph_height) - 100.) / float32(max_value-lowest_value)
		if selected != 0 && max_value-lowest_value == 0 {
			ratio = 0
		}

		start := len(stats.Stats[selected].Values) - min(int(max_graph_points), len(stats.Stats[selected].Values)) - int(graph_scroll)

		if selected != 0 {
			for i := start; i < len(stats.Stats[selected].Values)-int(graph_scroll); i++ {
				if i-len(stats.Stats[selected].Values)+1 == 0 {
					rl.DrawText("Today", 552+int32(i-start)*offset-rl.MeasureText("Today", 20)/2, 67+graph_height, 20, rl.White)
				} else {
					rl.DrawText(strconv.Itoa(i-len(stats.Stats[selected].Values)+1), 552+int32(i-start)*offset-rl.MeasureText(strconv.Itoa(i-len(stats.Stats[selected].Values)+1), 20)/2, 67+graph_height, 20, rl.White)
				}
			}
		}

		if mouse_x >= 502 && mouse_x <= width-103 && mouse_y >= 102 && mouse_y <= height-166 {
			hover_offset := (graph_width - 100) / (min(max_graph_points, int32(len(stats.Stats[selected].Values))) - 1)
			for i := start; i < len(stats.Stats[selected].Values)-1-int(graph_scroll); i++ {
				rl.DrawLineEx(rl.Vector2{X: float32(552 + int32(i-start)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i]-lowest_value)*ratio))},
					rl.Vector2{X: float32(552 + int32(i-start+1)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i+1]-lowest_value)*ratio))}, 2, ColorLine)

				rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(i-start+1)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i+1]-lowest_value)*ratio))}, 1, ColorLine)
			}
			rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(0)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[start]-lowest_value)*ratio))}, 1, ColorLine)

			found := false
			for i := len(stats.Stats[selected].Values) - 1 - int(graph_scroll); i >= start+1; i-- {
				if mouse_x > 552-offset/2+int32(i-start)*hover_offset {
					rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(i-start)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i]-lowest_value)*ratio))}, 4, ColorLine)
					rl.DrawText(strconv.Itoa(int(stats.Stats[selected].Values[i])), 552+int32(i-start)*offset-rl.MeasureText(strconv.Itoa(int(stats.Stats[selected].Values[i])), 20)/2, 52+graph_height-int32(float32(stats.Stats[selected].Values[i]-lowest_value)*ratio)-30, 20, ColorLine)
					found = true
					break
				}
			}
			if !found {
				rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(0)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[start]-lowest_value)*ratio))}, 4, ColorLine)
				rl.DrawText(strconv.Itoa(int(stats.Stats[selected].Values[start])), 552+int32(0)*offset-rl.MeasureText(strconv.Itoa(int(stats.Stats[selected].Values[start])), 20)/2, 52+graph_height-int32(float32(stats.Stats[selected].Values[start]-lowest_value)*ratio)-30, 20, ColorLine)
			}

			if wheel := rl.GetMouseWheelMove(); wheel > 0 && max_graph_points < int32(len(stats.Stats[selected].Values)) && graph_scroll < int32(len(stats.Stats[selected].Values))-max_graph_points {
				graph_scroll++
			} else if wheel < 0 && graph_scroll > 0 {
				graph_scroll--
			}
		} else {
			for i := start; i < len(stats.Stats[selected].Values)-1-int(graph_scroll); i++ {
				rl.DrawLineEx(rl.Vector2{X: float32(552 + int32(i-start)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i]-lowest_value)*ratio))},
					rl.Vector2{X: float32(552 + int32(i-start+1)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i+1]-lowest_value)*ratio))}, 2, ColorText)

				rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(i-start+1)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[i+1]-lowest_value)*ratio))}, 1, ColorText)
			}
			rl.DrawCircleV(rl.Vector2{X: float32(552 + int32(0)*offset), Y: float32(52 + graph_height - int32(float32(stats.Stats[selected].Values[start]-lowest_value)*ratio))}, 1, ColorText)
		}

		if mouse_x >= 102 && mouse_x <= 397 && mouse_y >= 102 && mouse_y <= height-103 {
			if mouse_y >= height-163 && selected != 0 {
				// Bottom buttons
				if mouse_x <= 162 {
					// -
					rl.DrawRectangle(115, height-135, 35, 5, ColorLine)

					if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
						stats.Add(selected, -10)
					}
				} else if mouse_x >= 337 {
					// +
					rl.DrawRectangle(350, height-135, 35, 5, ColorLine)
					rl.DrawRectangle(365, height-150, 5, 35, ColorLine)

					if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
						stats.Add(selected, 10)
					}
				} else if mouse_x >= 165 && mouse_x <= 334 {
					// Input field
					rl.DrawText(strconv.Itoa(int(stats.Get(selected))), 250-rl.MeasureText(strconv.Itoa(int(stats.Get(selected))), 40)/2, height-152, 40, ColorLine)

					k := rl.GetKeyPressed()
					if k == rl.KeyBackspace {
						value := stats.Get(selected)
						value /= 10
						stats.Set(selected, value)
					} else if c := rl.GetCharPressed(); c >= '0' && c <= '9' {
						value := stats.Get(selected)
						value = value*10 + int32(c-'0')
						stats.Set(selected, value)
					}

					if wheel := rl.GetMouseWheelMove(); wheel > 0 {
						stats.Add(selected, 1)
					} else if wheel < 0 {
						stats.Add(selected, -1)
					}
				}
			} else if mouse_y <= 162 {
				// Top buttons
				if mouse_x <= 162 {
					// -
					rl.DrawRectangle(115, 130, 35, 5, ColorLine)

					if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
						stats.Delete(input_text)
						input_text = ""
						if selected == len(stats.Stats) {
							selected--
						}
						if navigator_scroll != 0 && navigator_scroll == int32(len(stats.Stats)-int(navigator_height/buttonHeight)) {
							navigator_scroll--
						}
					}
				} else if mouse_x >= 337 {
					// +
					rl.DrawRectangle(350, 130, 35, 5, ColorLine)
					rl.DrawRectangle(365, 115, 5, 35, ColorLine)

					if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
						if stats.New(input_text) {
							selected = len(stats.Stats) - 1
							if int32(len(stats.Stats)) > (height-330)/buttonHeight+1 {
								navigator_scroll = int32(len(stats.Stats)) - navigator_height/buttonHeight - 1
							}
						}
						input_text = ""
					}
				} else if mouse_x >= 165 && mouse_x <= 334 {
					// Input field
					rl.DrawText(show_text, 250-rl.MeasureText(show_text, 20)/2, 122, 20, ColorLine)

					k := rl.GetKeyPressed()
					if k == rl.KeyBackspace {
						if len(input_text) > 0 {
							input_text = input_text[:len(input_text)-1]
						}
					}
					if c := rl.GetCharPressed(); c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == ' ' || c == '_' || c == '-' {
						if len(input_text) < 20 {
							input_text += string(c)
						}
					}

					if rl.IsMouseButtonPressed(rl.MouseRightButton) {
						input_text = ""
					}
				}
			} else if mouse_y <= 165+navigator_height {
				//  Navigator
				if wheel := rl.GetMouseWheelMove(); wheel < 0 && navigator_scroll+1 < int32(len(stats.Stats)-int(navigator_height/buttonHeight)) {
					navigator_scroll++
				} else if wheel > 0 && navigator_scroll > 0 {
					navigator_scroll--
				}

				for i := int(navigator_scroll); i < int(navigator_height/buttonHeight+navigator_scroll); i++ {
					if mouse_y >= 165+(int32(i)-navigator_scroll)*buttonHeight && mouse_y <= 165+buttonHeight+(int32(i)-navigator_scroll)*buttonHeight {
						rl.DrawText(stats.Stats[i+1].Name, 250-rl.MeasureText(stats.Stats[i+1].Name, fontSize)/2, 165+fontSize/2+(int32(i)-navigator_scroll)*buttonHeight, fontSize, ColorLine)
						if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
							selected = i + 1
							graph_scroll = 0
						}
					}
				}
			}
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
