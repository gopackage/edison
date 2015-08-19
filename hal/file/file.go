package file

import (
	"io/ioutil"
	"log"
	"os"
)

// DeviceFile abstracts access to device file drivers - most users will
// not need to implement or deal with DeviceFile objects directly.
//
// For unit testing we can create and verify the data in and out of the pin
// files. For troubleshooting we can log the data being read and written by
// introducing a proxy file (even at runtime). And for run-of-the-mill normal
// use, we can just use a simple writer.
//
// TODO Does i2c require random access? May need to expand the interface.
type DeviceFile interface {
	Write(path, data string) error
	Read(path string) (string, error)
}

// HardwareFile provides a real device file implementation for driving
// Edison hardware.
type HardwareFile struct {
}

// Write the given data to the file located at path.
func (e *HardwareFile) Write(path, data string) error {
	return ioutil.WriteFile(path, []byte(data), 0644)
}

// Read the file located at path, and return the data found or an error.
func (e *HardwareFile) Read(path string) (string, error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// LogFile logs all the the writes to the file and optionally passes the
// data to a child DeviceFile for further handling. If no proxy is specified,
// Read() calls always return the empty string and an error of os.ErrNotExists.
//
// TODO log should go into a log provider rather than directly to the system log
// We want to support people that use alternative loggers
type LogFile struct {
	Proxy DeviceFile // Proxy is an optional file the LogFile will pass data to.
}

// Write the given data to the file located at path.
func (l *LogFile) Write(path, data string) error {
	if l.Proxy == nil {
		log.Println(">>!", path, data)
		return nil
	}
	err := l.Proxy.Write(path, data)
	if err != nil {
		log.Println(">>!", path, data, err)
	} else {
		log.Println(">>", path, data)
	}
	return err
}

// Read the file located at path, and return the data found or an error.
func (l *LogFile) Read(path string) (string, error) {
	if l.Proxy == nil {
		log.Println("<<!", path)
		return "", os.ErrNotExist
	}
	body, err := l.Proxy.Read(path)
	if err != nil {
		log.Println("<<!", path, body, err)
	} else {
		log.Println("<<", path, body)
	}
	return body, err
}

// DeviceFileStats stores various statistics on the reads and writes
// that are performed on a DeviceFile.
type DeviceFileStats struct {
	Reads       int64
	ReadErrors  int64
	Read        int64
	Writes      int64
	WriteErrors int64
	Written     int64
}

// StatsFile collects statistics on reads and writes to a device file.
//
// TODO create and track stats on a per-path basis
type StatsFile struct {
	Proxy DeviceFile       // Proxy is an optional file the LogFile will pass data to.
	stats *DeviceFileStats // The current statistics accumulated on the file
}

// Write the given data to the file located at path.
func (l *StatsFile) Write(path, data string) error {
	l.stats.Writes++
	if l.Proxy == nil {
		return nil
	}
	err := l.Proxy.Write(path, data)
	if err != nil {
		l.stats.WriteErrors++
	} else {
		l.stats.Written += int64(len(data))
	}
	return err
}

// Read the file located at path, and return the data found or an error.
func (l *StatsFile) Read(path string) (string, error) {
	l.stats.Reads++
	if l.Proxy == nil {
		return "", os.ErrNotExist
	}
	body, err := l.Proxy.Read(path)
	if err != nil {
		l.stats.ReadErrors++
	} else {
		l.stats.Read += int64(len(body))
	}
	return body, err
}
