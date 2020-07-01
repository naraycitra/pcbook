package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// ImageStore is an interface to store laptop images
type ImageStore interface {
	// Save saves a new laptop image to the store
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore stores images on disk and its info on memory
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

// ImageInfo containts information of the Laptop Image
type ImageInfo struct {
	LaptopID string
	Type     string
	Path     string
}

// NewDiskImageStore return a new  DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

// Save save a new laptop image to the store
func (store *DiskImageStore) Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot create image id: %v", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %v", err)
	}

	_, err = imageData.WriteTo(file)

	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %v", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		LaptopID: laptopID,
		Path:     imagePath,
		Type:     imageType,
	}

	return imageID.String(), nil
}
