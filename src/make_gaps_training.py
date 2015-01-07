error "OBSOLETE!"

#! python
import billion
import sys

import argparse

parser = argparse.ArgumentParser(description='Converts corpus to "gaps training data"')
parser.add_argument('-i','--input', help='Input file name', required=True)
parser.add_argument(     '--vocab', help='Vocab file name', required=True)
parser.add_argument(     '--small', help='Number of "small words" to capture', required=False, default=32, type=int)
parser.add_argument('-o','--output', help='Output file name', required=True)

args = parser.parse_args()

inputfile = open(args.input)
outputfile = open(args.output, 'w')
small_limit = int(args.small)

regularize=billion.util.regularize
vocab_index = billion.util.load_vocab(args.vocab)

if False:
  for w in ['the', 'computer', 'investor', 'xNONEXISTENTx', ]:
    print w, ' -> ', vocab_index[w]

if False:
  missing_ones = ['the', 'and', 'of', 'to', 'for', 'a', 'an', 'on', 'in', 'at', 'by', 'from', ]
  if len(missing_ones)>30:
    print "missing_ones list too long to pack into integer"
    exit

  print(["%d=%s" % (vocab_index[w], w) for w in missing_ones ], "\n")

for l, line in enumerate(inputfile):  
  if 0 == l % 10000:
    billion.util.print_thousands("Line # ", l)
    # Skip this iter - since it's going into our 'holdout set'
    continue
    
  words = regularize(line)
  vocab_indices = [ vocab_index[w] for w in words ]
  
  for i in range(len(words)-2):
    if vocab_indices[i] is None or vocab_indices[i+1] is None or vocab_indices[i+2] is None:
      continue
      
    # Middle word missing
    out = [vocab_indices[i], vocab_indices[i+2], 1 ]
    
    #missing_word = words[i+1]
    #a=0
    #if missing_word in missing_ones:
    #  a = missing_ones.index(missing_word)
    
    # Pick out small words for 'easy' identification
    missing = vocab_indices[i+1]
    a = missing if missing<small_limit else 0 
    out.append(a)
    
    #print(out)
    outputfile.write("%d,%d,%d,%d\n" % (out[0], out[1], out[2], out[3]))
    
    
  for i in range(len(words)-1):
    if vocab_indices[i] is None or vocab_indices[i+1] is None:
      continue
    # Word not missing
    out = [vocab_indices[i], vocab_indices[i+1], 0 ]
    out.append(0)
    #print(out)
    outputfile.write("%d,%d,%d,%d\n" % (out[0], out[1], out[2], out[3]))
  
  #print(' '.join(out))
  #if l>5: break
  
inputfile.close()
outputfile.close()
