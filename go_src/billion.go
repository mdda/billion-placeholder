package main

// GOPATH=`pwd` go build billion.go  && ./billion

import (
	"fmt"
	"time"
	"math/rand"
	"flag"
)


const currently_running_version int = 1000

func main() {
	cmd:= flag.String("cmd", "", "Required : {}")
	cmd_type:= flag.String("type", "", "create:{fake_training_data|training_set_transitions|synthetic_transitions}, db:{test|insert_problems}, visualize:{data|ga}, submit:{kaggle|fakescore}")
	
	//delta := flag.Int("delta", 0, "Number of steps between start and end")
	seed  := flag.Int64("seed", 1, "Random seed to use")

	//id := flag.Int("id", 0, "Specific id to examine")
	//training_only := flag.Bool("training", false, "Act on training set (default=false, i.e. test set)")

	//count := flag.Int("count", 0, "Number of ids to process")

	
	flag.Parse()
	//fmt.Printf("CMD = %s\n", *cmd)
	
	//rand.Seed(time.Now().UnixNano()) 
	rand.Seed(*seed)
	
	
	if *cmd=="db" {
		/// ./billion -cmd=db -type=test
		if *cmd_type=="test" {
			test_open_db()
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

