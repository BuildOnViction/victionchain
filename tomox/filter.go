package tomox

import (
	"fmt"
	"sync"
	mrand "math/rand"
	crand "crypto/rand"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// Filter represents a TomoX message filter
type Filter struct {
	Topics     [][]byte          // Topics to filter messages with
	AllowP2P   bool              // Indicates whether this filter is interested in direct peer-to-peer messages
	id         string            // unique identifier

	Messages map[common.Hash]*ReceivedMessage
	mutex    sync.RWMutex
}

// Filters represents a collection of filters
type Filters struct {
	watchers map[string]*Filter

	topicMatcher     map[TopicType]map[*Filter]struct{} // map a topic to the filters that are interested in being notified when a message matches that topic
	allTopicsMatcher map[*Filter]struct{}               // list all the filters that will be notified of a new message, no matter what its topic is

	tomox *TomoX
	mutex   sync.RWMutex
}

// NewFilters returns a newly created filter collection
func NewFilters(t *TomoX) *Filters {
	return &Filters{
		watchers:         make(map[string]*Filter),
		topicMatcher:     make(map[TopicType]map[*Filter]struct{}),
		allTopicsMatcher: make(map[*Filter]struct{}),
		tomox:          t,
	}
}

// Install will add a new filter to the filter collection
func (fs *Filters) Install(watcher *Filter) (string, error) {
	if watcher.Messages == nil {
		watcher.Messages = make(map[common.Hash]*ReceivedMessage)
	}

	id, err := GenerateRandomID()
	if err != nil {
		return "", err
	}

	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.watchers[id] != nil {
		return "", fmt.Errorf("failed to generate unique ID")
	}

	watcher.id = id
	fs.watchers[id] = watcher
	fs.addTopicMatcher(watcher)
	return id, err
}

// Uninstall will remove a filter whose id has been specified from
// the filter collection
func (fs *Filters) Uninstall(id string) bool {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	if fs.watchers[id] != nil {
		fs.removeFromTopicMatchers(fs.watchers[id])
		delete(fs.watchers, id)
		return true
	}
	return false
}

// addTopicMatcher adds a filter to the topic matchers.
// If the filter's Topics array is empty, it will be tried on every topic.
// Otherwise, it will be tried on the topics specified.
func (fs *Filters) addTopicMatcher(watcher *Filter) {
	if len(watcher.Topics) == 0 {
		fs.allTopicsMatcher[watcher] = struct{}{}
	} else {
		for _, t := range watcher.Topics {
			topic := BytesToTopic(t)
			if fs.topicMatcher[topic] == nil {
				fs.topicMatcher[topic] = make(map[*Filter]struct{})
			}
			fs.topicMatcher[topic][watcher] = struct{}{}
		}
	}
}

// removeFromTopicMatchers removes a filter from the topic matchers
func (fs *Filters) removeFromTopicMatchers(watcher *Filter) {
	delete(fs.allTopicsMatcher, watcher)
	for _, topic := range watcher.Topics {
		delete(fs.topicMatcher[BytesToTopic(topic)], watcher)
	}
}

// getWatchersByTopic returns a slice containing the filters that
// match a specific topic
func (fs *Filters) getWatchersByTopic(topic TopicType) []*Filter {
	res := make([]*Filter, 0, len(fs.allTopicsMatcher))
	for watcher := range fs.allTopicsMatcher {
		res = append(res, watcher)
	}
	for watcher := range fs.topicMatcher[topic] {
		res = append(res, watcher)
	}
	return res
}

// Get returns a filter from the collection with a specific ID
func (fs *Filters) Get(id string) *Filter {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return fs.watchers[id]
}

// NotifyWatchers notifies any filter that has declared interest
// for the envelope's topic.
func (fs *Filters) NotifyWatchers(env *Envelope, p2pMessage bool) {
	var msg *ReceivedMessage

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	candidates := fs.getWatchersByTopic(env.Topic)
	for _, watcher := range candidates {
		if p2pMessage && !watcher.AllowP2P {
			log.Trace(fmt.Sprintf("msg [%x], filter [%s]: p2p messages are not allowed", env.Hash(), watcher.id))
			continue
		}

		var match bool
		if msg != nil {
			match = watcher.MatchMessage(msg)
		} else {
			match = watcher.MatchEnvelope(env)
			if match {
				msg = env.Open(watcher)
				if msg == nil {
					log.Trace("processing message: failed to open", "message", env.Hash().Hex(), "filter", watcher.id)
				}
			} else {
				log.Trace("processing message: does not match", "message", env.Hash().Hex(), "filter", watcher.id)
			}
		}

		if match && msg != nil {
			log.Trace("processing message: decrypted", "hash", env.Hash().Hex())
				watcher.Trigger(msg)
		}
	}
}

// Trigger adds a yet-unknown message to the filter's list of
// received messages.
func (f *Filter) Trigger(msg *ReceivedMessage) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, exist := f.Messages[msg.EnvelopeHash]; !exist {
		f.Messages[msg.EnvelopeHash] = msg
	}
}

// Retrieve will return the list of all received messages associated
// to a filter.
func (f *Filter) Retrieve() (all []*ReceivedMessage) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	all = make([]*ReceivedMessage, 0, len(f.Messages))
	for _, msg := range f.Messages {
		all = append(all, msg)
	}

	f.Messages = make(map[common.Hash]*ReceivedMessage) // delete old messages
	return all
}

// MatchMessage checks if the filter matches an already decrypted
// message (i.e. a Message that has already been handled by
// MatchEnvelope when checked by a previous filter).
// Topics are not checked here, since this is done by topic matchers.
func (f *Filter) MatchMessage(msg *ReceivedMessage) bool {
	return true
}

// MatchEnvelope checks if it's worth decrypting the message. If
// it returns `true`, client code is expected to attempt decrypting
// the message and subsequently call MatchMessage.
// Topics are not checked here, since this is done by topic matchers.
func (f *Filter) MatchEnvelope(envelope *Envelope) bool {
	return true
}

// GenerateRandomID generates a random string, which is then returned to be used as a key id
func GenerateRandomID() (id string, err error) {
	buf, err := generateSecureRandomData(keyIDSize)
	if err != nil {
		return "", err
	}
	if !validateDataIntegrity(buf, keyIDSize) {
		return "", fmt.Errorf("error in generateRandomID: crypto/rand failed to generate random data")
	}
	id = common.Bytes2Hex(buf)
	return id, err
}

// generateSecureRandomData generates random data where extra security is required.
// The purpose of this function is to prevent some bugs in software or in hardware
// from delivering not-very-random data. This is especially useful for AES nonce,
// where true randomness does not really matter, but it is very important to have
// a unique nonce for every message.
func generateSecureRandomData(length int) ([]byte, error) {
	x := make([]byte, length)
	y := make([]byte, length)
	res := make([]byte, length)

	_, err := crand.Read(x)
	if err != nil {
		return nil, err
	} else if !validateDataIntegrity(x, length) {
		return nil, errors.New("crypto/rand failed to generate secure random data")
	}
	_, err = mrand.Read(y)
	if err != nil {
		return nil, err
	} else if !validateDataIntegrity(y, length) {
		return nil, errors.New("math/rand failed to generate secure random data")
	}
	for i := 0; i < length; i++ {
		res[i] = x[i] ^ y[i]
	}
	if !validateDataIntegrity(res, length) {
		return nil, errors.New("failed to generate secure random data")
	}
	return res, nil
}
