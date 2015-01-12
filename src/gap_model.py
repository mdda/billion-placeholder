#! python
from __future__ import print_function

import billion

import numpy as np

import theano
import theano.tensor as T

import lasagne

import argparse
import itertools

import hickle 

import warnings
warnings.simplefilter("error", RuntimeWarning)

parser = argparse.ArgumentParser(description='Converts corpus to "gaps training data"')
parser.add_argument('-m','--mode',  help='{train|test}', required=True)

parser.add_argument(     '--vocab', help='Vocab file name', required=True)
parser.add_argument(     '--vectors', help='Word Embedding Vectors file name', required=True)
parser.add_argument(     '--small', help='Number of "small words" to capture', required=False, default=32, type=int)

parser.add_argument(     '--train', help='Training text file name', )
parser.add_argument(     '--valid', help='Validation text file name', )
parser.add_argument(     '--epochs', help='Number of Epochs', required=False, default=10, type=int)

parser.add_argument(     '--test',  help='Test text file name', )
parser.add_argument(     '--output',  help='Submission file name to write', )

parser.add_argument(     '--load',  help='File to load model from', )
parser.add_argument(     '--save',  help='File to save model to', )

args = parser.parse_args()


# The vectors will be stored on GPU all the time
# Blocks of training data will be 'mini-batched' and also paged in
# in units of 'BULK_SIZE'
BULK_SIZE = 1000*1000  # Training Records to read in blocks off disk

# Memory usage = (ints for embedding index + byte for answer) * BULK_SIZE
#              = (CONTEXT_LENGTH * 4 + 1) * BULK_SIZE

## The examples are being created dynamically by generators
## since the files (for ALL in particular) will get stupidly large

# These are the mini-batches over which SGD takes place
MINIBATCH_SIZE = 500

NUM_HIDDEN_UNITS = 240

CONTEXT_LENGTH = 2

# This will use ADAGRAD, rather than momentum, etc


def load_language(vocab, vectors, small):
    print("Vectors file = %s" % (vectors,))
    d = hickle.load(vectors)
    vectors = theano.shared(lasagne.utils.floatX(d['vectors']))
    
    print("  Vectors.nbytes \t= \t", billion.util.comma_000(d['vectors'].nbytes))
    
    return dict(
      vectors = vectors,
      vocab_size = d['vectors'].shape[0],
      vector_width = d['vectors'].shape[1],
      
      gaps = billion.gaps.Gaps(vocab, small),
    )


def data_loader(filename, gaps, comment):
    inputfile = open(filename)
    for l, line in enumerate(inputfile):  
        if 0 == l % 10000:
            billion.util.print_thousands(comment+" Line # ", l)
        for p in gaps.generate_training(line):
            yield p
    billion.util.print_thousands(comment+" Line # ", l, s_after="    \n")
    inputfile.close()
    
def reset_training_set_loader(training_set, gaps):
    print("Resetting TrainingSet.loader")
    training_set['loader'] = data_loader(training_set['filename'], gaps, "Training Data")
    
def create_training_set(train, gaps):  # BULK_SIZE
    # These are just 'sized' - will be loaded dynamically due to GPU size constraints
    X = np.empty( (BULK_SIZE, CONTEXT_LENGTH), dtype=np.int32)
    Y = np.empty( (BULK_SIZE), dtype=np.int8)

    print("  Training.X.nbytes \t= \t", billion.util.comma_000(X.nbytes))
    print("  Training.Y.nbytes \t= \t", billion.util.comma_000(Y.nbytes))
    
    return dict(
        filename = train,
        #loader = data_loader(train, gaps, "Training Data"),
        
        X = theano.shared(X),
        Y = T.cast(theano.shared(Y), dtype='int8'),
        
        num_examples=X.shape[0],
        
        input_dim = X.shape[1],
        output_dim = 2 + gaps.small_limit,
	)

def load_training_set_inplace(training_set):  
    """ 
    Load in next piece of 'HUGE' dataset - returns False if insufficient data remains in file
    """
    g = training_set['loader']
    try:
        arr = [ g.next() for i in range(0, BULK_SIZE) ]
    except StopIteration:
        return False # This is a failure to load
    
    X = np.array([x for (x,y) in arr], dtype=np.int32)
    Y = np.array([y for (x,y) in arr], dtype=np.int8)
    
    #print(X[0:60])
    #print(Y[0:60])
    
    training_set['X'].set_value(X, borrow=True)
    training_set['Y'].set_value(Y, borrow=True)

    return True
    
