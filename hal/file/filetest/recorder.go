package filetest

import (
	"os"
	pth "path"
	"strings"
	"sync"
	"time"
)

// FileRecorder records data written to files and can generate file data.
// The FileRecorder is intended for testing device file interactions.
type FileRecorder struct {
	records    []*Record            // records stores a list of interaction Record items
	responders map[string]Responder // The responders for various paths
}

// Write records the path and text sent to the FileRecorder.
// An error is returned if there is a write error returned by the responder
// registered for the provided path or os.ErrNotExist if there is no responder.
func (r *FileRecorder) Write(path, text string) error {
	path = AbsPath(path)
	responder, ok := r.responders[path]
	var err error
	if ok {
		err = responder.Write(text)
	} else {
		err = os.ErrNotExist
	}
	r.records = append(r.records, NewWriteRecord(path, text, err))
	return err
}

// Read the path read from and any text that was returned to the caller.
// The data (or error) for the read is provided by a write responder if
// registered, or an os.ErrNotExist error is returned otherwise.
func (r *FileRecorder) Read(path string) (string, error) {
	path = AbsPath(path)
	responder, ok := r.responders[path]
	var text string
	var err error
	if ok {
		text, err = responder.Read()
	} else {
		err = os.ErrNotExist
	}
	r.records = append(r.records, NewReadRecord(path, text, err))
	return text, err
}

// Records retrieves a snapshot of the records in the FileRecorder.
func (r *FileRecorder) Records() []*Record {
	list := make([]*Record, len(r.records))
	copied := copy(list, r.records)
	// We don't protect the records slice from multi-threaded mutations
	// so we are careful to return the actual snapshot items copied rather
	// than whatever is in the list which may be longer than what was copied.
	return list[:copied]
}

// Reset the records stored in the FileRecorder.
func (r *FileRecorder) Reset() {
	r.records = nil
}

// Respond adds a Responder for a particular path. Adding a responder to
// a path that already has a responder replaces the existing responder.
func (r *FileRecorder) Respond(path string, responder Responder) {
	path = AbsPath(path)
	if r.responders == nil {
		r.responders = map[string]Responder{}
	}
	r.responders[path] = responder
}

// Responder is implemented by objects that respond to reads and writes
// at a certain file path.
type Responder interface {
	// Write receives text and can respond with an error.
	Write(text string) error
	// Read responds to a read request by returning text or an error.
	Read() (string, error)
}

// StaticResponder allows any writes and returns a single static text value.
type StaticResponder struct {
	Text string // The text value returned by Read calls
}

// Write always returns no error (any/all writes are allowed).
func (s *StaticResponder) Write(text string) error {
	return nil
}

// Read returns the static Text property of the responder.
func (s *StaticResponder) Read() (string, error) {
	return s.Text, nil
}

// ListResponder supports a list of text to respond to read requests. Each
// Read() removes and returns the last item in the list.
type ListResponder struct {
	texts   []string   // texts stores the text values returned by Read calls
	readMux sync.Mutex // readMux protects against multi-threaded reads
}

// Add adds another text entry to the list of Read responses.
func (s *ListResponder) Add(text string) {
	s.readMux.Lock()
	defer s.readMux.Unlock()
	s.texts = append(s.texts, text)
}

// Write always returns no error (any/all writes are allowed).
func (s *ListResponder) Write(text string) error {
	return nil
}

// Read returns the next texts item.
func (s *ListResponder) Read() (string, error) {
	s.readMux.Lock()
	defer s.readMux.Unlock()

	if len(s.texts) == 0 {
		return "", os.ErrNotExist
	}
	head := s.texts[0]
	s.texts = s.texts[1:]
	return head, nil
}

// Record stores the information about a single FileRecorder interaction.
type Record struct {
	Stamp time.Time // Stamp is a timestamp for the record
	Err   error     // Err is any error associated with the record (if any)
	Path  string    // Path is the file path associated with the record
	Text  string    // Text is the string either read or written
	Write bool      // True if this records a write operation
}

// NewReadRecord creates a new Record object with the given path, text and error
// information.
func NewReadRecord(path, text string, err error) *Record {
	return NewRecord(false, path, text, err)
}

// NewWriteRecord creates a new Record object with the given path, text and error
// information.
func NewWriteRecord(path, text string, err error) *Record {
	return NewRecord(true, path, text, err)
}

// NewRecord creates a new Record object with the given path, text and error
// information.
func NewRecord(write bool, path, text string, err error) *Record {
	return &Record{
		Stamp: time.Now(),
		Write: write,
		Path:  path,
		Text:  text,
		Err:   err,
	}
}

// AbsPath creates a cleaned, absolute path from the given path regardless of current
// working directory (it's always relative to the root directory).
func AbsPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return pth.Clean(path)
}
