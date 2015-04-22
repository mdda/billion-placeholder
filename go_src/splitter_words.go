package main

import (
	"fmt"
)

import (
	"os"
	"encoding/csv"
	"io"
	"strconv"
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
    sv[pl[i].Key] = SVAtom{ Vocab{}, Vocab{} }
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
    //break
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return sv
	}

	return sv
}

func (self *Vocab) to_single_string() string {
  pl := sortMapByValue(self)
  
  s := []string{}
  for _, p := range pl {
    new_key := strings.Replace(strings.Replace(p.Key, ",", "#COMMA#",-1), ":","#COLON",-1)
    s = append(s, fmt.Sprintf("%s:%d", new_key,p.Value))
  }
  return "{"+strings.Join(s, ",")+"}"
}

func to_vocab(s string) Vocab {
  vocab := Vocab{}
  
  if s[0] == '{' && s[len(s)-1] == '}' {  // Strip off the surrounding {}
    s = s[1:(len(s)-1)]
  }
  //fmt.Printf("SplitterVocab='%s'\n", s)
  if len(s)==0 {
    return vocab // Nothing
  }
  for _,pair := range strings.Split(s, ",") {
    piece := strings.Split(pair, ":")
    k :=  piece[0]
    v, err := strconv.Atoi(piece[1])
    if err != nil  {
      fmt.Printf("SplitterVocab.pair='%s'\n", pair)
      continue
    }
    new_key := strings.Replace(strings.Replace(k, "#COMMA#",",", -1), "#COLON",":", -1)
    vocab[new_key] = v
    if v>65*1000 {
      fmt.Printf("size=%8d for SplitterVocab[...][%s]\n", v, new_key)
    }
  }
  return vocab
}

func (sv SplitterVocab) Save(filename string, pl PairList) {  
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return 
	}
	defer file.Close()
	writer := csv.NewWriter(file)
  
	//for w, sa := range sv {
	for _,p := range pl {
    w := p.Key 
    sa := sv[w]
    writer.Write( []string{w, sa.Together.to_single_string(), sa.Separate.to_single_string()} )
  }
  writer.Flush()
}

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
    together := record[1]
    separate := record[2]
    
    self[word] = SVAtom{ to_vocab(together),to_vocab(separate) }
  }
}

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
	writer := bufio.NewWriter(file_out)
  
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

    id  := record[0] // Don't really care about content
		txt := record[1]

    best_i := -1
    best_v := 20.0 + float64(hyper[0])  // Must beat this to register at all
    if best_v<0 {
			best_v = -best_v
		}
    if currently_running_version>1000 {
      //best_v = 0
    }

		words := strings.Split(txt, " ")
    //words[0] = strings.ToLower(words[0])
    for i := 0; i < len(words)-1; i++ {
      word := words[i]
      
      // How frequent the words are stand-alone
      v0 := float64((*vocab)[word])
      v1 := float64((*vocab)[words[i+1]])
      
      // Stats for the size of the expected word lists
      sa := self[word]
      together := float64(len(sa.Together))
      separate := float64(len(sa.Separate))
      tot := together+separate
      if tot<1.0 {
         tot=1.0
      }
      
      max_prop := 0.0
      if v0>0.0 && v1>0.0 {
        max_prop = 100.0 * tot/v0
        if max_prop < 100.0 * tot/v1 {
          max_prop = 100.0 * tot/v1
        }
      }
      
      v := (max_prop * separate)/tot
      
      if separate<( 10.0 + float64(hyper[1]) ) { // Evidence either way is very poor
        v=0.0
      }
      if 100.0*separate/tot<( 90.0 + float64(hyper[2]) ) { // poor percentage suggesting split
        v=0.0
      }
      
      if v>best_v {
        best_i=i
        best_v=v
      }
      
      if(skip_check>0) {
        fmt.Printf("%20s - %20s :: [%7d,%7d] :: %7d %3d%% :: vocab:(%8d,%8d)=(%3d%%,%3d%%) -> %3d%%\n", words[i], words[i+1], 
          int(together), int(separate), int(together+separate), int(100.0*separate/tot),
          int(v0), int(v1), int(100.0*tot/v0), int(100.0*tot/v1),
          int(v) )
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