def load_validation_set(valid, gaps):  # Will load all
    arr = [ p for p in data_loader(valid, gaps, "Validation Data") ]
    
    X = np.array([x for (x,y) in arr], dtype=np.int32)
    Y = np.array([y for (x,y) in arr], dtype=np.int8)

    print("  Valid.X.nbytes \t= \t", billion.util.comma_000(X.nbytes))
    print("  Valid.Y.nbytes \t= \t    ", billion.util.comma_000(Y.nbytes))
    
    return dict(
        X = theano.shared(X),
        Y = T.cast(theano.shared(Y), dtype='int8'),
        
        num_examples=X.shape[0],
        
        input_dim = X.shape[1],
        output_dim = 2 + gaps.small_limit,
	)


def build_model(input_dim, output_dim,
                batch_size=MINIBATCH_SIZE, 
                num_hidden_units=NUM_HIDDEN_UNITS):

    # Need to understand InputLayer structure 
    # (how does l_out keep a reference to it? = This is tracked through whole network)
    # And then need to take out [int32] and convert it into concatinated embedding vectors
    
    # input_dim = CONTEXT_LENGTH # of int32
    # processed_input_dim = CONTEXT_LENGTH * language['vector_width'] # of floatX

    l_in = lasagne.layers.InputLayer(
        shape=(batch_size, input_dim),
	)
    
    l_hidden1 = lasagne.layers.DenseLayer(
        l_in,
        num_units=num_hidden_units,
		nonlinearity=lasagne.nonlinearities.rectify,
	)
    if False:
        l_hidden1_dropout = lasagne.layers.DropoutLayer(
            l_hidden1,
            p=0.5,
        )
        l_hidden2 = lasagne.layers.DenseLayer(
            l_hidden1_dropout,
            num_units=num_hidden_units,
            nonlinearity=lasagne.nonlinearities.rectify,
        )
        l_hidden2_dropout = lasagne.layers.DropoutLayer(
            l_hidden2,
            p=0.5,
        )
    l_out = lasagne.layers.DenseLayer(
        #l_hidden2_dropout,
        l_hidden1,
        num_units=output_dim,
        nonlinearity=lasagne.nonlinearities.softmax,
	)
    
    # Perhaps this need a different output layer
    # a = NotMissing
    # b = Missing (complex + simple)
    # c-x = Missing a simple word (take shift into account)
    
    # But this (for the first runs) easy to model as a soft-max thing 
    # from 0=(nogap), 1=(complex), 2..(small_limit+2)=small-word
    
    return l_out


def create_iter_functions(dataset, output_layer,
                          X_tensor_type=T.imatrix,
                          batch_size=MINIBATCH_SIZE
                         ):
    batch_index = T.iscalar('batch_index')
    X_batch = X_tensor_type('x')
    
    # See http://stackoverflow.com/questions/25166657/index-gymnastics-inside-a-theano-function
    vectors = dataset['language']['vectors']
    X_batch_flat_vectors =  vectors[X_batch].reshape( (X_batch.shape[0], -1) )
    
    #Y_batch = T.ivector('y') 
    Y_batch = T.bvector('y') # This is smaller...
    batch_slice = slice(
        batch_index * batch_size, (batch_index + 1) * batch_size
    )

    def loss(output):
        return -T.mean(T.log(output)[T.arange(Y_batch.shape[0]), Y_batch])

    loss_train = loss(output_layer.get_output(X_batch_flat_vectors))
    loss_eval  = loss(output_layer.get_output(X_batch_flat_vectors, deterministic=True))

    pred = T.argmax(
        output_layer.get_output(X_batch_flat_vectors, deterministic=True), axis=1
    )
    accuracy = T.mean(T.eq(pred, Y_batch))

    all_params = lasagne.layers.get_all_params(output_layer)
    
    #updates = lasagne.updates.nesterov_momentum(
    #    loss_train, all_params, learning_rate, momentum
    #)
    
    #def adagrad(loss, all_params, learning_rate=1.0, epsilon=1e-6):
    updates = lasagne.updates.adagrad(
        loss_train, all_params #, learning_rate, momentum
    )

    iters={}
    
    if 'train' in dataset:
        d=dataset['train']
        iters['train'] = theano.function(
            [batch_index], loss_train,
            updates=updates,
            givens={
                X_batch: d['X'][batch_slice],
                Y_batch: d['Y'][batch_slice],
            },
        )

    if 'valid' in dataset:
        d=dataset['valid']
        iters['valid'] = theano.function(
            [batch_index], [loss_eval, accuracy],
            givens={
                X_batch: d['X'][batch_slice],
                Y_batch: d['Y'][batch_slice],
            },
        )

    if 'test' in dataset:
        d=dataset['test']
        iters['test'] = theano.function(
            [batch_index], [loss_eval, accuracy],
            givens={
                X_batch: d['X'][batch_slice],
                Y_batch: d['Y'][batch_slice],
            },
        )

    return iters

