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

	for i,p := range pl {
		fmt.Printf("%7d -> %8d %s\n", i, p.Value, p.Key)
    if(i>=20) {
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
	for _,p := range pl {
    writer.Write( []string{p.Key, strconv.Itoa(p.Value)} )
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

	return splitter
}

func (self Splitter) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return 
	}
	defer file.Close()
	writer := csv.NewWriter(file)
  
	for w, sa := range self {
    //line := fmt.Sprintf("\"%s\",%d,%d\n", strings.Replace(w, "\"", "\"\"", -1), sa.Together, sa.Separate)
    writer.Write( []string{w, strconv.Itoa(sa.Together), strconv.Itoa(sa.Separate)} )
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

func (self Splitter) CreateSubmission(filename_test string, filename_submit string, vocab *Vocab, skip_check int) {
	file_in, err := os.Open(filename_test)
	if err != nil {
		fmt.Println("Error file_in:", err)
		return
	}
	defer file_in.Close()
	reader := csv.NewReader(file_in)
  
	file_out, err := os.Create(filename_submit)
	if err != nil {
		fmt.Println("Error file_out:", err)
		return 
	}
	defer file_out.Close()
	writer := bufio.NewWriter(file_out)  // since csv.Writer doesn't allow force quoting
  
	//fmt.Printf("filename_submit = %s", filename_submit)

	// First line different
	header, err := reader.Read()
	if header[0] != "id" {
		fmt.Println("Bad Header", err)
		return
	}
  writer.WriteString("\"id\",\"sentence\"\n");

	for {
    if(skip_check>0) {
      // Waste a few lines...  (3032 lines in heldout.txt.csv)
      for i:=0; i<skip_check-1; i++ {
        reader.Read()
      }
    }

		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record is []string

		//id, _ := strconv.ParseFloat(record[0], 32)
    id  := record[0] // Don't really care about content
		txt := record[1]

    //fmt.Printf("%6.4f\n", id)

    best_i := -1
    best_v := -1

		words := strings.Split(txt, " ")
    words[0] = strings.ToLower(words[0])
    for i := 0; i < len(words)-1; i++ {
      word := words[i] + "|" + words[i+1]

      v0 := (*vocab)[words[i]]
      v1 := (*vocab)[words[i+1]]
      
      // Let's print the word, and its corresponding stats
      sa := self[word]
      tot := sa.Together+sa.Separate
      if 0==tot {
         tot=1
      }
      
      max_prop := 0
      if v0>0 && v1>0 {
        max_prop = (tot*100)/v0
        if max_prop < (tot*100)/v1 {
          max_prop = (tot*100)/v1
        }
      }
      
      v := (max_prop * sa.Separate)/tot
      //if (sa.Separate*100)/tot<50 {
      if sa.Separate<50 { // No evidence
        v=0 //
      }
      if (sa.Separate*100)/tot<75 { // poor percentage suggesting split
        v=0
      }
      
      if v>best_v {
        best_i=i
        best_v=v
      }
      
      if(skip_check>0) {
        fmt.Printf("%12s - %12s :: [%7d,%7d] :: %7d %3d%% :: vocab:(%8d,%8d)=(%3d%%,%3d%%) -> %3d%%\n", words[i], words[i+1], 
          sa.Together, sa.Separate, sa.Together+sa.Separate, (sa.Separate*100)/tot,
          v0, v1, (tot*100)/v0, (tot*100)/v1,
          v )
      }
    }
    
    //line := fmt.Sprintf("\"%s\",%d,%d\n", strings.Replace(w, "\"", "\"\"", -1), sa.Together, sa.Separate)
    //writer.Write( []string{id, strings.Join( words_verbatim, " ")} )
    
    words_output := strings.Split(txt, " ")
    
    if best_v>0 { // Insert it into the list of words
      i := best_i+1
      words_output = append(words_output, "")
      copy(words_output[i+1:], words_output[i:])
      words_output[i] = "" // Insert an empty word...
    }
    
    txt_out := strings.Join( words_output, " ")
    writer.WriteString( fmt.Sprintf("%s,\"%s\"\n", id, strings.Replace(txt_out, "\"", "\"\"", -1)) )
    
    if(skip_check>0) {
      break
    }
	}
  writer.Flush()
}


const currently_running_version int = 1000

func main() {
	cmd := flag.String("cmd", "", "Required : {size}")
	cmd_type := flag.String("type", "", "size:{vocab|bigrams}")

	file_save := flag.String("save", "", "filename")
	file_load := flag.String("load", "", "filename")

	//delta := flag.Int("delta", 0, "Number of steps between start and end")
	seed := flag.Int64("seed", 1, "Random seed to use")
  
	skip := flag.Int("skip", 0, "Debugging aid")

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
		/// ./billion -cmd=size -type=vocab -save=0-vocab.csv
		if *cmd_type == "vocab" {
      vocab := Vocab{}
      
			// Read in the vocab for test file - counts will be disgarded
			vocab.ReadTestNGram(fname_test, 1)
      //vocab.MakeSortedPairList()
      
			// Read in the vocab for holdout file (additional, for validation)
			vocab.ReadTestNGram(fname_heldout, 1)
      
      test_pairs := vocab.MakeSortedPairList()
      
			fmt.Printf("Billion elapsed : %s\n", time.Since(start))

      vocab_train := Vocab{}
			vocab_train.ReadTrainingFile(fname_train)
      train_pairs := vocab_train.MakeSortedPairList()
      
			fmt.Printf("Billion elapsed : %s\n", time.Since(start))

			// Create an empty test vocab from the test set 'skeleton'
			test_vocab := Vocab{}
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
      
      if len(*file_save)>0 {
        test_vocab.Save(*file_save)
      }
		}

		/// ./billion -cmd=size -type=bigrams -save=0-bigrams.csv
		if *cmd_type == "bigrams" {
      vocab := Vocab{}
      
			// Read in the vocab for test file
			vocab.ReadTestNGram(fname_test, 2)
			// Read in the vocab for holdout file (additional, for validation)
			vocab.ReadTestNGram(fname_heldout, 2)

      pl := vocab.MakeSortedPairList()
      
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
      vocab := Vocab{}
      vocab.Load("0-vocab.csv") // Hard coded for now
      
      splitter := Splitter{}
      splitter.Load(*file_load)

      splitter.CreateSubmission(fname_validation, "1-valid"+*file_save, &vocab, *skip)
      //splitter.CreateSubmission(fname_test, "1-test"+*file_save, &vocab, *skip)
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
