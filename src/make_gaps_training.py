#! python

import billion
 
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

for w in ['the', 'computer', 'investor', 'xNONEXISTENTx', ]:
  print w, ' -> ', vocab_index[w]

"""
for l, line in enumerate(inputfile):  
  words = regularize(line)
  
  outputfile.write(' '.join(words))
  outputfile.write("\n")
  
  #print(' '.join(words))
  if l>args.lines: break
"""
  
inputfile.close()
outputfile.close()
