package main

// GOPATH=`pwd` go build billion.go  && ./billion

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
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
func sortMapByValue(m *map[string]int) PairList {
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
func read_test_ngram(filename string, vocab *map[string]int, n int) {
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

func vocab_to_pairslist(vocab *map[string]int) PairList {
	pl := sortMapByValue(vocab)

	l := len(pl)
	fmt.Printf("Vocab size : %d\n", l)

	if l > 25 {
		l = 25
	}
	for i := 0; i < l; i++ {
		fmt.Printf("%7d -> %7d %s\n", i, pl[i].Value, pl[i].Key)
	}

	return pl
}

func read_train(filename string, vocab *map[string]int) {
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


type SplitterAtom struct {
  Together int
  Separate int
}
type Splitter map[string]SplitterAtom

func get_train_ngrams(filename string, pl PairList) Splitter {
	splitter := Splitter{}
  for i:=0; i<len(pl); i++ {
    splitter[pl[i].Key] = SplitterAtom{0,0}
  }

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return splitter
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		txt := scanner.Text()
    words := strings.Split(txt, " ")
    words[0] = strings.ToLower(words[0])
    
    for i := 0; i < len(words)-1; i++ {
      word := words[i] + "|" + words[i+1]
      if sa, ok := splitter[word]; ok {
        sa.Together++
        splitter[word] = sa
      }
    }
    for i := 0; i < len(words)-2; i++ {
      word := words[i] + "|" + words[i+2]
      if sa, ok := splitter[word]; ok {
        sa.Separate++
        splitter[word] = sa
      }
    }
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return splitter
	}

  // Stuff here

	return splitter
}

func (self Splitter) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return 
	}
	defer file.Close()
  
	writer := bufio.NewWriter(file)
  
	for w, sa := range self {
    line := fmt.Sprintf("\"%s\",%d,%d\n", strings.Replace(w, "\"", "\"\"", -1), sa.Together, sa.Separate)
    writer.WriteString(line)
  }
  writer.Flush()
}

func (self Splitter) Load(filename string) {
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
    together, _ := strconv.Atoi(record[1])
    separate, _ := strconv.Atoi(record[2])
    
    self[word] = SplitterAtom{together,separate}
  }
}

func (self Splitter) CreateSubmission(filename_test string, filename_submit string) {
	file, err := os.Open(filename_test)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
  
	fmt.Printf("filename_submit = %s", filename_submit)

	// First line different
	header, err := reader.Read()
	if header[0] != "id" {
		fmt.Println("Bad Header", err)
		return
	}

	for {
    // Waste a few lines...
    for i:=0; i<100; i++ {
      reader.Read()
    }

		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record is []string

		id, _ := strconv.ParseFloat(record[0], 32)
		txt := record[1]

		//fmt.Println("id,txt:", id, txt)
    fmt.Printf("%6.4f\n", id)

		words := strings.Split(txt, " ")
    words[0] = strings.ToLower(words[0])
    for i := 0; i < len(words)-1; i++ {
      word := words[i] + "|" + words[i+1]
      
      // Let's print the word, and its corresponding stats
      sa := self[word]
      tot := sa.Together+sa.Separate
      if 0==tot {
         tot=1
      }
      fmt.Printf("%20s - %20s :: [%7d,%7d] :: %7d %3d%%\n", words[i], words[i+1], 
        sa.Together, sa.Separate, sa.Together+sa.Separate, 100-(sa.Together*100)/tot)
      
    }
    break
	}
}


const currently_running_version int = 1000

