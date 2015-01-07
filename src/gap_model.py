#! python

#from __future__ import print_function

import billion
import sys

import hickle 

#import cPickle as pickle
#import gzip
import itertools
#import urllib

import numpy as np

import lasagne

import theano
import theano.tensor as T

NUM_EPOCHS = 10

import argparse

parser = argparse.ArgumentParser(description='Converts corpus to "gaps training data"')
parser.add_argument('-m','--mode',  help='{train|test}', required=True)

parser.add_argument(     '--vocab', help='Vocab file name', required=True)
parser.add_argument(     '--vectors', help='Word Embedding Vectors file name', required=True)
parser.add_argument(     '--small', help='Number of "small words" to capture', required=False, default=32, type=int)

parser.add_argument(     '--train', help='Training text file name', )
parser.add_argument(     '--valid', help='Validation text file name', )

parser.add_argument(     '--test',  help='Test text file name', )
parser.add_argument(     '--output',  help='Submission file name to write', )

parser.add_argument(     '--load',  help='File to load model from', )
parser.add_argument(     '--save',  help='File to save model to', )

args = parser.parse_args()


# The vectors will be stored on GPU all the time
# Blocks of training data will be 'mini-batched' and also paged in
# in units of 'BULK_SIZE'
BULK_SIZE = 1000*1000  # Training Records to read in blocks off disk

## Hmm : Maybe the examples should be created dynamically by generators
## since the files (for ALL in particular) will get stupidly large

# These are the mini-batches over which SGD takes place
MINIBATCH_SIZE = 500

NUM_HIDDEN_UNITS = 240

# This will use ADAGRAD, rather than momentum, etc


def load_language(vocab, vectors, small):
    print("Vectors file = %s" % (vectors,))
    d = hickle.load(vectors)
    vectors = theano.shared(lasagne.utils.floatX(d['vectors']))
    
    return dict(
      vectors = vectors,
      gaps = billion.gaps.Gaps(vocab, small),
    )


def create_training_set(train, gaps):  # BULK_SIZE
    # These are just 'sized' - will be loaded dynamically due to GPU size constraints
    X = np.zeroes( (2, BULK_SIZE), dtype='int32')
    Y = np.zeroes( (1, BULK_SIZE), dtype='int8' )
    
    return dict(
        X = theano.shared(X),
        Y = T.cast(theano.shared(Y), dtype='int8'),
        
        num_examples=X.shape[0],
        
        input_dim = X.shape[1],
        output_dim = 1 + gaps.small_limit,
	)

def load_validation_set(valid, gaps):  # Will load all
    def g(f, comment):
        inputfile = open(f)
        for l, line in enumerate(inputfile):  
            if 0 == l % 10000:
                billion.util.print_thousands(comment+" Line # ", l)
            for p in gaps.generate_training(line):
                yield p
        inputfile.close()

    arr = [ p for p in g(valid, "Validation") ]
    
    X = np.array([x for (x,y) in arr], dtype=np.int)
    Y = np.array([y for (x,y) in arr], dtype=np.int8)
    
    print X[0:60]
    print Y[0:60]
    
    #print X.shape[0]

    return dict(
        X = theano.shared(X),
        Y = T.cast(theano.shared(Y), dtype='int8'),
        
        num_examples=X.shape[0],
        
        input_dim = X.shape[1],
        output_dim = 1 + gaps.small_limit,
	)


def build_model(input_dim, output_dim,
                batch_size=MINIBATCH_SIZE, num_hidden_units=NUM_HIDDEN_UNITS):

    l_in = lasagne.layers.InputLayer(
        shape=(batch_size, input_dim),
	)
    l_hidden1 = lasagne.layers.DenseLayer(
        l_in,
        num_units=num_hidden_units,
		nonlinearity=lasagne.nonlinearities.rectify,
	)
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
        l_hidden2_dropout,
        num_units=output_dim,
        nonlinearity=lasagne.nonlinearities.softmax,
	)
    return l_out


def create_iter_functions(dataset, output_layer,
                          X_tensor_type=T.matrix,
                          batch_size=MINIBATCH_SIZE
                         ):
    batch_index = T.iscalar('batch_index')
    X_batch = X_tensor_type('x')
    y_batch = T.ivector('y')
    batch_slice = slice(
        batch_index * batch_size, (batch_index + 1) * batch_size
    )

    def loss(output):
        return -T.mean(T.log(output)[T.arange(y_batch.shape[0]), y_batch])

    loss_train = loss(output_layer.get_output(X_batch))
    loss_eval = loss(output_layer.get_output(X_batch, deterministic=True))

    pred = T.argmax(
        output_layer.get_output(X_batch, deterministic=True), axis=1
    )
    accuracy = T.mean(T.eq(pred, y_batch))

    all_params = lasagne.layers.get_all_params(output_layer)
    updates = lasagne.updates.nesterov_momentum(
        loss_train, all_params, learning_rate, momentum
    )

    iter_train = theano.function(
        [batch_index], loss_train,
        updates=updates,
        givens={
            X_batch: dataset['X_train'][batch_slice],
            y_batch: dataset['y_train'][batch_slice],
		},
	)

    iter_valid = theano.function(
        [batch_index], [loss_eval, accuracy],
        givens={
            X_batch: dataset['X_valid'][batch_slice],
            y_batch: dataset['y_valid'][batch_slice],
		},
	)

    iter_test = theano.function(
        [batch_index], [loss_eval, accuracy],
        givens={
            X_batch: dataset['X_test'][batch_slice],
            y_batch: dataset['y_test'][batch_slice],
		},
	)

    return dict(
        train=iter_train,
        valid=iter_valid,
        test=iter_test,
	)


def train(iter_funcs, dataset, batch_size=MINIBATCH_SIZE):
    num_batches_train = dataset['num_examples_train'] // batch_size
    num_batches_valid = dataset['num_examples_valid'] // batch_size
    num_batches_test = dataset['num_examples_test'] // batch_size

    for epoch in itertools.count(1):
        batch_train_losses = []
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


def main(num_epochs=NUM_EPOCHS):
    dataset = load_data()
    output_layer = build_model(
        input_dim=dataset['input_dim'],
        output_dim=dataset['output_dim'],
    )
    iter_funcs = create_iter_functions(dataset, output_layer)

    print("Starting training...")
    for epoch in train(iter_funcs, dataset):
        print("Epoch %d of %d" % (epoch['number'], num_epochs))
        print("  training loss:\t\t%.6f" % epoch['train_loss'])
        print("  validation loss:\t\t%.6f" % epoch['valid_loss'])
        print("  validation accuracy:\t\t%.2f %%" %
              (epoch['valid_accuracy'] * 100))

        if epoch['number'] >= num_epochs:
            break

    return output_layer


if __name__ == '__main__':
    if args.mode != 'train' and args.mode != 'test':
        args.print_help()
        exit(1)
        
    language = load_language(args.vocab, args.vectors, args.small)
    
    if args.mode == 'train':
        validation = load_validation_set(args.valid, language['gaps'])
        #main()
    if args.mode == 'test':
        pass
