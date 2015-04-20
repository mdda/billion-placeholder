package main

// GOPATH=`pwd` go build billion.go vocab.go levenshtein.go && ./billion

import (
	"fmt"
)

import (
	"os"
	"encoding/csv"
//	"io"
//	"strconv"
	"strings"
)

import (
	"bufio"
//	"math"
)


type SVAtom struct {
  Together Vocab
  Separate Vocab
}
type SplitterVocab map[string]SVAtom

func get_train_splittervocab(filename string, pl PairList) SplitterVocab {
	sv := SplitterVocab{}
  for i:=0; i<len(pl); i++ {
    sv[pl[i].Key] = SVAtom{}
  }

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return sv
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		txt := scanner.Text()
    words := strings.Split(txt, " ")
    words[0] = strings.ToLower(words[0])
    
    for i := 0; i < len(words)-1; i++ {
      word := words[i]
      if sa, ok := sv[word]; ok {
        sa.Together[words[i+1]]++
        sv[word] = sa
      }
    }
    for i := 0; i < len(words)-2; i++ {
      word := words[i]
      if sa, ok := sv[word]; ok {
        sa.Separate[words[i+2]]++
        sv[word] = sa
      }
    }
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return sv
	}

	return sv
}

func (v Vocab) to_single_string() string {
  s := []string{}
  for k,v := range v {
    s = append(s, fmt.Sprintf("%d:%s", k,v))
  }
  return "{"+strings.Join(s, ",")+"}"
}

func (sv SplitterVocab) Save(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return 
	}
	defer file.Close()
	writer := csv.NewWriter(file)
  
	for w, sa := range sv {
    writer.Write( []string{w, sa.Together.to_single_string(), sa.Separate.to_single_string()} )
  }
  writer.Flush()
}

/*
func (self SplitterVocab) Load(filename string) {
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
*/

/*
func (self SplitterVocab) CreateSubmission(filename_test string, filename_submit string, vocab *Vocab, skip_check int, hyper []int) {
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

  line_num := 0
	for {
    if(skip_check>0) {
      // Waste a few lines...  (3032 lines in heldout.txt.csv)
      for i:=0; i<skip_check-1; i++ {
        reader.Read()
        line_num++
      }
    }

		record, err := reader.Read()
    line_num++
    
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
    best_v := 20 + hyper[0]  // Must beat this to register at all
    if best_v<0 {
			best_v = -best_v
		}
    if currently_running_version>1000 {
      //best_v = 0
    }

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
      if currently_running_version>1000 {
        v = (max_prop * sa.Separate)/tot
      }
      
      //if (sa.Separate*100)/tot<50 {
      if sa.Separate<( 10 + hyper[1] ) { // Evidence either way is very poor
        v=0 //
      }
      if (sa.Separate*100)/tot<( 90 + hyper[2] ) { // poor percentage suggesting split
        v=0
      }
      
      if v>best_v {
        best_i=i
        best_v=v
      }
      
      if(skip_check>0) {
        fmt.Printf("%20s - %20s :: [%7d,%7d] :: %7d %3d%% :: vocab:(%8d,%8d)=(%3d%%,%3d%%) -> %3d%%\n", words[i], words[i+1], 
          sa.Together, sa.Separate, sa.Together+sa.Separate, (sa.Separate*100)/tot,
          v0, v1, (tot*100)/v0, (tot*100)/v1,
          v )
      }
    }
    
    //line := fmt.Sprintf("\"%s\",%d,%d\n", strings.Replace(w, "\"", "\"\"", -1), sa.Together, sa.Separate)
    //writer.Write( []string{id, strings.Join( words_verbatim, " ")} )
    
    words_output := strings.Split(txt, " ")
    
    if best_i>=0 { // Insert it into the list of words
      i := best_i+1
      words_output = append(words_output, "")
      copy(words_output[i+1:], words_output[i:])
      words_output[i] = "" // Insert an empty word...
      
      if(false || skip_check>0) {
        words_highlight := append(strings.Split(txt, " "), "")
        copy(words_highlight[i+1:], words_highlight[i:])
        words_highlight[i] = "***"
        fmt.Printf("%d - %s\n", line_num, strings.Join( words_highlight, " "))
      }
    }
    
    txt_out := strings.Join( words_output, " ")
    writer.WriteString( fmt.Sprintf("%s,\"%s\"\n", id, strings.Replace(txt_out, "\"", "\"\"", -1)) )
    
    if(skip_check>0) {
      break
    }
	}
  writer.Flush()
}

*/