func main() {
	cmd := flag.String("cmd", "", "Required : {size}")
	cmd_type := flag.String("type", "", "size:{vocab|bigrams}")

	file_save := flag.String("save", "", "filename")
	file_load := flag.String("load", "", "filename")

	//delta := flag.Int("delta", 0, "Number of steps between start and end")
	seed := flag.Int64("seed", 1, "Random seed to use")

	//id := flag.Int("id", 0, "Specific id to examine")
	//training_only := flag.Bool("training", false, "Act on training set (default=false, i.e. test set)")

	//count := flag.Int("count", 0, "Number of ids to process")

	flag.Parse()
	//fmt.Printf("CMD = %s\n", *cmd)

	//rand.Seed(time.Now().UnixNano())
	rand.Seed(*seed)

	fmt.Printf("Billion Start : %s\n", time.Now().Format("2006-01-02_15h04m05s"))
	start := time.Now()

	fname_test := "../data/0-orig/test_v2.txt"
	fname_heldout := "../data/1-holdout/heldout.txt.csv"  // This was the data not included in the training set (same format as test)
	fname_validation := "../data/1-holdout/valid.txt"     // This is a test set for which we have perfect comparison ('truth.txt')
  
	//fname_train := "../data/0-orig/train_v2.txt"
	fname_train := "../data/1-holdout/train.txt"

	if *cmd == "size" {
		/// ./billion -cmd=size -type=vocab
		if *cmd_type == "vocab" {
      vocab := map[string]int{}
      
			// Read in the vocab for test file
			read_test_ngram(fname_test, &vocab, 1)
      vocab_to_pairslist(&vocab)
      
			// Read in the vocab for holdout file (additional, for validation)
			read_test_ngram(fname_heldout, &vocab, 1)
      
      test_pairs := vocab_to_pairslist(&vocab)
      
			fmt.Printf("Billion elapsed : %s\n", time.Since(start))

      vocab_train := map[string]int{}
			read_train(fname_train, &vocab_train)
      train_pairs := vocab_to_pairslist(&vocab_train)
      
			fmt.Printf("Billion elapsed : %s\n", time.Since(start))

			// Create an empty test vocab
			test_vocab := map[string]int{}
			for i := 0; i < len(test_pairs); i++ {
				p := test_pairs[i]
				test_vocab[p.Key] = 0
			}

			// The fill it with train word freqs (where applicable)
			for i := 0; i < len(train_pairs); i++ {
				p := train_pairs[i]
				if _, ok := test_vocab[p.Key]; ok {
					test_vocab[p.Key] = p.Value
				}
			}

			// Count the non-zero-freq test words
			nonzero := 0

			hist_max := 10
			hist := make([]int, hist_max)

			for _, v := range test_vocab {
				if v > 0 {
					nonzero++
				}
				if v < hist_max {
					hist[v]++
				}
			}

			fmt.Printf("NonZero   test vocab words : %d\n", nonzero)
			for i := 0; i < hist_max; i++ {
				fmt.Printf("%2d occurences in train : %d\n", i, hist[i])
			}
		}

		/// ./billion -cmd=size -type=bigrams -save=0-bigrams.csv
		if *cmd_type == "bigrams" {
      vocab := map[string]int{}
      
			// Read in the vocab for test file
			read_test_ngram(fname_test, &vocab, 2)
			// Read in the vocab for holdout file (additional, for validation)
			read_test_ngram(fname_heldout, &vocab, 2)

      pl := vocab_to_pairslist(&vocab)
      
      // Now, go through the training set, building up a picture of '' or 'something' for each found bigram
      splitter := get_train_ngrams(fname_train, pl)
      
      if len(*file_save)>0 {
        splitter.Save(*file_save)
      }
		}
	}

	if *cmd == "validate" {
		/// ./billion -cmd=validate -type=bigrams -load=0-bigrams.csv -save=.bigram_01.csv
		if *cmd_type == "bigrams" {
      //vocab := map[string]int{}
      splitter := Splitter{}
      splitter.Load(*file_load)

      splitter.CreateSubmission(fname_validation, "1-valid"+*file_save)
      //splitter.CreateSubmission(fname_test, "1-test"+*file_save)
    }
  }
	fmt.Printf("Billion elapsed : %s\n", time.Since(start))

	/*
		if *cmd=="create" {
			/// ./reverse-gol -cmd=create -type=fake_training_data
			if *cmd_type=="fake_training_data" {
				if *seed==1 {
					fmt.Println("Must not have seed same as one used to generate Synthetic Transitions!")
					flag.Usage()
					return
				}
				main_create_fake_training_data()

				// Prevent solving of actual training set (since this is where our state came from, so it's not particularly helpful
				// UPDATE problems SET solution_count=100 WHERE id>-60000  and id<0
				// UPDATE problems SET solution_count=0 WHERE id>-100000 and id<-60000
			}

			/// ./reverse-gol -cmd=create -type=training_set_transitions
			if *cmd_type=="training_set_transitions" {
				main_create_stats_all(true)
				//main_read_stats(1)
			}

			/// ./reverse-gol -cmd=create -type=synthetic_transitions -delta=1
			/// ./reverse-gol -cmd=create -type=synthetic_transitions -delta=2
			/// ./reverse-gol -cmd=create -type=synthetic_transitions -delta=3
			/// ./reverse-gol -cmd=create -type=synthetic_transitions -delta=4
			/// ./reverse-gol -cmd=create -type=synthetic_transitions -delta=5
			if *cmd_type=="synthetic_transitions" {
				if *delta<=0 {
					fmt.Println("Need to specify '-delta=%d' to identify which stats to generate")
					flag.Usage()
					return
				}
				main_create_stats(*delta, false)
				//main_read_stats(1)
			}

		}

		if *cmd=="visualize" {
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=50
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=60001
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=60201
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=60401
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=60601
			/// ./reverse-gol -cmd=visualize -type=data -training=true -id=60801
			if *cmd_type=="data" {
				if *id<=0 {
					fmt.Println("Need to specify '-id=%d' as base id to view (will also show 9 following)")
					flag.Usage()
					return
				}
				if !*training_only {
					fmt.Println("Need to specify '-training=true' (don't know start boards for test...)")
					flag.Usage()
					return
				}
				main_verify_training_examples(*id)
			}

			/// ./reverse-gol -cmd=visualize -type=ga -training=true -id=58
			///
			if *cmd_type=="ga" {
				if *id<=0 {
					fmt.Println("Need to specify '-id=%d'")
					flag.Usage()
					return
				}
				main_population_score(*training_only, *id)
			}
		}
	*/

}
