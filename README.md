billion-placeholder
===================

*Placeholder while working on Kaggle HowTo* 

The intention was to put a working solution 
to the [Billion Word Imputation challenge](http://www.kaggle.com/c/billion-word-imputation/) 
into the pubic domain during Jan-2015 (i.e. well prior to the end of the competition).
My thinking was that many would-be competitors were holding back because
the task was so daunting, but would be encouraged to 'have a go' if there
was some known-decent starting point.

However, the real world intervened, putting the Neural Network (Theano-based)
implementation on-hold - to such an extent that the clock started to run out.

Just to put in some kind of showing, I created a wholely new approach, based
on bigrams only using the Go language, over a weekend, for fun.

This solution was (of course) much worse than would have been acheivable 
using the NN method(s) originally targetted, but had the advantage of 
being ready to submit quickly.

So, somewhat resigned to not acheiving the original goals, I just submitted 
a basic bigram approach-based submission, plus another couple that were
based on optimising a couple of hyperparameters too.  That was 
good enough for 12th out of 87.


And so it ends...
-------------------

The only submittable results were produced by code in ```<THIS-REPO>/go_src```, 
which took about a weekend to write.

There's probably enough 'juice' in the pure bigram approach (which is all
that is done in go_src) to get to #8 from current #12, 
but no futher, since there is a step-wise change between #8 and #7.

But time would be better spent (i) sleeping; (ii) working on something 
that can be submitted to a conference instead.