def set_up_complete_model(dataset):
    output_layer = build_model(
        CONTEXT_LENGTH * dataset['language']['vector_width'],  # input_dim
        dataset['language']['gaps'].small_limit + 2,           # output_dim
    )
    print("Creating IterFunctions...")
    iter_funcs = create_iter_functions(dataset, output_layer)
    
    return iter_funcs

def train_and_validate(iter_funcs, dataset, batch_size=MINIBATCH_SIZE):
    num_batches_train = dataset['train']['num_examples'] // batch_size
    num_batches_valid = dataset['valid']['num_examples'] // batch_size
    #num_batches_test  = dataset['test' ]['num_examples'] // batch_size

    for epoch in itertools.count(1):  # This just allows us to enumerate epoch_results
        batch_train_losses = []
        
        reset_training_set_loader(dataset['train'], dataset['language']['gaps'])
        
        while True:  # Loop for loading additional training data
            loaded = load_training_set_inplace(dataset['train'])
            print(" full = ", loaded)
            if not loaded: # There wasn't enough data for a full 'BULK' so ditch attempt
                break
            
            for b in range(num_batches_train):
                batch_train_loss = iter_funcs['train'](b)
                batch_train_losses.append(batch_train_loss)

        avg_train_loss = np.mean(batch_train_losses)

        batch_valid_losses = []
        batch_valid_accuracies = []
        for b in range(num_batches_valid):
            batch_valid_loss, batch_valid_accuracy = iter_funcs['valid'](b)
            batch_valid_losses.append(batch_valid_loss)
            batch_valid_accuracies.append(batch_valid_accuracy)

        avg_valid_loss = np.mean(batch_valid_losses)
        avg_valid_accuracy = np.mean(batch_valid_accuracies)

        yield {
            'number': epoch,
            'train_loss': avg_train_loss,
            'valid_loss': avg_valid_loss,
            'valid_accuracy': avg_valid_accuracy,
		}

def train_and_validate_all(iter_funcs, dataset, num_epochs):
    print("Starting training...")
    for epoch_results in train_and_validate(iter_funcs, dataset):
        print("Epoch %d of %d" % (epoch_results['number'], num_epochs))
        print("  training loss:\t\t%.6f" % epoch_results['train_loss'])
        print("  validation loss:\t\t%.6f" % epoch_results['valid_loss'])
        print("  validation accuracy:\t\t%.2f %%" %
              (epoch_results['valid_accuracy'] * 100))

        if epoch_results['number'] >= num_epochs:
            break


if __name__ == '__main__':
    if args.mode != 'train' and args.mode != 'test':
        args.print_help()
        exit(1)

    language = load_language(args.vocab, args.vectors, args.small)
    dataset = dict( language = language )
    
    if args.mode == 'train':
        # Training data loads progressively, since it is so large
        dataset['train'] = create_training_set(args.train, language['gaps'])  
        
        # Validation data loads immediately, since it is fairly small
        dataset['valid'] = load_validation_set(args.valid, language['gaps'])  
        
        iter_funcs = set_up_complete_model(dataset)
        train_and_validate_all(iter_funcs, dataset, num_epochs=args.epochs)
        
        
    if args.mode == 'test':
        pass
