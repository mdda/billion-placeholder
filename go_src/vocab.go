package main

import (
	"fmt"
)

import (
	"sort"
)

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

import (
	"bufio"
)

type Vocab map[string]int

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value } // Sort DESC

// A function to turn a map into a PairList, then sort and return it.
//func sortMapByValue(m *map[string]int) PairList {
func sortMapByValue(m *Vocab) PairList {
	p := make(PairList, len(*m))
	i := 0
	for k, v := range *m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

// This adds the results of the read into the vocab passed as a pointer
func (vocab *Vocab) ReadTestNGram(filename string, n int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// First line different
	header, err := reader.Read()
	if header[0] != "id" {
		fmt.Println("Bad Header", err)
		return
	}

	/*
		id_max := 0
		id_map := make(map[int]bool)
		for _, id := range id_list {
			id_map[id] = true
			if id_max < id {
				id_max = id
			}
		}
	*/

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record is []string

		//id, _ := strconv.Atoi(record[0])
		txt := record[1]

		//fmt.Println("id,txt:", id, txt)

		words := strings.Split(txt, " ")
		if 1 == n {
			for i := 0; i < len(words); i++ {
				word := words[i]
				//fmt.Println("word:", word)
				(*vocab)[word]++
			}
		}

		if 2 == n {
			words[0] = strings.ToLower(words[0])
			for i := 0; i < len(words)-1; i++ {
				word := words[i] + "|" + words[i+1]
				//fmt.Println("word:", word)
				(*vocab)[word]++
			}
		}

		/*
			if id_map[id] {
				//fmt.Println(record) // record has the type []string

				steps:=0
				var data []string

				if has_steps {
					steps, _ = strconv.Atoi(record[1])
					data = record[2:]
				} else {
					data = record[1:]
				}

				start := NewBoard_BoolPacked(board_width, board_height)
				end := NewBoard_BoolPacked(board_width, board_height)
				if is_training {
					start.LoadArray(data[0:400])
					end.LoadArray(data[400:800])
				} else {
					end.LoadArray(data[0:400])
				}

				s.problem[id] = LifeProblem{
					id:    id,
					start: start,
					end:   end,
					steps: steps,
				}
				fmt.Printf("Loaded problem[%d] : steps=%d\n", id, steps)
				//fmt.Print(s.problem[id].start)
			}
			if id > id_max {
				return // fact-of-life : ids are ascending order, so can quit reading early
			}
		*/
	}
}

func (vocab *Vocab) MakeSortedPairList() PairList {
	pl := sortMapByValue(vocab)

	fmt.Printf("Vocab size : %d\n", len(pl))

	for i, p := range pl {
		fmt.Printf("%7d -> %8d %s\n", i, p.Value, p.Key)
		if i >= 20 {
			break
		}
	}

	return pl
}

func (vocab *Vocab) ReadTrainingFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		//  fmt.Println(scanner.Text())
		txt := scanner.Text()
		for _, word := range strings.Split(txt, " ") {
			(*vocab)[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		//log.Fatal(err)
		fmt.Println(err)
		return
	}
}

func (self *Vocab) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)

	pl := sortMapByValue(self) // Far prefer this to be ordered...
	for _, p := range pl {
		writer.Write([]string{p.Key, strconv.Itoa(p.Value)})
	}
	writer.Flush()
}

func (self Vocab) Load(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record is []string
		word := record[0]
		count, _ := strconv.Atoi(record[1])

		self[word] = count
	}
}
