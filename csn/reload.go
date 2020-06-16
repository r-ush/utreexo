package csn

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/mit-dci/utreexo/accumulator"
	"github.com/mit-dci/utreexo/util"
)

// restorePollard restores the pollard from disk to memory.
// If starting anew, it just returns a empty pollard.
func restorePollard() (p accumulator.Pollard, err error) {
	// Restore Pollard
	pollardFile, err := os.OpenFile(
		util.PollardFilePath, os.O_RDWR, 0600)
	if err != nil {
		return
	}
	err = p.RestorePollard(pollardFile)
	if err != nil {
		fmt.Printf("restore error\n")
		return
	}

	return
}

// restorePollardHeight restores the current height that pollard is synced to
// Not to be confused with the height variable for genproofs
func restorePollardHeight() (height int32, err error) {

	var pHeightFile *os.File
	// Restore height
	pHeightFile, err = os.OpenFile(
		util.PollardHeightFilePath, os.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}

	err = binary.Read(pHeightFile, binary.BigEndian, &height)
	if err != nil {
		return 0, err
	}

	return
}

// saveIBDsimData saves the state of ibdsim so that when the
// user restarts, they'll be able to resume.
// Saves height for ibdsim and pollard itself
func saveIBDsimData(height int32, p accumulator.Pollard) error {

	var fileWait sync.WaitGroup

	fileWait.Add(1)
	go func() error {
		pHeightFile, err := os.OpenFile(
			util.PollardHeightFilePath, os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		// write to the heightfile
		err = binary.Write(pHeightFile, binary.BigEndian, height)
		if err != nil {
			return err
		}
		fileWait.Done()
		return nil
	}()

	fileWait.Add(1)
	go func() error {

		pollardFile, err := os.OpenFile(
			util.PollardFilePath, os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		err = p.WritePollard(pollardFile)
		if err != nil {
			return err
		}
		fileWait.Done()
		return nil
	}()

	fileWait.Wait()
	return nil
}
