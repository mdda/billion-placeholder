package main

// GOPATH=`pwd` go build billion.go  && ./billion

import (
	"fmt"
	"time"
	"flag"
	"math/rand"
)

import (
	"sort"
)

import (
  "os"
  "io"
	"encoding/csv"
//	"strconv"
  "strings"
)

import (
	"bufio"
)

// A data structure to hold a key/value pair.
type Pair struct {
  Key string
  Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value } // Sort DESC

// A function to turn a map into a PairList, then sort and return it. 
func sortMapByValue(m map[string]int) PairList {
   p := make(PairList, len(m))
   i := 0
   for k, v := range m {
      p[i] = Pair{k, v}
      i++
   }
   sort.Sort(p)
   return p
}

func read_test(filename string) PairList {
	vocab := map[string]int{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return PairList{}
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// First line different
	header, err := reader.Read()
	if header[0] != "id" {
		fmt.Println("Bad Header", err)
		return PairList{}
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
      return PairList{}
		}
		// record is []string
    
		//id, _ := strconv.Atoi(record[0])
    txt   := record[1]
    
		//fmt.Println("id,txt:", id, txt)
    
    for _, word := range strings.Split(txt, " ") {
      //fmt.Println("word:", word)
      vocab[word]++
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
  
  pl := sortMapByValue(vocab)
  
  l := len(pl)
  fmt.Printf("Test Vocab size : %d\n", l)
  
  if(l>25) { l=25 }
  for i := 0; i<l; i++ {
    fmt.Printf("%7d -> %7d %s\n", i, pl[i].Value, pl[i].Key)
  }  
  
  return pl
}

func read_train(filename string) PairList {
	vocab := map[string]int{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return PairList{}
	}
	defer file.Close()
  
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    //  fmt.Println(scanner.Text())
    txt := scanner.Text()
    for _, word := range strings.Split(txt, " ") {
      vocab[word]++
    }
  }

  if err := scanner.Err(); err != nil {
    //log.Fatal(err)
    fmt.Println(err)
		return PairList{}
  }

  pl := sortMapByValue(vocab)
  
  l := len(pl)
  fmt.Printf("Train Vocab size : %d\n", l)
  
  if(l>25) { l=25 }
  for i := 0; i<l; i++ {
    fmt.Printf("%7d -> %7d %s\n", i, pl[i].Value, pl[i].Key)
  }  
  
  return pl
}


const currently_running_version int = 1000

func main() {
	cmd := flag.String("cmd", "", "Required : {}")
	cmd_type := flag.String("type", "", "create:{fake_training_data|training_set_transitions|synthetic_transitions}, db:{test|insert_problems}, visualize:{data|ga}, submit:{kaggle|fakescore}")

	//delta := flag.Int("delta", 0, "Number of steps between start and end")
	seed := flag.Int64("seed", 1, "Random seed to use")

	//id := flag.Int("id", 0, "Specific id to examine")
	//training_only := flag.Bool("training", false, "Act on training set (default=false, i.e. test set)")

	//count := flag.Int("count", 0, "Number of ids to process")

	flag.Parse()
	//fmt.Printf("CMD = %s\n", *cmd)

	//rand.Seed(time.Now().UnixNano())
	rand.Seed(*seed)

	fmt.Printf("Billion Start : %s\n", time.Now().Format("2006-01-02_15-04_05:06:07"))
  start := time.Now()

	fname_test  := "../data/0-orig/test_v2.txt"
	fname_train := "../data/0-orig/train_v2.txt"
  
  // Read in the vocab for test
  test_pairs  := read_test(fname_test)
	fmt.Printf("Billion elapsed : %s\n", time.Since(start))
  
  train_pairs := read_train(fname_train)
	fmt.Printf("Billion elapsed : %s\n", time.Since(start))
  
  // Create an empty test vocab
	test_vocab := map[string]int{}
  for i:=0; i<len(test_pairs); i++ {
    p := test_pairs[i]
    test_vocab[p.Key] = 0
  }
  
  // The fill it with train word freqs (where applicable)
  for i:=0; i<len(train_pairs); i++ {
    p := train_pairs[i]
    if _, ok := test_vocab[p.Key]; ok {
      test_vocab[p.Key] = p.Value
    }
  }
  
  // Count the non-zero-freq test words
  nonzero :=0
  
  hist_max := 10
  hist := make([]int,hist_max)
  
  for _,v := range test_vocab {
    if v>0 {
      nonzero++
    }
    if v<hist_max {
      hist[v]++
    }
  }

	fmt.Printf("NonZero   test vocab words : %d\n", nonzero)
  for i:=0; i<hist_max; i++ {
    fmt.Printf("%2d occurences in train : %d\n", i, hist[i])
  }
  
	fmt.Printf("Billion elapsed : %s\n", time.Since(start))

	if *cmd == "db" {
		/// ./billion -cmd=db -type=test
		if *cmd_type == "test" {
			//test_open_db()
		}

	}

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
