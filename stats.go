package main

import (
	"os"
	"time"
	"unsafe"
)

type Stats struct {
	Folder          string
	LastDayRecorded int
	Stats           []Stat
}

type Stat struct {
	Name          string
	Values        []int32
	Total         int32
	LongestStreak int32
	Max           int32
}

func (s *Stats) Save() {
	for i := 1; i < len(s.Stats); i++ {
		data := make([]byte, len(s.Stats[i].Values)*4)
		for j := range len(s.Stats[i].Values) * 4 {
			data[j] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&s.Stats[i].Values[0])) + uintptr(j)))
		}

		err := os.WriteFile(s.Folder+s.Stats[i].Name, data, 0644)
		if err != nil {
			panic(err)
		}
	}

	os.WriteFile(s.Folder+".ldr", (*(*[unsafe.Sizeof(s.LastDayRecorded)]byte)(unsafe.Pointer(&s.LastDayRecorded)))[:], 0644)
}

func (s *Stats) Load() {
	_, err := os.Stat(s.Folder)
	if os.IsNotExist(err) {
		os.MkdirAll(s.Folder, 0755)
	}

	_, err = os.Stat(s.Folder + ".ldr")
	if os.IsNotExist(err) {
		day_now := time.Now().Local().YearDay()
		for i := 2025; i < time.Now().Local().Year(); i++ {
			if i%4 == 0 && (i%100 != 0 || i%400 == 0) {
				day_now += 366
			} else {
				day_now += 365
			}
		}
		s.LastDayRecorded = day_now
		s.Save()
		return
	}

	files, err := os.ReadDir(s.Folder)
	if err != nil {
		panic(err)
	}

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		data, err := os.ReadFile(s.Folder + files[i].Name())
		if err != nil {
			panic(err)
		}

		if files[i].Name() == ".ldr" {
			s.LastDayRecorded = *(*int)(unsafe.Pointer(&data[0]))
			continue
		}

		s.Stats = append(s.Stats, Stat{
			Name:          files[i].Name(),
			Values:        make([]int32, len(data)/4),
			Total:         0,
			LongestStreak: 0,
			Max:           0,
		})

		for i := range len(data) / 4 {
			s.Stats[len(s.Stats)-1].Values[i] = *(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&data[0])) + uintptr(i*4)))
		}

		s.CheckLastDayRecorded(len(s.Stats) - 1)

		last_time := false
		streak := int32(0)

		for i := 0; i < len(s.Stats[len(s.Stats)-1].Values)-1; i++ {
			if s.Stats[len(s.Stats)-1].Values[i] > s.Stats[len(s.Stats)-1].Max {
				s.Stats[len(s.Stats)-1].Max = s.Stats[len(s.Stats)-1].Values[i]
			}
			s.Stats[len(s.Stats)-1].Total += s.Stats[len(s.Stats)-1].Values[i]

			if s.Stats[len(s.Stats)-1].Values[i] > 0 {
				if last_time {
					streak++
				} else {
					last_time = true
					streak = 1
				}
				if streak > s.Stats[len(s.Stats)-1].LongestStreak {
					s.Stats[len(s.Stats)-1].LongestStreak = streak
				}
			} else {
				last_time = false
			}
		}
	}

	day_now := time.Now().Local().YearDay()
	for i := 2025; i < time.Now().Local().Year(); i++ {
		if i%4 == 0 && (i%100 != 0 || i%400 == 0) {
			day_now += 366
		} else {
			day_now += 365
		}
	}
	s.LastDayRecorded = day_now
	s.Save()
}

func (s *Stats) CheckLastDayRecorded(i int) {
	// Calculate the number of days since 2025
	day_now := time.Now().Local().YearDay()
	for i := 2025; i < time.Now().Local().Year(); i++ {
		if i%4 == 0 && (i%100 != 0 || i%400 == 0) {
			day_now += 366
		} else {
			day_now += 365
		}
	}

	if day_now <= s.LastDayRecorded {
		return
	}

	ldr := s.LastDayRecorded

	// If the current day is greater than the last recorded day, update stats
	for day_now > ldr {
		s.Stats[i].Values = append(s.Stats[i].Values, 0)
		if len(s.Stats[i].Values) == 3 && s.Stats[i].Values[0] == 0 {
			s.Stats[i].Values = s.Stats[i].Values[1:]
		}
		ldr++
	}

	s.Save()
}

func (s *Stats) New(name string) bool {
	if name == ".ldr" || name == "0" || name == "" {
		return false
	}

	for i := range s.Stats {
		if s.Stats[i].Name == name {
			return false
		}
	}

	s.Stats = append(s.Stats, Stat{
		Name:   name,
		Values: []int32{0, 0},
	})

	s.Save()

	return true
}

func (s *Stats) Delete(name string) {
	if name == ".ldr" || name == "0" || name == "" {
		return
	}

	for i := range s.Stats {
		if s.Stats[i].Name == name {
			os.Remove(s.Folder + s.Stats[i].Name)
			if i == len(s.Stats)-1 {
				s.Stats = s.Stats[:len(s.Stats)-1]
			} else {
				s.Stats = append(s.Stats[:i], s.Stats[i+1:]...)
			}
			return
		}
	}
}

func (s *Stats) Set(i int, value int32) {
	if i < 0 || i >= len(s.Stats) {
		return
	}

	s.Stats[i].Values[len(s.Stats[i].Values)-1] = value
	if s.Stats[i].Values[len(s.Stats[i].Values)-1] < 0 {
		s.Stats[i].Values[len(s.Stats[i].Values)-1] = 0
	}
	s.Save()
}

func (s *Stats) Add(i int, value int32) {
	if i < 0 || i >= len(s.Stats) {
		return
	}

	s.Stats[i].Values[len(s.Stats[i].Values)-1] += value
	if s.Stats[i].Values[len(s.Stats[i].Values)-1] < 0 {
		s.Stats[i].Values[len(s.Stats[i].Values)-1] = 0
	}
	s.Save()
}

func (s *Stats) Get(i int) int32 {
	if i < 0 || i >= len(s.Stats) {
		return 0
	}

	return s.Stats[i].Values[len(s.Stats[i].Values)-1]
}
