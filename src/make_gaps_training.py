#! python
import billion
import sys

import argparse

parser = argparse.ArgumentParser(description='Parses training file in to Corpus (for GloVe)')
parser.add_argument('-i','--input', help='Input file name', required=True)
parser.add_argument('--vocab', help='Vocab file name', required=True)
parser.add_argument('-o','--output', help='Output file name', required=True)

args = parser.parse_args()

inputfile = open(args.input)
outputfile = open(args.output, 'w')

regularize=billion.util.regularize
vocab_index = billion.util.load_vocab(args.vocab)

#for w in ['the', 'computer', 'investor', 'xNONEXISTENTx', ]:
#  print w, ' -> ', vocab_index[w]

missing_ones = ['the', 'and', 'of', 'to', 'for', 'a', 'an', 'on', 'in', 'at', 'by', 'from', ]
if len(missing_ones)>30:
  print "missing_ones list too long to pack into integer"
  exit

print(missing_ones, "\n")

for l, line in enumerate(inputfile):  
  if 0 == l % 10000:
    print '\x1b[0G', 'Line : ', l, # Nice over-writing (no newline)
    sys.stdout.flush()

    
  words = regularize(line)
  vocab_indices = [ vocab_index[w] for w in words ]
  
  for i in range(len(words)-2):
    if vocab_indices[i] is None or vocab_indices[i+1] is None or vocab_indices[i+2] is None:
      continue
    # Middle word missing
    out = [vocab_indices[i], vocab_indices[i+2], 1 ]
    missing = words[i+1]
    if missing in missing_ones:
      a = 1 << missing_ones.index(missing)
      out.append(a)
    else:
      out.append(0)
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
