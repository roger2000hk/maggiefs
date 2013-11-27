package client

import (
	"fmt"
	"github.com/jbooth/maggiefs/maggiefs"
	"time"
	"sync"
)

type InodeWriter struct {
	// instantiated vars
	inodeid   uint64
	leases    maggiefs.LeaseService
	names     maggiefs.NameService
	datas     maggiefs.DataService
	localDnId *uint32
	
	// state vars, these are protected by lock
	l *sync.Mutex
	currLease maggiefs.WriteLease
	currInode *maggiefs.Inode
	currBlock maggiefs.Block
	currPipeline maggiefs.BlockWriter
	leaseCreatedTime time.Time
	lastWriteTime time.Time
	closed bool
}

var (
	MAX_LEASE_HOLD=time.Duration(10*time.Second) // close/reacquire lease if held for 10 seconds (10,000 millis)
	MAX_LEASE_IDLE=time.Duration(1 * time.Second) // close lease if idle for 1 second (1,000 millis)
)

func NewInodeWriter(inodeid uint64, leases maggiefs.LeaseService, names maggiefs.NameService, datas maggiefs.DataService, localDnId *uint32) (w *InodeWriter, err error) {
	if err != nil {
		return nil, err
	}

	ret := &InodeWriter{inodeid, leases, names, datas, localDnId, new(sync.Mutex),nil,nil,maggiefs.Block{},nil,time.Time{}, time.Time{},false}
	go ret.checkLease()
	return ret, nil
}

// checks every second if we haven't done a write in more than a second
func (w *InodeWriter) checkLease() {
	for {
		// sleep 1 second in between checks
		time.Sleep(1 * time.Second)
		w.l.Lock()
		now := time.Now()
		// if lease is held and we've held it too long
		if w.currLease != nil && (now.After(w.leaseCreatedTime.Add(MAX_LEASE_HOLD)) || now.After(w.lastWriteTime.Add(MAX_LEASE_IDLE))) {
			err := w.expireLease()
			if err != nil {
				// already marked close, just bail this goroutine
				fmt.Printf("Error expiring lease for writer on inode %d!, closing..\n",w.inodeid)
				w.l.Unlock()
				return
			}
			
		}
		w.l.Unlock()
	}
}

// commits all changes and drops our lease
// expects lock held
func (w *InodeWriter) expireLease() (err error) {
	if w.currInode != nil {
		err = w.names.SetInode(w.currInode)
		if err != nil {
			w.closed = true
			return
		}
	}
	if w.currPipeline != nil {
		err = w.currPipeline.SyncAndClose()
		if err != nil {
			w.closed = true
			return
		}
	}
	if w.currLease != nil {
		err = w.currLease.Release()
		if err != nil {
			w.closed = true
			return
		}
	}
	w.currLease = nil
	w.currInode = nil
	w.currPipeline = nil
	return
}



// calls out to name service to truncate this file by repeatedly shrinking blocks
func (w *InodeWriter) Truncate(length uint64) error {
	// expire/flush any previous lease
	err := w.expireLease()
	if err != nil {
		w.closed = true
		return err
	}
	// pick up lease
	lease, err := w.leases.WriteLease(w.inodeid)
	if err != nil {
		w.closed = true
		if lease != nil { lease.Release() }
		return err
	}	
	defer lease.Release()
	return w.names.Truncate(w.inodeid, length)
}

//io.InodeWriter
func (w *InodeWriter) WriteAt(p []byte, off uint64, length uint32) (nWritten uint32, err error) {
	w.l.Lock()
	defer w.l.Unlock()
	
	
	// check lease
	if w.currLease != nil && w.currInode != nil && w.currPipeline != nil {
		blockNeeded,err := blockForPos(off,w.currInode)
		if err != nil {
			return 0,err
		}
		// if the write (or a portion) can be served using current writer
		if blockNeeded.Id == w.currBlock.Id {
			// serve from local pipeline
			posInBlock := off - w.currBlock.StartPos
		}
	}
	
	// if not, release lock and serve from new pipeline
	
	// lastly, reacquire lock and cache our new pipeline 


	// get our copy of buffer
	return length, nil
}

type blockwrite struct {
	p []byte
	posInBlock uint64
}

// gets the list of block writes
func (w *InodeWriter) blockwrites(p []byte, off uint64, length uint32) []blockwrite {
	return nil
}

func (w *InodeWriter) doBlockWrite(b blockwrite) error {
	
}

// assumes lease held and state (inode, writer) is up to date
func (w *InodeWriter) servableFromCurrWriter(off uint64, length uint32) {
	if off >= w.currBlock.StartPos {
		if (off + uint64(length)) < w.currBlock.EndPos {
			return true
		}
		// if it's the last block, it can be extended up to blocklength
		
	}
}

func (w *InodeWriter) Fsync() (err error) {
	return w.expireLease()
}

