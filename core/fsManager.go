package core

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/mdlayher/wavepipe/core/models"

	"github.com/mdlayher/goset"
	"github.com/wtolson/go-taglib"
)

// validSet is a set of valid file extensions which we should scan as media, as they are the ones
// which TagLib is capable of reading
var validSet = set.New(".ape", ".flac", ".m4a", ".mp3", ".mpc", ".ogg", ".wma", ".wv")

// fsManager handles fsWalker processes, and communicates back and forth with the manager goroutine
func fsManager(mediaFolder string, fsKillChan chan struct{}) {
	log.Println("fs: starting...")

	// Trigger a filesystem walk, which can be halted via channel
	walkCancelChan := make(chan struct{})
	errChan := fsWalker(mediaFolder, walkCancelChan)

	// Trigger events via channel
	for {
		select {
		// Stop filesystem manager
		case <-fsKillChan:
			// Halt any in-progress walks
			walkCancelChan <- struct{}{}

			// Inform manager that shutdown is complete
			log.Println("fs: stopped!")
			fsKillChan <- struct{}{}
			return
		// Error return channel
		case err := <-errChan:
			// Check if error occurred
			if err == nil {
				break
			}

			// Report walk errors
			log.Println(err)
		}
	}
}

// fsWalker scans for media files in a specified path, and queues them up for inclusion
// in the wavepipe database
func fsWalker(mediaFolder string, walkCancelChan chan struct{}) (chan error) {
	// Return errors on channel
	errChan := make(chan error)

	// Halt walk if needed
	var mutex sync.RWMutex
	haltWalk := false
	go func() {
		// Wait for signal
		<-walkCancelChan

		// Halt!
		mutex.Lock()
		haltWalk = true
		mutex.Unlock()
	}()

	// Invoke walker goroutine
	go func() {
		// Keep sets of unique artists, albums, and songs encountered
		artistStringSet := set.New()
		albumStringSet := set.New()
		artistSet := set.New()
		albumSet := set.New()
		songSet := set.New()

		// Invoke a recursive file walk on the given media folder, passing closure variables into
		// walkFunc to enable additional functionality
		err := filepath.Walk(mediaFolder, func(currPath string, info os.FileInfo, err error) error {
			// Stop walking immediately if needed
			mutex.RLock()
			if haltWalk {
				return errors.New("walk: halted by channel")
			}
			mutex.RUnlock()

			// Make sure path is actually valid
			if info == nil {
				return errors.New("walk: invalid path: " + currPath)
			}

			// Ignore directories for now
			if info.IsDir() {
				return nil
			}

			// Check for a valid media extension
			if !validSet.Has(path.Ext(currPath)) {
				return nil
			}

			// Attempt to scan media file with taglib
			file, err := taglib.Read(currPath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Generate a song model from the file
			// TODO: insert song into database, and get ID
			song, err := models.SongFromFile(file)
			song.ID = int64(songSet.Size()) + 1
			if err != nil {
				return err
			}

			// Check for new artist
			if artistStringSet.Add(song.Artist) {
				// Generate the artist model from this song's metadata
				// TODO: insert artist into database, and get ID
				artist := models.ArtistFromSong(song)
				artist.ID = int64(artistSet.Size()) + 1

				// Add artist to set
				log.Printf("New artist: [%02d] %s", artist.ID, artist.Title)
				artistSet.Add(artist)
			}

			// Check for new artist/album combination
			if albumStringSet.Add(song.Artist + "-" + song.Album) {
				// Generate the album model from this song's metadata
				// TODO: insert album into database, and get ID, as well as artist ID
				album := models.AlbumFromSong(song)
				album.ArtistID = int64(artistSet.Size())
				album.ID = int64(albumSet.Size()) + 1

				// Add album to set
				log.Printf("New album: [%02d] %s- %s", album.ID, album.Artist, album.Title)
				albumSet.Add(album)
			}

			// Check for new song (struct, no need to worry about name overlap)
			if songSet.Add(song) {
				log.Printf("Song: [%02d] %s - %s - %s", song.ID, song.Artist, song.Album, song.Title)
			}

			return nil
		})

		// Check for filesystem walk errors
		if err != nil {
			errChan <- err
		}

		log.Println(artistSet)
		log.Println(albumSet)
		log.Println(songSet)

		// No errors
		errChan <- nil
	}()

	// Return communication channel
	return errChan
}