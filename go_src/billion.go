package main

// GOPATH=`pwd` go build billion.go vocab.go levenshtein.go && ./billion

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

import (
	"math"
)

func get_validation_score(filename_truth string, filename_attempt string) float32 {
	file_truth, err := os.Open(filename_truth)
	if err != nil {
		fmt.Println("Error file_truth:", err)
		return 0.0
	}
	defer file_truth.Close()
	reader_a := csv.NewReader(file_truth)

	file_attempt, err := os.Open(filename_attempt)
	if err != nil {
		fmt.Println("Error file_attempt:", err)
		return 0.0
	}
	defer file_attempt.Close()
	reader_b := csv.NewReader(file_attempt)

	// First line different
	header_a, err := reader_a.Read()
	if header_a[0] != "id" {
		fmt.Println("Bad Header", err)
		return 0.0
	}
	header_b, err := reader_b.Read()
	if header_b[0] != "id" {
		fmt.Println("Bad Header", err)
		return 0.0
	}
	
	total, total2, cnt := 0,0,0
	for {
		record_a, err := reader_a.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error ReaderA:", err)
			return 0.0
		}
		record_b, err := reader_b.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error ReaderB:", err)
			return 0.0
		}
		
		if record_a[0] != record_b[0] {
			fmt.Printf("LineID mismatch %s != %s\n", record_a[0],record_b[0])
			break
		}
		
		dist := LevenshteinDistance(record_a[1], record_b[1])
		total += dist
		total2 += dist*dist
		cnt ++
	}
	if 0==cnt {
		return 0.0
	}
	
	count := float64(cnt)
	mean := float64(total) / count
  sd   := float64(total2)/ count - mean*mean
	conf := math.Sqrt(sd/count)

	fmt.Printf("Count : %d,  av: %7.5f (%7.5f,%7.5f)\n", 
		cnt, mean,
		mean-conf, mean+conf )
	return float32(mean)
}

//const currently_running_version int = 1000  // Used for submissions 1,2,3
const currently_running_version int = 1001  //

func main() {
	cmd := flag.String("cmd", "", "Required : {size}")
	cmd_type := flag.String("type", "", "size:{vocab|bigrams}")

	file_save := flag.String("save", "", "filename")
	file_load := flag.String("load", "", "filename")

	//delta := flag.Int("delta", 0, "Number of steps between start and end")
	seed := flag.Int64("seed", 1, "Random seed to use")
  
	skip := flag.Int("skip", 0, "Debugging aid")
	submit := flag.Int("submit", 0, "Build the submissions file too")
	search := flag.Int("search", 0, "Number of search iterations")

	hyper_str := flag.String("hyper", "0,0,0,0,0,0,0,0,0,0", "integer,comma-separated hyperparameters (up to 10)")

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
  fname_truth      := "../data/1-holdout/truth.txt"     // This is the perfect comparison file for the valid.txt test set
  
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
			for _, p := range test_pairs {
				test_vocab[p.Key] = 0
			}

			// The fill it with train word freqs (where applicable)
			for _, p := range train_pairs {
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
    
		/// ./billion -cmd=size -type=sv -save=0-sv.csv
		if *cmd_type == "sv" {
      vocab := Vocab{}
      
			// Read in the vocab for test file
			vocab.ReadTestNGram(fname_test, 1)
			// Read in the vocab for holdout file (additional, for validation)
			vocab.ReadTestNGram(fname_heldout, 1)

      pl := vocab.MakeSortedPairList()
      
      // Now, go through the training set, building up a picture of '' or 'something' for each found bigram
      sv := get_train_splittervocab(fname_train, pl)
      
      if len(*file_save)>0 {
        sv.Save(*file_save)
      }
    }
	}

	if *cmd == "validate" {
		/// ./billion -cmd=validate -type=bigrams -load=0-bigrams.csv -save=.bigram_02.csv -hyper=1,2,3,5,6 -search=0
		/// ./billion -cmd=validate -type=sv -load=0-sv.csv -save=.bigram_05.csv -hyper=1,2,3,5,6 -search=0 -skip=1555
		if *cmd_type == "bigrams" || *cmd_type == "sv" {
      vocab := Vocab{}
      vocab.Load("0-vocab.csv") // Hard coded for now

			hyper := get_hyper(*hyper_str)
			
			hyper_best := make([]int, 10)
			copy(hyper_best, hyper)  // dest, src
			hyper_best_score := float32(-1.0) // Initialisation

      splitter := Splitter{} // empty
      if *cmd_type == "bigrams" {
        splitter.Load(*file_load)
      }
      sv := SplitterVocab{}  // empty
      if *cmd_type == "sv" {
        sv.Load(*file_load)
        fmt.Printf("Loaded SplitterVocab : size=%d\n", len(sv))
      }

			for {
        if *cmd_type == "bigrams" {
          splitter.CreateSubmission(fname_validation, "1-valid"+*file_save, &vocab, *skip, hyper)
          if *submit>0 {
            splitter.CreateSubmission(fname_test, "1-test"+*file_save, &vocab, *skip, hyper)
          }
        }
        if *cmd_type == "sv" {
          break
          /*
          sv.CreateSubmission(fname_validation, "1-valid"+*file_save, &vocab, *skip, hyper)
          if *submit>0 {
            sv.CreateSubmission(fname_test, "1-test"+*file_save, &vocab, *skip, hyper)
          }
          */
        }
        
        if *skip>0 {
          break // all done
        }
				score := get_validation_score(fname_truth, "1-valid"+*file_save)
				hyper_asstrings := make([]string, len(hyper))
				for i, v:= range hyper {
					hyper_asstrings[i] = fmt.Sprintf("%+d", v);
				}
				fmt.Printf("  Levenshtein score : %7.5f :: hyper=%s\n", score, strings.Join(hyper_asstrings, ","))
				
				// If this is the best so far, re-base the best
				if hyper_best_score<0 || score<hyper_best_score { // new best score
					copy(hyper_best, hyper)
					hyper_best_score = score
					fmt.Printf("    New Best score established\n")
				}
				
				*search--
				if *search<0 {
					break
				}

				// Create next hyper as a variation of hyper_best ...
				for i,v := range hyper_best {
					hyper[i] = v + (rand.Intn(5)-2) // {-2,-1,0,1,2}
				}
				// And loop around...
			}
    }
		/// ./billion -cmd=validate -type=score -load=.bigram_02.csv
		if *cmd_type == "score" {
			score := get_validation_score(fname_truth, "1-valid"+*file_load)
			fmt.Printf("  Levenshtein score : %7.5f\n", score)
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