func (w *InodeWriter) Close() (err error) {
	err = w.expireLease()
	w.closed = true
	return
}



// assumes lease is held and passed in inode is up to date -- may flush inode state to nameserver if we're adding a new block
// but will not most of the time, will modify passed in inode instead
func (w *InodeWriter) doWrite(ino *maggiefs.Inode, p []byte, off uint64, length uint32) (uint32, error) {
	// check if we need to add blocks, note we support sparse files
	// so it's ok to add blocks that won't be written to right away,
	// on read they'll return 0 for all bytes in empty sections, and
	// won't read past
	if off+uint64(length) > w.inode.Length {
		err := w.addBlocksForFileWrite(off, length)
		if err != nil {
			return 0, err
		}
		//fmt.Printf("Added blocks for filewrite to ino %+v\n", inode)
	}
	// now write bytes
	nWritten := 0
	endOfWritePos := off + uint64(length) - 1
	startOfWritePos := off + uint64(nWritten)
	for _, b := range w.inode.Blocks {
		//fmt.Printf("evaluating block %+v for writeStartPos %d endofWritePos %d\n", b, startOfWritePos, endOfWritePos)
		// TODO do we need that last endOfWritePos <= b.EndPos here?
		if (b.StartPos <= startOfWritePos && b.EndPos > startOfWritePos) || (b.StartPos < endOfWritePos && endOfWritePos <= b.EndPos) {

			posInBlock := uint64(0)
			if b.StartPos < off {
				posInBlock += off - b.StartPos
			}
			//fmt.Printf("nWritten %d off %d len %d endofWritePos %d block %+v posInBlock %d\n", nWritten, off, length, endOfWritePos, b, posInBlock)
			writeLength := int(length) - nWritten
			if b.EndPos < endOfWritePos {
				writeLength = int(b.Length() - posInBlock)
			}
			//fmt.Printf("startIdx %d writeLength %d\n", nWritten, writeLength)
			startIdx := nWritten
			endIdx := startIdx + writeLength
			//fmt.Printf("Writing %d bytes to block %+v pos %d startIdx %d endIdx %d\n", b, posInBlock, startIdx, endIdx)
			err := w.datas.Write(b, p[startIdx:endIdx], posInBlock)
			if err != nil {
				return 0, err
			}
			//			fmt.Printf("Wrote %d bytes to block %+v\n", endIdx-startIdx, b)
			nWritten += writeLength
			//			fmt.Printf("Wrote %d, nWritten total %d", writeLength, nWritten)
		}
	}
	//fmt.Printf("Returning from write, nWritten %d",nWritten)
	return uint32(nWritten), nil
}

// acquires lease, then adds the blocks to the namenode,
// patching up the referenced inode to match
//
// only sets inode with nameservice IFF we had to add a block.
// in that case, we first set the inode so that we're consistent, then ask the
// namenode to add a new block of the required length
func (w *InodeWriter) addBlocksForFileWrite(off uint64, length uint32) error {
	//fmt.Printf("Adding/extending blocks for write at off %d length %d\n", off, length)
	newEndPos := off + uint64(length)
	if newEndPos > w.inode.Length {
		// if we have a last block and it's less than max length,
		// extend last block to max block length first
		//fmt.Printf("Adding/extending blocks for file write to inode %+v\n", inode)
		if w.inode.Blocks != nil && len(w.inode.Blocks) > 0 {
			idx := int(len(w.inode.Blocks) - 1)
			lastBlock := w.inode.Blocks[idx]
			if lastBlock.Length() < BLOCKLENGTH {
				extendLength := BLOCKLENGTH - lastBlock.Length()
				if lastBlock.EndPos+extendLength > off+uint64(length) {
					// only extend as much as we need to
					extendLength = off + uint64(length-1) - lastBlock.EndPos
				}
				lastBlock.EndPos = lastBlock.EndPos + extendLength
				w.inode.Blocks[idx] = lastBlock
				w.inode.Length += extendLength
				//fmt.Printf("Extended block %v on inode %v\n",lastBlock,inode)
			}
		}
		// and add new blocks as necessary
		for newEndPos > w.inode.Length {
			//fmt.Printf("New end pos %d still greater than inode length %d\n", newEndPos, inode.Length)
			// make suer
			newBlockLength := newEndPos - w.inode.Length
			if newBlockLength > BLOCKLENGTH {
				newBlockLength = BLOCKLENGTH
			}
			// write prev inode block state so we stay consistent with master while adding a block
			w.inode.Mtime = time.Now().Unix()
			err := w.names.SetInode(w.inode)
			if err != nil {
				return err
			}
			newBlock, err := w.names.AddBlock(w.inode.Inodeid, uint32(newBlockLength), w.localDnId)
			if err != nil {
				return err
			}
			// patch up local copy
			w.inode.Blocks = append(w.inode.Blocks, newBlock)
			w.inode.Length += newBlockLength
		}
	}
	return nil
}
