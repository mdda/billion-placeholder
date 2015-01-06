import sys

import billion

regularize=billion.util.regularize

class Gaps:
	def __init__(vocab, small_limit):
		self.vocab_index = billion.util.load_vocab(args.vocab)
		self.small_limit = small_limit

		if False:
		  for w in ['the', 'computer', 'investor', 'xNONEXISTENTx', ]:
			print w, ' -> ', self.vocab_index[w]

		if False:
		  self.missing_ones = ['the', 'and', 'of', 'to', 'for', 'a', 'an', 'on', 'in', 'at', 'by', 'from', ]
		  if len(self.missing_ones)>30:
			print "missing_ones list too long to pack into integer"
			exit

		  print(["%d=%s" % (self.vocab_index[w], w) for w in self.missing_ones ], "\n")

	def generate_training(self, line):
	  words = regularize(line)
	  vocab_indices = [ self.vocab_index[w] for w in words ]
	  
	  for i in range(len(words)-2):
		if vocab_indices[i] is None or vocab_indices[i+1] is None or vocab_indices[i+2] is None:
		  continue
		  
		# Middle word missing
		x = [ vocab_indices[i], vocab_indices[i+2] ]
			
		#missing_word = words[i+1]
		#a=0
		#if missing_word in missing_ones:
		#  a = self.missing_ones.index(missing_word)
		
		# Pick out small words for 'easy' identification
		missing = vocab_indices[i+1]
		a = missing if missing<small_limit else 0
		
		yield (x, [1,a])
		
	  for i in range(len(words)-1):
		if vocab_indices[i] is None or vocab_indices[i+1] is None:
		  continue
		# Word not missing
		x = [ vocab_indices[i], vocab_indices[i+1] ]
		
		yield (x, [0,0])
	  
	  #print(' '.join(out))
	  #if l>5: break
		
"""
for l, line in enumerate(inputfile):  
  if 0 == l % 10000:
    billion.util.print_thousands("Line # ", l)
    # Skip this iter - since it's going into our 'holdout set'
    continue
    
"""
